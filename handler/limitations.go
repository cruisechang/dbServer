package handler

import (
	"net/http"
	"database/sql"
	"github.com/cruisechang/dbex"
	"fmt"
)

func NewLimitationsHandler(base baseHandler) *limitationsHandler {
	return &limitationsHandler{
		baseHandler: base,
	}
}

type limitationsHandler struct {
	baseHandler
}

func (h *limitationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "limitationsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT limitation_id,limitation FROM limitation "
		//h.get(w, r,logPrefix,queryString,h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *limitationsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *limitationsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := limitationDB{}
		resData := []limitationData{}

		for rows.Next() {
			err := rows.Scan(&ud.limitation_id, &ud.limitation)
			if err == nil {
				count ++
				resData = append(resData,
					limitationData{
						ud.limitation_id,
						ud.limitation})
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
func (h *limitationsHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		ud := limitationDB{}
		resData := []limitationData{}

		for rows.Next() {
			err := rows.Scan(&ud.limitation_id, &ud.limitation)
			if err == nil {
				count += 1
				resData = append(resData,
					limitationData{
						ud.limitation_id,
						ud.limitation})
			}
		}
		return resData, count
	}
}
*/
