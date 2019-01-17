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

//NewOfficialCMSManagerIDHandler returns structure of official cms manager
func NewOfficialCMSManagerIDHandler(base baseHandler) *OfficialCMSManagerIDHandler {
	return &OfficialCMSManagerIDHandler{
		baseHandler: base,
	}
}

//OfficialCMSManagerIDHandler handles request for official cms manager
type OfficialCMSManagerIDHandler struct {
	baseHandler
}

func (h *OfficialCMSManagerIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSManagerIDHandler"

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

		queryString := "SELECT manager_id,account,role_id,active, login,create_date FROM official_cms_manager where manager_id = ? LIMIT 1"
		//h.getTargetRow(w, r, logPrefix, ID, queryString, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM official_cms_manager  where manager_id = ? LIMIT 1"
		h.delete(w, r, logPrefix, ID, queryString, h.returnIDResData)
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
		//if !strings.Contains(r.URL.Path, "patch") {
		//	h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler patch path error :%s", logPrefix, r.URL.Path))
		//	h.writeError(w, http.StatusOK, CodeRequestPathError, "")
		//	return
		//}

		//umarshal request body
		patchData, err := h.getPatchData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "UPDATE official_cms_manager set password = ? , role_id =? , active =?  WHERE manager_id = ? LIMIT 1"

		h.patch(w, r, logPrefix, ID, queryString, patchData, h.patchExec, h.returnIDResData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *OfficialCMSManagerIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *OfficialCMSManagerIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		d := officialCMSManagerDB{}
		resData := []officialCMSManagerData{}

		for rows.Next() {
			err := rows.Scan(&d.manager_id, &d.account, &d.role_id, &d.active, &d.login, &d.create_date)
			if err == nil {
				count++
				resData = append(resData,
					officialCMSManagerData{
						d.manager_id,
						d.account,
						d.active,
						d.role_id,
						d.login,
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

/*
func (h *officialCMSManagerIDHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		d := officialCMSManagerDB{}
		resData := []officialCMSManagerData{}

		for rows.Next() {
			err := rows.Scan(&d.manager_id, &d.account, &d.role_id, &d.active, &d.login, &d.create_date)
			if err == nil {
				count += 1
				resData = append(resData,
					officialCMSManagerData{
						d.manager_id,
						d.account,
						d.active,
						d.role_id,
						d.login,
						d.create_date})
			}
		}
		return resData, count
	}
}
*/
func (h *OfficialCMSManagerIDHandler) getPatchData(body []byte) (interface{}, error) {

	d := &officialCMSManagerPatchParam{}
	err := json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}
	if len(d.Password) == 0 {
		return nil, errors.New("patch data password unmarshal error")
	} else if d.RoleID < 0 {
		return nil, errors.New("patch data roleID unmarshal error")
	} else if d.Active < 0 {
		return nil, errors.New("patch data active unmarshal error")

	}

	return d, nil

}

func (h *OfficialCMSManagerIDHandler) patchExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	//檢查參數是否合法
	if p, ok := param.(*officialCMSManagerPatchParam); ok {

		return stmt.Exec(p.Password, p.RoleID, p.Active, ID)
	}
	return nil, errors.New("parsing param error")

}
func (h *OfficialCMSManagerIDHandler) returnIDResData(ID uint64) interface{} {

	return []managerIDData{
		{
			uint(ID),
		},
	}
}

func (h *OfficialCMSManagerIDHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSManagerPostParam); ok {

		return stmt.Exec(p.Account, p.Password, p.RoleID)
	}
	return nil, errors.New("")

}
func (h *OfficialCMSManagerIDHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []managerIDData{
		{
			uint(lastID),
		},
	}
}
