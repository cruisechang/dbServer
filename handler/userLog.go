package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"time"
	"database/sql"
	"errors"
)

func NewUserLogHandler(base baseHandler) *userLogHandler {
	return &userLogHandler{
		baseHandler: base,
	}
}

type userLogHandler struct {
	baseHandler
}

func (h *userLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "userLog"

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

	vars := mux.Vars(r)
	mid, ok := vars["id"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get id not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	ID, err := strconv.ParseUint(mid, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//check time format
	param := &timeParam{}
	err = json.Unmarshal(body, param)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get data unmarshal error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
		return
	}

	_, err = time.Parse(timeFormat, param.BeginDate)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get data parse beginDate error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeRequestDataError, "parse beginDate error")
		return
	}

	_, err = time.Parse(timeFormat, param.EndDate)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s hander get data parse endDate error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeRequestDataError, "parse endDate error")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		queryString, queryArgs := h.getQueryStringArgs(param)
		h.dbQuery(w, r, logPrefix, ID, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
//query
func (h *userLogHandler) getQueryStringArgs(param *timeParam) (queryString string, queryArgs []interface{}) {

	queryString = "SELECT user_log.log_id,user_log.user_id,user_log.category,user_log.ip,user_log.platform,user_log.create_date ,user.account,user.name from user_log LEFT JOIN user on user_log.user_id=user.user_id WHERE user_log.user_id = ? AND user_log.create_date BETWEEN ? AND ?"

	queryArgs=append(queryArgs,param.BeginDate)
	queryArgs=append(queryArgs,param.EndDate)
	return
}
func (h *userLogHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	args, ok := param.([]interface{})
	if !ok {
		return nil, errors.New("args error")
	}

	return stmt.Query(IDOrAccount, args[0], args[1])
}
func (h *userLogHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		resData := []userLogData{}
		for rows.Next() {
			ud := userLogDB{}
			err := rows.Scan(&ud.log_id, &ud.user_id, &ud.category, &ud.ip, &ud.platform, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count ++
				resData = append(resData,
					userLogData{
						ud.log_id,
						ud.user_id,
						ud.account,
						ud.name,
						ud.category,
						ud.ip,
						ud.platform,
						ud.create_date,})
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


