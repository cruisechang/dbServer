package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)
//NewOfficialCMSManagerAccountHandler returns OfficialCMSManagerAccountHandler structure
func NewOfficialCMSManagerAccountHandler(base baseHandler) *OfficialCMSManagerAccountHandler {
	return &OfficialCMSManagerAccountHandler{
		baseHandler: base,
	}
}

//OfficialCMSManagerAccountHandler selects data from DB by account
type OfficialCMSManagerAccountHandler struct {
	baseHandler
}

func (h *OfficialCMSManagerAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "officialCMSManagerAccountHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get account not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if len(account) < 4 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s account length =%d ", logPrefix, len(account)))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "login") {
			queryString := "SELECT manager_id,password, active  FROM official_cms_manager where account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "login", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			return
		}
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *OfficialCMSManagerAccountHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *OfficialCMSManagerAccountHandler) returnTargetColumnResponseData() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "login":
			resData := []checkLoginData{}
			count := 0
			var managerID uint
			var password string
			var active uint
			for rows.Next() {
				err := rows.Scan(&managerID, &password, &active)
				if err == nil {
					count++
					resData = append(resData,
						checkLoginData{
							managerID,
							password,
							active,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		default:
			return &responseData{}
		}
	}
}
