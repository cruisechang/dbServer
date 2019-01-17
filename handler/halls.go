package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewHallsHandler(base baseHandler) *hallsHandler {
	return &hallsHandler{
		baseHandler: base,
	}
}

type hallsHandler struct {
	baseHandler
}

func (h *hallsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "hallsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT hall_id,name,active,create_date FROM hall "
		//h.get(w, r,"halls",queryString,h.returnResDataFunc)
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

		param := &hallPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.HallID < 0 || len(param.Name) < 3 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		queryString := "INSERT  INTO hall (hall_id,name) values (? ,?)"
		h.post(w, r, logPrefix, uint64(param.HallID), queryString, param, h.sqlExec, h.returnPostResData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *hallsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *hallsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := hallDB{}
		resData := []hallData{}

		for rows.Next() {
			err := rows.Scan(&ud.hall_id, &ud.name, &ud.active, &ud.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					hallData{
						ud.hall_id,
						ud.name,
						ud.active,
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

func (h *hallsHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*hallPostParam); ok {

		return stmt.Exec(p.HallID, p.Name)
	}
	return nil, errors.New("")

}
func (h *hallsHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []hallIDData{
		{
			uint(ID),
		},
	}
}
