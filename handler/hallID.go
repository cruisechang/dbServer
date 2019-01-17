package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"database/sql"
)

func NewHallIDHandler(base baseHandler) *hallIDHandler {
	return &hallIDHandler{
		baseHandler: base,
	}
}

type hallIDHandler struct {
	baseHandler
}

func (h *hallIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "hallIDHandler"

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
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	id, err := strconv.ParseUint(mid, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if id == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT hall_id,name,active,create_date FROM hall where hall_id = ? LIMIT 1"
		//h.getTargetRow(w, r, logPrefix, id, queryString, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, id, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
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
		patchData, err := h.getPatchData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, err.Error())
			return
		}
		h.patch(w, r, logPrefix, id, queryString, patchData, h.patchExec, h.returnIDResData)
		return
	}

	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM hall  where hall_id = ? LIMIT 1"
		h.delete(w, r, logPrefix, id, queryString, h.returnIDResData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *hallIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *hallIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := hallDB{}
		resData := []hallData{}

		for rows.Next() {
			err := rows.Scan(&ud.hall_id, &ud.name, &ud.active, &ud.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					hallData{
						ud.hall_id,
						ud.name,
						ud.active,
						ud.create_date,})
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

//func (h *hallIDHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {
//
//	return func(rows *sql.Rows) (interface{}, int) {
//		count := 0
//		ud := hallDB{}
//		resData := []hallData{}
//
//		for rows.Next() {
//			err := rows.Scan(&ud.hall_id, &ud.name, &ud.active, &ud.create_date)
//			if err == nil {
//				count += 1
//				resData = append(resData,
//					hallData{
//						ud.hall_id,
//						ud.name,
//						ud.active,
//						ud.create_date,})
//			}
//		}
//		return resData, count
//	}
//}
func (h *hallIDHandler) getPatchData(body []byte) (interface{}, error) {
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
func (h *hallIDHandler) patchExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*hallPatchParam); ok {
		return stmt.Exec(p.HallID, p.Name, p.Active, ID)
	}

	return nil, errors.New("parsing param error")

}

func (h *hallIDHandler) returnIDResData(ID uint64) interface{} {
	return []hallIDData{
		{
			uint(ID),
		},
	}
}