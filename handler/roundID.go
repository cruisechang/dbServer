package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

//NewRoundIDHandler returns a handler
func NewRoundIDHandler(base baseHandler) *RoundIDHandler {
	return &RoundIDHandler{
		baseHandler: base,
	}
}

type RoundIDHandler struct {
	baseHandler
}

func (h *RoundIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "roundIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s ", logPrefix))

	vars := mux.Vars(r)
	var ID uint64
	mid, ok := vars["id"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	ID, err := strconv.ParseUint(mid, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID < 1 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {
		queryString := "select round.round_id,round.hall_id,round.room_id,round.room_type,round.brief,round.record,round.status,round.create_date, round.end_date,room.name from round LEFT JOIN room on round.room_id=room.room_id where round.round_id= ? "
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s patch checkBody error=%s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		//check patch
		if !strings.Contains(r.URL.Path, "patch") {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch path error :%s", logPrefix, r.URL.Path))
			h.writeError(w, http.StatusOK, CodeRequestPathError, "")
			return
		}

		//unmarshal request body
		param, err := h.getPatchData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "UPDATE round set brief = ? , record =? , status =?  WHERE round_id = ? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *RoundIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *RoundIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := roundDB{}
		resData := []roundData{}

		for rows.Next() {
			err := rows.Scan(&ud.round_id, &ud.hall_id, &ud.room_id, &ud.room_type, &ud.brief, &ud.record, &ud.status, &ud.create_date, &ud.end_datea, &ud.name)
			if err == nil {
				count++
				resData = append(resData,
					roundData{
						ud.round_id,
						ud.hall_id,
						ud.room_id,
						ud.room_type,
						ud.brief,
						ud.record,
						ud.status,
						ud.create_date,
						ud.end_datea,
						ud.name,
					})
			}
		}

		return &responseData{
			Code:    CodeSuccess,
			Count:   count,
			Message: "",
			Data:    resData,
		}
	}
}

//patch
func (h *RoundIDHandler) getPatchData(body []byte) (interface{}, error) {

	d := &roundPatchParam{}
	err := json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}
	if len(d.Record) == 0 {
		return nil, errors.New("patch data result unmarshal error")
	} else if len(d.Brief) == 0 {
		return nil, errors.New("patch data brief unmarshal error")
	} else if d.Status < 0 {
		return nil, errors.New("patch data status unmarshal error")

	}

	return d, nil

}
func (h *RoundIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	//檢查參數是否合法
	if p, ok := param.(*roundPatchParam); ok {

		return stmt.Exec(p.Brief, p.Record, p.Status, IDOrAccount)
	}
	return nil, errors.New("parsing param error")
}

func (h *RoundIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*roundIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint64)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*roundIDData{
			{
				ID,
			},
		},
	}
}
