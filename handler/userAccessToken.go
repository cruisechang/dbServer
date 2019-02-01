package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

//NewUserAccessTokenHandler returns NewUserAccessTokenHandler structure
func NewUserAccessTokenHandler(base baseHandler) *UserAccessTokenHandler {
	return &UserAccessTokenHandler{
		baseHandler: base,
	}
}

//UserAccessTokenHandler presents structure of user account handler
//user使用access token 登入，查詢db是否有此access token
type UserAccessTokenHandler struct {
	baseHandler
}

func (h *UserAccessTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "userAccessTokenHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	token, ok := vars["accessToken"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  accessToken variable not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if len(token) < 5 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  accessToken variable length =%d ", logPrefix, len(token)))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "tokenData") {

			queryString := "SELECT user_id, account , credit, name, partner_id, active, access_token_expire from user WHERE access_token = ? LIMIT 1"
			param := &userAccessTokenGetParam{token}
			h.dbQuery(w, r, logPrefix, 0, "", queryString, param, h.sqlQuery, h.returnResDataFunc)
			return

		}
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get path error path=%s ", logPrefix, r.RequestURI))
		h.writeError(w, http.StatusOK, CodeRequestPathError, "")
		return

	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *UserAccessTokenHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	//檢查參數是否合法
	if p, ok := param.(*userAccessTokenGetParam); ok {

		return stmt.Query(p.AccessToken)
	}
	return nil, errors.New("parsing param error")

}

func (h *UserAccessTokenHandler) returnResDataFunc() func(IDOrAccount interface{}, column string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, column string, rows *sql.Rows) *responseData {
		count := 0
		d := userAccessTokenDB{}
		resData := []userAccessTokenGetData{}

		for rows.Next() {
			err := rows.Scan(&d.user_id, &d.account, &d.credit, &d.name, &d.partner_id, &d.active, &d.access_token_expire)
			if err == nil {
				count++
				resData = append(resData,
					userAccessTokenGetData{
						d.user_id,
						d.account,
						d.credit,
						d.name,
						d.partner_id,
						d.active,
						d.access_token_expire})
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
