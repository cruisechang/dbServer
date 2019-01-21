package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbex"
)

//NewLimitationsHandler returns LimitationsHandler structure
func NewLimitationsHandler(base baseHandler) *LimitationsHandler {
	return &LimitationsHandler{
		baseHandler: base,
	}
}

//LimitationsHandler selects limitation data  from db
type LimitationsHandler struct {
	baseHandler
}

func (h *LimitationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "limitationsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT limitation_id,limitation FROM limitation "
		h.dbQuery(w, r, logPrefix, 0, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *LimitationsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *LimitationsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := limitationDB{}
		resData := []limitationData{}

		for rows.Next() {
			err := rows.Scan(&ud.limitation_id, &ud.limitation)
			if err == nil {
				count++
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
