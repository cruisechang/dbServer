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

//NewOfficialCMSRoleIDHandler returns official cms role structure
func NewOfficialCMSRoleIDHandler(base baseHandler) *OfficialCMSRoleIDHandler {
	return &OfficialCMSRoleIDHandler{
		baseHandler: base,
	}
}
//OfficialCMSRoleIDHandler handles request of official cms role
type OfficialCMSRoleIDHandler struct {
	baseHandler
}

func (h *OfficialCMSRoleIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSRoleIDHandler"

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
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", mid, logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {

		queryString := "SELECT role_id,permission,create_date FROM official_cms_role where role_id = ? LIMIT 1"
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	//if r.Method == "DELETE" || r.Method == "delete" {
	//	queryString := "DELETE FROM official_cms_role  where role_id = ? LIMIT 1"
	//	h.delete(w, r, logPrefix, ID, queryString, h.returnIDResData)
	//	return
	//}

	if r.Method == "PATCH" || r.Method == "patch" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s patch checkBody error=%s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		//unmarshal request body
		param, err := h.getPatchData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "UPDATE official_cms_role set permission = ?   WHERE role_id = ? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *OfficialCMSRoleIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *OfficialCMSRoleIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		d := officialCMSRoleDB{}
		resData := []officialCMSRoleData{}

		for rows.Next() {
			err := rows.Scan(&d.role_id, &d.permission, &d.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					officialCMSRoleData{
						d.role_id,
						d.permission,
						d.create_date})
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
func (h *OfficialCMSRoleIDHandler) getPatchData(body []byte) (interface{}, error) {

	d := &officialCMSRolePatchParam{}
	err := json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}

	//param.Permission is json string []int
	pa := &[]int{}
	err = json.Unmarshal([]byte(d.Permission), pa)
	if err != nil {
		return nil, err
	}


	if len(d.Permission) <2 {
		return nil, errors.New("patch data password unmarshal error")
	}

	return d, nil

}
func (h *OfficialCMSRoleIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSRolePatchParam); ok {

		return stmt.Exec(p.Permission, IDOrAccount)
	}
	return nil, errors.New("parsing param error")
}

func (h *OfficialCMSRoleIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*roleIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*roleIDData{
			{
				ID,
			},
		},
	}
}


