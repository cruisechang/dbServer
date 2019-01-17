package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"strings"
	"errors"
	"database/sql"
)

func NewBetIDHandler(base baseHandler) *betIDHandler {
	return &betIDHandler{
		baseHandler: base,
	}
}

type betIDHandler struct {
	baseHandler
}

func (h *betIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "betIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	var id uint64
	mid, ok := vars["id"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get id not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}
	id, err := strconv.ParseUint(mid, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if id < 1 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {

		queryString := "select bet.bet_id,bet.partner_id,bet.user_id,bet.room_id,bet.room_type,bet.round_id,bet.seat_id,bet.bet_credit,bet.active_credit,bet.prize_credit,bet.result_credit,bet.balance_credit,bet.original_credit,bet.record,bet.status,bet.create_date, user.account, user.name from bet LEFT JOIN user on bet.user_id=user.user_id where bet.bet_id = ?"
		//h.getTargetRow(w, r, logPrefix, id, queryString, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, id, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {

		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s handler  patch checkBody error=%s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		columm := ""
		queryString := ""

		if strings.Contains(r.URL.Path, "status") {
			columm = "status"
		} else {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler patch path error :%s", logPrefix, r.URL.Path))
			h.writeError(w, http.StatusOK, CodeRequestPathError, "")
			return
		}
		queryString = "UPDATE bet set " + columm + " = ?  WHERE bet_id = ? LIMIT 1"

		//umarshal request body
		patchData, err := h.getPatchData(columm, body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}
		h.patch(w, r, logPrefix, id, queryString, patchData, h.patchExec, h.returnIDResData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *betIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *betIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := betDB{}
		resData := []betData{}

		for rows.Next() {
			err := rows.Scan(&ud.bet_id, &ud.partner_id, &ud.user_id, &ud.room_id, &ud.room_type, &ud.round_id, &ud.seat_id, &ud.bet_credit, &ud.active_credit, &ud.prize_credit, &ud.result_credit, &ud.balance_credit, &ud.original_credit, &ud.partner_id, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count ++
				resData = append(resData,
					betData{
						ud.bet_id,
						ud.partner_id,
						ud.user_id,
						ud.room_id,
						ud.room_type,
						ud.round_id,
						ud.seat_id,
						ud.bet_credit,
						ud.active_credit,
						ud.prize_credit,
						ud.result_credit,
						ud.balance_credit,
						ud.original_credit,
						ud.record,
						ud.status,
						ud.create_date,
						ud.account,
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

/*
func (h *betIDHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		ud := betDB{}
		resData := []betData{}

		for rows.Next() {
			err := rows.Scan(&ud.bet_id, &ud.partner_id, &ud.user_id, &ud.room_id, &ud.room_type, &ud.round_id, &ud.seat_id, &ud.bet_credit, &ud.active_credit, &ud.prize_credit, &ud.result_credit, &ud.balance_credit, &ud.original_credit, &ud.partner_id, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count ++
				resData = append(resData,
					betData{
						ud.bet_id,
						ud.partner_id,
						ud.user_id,
						ud.room_id,
						ud.room_type,
						ud.round_id,
						ud.seat_id,
						ud.bet_credit,
						ud.active_credit,
						ud.prize_credit,
						ud.result_credit,
						ud.balance_credit,
						ud.original_credit,
						ud.record,
						ud.status,
						ud.create_date,
						ud.account,
						ud.name,
					})
			}
		}

		return resData, count
	}
}
*/

func (h *betIDHandler) getPatchData(column string, body []byte) (interface{}, error) {
	switch(column) {
	case "status":
		dt := &statusData{}
		err := json.Unmarshal(body, dt)
		if err != nil {
			return nil, err
		}
		return dt, nil
	default:
		return nil, errors.New("column error")
	}
}
func (h *betIDHandler) patchExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	//檢查參數是否合法
	if p, ok := param.(*statusData); ok {

		return stmt.Exec(p.Status, ID)
	}
	return nil, errors.New("parsing param error")

}

func (h *betIDHandler) returnIDResData(ID uint64) interface{} {
	return []betIDData{
		{
			ID,
		},
	}
}
