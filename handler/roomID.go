package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)
//NewRoomIDHandler returns RoomIDHandler structure
func NewRoomIDHandler(base baseHandler) *RoomIDHandler {
	return &RoomIDHandler{
		baseHandler: base,
	}
}
//RoomIDHandler does select, update, delete data from DB by ID
type RoomIDHandler struct {
	baseHandler
}

func (h *RoomIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "roomIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	var ID uint64
	mid, ok := vars["id"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	ID, err := strconv.ParseUint(mid, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", logPrefix, mid))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT room_id,hall_id,name,room_type,active,hls_url,boot,round_id,status,bet_countdown,dealer_id,limitation_id ,create_date FROM room where room_id = ? LIMIT 1"
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s patch checkBody error=%s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		column := ""
		queryString := ""

		if strings.Contains(r.URL.Path, "name") {
			column = "name"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "active") {
			column = "active"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "hlsURL") {
			column = "hls_url"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "boot") {
			column = "boot"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "round") {
			column = "round_id"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "status") {
			column = "status"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "betCountdown") {
			column = "bet_countdown"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "dealerID") {
			column = "dealer_id"
			queryString = "UPDATE room set " + column + "= ?  WHERE room_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "newRound") {
			column = "newRound"
			queryString = "UPDATE room set boot = ?,round_id=?,status=?  WHERE room_id = ? LIMIT 1"
		} else {
			queryString = "UPDATE room set  room_id = ? , hall_id = ? , name = ? , room_type = ? , active = ? , hls_url= ? , bet_countdown= ? , limitation_id= ?   WHERE room_id = ? LIMIT 1"
		}

		//unmarshal request body
		param, err := h.getPatchData(column, body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, err.Error())
			return
		}

		h.dbExec(w, r, logPrefix, ID, column, queryString, param, h.sqlPatch, h.returnExecResponseData)
		return

	}
	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM room  where room_id = ? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, nil, h.sqlDelete, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *RoomIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *RoomIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

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

//delete
func (h *RoomIDHandler) sqlDelete(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	return stmt.Exec(IDOrAccount)
}

//patch
func (h *RoomIDHandler) getPatchData(column string, body []byte) (interface{}, error) {
	switch column {
	case "newRound":
		d := &roomNewRoundPatchParam{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "name":
		d := &nameData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "active":
		d := &activeData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "hls_url":
		d := &hlsURLData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "boot":
		d := &bootData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "round_id":
		d := &roundIDData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "status":
		d := &statusData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "bet_countdown":
		d := &betCountdownData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "dealer_id":
		d := &dealerIDData{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	case "":
		d := &roomPatchParam{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		if len(d.HLSURL) == 0 || len(d.Name) == 0 || d.BetCountdown < 5 {
			return nil, errors.New("param marshal error")
		}
		return d, nil
	default:
		return nil, errors.New("column error")
	}
}
func (h *RoomIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*roomNewRoundPatchParam); ok {
		return stmt.Exec(p.Boot, p.RoundID, p.Status, IDOrAccount)
	}
	if p, ok := param.(*nameData); ok {

		return stmt.Exec(p.Name, IDOrAccount)
	}
	if p, ok := param.(*activeData); ok {
		return stmt.Exec(p.Active, IDOrAccount)
	}
	if p, ok := param.(*hlsURLData); ok {
		return stmt.Exec(p.HLSURL, IDOrAccount)
	}
	if p, ok := param.(*bootData); ok {
		return stmt.Exec(p.Boot, IDOrAccount)
	}
	if p, ok := param.(*roundIDData); ok {
		return stmt.Exec(p.Round, IDOrAccount)
	}
	if p, ok := param.(*statusData); ok {
		return stmt.Exec(p.Status, IDOrAccount)
	}
	if p, ok := param.(*betCountdownData); ok {
		return stmt.Exec(p.BetCountdown, IDOrAccount)
	}
	if p, ok := param.(*dealerIDData); ok {
		return stmt.Exec(p.DealerID, IDOrAccount)
	}
	if p, ok := param.(*roomPatchParam); ok {
		return stmt.Exec(p.RoomID, p.HallID, p.Name, p.RoomType, p.Active, p.HLSURL, p.BetCountdown, p.LimitationID, IDOrAccount)
	}
	return nil, errors.New("parsing param error")
}

func (h *RoomIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*roomIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*roomIDData{
			{
				ID,
			},
		},
	}
}
