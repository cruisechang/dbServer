package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

//NewPartnerLogHandler returns PartnerLogHandler structure
func NewPartnerLogHandler(base baseHandler) *PartnerLogHandler {
	return &PartnerLogHandler{
		baseHandler: base,
	}
}

//PartnerLogHandler does select and insert
type PartnerLogHandler struct {
	baseHandler
}

func (h *PartnerLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "partnerLog"

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

func (h *PartnerLogHandler) getQueryStringArgs(param *timeParam) (queryString string, queryArgs []interface{}) {

	queryString = "SELECT partner_log.log_id,partner_log.partner_id,partner_log.category,partner_log.create_date ,partner.account,partner.name from partner_log LEFT JOIN partner on partner_log.partner_id=partner.partner_id WHERE partner_log.partner_id = ? AND partner_log.create_date BETWEEN ? AND ?"

	queryArgs = append(queryArgs, param.BeginDate)
	queryArgs = append(queryArgs, param.EndDate)
	return
}

func (h *PartnerLogHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	args, ok := param.([]interface{})
	if !ok {
		return nil, errors.New("args error")
	}

	return stmt.Query(IDOrAccount, args[0], args[1])
}
func (h *PartnerLogHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		resData := []partnerLogData{}
		for rows.Next() {
			ud := partnerLogDB{}
			err := rows.Scan(&ud.log_id, &ud.partner_id, &ud.category, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count += 1
				resData = append(resData,
					partnerLogData{
						ud.log_id,
						ud.partner_id,
						ud.account,
						ud.name,
						ud.category,
						ud.create_date})
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
