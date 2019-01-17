package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewRoomsHandler(base baseHandler) *roomsHandler {
	return &roomsHandler{
		baseHandler: base,
	}
}

type roomsHandler struct {
	baseHandler
}

func (h *roomsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "roomsHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT room_id,hall_id,name,room_type,active,hls_url,boot,round_id,status,bet_countdown,dealer_id,limitation_id ,create_date FROM room "
		//h.get(w, r, "rooms", queryString, h.returnResDataFunc)
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
		h.post(w, r, logPrefix, uint64(param.RoomID), queryString, param, h.sqlExec, h.returnPostResData)

		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *roomsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *roomsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := roomDB{}
		resData := []roomData{}

		for rows.Next() {
			err := rows.Scan(&ud.room_id, &ud.hall_id, &ud.name, &ud.room_type, &ud.active, &ud.hls_url, &ud.boot, &ud.round_id, &ud.status, &ud.bet_countdown, &ud.dealer_id, &ud.limitation_id, &ud.create_date)
			if err == nil {
				count ++
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

/*
func (h *roomsHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		ud := roomDB{}
		resData := []roomData{}

		for rows.Next() {
			err := rows.Scan(&ud.room_id, &ud.hall_id, &ud.name, &ud.room_type, &ud.active, &ud.hls_url, &ud.boot, &ud.round_id, &ud.status, &ud.bet_countdown, &ud.dealer_id, &ud.limitation_id, &ud.create_date)
			if err == nil {
				count += 1
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
		return resData, count
	}
}
*/

func (h *roomsHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*roomPostParam); ok {

		return stmt.Exec(p.RoomID, p.HallID, p.Name, p.RoomType, p.HLSURL, p.BetCountdown,1, p.LimitationID)
	}
	return nil, errors.New("")

}
func (h *roomsHandler) returnPostResData(ID, lastID uint64) interface{} {
	return []roomIDData{
		{
			uint(ID),
		},
	}
}
