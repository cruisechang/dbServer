package handler

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/cruisechang/dbex"
	"github.com/juju/errors"
	"github.com/cruisechang/dbServer/util"
)

func NewRoundsHandler(base baseHandler) *roundsHandler {
	return &roundsHandler{
		baseHandler: base,
	}
}

type roundsHandler struct {
	baseHandler
}

func (h *roundsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "roundsHandler"

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

	if r.Method == "GET" {
		//umarshal request body
		param := &roundGetParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString, queryArgs := h.getQueryStringArgs(param)
		//h.getByFilter(w, r, logPrefix, queryString, queryArgs, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "POST" {
		param := &roundPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.HallID < 0 || param.RoomID < 0 || param.RoomType < 0 || param.Status < 0 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		ID, err := util.GetUniqueID()
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
		}

		queryString := "INSERT  INTO round (round_id,hall_id,room_id,room_type,brief,record,status) VALUE (?,?,?,?,?,?,?)"

		h.post(w, r, logPrefix, ID, queryString, param, h.sqlExec, h.returnPostResData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}


func (h *roundsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	args, ok := param.([]interface{})
	if !ok {
		return nil, errors.New("args error")
	}

	switch len(args) {
	case 0:
		return stmt.Query()
	case 1:
		return stmt.Query(args[0])
	case 2:
		return stmt.Query(args[0], args[1])
	case 3:
		return stmt.Query(args[0], args[1], args[2])
	case 4:
		return stmt.Query(args[0], args[1], args[2], args[3])
	case 5:
		return stmt.Query(args[0], args[1], args[2], args[3], args[4])
	case 6:
		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5])
	case 7:
		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5], args[6])
	case 8:
		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7])

	}
	return nil, errors.New("args error")
}
func (h *roundsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := roundDB{}
		resData := []roundData{}

		for rows.Next() {
			err := rows.Scan(&ud.round_id, &ud.hall_id, &ud.room_id, &ud.room_type, &ud.brief, &ud.record, &ud.status, &ud.create_date, &ud.end_datea, &ud.name)
			if err == nil {
				count ++
				resData = append(resData,
					roundData{
						ud.round_id,
						ud.hall_id,
						ud.room_id,
						ud.room_type,
						ud.brief,
						ud.record,
						ud.status,
						ud.create_date,
						ud.end_datea,
						ud.name,
					})
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
func (h *roundsHandler) getQueryStringArgs(param *roundGetParam) (queryString string, queryArgs []interface{}) {

	hasFilter := false
	queryString = "select round.round_id,round.hall_id,round.room_id,round.room_type,round.brief,round.record,round.status,round.create_date, round.end_date,room.name from round LEFT JOIN room on round.room_id=room.room_id where "

	if param.HallID > -1 {
		queryString += " round.hall_id = ?  "
		queryArgs = append(queryArgs, param.HallID)
		hasFilter = true
	}

	if param.RoomID > -1 {
		if hasFilter {
			queryString += " and round.room_id = ?  "

		} else {
			queryString += " round.room_id = ?  "
		}
		hasFilter = true
		queryArgs = append(queryArgs, param.RoomID)
	}

	if param.RoomType > -1 {

		if hasFilter {
			queryString += " AND round.room_type = ? "
		} else {
			queryString += " round.room_type = ? "
		}
		hasFilter = true
		queryArgs = append(queryArgs, param.RoomType)
	}

	if param.Status > -1 {
		if hasFilter {
			queryString += " AND round.status = ? "
		} else {
			queryString += " round.status = ? "

		}
		hasFilter = true
		queryArgs = append(queryArgs, param.Status)
	}

	if hasFilter {
		queryString += "AND round.create_date BETWEEN ? AND ?"
	} else {
		queryString += " round.create_date BETWEEN ? AND ?"
	}
	queryArgs = append(queryArgs, param.BeginDate, param.EndDate)

	return
}
//func (h *roundsHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {
//
//	return func(rows *sql.Rows) (interface{}, int) {
//		count := 0
//		ud := roundDB{}
//		resData := []roundData{}
//
//		for rows.Next() {
//			err := rows.Scan(&ud.round_id, &ud.hall_id, &ud.room_id, &ud.room_type, &ud.brief, &ud.record, &ud.status, &ud.create_date, &ud.end_datea, &ud.name)
//			if err == nil {
//				count ++
//				resData = append(resData,
//					roundData{
//						ud.round_id,
//						ud.hall_id,
//						ud.room_id,
//						ud.room_type,
//						ud.brief,
//						ud.record,
//						ud.status,
//						ud.create_date,
//						ud.end_datea,
//						ud.name,
//					})
//			}
//		}
//
//		return resData, count
//	}
//}

func (h *roundsHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*roundPostParam); ok {

		return stmt.Exec(ID, p.HallID, p.RoomID, p.RoomType, p.Brief, p.Record, p.Status)
	}
	return nil, errors.New("parsing param error")

}
func (h *roundsHandler) returnPostResData(ID, lastID uint64) interface{} {

	return []roundIDData{
		{
			ID,
		},
	}
}