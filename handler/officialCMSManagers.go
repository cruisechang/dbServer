package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbex"
)

//NewOfficialCMSManagersHandler returns OfficialCMSManagersHandler structure
func NewOfficialCMSManagersHandler(base baseHandler) *OfficialCMSManagersHandler {
	return &OfficialCMSManagersHandler{
		baseHandler: base,
	}
}

//OfficialCMSManagersHandler select and insert new data
type OfficialCMSManagersHandler struct {
	baseHandler
}

func (h *OfficialCMSManagersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSManagersHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT manager_id,account,active, role_id,login,create_date FROM official_cms_manager "
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

		param := &officialCMSManagerPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.Account) < 5 || len(param.Password) < 5 || param.RoleID == 0 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		queryString := "INSERT  INTO official_cms_manager (account,password,role_id ) values (? ,?,?)"
		h.dbExec(w, r, logPrefix, 0, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *OfficialCMSManagersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *OfficialCMSManagersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		d := officialCMSManagerDB{}
		resData := []officialCMSManagerData{}

		for rows.Next() {
			err := rows.Scan(&d.manager_id, &d.account, &d.active, &d.role_id, &d.login, &d.create_date)
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

//post
func (h *OfficialCMSManagersHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSManagerPostParam); ok {

		return stmt.Exec(p.Account, p.Password, p.RoleID)
	}
	return nil, errors.New("")

}

//id 自動產生
func (h *OfficialCMSManagersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*managerIDData{{}},
		}

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecLastIDError,
			Count:   0,
			Message: "",
			Data:    []*managerIDData{{}},
		}
	}

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*managerIDData{
			{
				uint(lastID),
			},
		},
	}
}


func (h *OfficialCMSManagersHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSManagerPostParam); ok {

		return stmt.Exec(p.Account, p.Password, p.RoleID)
	}
	return nil, errors.New("")

}
func (h *OfficialCMSManagersHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []managerIDData{
		{
			uint(lastID),
		},
	}
}
