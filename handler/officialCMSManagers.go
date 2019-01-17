package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewOfficialCMSManagersHandler(base baseHandler) *officialCMSManagersHandler {
	return &officialCMSManagersHandler{
		baseHandler: base,
	}
}

type officialCMSManagersHandler struct {
	baseHandler
}

func (h *officialCMSManagersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSManagersHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT manager_id,account,active, role_id,login,create_date FROM official_cms_manager "
		//h.get(w, r, logPrefix, queryString, h.returnResDataFunc)
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
		h.post(w, r, "dealers", 0, queryString, param, h.sqlExec, h.returnPostResData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *officialCMSManagersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *officialCMSManagersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		d := officialCMSManagerDB{}
		resData := []officialCMSManagerData{}

		for rows.Next() {
			err := rows.Scan(&d.manager_id, &d.account, &d.active, &d.role_id, &d.login, &d.create_date)
			if err == nil {
				count ++
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

//func (h *officialCMSManagersHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {
//
//	return func(rows *sql.Rows) (interface{}, int) {
//		count := 0
//		d := officialCMSManagerDB{}
//		resData := []officialCMSManagerData{}
//
//		for rows.Next() {
//			err := rows.Scan(&d.manager_id, &d.account, &d.active, &d.role_id, &d.login, &d.create_date)
//			if err == nil {
//				count += 1
//				resData = append(resData,
//					officialCMSManagerData{
//						d.manager_id,
//						d.account,
//						d.active,
//						d.role_id,
//						d.login,
//						d.create_date})
//			}
//		}
//		return resData, count
//	}
//}

func (h *officialCMSManagersHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*officialCMSManagerPostParam); ok {

		return stmt.Exec(p.Account, p.Password, p.RoleID)
	}
	return nil, errors.New("")

}
func (h *officialCMSManagersHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []managerIDData{
		{
			uint(lastID),
		},
	}
}