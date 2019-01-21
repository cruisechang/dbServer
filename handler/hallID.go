package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)
//NewHallIDHandler returns HallIDHandler structure
func NewHallIDHandler(base baseHandler) *HallIDHandler {
	return &HallIDHandler{
		baseHandler: base,
	}
}

//HallIDHandler do select, delete, patch by ID
type HallIDHandler struct {
	baseHandler
}

func (h *HallIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "hallIDHandler"

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
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT hall_id,name,active,create_date FROM hall where hall_id = ? LIMIT 1"
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM hall  where hall_id = ? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, nil, h.sqlDelete, h.returnExecResponseData)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s patch checkBody error=%s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		queryString := "UPDATE hall SET hall_id=?, name=?, active= ? WHERE hall_id = ? LIMIT 1"

		//unmarshal request body
		param, err := h.getPatchData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, err.Error())
			return
		}
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *HallIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *HallIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := hallDB{}
		resData := []hallData{}

		for rows.Next() {
			err := rows.Scan(&ud.hall_id, &ud.name, &ud.active, &ud.create_date)
			if err == nil {
				count++
				resData = append(resData,
					hallData{
						ud.hall_id,
						ud.name,
						ud.active,
						ud.create_date})
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

//delete
func (h *HallIDHandler) sqlDelete(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	return stmt.Exec(IDOrAccount)
}

//patch
func (h *HallIDHandler) getPatchData(body []byte) (interface{}, error) {
	d := &hallPatchParam{}
	err := json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}
	if len(d.Name) < 3 {
		return nil, errors.New("patch param error")
	}
	return d, nil
}

func (h *HallIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*hallPatchParam); ok {
		return stmt.Exec(p.HallID, p.Name, p.Active, IDOrAccount)
	}

	return nil, errors.New("parsing param error")
}

func (h *HallIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*hallIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*hallIDData{
			{
				ID,
			},
		},
	}
}

