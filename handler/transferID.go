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

func NewTransferIDHandler(base baseHandler) *transferIDHandler {
	return &transferIDHandler{
		baseHandler: base,
	}
}

type transferIDHandler struct {
	baseHandler
}

func (h *transferIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "transferIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

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

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {

		queryString := "select transfer.transfer_id, transfer.partner_transfer_id, transfer.partner_id, transfer.user_id, transfer.category, transfer.transfer_credit, transfer.credit, transfer.status, transfer.create_date, user.account, user.name from transfer LEFT JOIN user on transfer.user_id=user.user_id  WHERE transfer.transfer_id = ? "
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

		column := ""
		queryString := ""

		if strings.Contains(r.URL.Path, "status") {
			column = "status"
		} else {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch path error :%s", logPrefix, r.URL.Path))
			h.writeError(w, http.StatusOK, CodeRequestPathError, "")
			return
		}
		queryString = "UPDATE transfer set " + column + " = ?  WHERE transfer_id = ? LIMIT 1"

		//unmarshal request body
		param, err := h.getPatchData(column, body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}
		//h.patch(w, r, logPrefix, id, queryString, patchData, h.patchExec, h.returnIDResData)
		h.dbExec(w, r, logPrefix, ID, column, queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *transferIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *transferIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := transferDB{}
		resData := []transferData{}

		for rows.Next() {
			err := rows.Scan(&ud.transfer_id, &ud.partner_transfer_id, &ud.partner_id, &ud.user_id, &ud.category, &ud.transfer_credit, &ud.credit, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count++
				resData = append(resData,
					transferData{
						ud.transfer_id,
						ud.partner_transfer_id,
						ud.partner_id,
						ud.user_id,
						ud.category,
						ud.transfer_credit,
						ud.credit,
						ud.status,
						ud.create_date,
						ud.account,
						ud.name})
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

func (h *transferIDHandler) getPatchData(column string, body []byte) (interface{}, error) {
	switch column {
	case "status":
		ug := &statusData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	default:
		return nil, errors.New("column error")
	}
}

//patch
func (h *transferIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*statusData); ok {

		return stmt.Exec(p.Status, IDOrAccount)
	}
	return nil, errors.New("parsing param error")
}

func (h *transferIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*transferIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint64)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*transferIDData{
			{
				ID,
			},
		},
	}
}
