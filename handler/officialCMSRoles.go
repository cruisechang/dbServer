package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewOfficialCMSRolesHandler(base baseHandler) *OfficialCMSRolesHandler {
	return &OfficialCMSRolesHandler{
		baseHandler: base,
	}
}

type OfficialCMSRolesHandler struct {
	baseHandler
}

func (h *OfficialCMSRolesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSRolesHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT role_id, permission, create_date FROM official_cms_role"
		h.dbQuery(w, r, logPrefix, 0, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		//h.get(w, r, logPrefix, queryString, h.returnResDataFunc)
		return
	}
	if r.Method == "POST" || r.Method == "post" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, errMsg)
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		postData,err:=h.getPostData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "INSERT  INTO official_cms_role ( permission ) values (?)"
		h.post(w, r, logPrefix, 0, queryString, postData, h.sqlExec, h.returnPostResData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *OfficialCMSRolesHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *OfficialCMSRolesHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

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

/*
func (h *officialCMSRolesHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		d := officialCMSRoleDB{}
		resData := []officialCMSRoleData{}

		for rows.Next() {
			err := rows.Scan(&d.role_id, &d.permission, &d.create_date)
			if err == nil {
				count += 1
				resData = append(resData,
					officialCMSRoleData{
						d.role_id,
						d.permission,
						d.create_date})
			}
		}
		return resData, count
	}
}
*/

func (h *OfficialCMSRolesHandler) getPostData(body []byte) (interface{}, error) {

	d := &officialCMSRolePostParam{}
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

func (h *OfficialCMSRolesHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSRolePostParam); ok {

		return stmt.Exec(p.Permission)
	}
	return nil, errors.New("")

}
func (h *OfficialCMSRolesHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []roleIDData{
		{
			uint(lastID),
		},
	}
}
