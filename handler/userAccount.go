package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
	uuid "github.com/satori/go.uuid"
)

//NewUserAccountHandler returns userAccountHandler structure
func NewUserAccountHandler(base baseHandler) *userAccountHandler {
	return &userAccountHandler{
		baseHandler: base,
	}
}

//UserAccountHandler does mysql select by account
type userAccountHandler struct {
	baseHandler
}

func (h *userAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "userAccountHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	body, errCode, errMsg := h.checkBody(w, r)
	if errCode != CodeSuccess {
		h.logger.Log(dbex.LevelError, errMsg)
		h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
		return
	}

	param := &userAccountGetParam{}
	err := json.Unmarshal(body, param)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get data unmarshal error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
		return
	}

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get account not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if len(account) < 4 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  account length =%d ", logPrefix, len(account)))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "password") {

			queryString := "SELECT password from user WHERE partner_id = ? AND account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "password", queryString, param, h.sqlQuery, h.returnTargetColumnResponseData)
			return

		} else if strings.Contains(r.URL.Path, "id") {
			queryString := "SELECT user_id from user WHERE partner_id = ? AND account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "user_id", queryString, param, h.sqlQuery, h.returnTargetColumnResponseData)
			return

		} else if strings.Contains(r.URL.Path, "accessToken") {
			h.handleAccessToken(w, r, logPrefix, account, param)
			return

		} else {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get path error path=%s ", logPrefix, r.RequestURI))
			h.writeError(w, http.StatusOK, CodeRequestPathError, "")
			return
		}
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *userAccountHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	if p, ok := param.(*userAccountGetParam); ok {
		return stmt.Query(p.PartnerID, IDOrAccount)
	}
	return nil, errors.New("assertion error")
}
func (h *userAccountHandler) returnTargetColumnResponseData() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "password":
			resData := []passwordData{}
			count := 0
			var password string
			for rows.Next() {
				err := rows.Scan(&password)
				if err == nil {
					count += 1
					resData = append(resData,
						passwordData{
							password,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "user_id":
			resData := []userIDData{}
			count := 0
			var id uint64
			for rows.Next() {
				err := rows.Scan(&id)
				if err == nil {
					count += 1
					resData = append(resData,
						userIDData{
							id,
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

func (h *userAccountHandler) handleAccessToken(w http.ResponseWriter, r *http.Request, logPrefix string, account string, param *userAccountGetParam) {
	count := 0
	active := -1
	token := ""
	code := CodeSuccess
	found := 0

	queryString := "SELECT active from user WHERE partner_id = ? AND account = ? AND password = ? "
	sqlDB := h.db.GetSQLDB()
	row := sqlDB.QueryRow(queryString, param.PartnerID, account, param.Password)
	row.Scan(&active)

	//found
	if active != -1 {
		found = 1

		//啟用
		if active == 1 {
			count = 1
			//產生token
			u, _ := uuid.NewV4()
			token = u.String()

		} else {
			active = 0
		}

	}

	//有找照，且active==1
	if active == 1 {
		updateString := "UPDATE user set access_token = ? , access_token_expire =?   WHERE partner_id =? AND account  = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(updateString)
		defer stmt.Close()
		t := time.Now().Add(time.Duration(10) * time.Minute)
		stmt.Exec(token, t.Format("2006-01-02 15:04:05"), param.PartnerID, account)
	}

	rd := responseData{
		Code:    code,
		Count:   count,
		Message: "",
		Data: []userAccessTokenData{{
			found,
			active,
			token}},
	}
	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  exec res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s  exec response data=%s", logPrefix, resStr))
}
