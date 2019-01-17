package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewBroadcastsHandler(base baseHandler) *BroadcastsHandler {
	return &BroadcastsHandler{
		baseHandler: base,
	}
}

type BroadcastsHandler struct {
	baseHandler
}

func (h *BroadcastsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "BroadcastsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT broadcast_id, content, internal, repeat_times, active, create_date  FROM broadcast "
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

		param := &broadcastPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.Internal < 1 || param.RepeatTimes < 1 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		queryString := "INSERT  INTO broadcast (content,internal,repeat_times,active) values (? ,?, ?, ?)"
		h.dbExec(w, r, logPrefix, 0, "", queryString, param, h.sqlExec, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//get
func (h *BroadcastsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *BroadcastsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := broadcastDB{}
		resData := []broadcastData{}

		for rows.Next() {
			err := rows.Scan(&ud.broadcast_id, &ud.content, &ud.internal, &ud.repeat_times, &ud.active, &ud.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					broadcastData{
						ud.broadcast_id,
						ud.content,
						ud.internal,
						ud.repeat_times,
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

//post
func (h *BroadcastsHandler) sqlExec(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*broadcastPostParam); ok {

		return stmt.Exec(p.Content, p.Internal, p.RepeatTimes, p.Active)
	}
	return nil, errors.New("")

}
func (h *BroadcastsHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []broadcastIDData{{}},
		}
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecLastIDError,
			Count:   0,
			Message: "",
			Data:    []broadcastIDData{{}},
		}
	}

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []broadcastIDData{
			{
				uint64(lastID),
			},
		},
	}
}
