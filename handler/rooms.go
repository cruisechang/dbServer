package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbex"
)

//NewRoomsHandler returns RoomsHandler structure
func NewRoomsHandler(base baseHandler) *RoomsHandler {
	return &RoomsHandler{
		baseHandler: base,
	}
}

//RoomsHandler does select and insert
type RoomsHandler struct {
	baseHandler
}

func (h *RoomsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "roomsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT room_id,hall_id,name,room_type,active,hls_url,boot,round_id,status,bet_countdown,dealer_id,limitation_id ,create_date FROM room "
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

		param := &roomPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.RoomID < 1 || param.HallID < 1 || param.RoomType < 0 || param.BetCountdown < 1 || param.LimitationID < 0 || len(param.Name) < 3 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}

		queryString := "INSERT  INTO room (room_id,hall_id,name,room_type,hls_url,bet_countdown,dealer_id,limitation_id) values (? ,?,?,?,?,?,?,?)"
		h.dbExec(w, r, logPrefix, param.RoomID, "", queryString, param, h.sqlPost, h.returnPostResponseData)

		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//get
func (h *RoomsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *RoomsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := roomDB{}
		resData := []roomData{}

		for rows.Next() {
			err := rows.Scan(&ud.room_id, &ud.hall_id, &ud.name, &ud.room_type, &ud.active, &ud.hls_url, &ud.boot, &ud.round_id, &ud.status, &ud.bet_countdown, &ud.dealer_id, &ud.limitation_id, &ud.create_date)
			if err == nil {
				count++
				resData = append(resData,
					roomData{
						ud.room_id,
						ud.hall_id,
						ud.name,
						ud.room_type,
						ud.active,
						ud.hls_url,
						ud.boot,
						ud.round_id,
						ud.status,
						ud.bet_countdown,
						ud.dealer_id,
						ud.limitation_id,
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

//post
func (h *RoomsHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*roomPostParam); ok {

		return stmt.Exec(p.RoomID, p.HallID, p.Name, p.RoomType, p.HLSURL, p.BetCountdown, 1, p.LimitationID)
	}
	return nil, errors.New("")

}

//id預先產生
func (h *RoomsHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*roomIDData{{}},
		}

	}

	if id, ok := IDOrAccount.(uint); ok {
		return &responseData{
			Code:    CodeSuccess,
			Count:   int(affRow),
			Message: "",
			Data: []*roomIDData{
				{
					id,
				},
			},
		}
	}

	//error
	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data:    []*roomIDData{{}},
	}
}
