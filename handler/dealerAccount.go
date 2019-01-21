package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

//NewDealerAccountHandler returns DealerAccountHandler structure
func NewDealerAccountHandler(base baseHandler) *DealerAccountHandler {
	return &DealerAccountHandler{
		baseHandler: base,
	}
}

//DealerAccountHandler selects dealer value by account
type DealerAccountHandler struct {
	baseHandler
}

func (h *DealerAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "dealerAccount"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get account not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if len(account) < 3 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s path error %s ", logPrefix, r.RequestURI))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "password") {
			queryString := "SELECT password  FROM dealer where account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "password", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			return
		}
		if strings.Contains(r.URL.Path, "id") {
			queryString := "SELECT dealer_id  FROM dealer where account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "dealer_id", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			return
		}

		if strings.Contains(r.URL.Path, "login") {

			body, errCode, errMsg := h.checkBody(w, r)
			if errCode != CodeSuccess {
				h.logger.Log(dbex.LevelError, errMsg)
				h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
				return
			}

			param := &dealerAccountGetParam{}
			err := json.Unmarshal(body, param)
			if err != nil {
				h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get data unmarshal error=%s", logPrefix, err.Error()))
				h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
				return
			}
			h.handleLogin(w, r, logPrefix, account, param)

			return
		}
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *DealerAccountHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}

func (h *DealerAccountHandler) returnTargetColumnResponseData() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "password":
			resData := []passwordData{}
			count := 0
			var password string
			for rows.Next() {
				err := rows.Scan(&password)
				if err == nil {
					count ++
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
		case "dealer_id":
			resData := []dealerIDData{}
			count := 0
			var id uint
			for rows.Next() {
				err := rows.Scan(&id)
				if err == nil {
					count ++
					resData = append(resData,
						dealerIDData{
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

func (h *DealerAccountHandler) handleLogin(w http.ResponseWriter, r *http.Request, logPrefix string, account string, param *dealerAccountGetParam) {
	count := 0
	active := -1
	dealer_id := -1  //這個是db column name 所以用底線
	code := CodeSuccess

	queryString := "SELECT dealer_id ,active from dealer WHERE  account = ? AND password = ? "

	sqlDB := h.db.GetSQLDB()
	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  handleLogin sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, fmt.Sprintf("%s  handleLogin sqlDB prepare error=%s", logPrefix, err.Error()))
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(account, param.Password)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  handleLogin sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, fmt.Sprintf("%s  handleLogin sqlDB query error=%s", logPrefix, err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dealer_id, &active)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  handleLogin sqlDB query error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeDBScanError, fmt.Sprintf("%s  handleLogin sqlDB scan error=%s", logPrefix, err.Error()))
			return
		}
		break
	}

	//found
	if active != -1 && dealer_id != -1 {

		//啟用
		if active == 1 {
			count = 1
		}
	} else {
		active = 0
		dealer_id = 0
	}

	rd := responseData{
		Code:    code,
		Count:   count,
		Message: "",
		Data: []dealerLoginData{{
			dealer_id,
			active,
		}},
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
