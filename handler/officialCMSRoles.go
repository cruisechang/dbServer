package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbex"
)
//NewOfficialCMSRolesHandler returns OfficialCMSRolesHandler strurcture
func NewOfficialCMSRolesHandler(base baseHandler) *OfficialCMSRolesHandler {
	return &OfficialCMSRolesHandler{
		baseHandler: base,
	}
}
//OfficialCMSRolesHandler select and insert new data
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
		return
	}
	if r.Method == "POST" || r.Method == "post" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, errMsg)
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		param, err := h.getPostData(body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "INSERT  INTO official_cms_role ( permission ) values (?)"
		h.dbExec(w, r, logPrefix, 0, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//get
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
				count++
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

//post
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

	if len(d.Permission) < 2 {
		return nil, errors.New("patch data password unmarshal error")
	}

	return d, nil

}
func (h *OfficialCMSRolesHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSRolePostParam); ok {

		return stmt.Exec(p.Permission)
	}
	return nil, errors.New("")

}

//id 自動產生
func (h *OfficialCMSRolesHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*roleIDData{{}},
		}

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecLastIDError,
			Count:   0,
			Message: "",
			Data:    []*roleIDData{{}},
		}
	}

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*roleIDData{
			{
				uint(lastID),
			},
		},
	}
}

//
//func (h *OfficialCMSRolesHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {
//
//	if p, ok := param.(*officialCMSRolePostParam); ok {
//
//		return stmt.Exec(p.Permission)
//	}
//	return nil, errors.New("")
//
//}
//func (h *OfficialCMSRolesHandler) returnPostResData(ID, lastID uint64) interface{} {
//	return []roleIDData{
//		{
//			uint(lastID),
//		},
//	}
//}
