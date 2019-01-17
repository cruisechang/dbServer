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

func NewBetsHandler(base baseHandler) *betsHandler {
	return &betsHandler{
		baseHandler: base,
	}
}

type betsHandler struct {
	baseHandler
}

func (h *betsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "bets"

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
		//unmarshal request body
		param := &betGetParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString, queryArgs := h.getQueryStringArgs(param)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		//h.getByFilter(w, r, logPrefix, queryString, queryArgs, h.returnResDataFunc)
		return
	}

	if r.Method == "POST" {
		param := &betPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.BetCredit < 0 || param.ActiveCredit < 0 || param.PrizeCredit < 0 || param.OriginalCredit < 0 || param.BalanceCredit < 0 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		betID, err := util.GetUniqueID()
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
		}

		queryString := "INSERT  INTO bet (bet_id,partner_id,user_id,room_id,room_type,round_id,seat_id,bet_credit,active_credit,prize_credit,result_credit,balance_credit,original_credit,record,status) VALUE (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		h.post(w, r, logPrefix, betID, queryString, param, h.sqlExec, h.returnPostResData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *betsHandler) getQueryStringArgs(param *betGetParam) (queryString string, queryArgs []interface{}) {

	hasFilter := true
	queryString = "select bet.bet_id,bet.partner_id,bet.user_id,bet.room_id,bet.room_type,bet.round_id,bet.seat_id,bet.bet_credit,bet.active_credit,bet.prize_credit,bet.result_credit,bet.balance_credit,bet.original_credit,bet.record,bet.status,bet.create_date, user.account, user.name from bet LEFT JOIN user on bet.user_id=user.user_id where "

	if param.PartnerID > -1 && param.UserID > -1 {
		queryString += " bet.partner_id = ?  AND bet.user_id = ? "
		queryArgs = append(queryArgs, param.PartnerID, param.UserID)
	} else if param.PartnerID > -1 {
		queryString += " bet.partner_id = ?  "
		queryArgs = append(queryArgs, param.PartnerID)
	} else if param.UserID > -1 {
		queryString += " bet.user_id = ?  "
		queryArgs = append(queryArgs, param.UserID)
	} else {
		hasFilter = false
	}

	if param.RoomID > -1 {
		if hasFilter {
			queryString += " AND bet.room_id = ? "
		} else {
			queryString += " bet.room_id = ? "
			hasFilter = true
		}
		queryArgs = append(queryArgs, param.RoomID)
	}
	if param.RoomType > -1 {
		if hasFilter {
			queryString += " AND bet.room_type = ? "
		} else {
			queryString += " bet.room_type = ? "
			hasFilter = true
		}
		queryArgs = append(queryArgs, param.RoomType)
	}

	if param.RoundID > -1 {
		if hasFilter {
			queryString += " AND bet.round_id = ? "

		} else {
			queryString += " bet.round_id = ? "
			hasFilter = true
		}
		queryArgs = append(queryArgs, param.RoundID)
	}

	if param.Status > -1 {
		if hasFilter {
			queryString += " AND bet.status = ? "
		} else {
			queryString += " bet.status = ? "
			hasFilter = true
		}
		queryArgs = append(queryArgs, param.Status)
	}

	if hasFilter {
		queryString += "AND bet.create_date BETWEEN ? AND ?"
	} else {
		queryString += " bet.create_date BETWEEN ? AND ?"
	}
	queryArgs = append(queryArgs, param.BeginDate, param.EndDate)

	return
}

func (h *betsHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	//return stmt.Query(IDOrAccount)

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
func (h *betsHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := betDB{}
		resData := []betData{}

		for rows.Next() {
			err := rows.Scan(&ud.bet_id, &ud.partner_id, &ud.user_id, &ud.room_id, &ud.room_type, &ud.round_id, &ud.seat_id, &ud.bet_credit, &ud.active_credit, &ud.prize_credit, &ud.result_credit, &ud.balance_credit, &ud.original_credit, &ud.partner_id, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count ++
				resData = append(resData,
					betData{
						ud.bet_id,
						ud.partner_id,
						ud.user_id,
						ud.room_id,
						ud.room_type,
						ud.round_id,
						ud.seat_id,
						ud.bet_credit,
						ud.active_credit,
						ud.prize_credit,
						ud.result_credit,
						ud.balance_credit,
						ud.original_credit,
						ud.record,
						ud.status,
						ud.create_date,
						ud.account,
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

//func (h *betsHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {
//
//	return func(rows *sql.Rows) (interface{}, int) {
//		count := 0
//		ud := betDB{}
//		resData := []betData{}
//
//		for rows.Next() {
//			err := rows.Scan(&ud.bet_id, &ud.partner_id, &ud.user_id, &ud.room_id, &ud.room_type, &ud.round_id, &ud.seat_id, &ud.bet_credit, &ud.active_credit, &ud.prize_credit, &ud.result_credit, &ud.balance_credit, &ud.original_credit, &ud.partner_id, &ud.status, &ud.create_date, &ud.account, &ud.name)
//			if err == nil {
//				count += 1
//				resData = append(resData,
//					betData{
//						ud.bet_id,
//						ud.partner_id,
//						ud.user_id,
//						ud.room_id,
//						ud.room_type,
//						ud.round_id,
//						ud.seat_id,
//						ud.bet_credit,
//						ud.active_credit,
//						ud.prize_credit,
//						ud.result_credit,
//						ud.balance_credit,
//						ud.original_credit,
//						ud.record,
//						ud.status,
//						ud.create_date,
//						ud.account,
//						ud.name,
//					})
//			}
//		}
//
//		return resData, count
//	}
//}

func (h *betsHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	if p, ok := param.(*betPostParam); ok {

		return stmt.Exec(ID, p.PartnerID, p.UserID, p.RoomID, p.RoomType, p.RoundID, p.SeatID, p.BetCredit, p.ActiveCredit, p.PrizeCredit, p.ResultCredit, p.BalanceCredit, p.OriginalCredit, p.Record, p.Status)
	}
	return nil, errors.New("parsing param error")

}
func (h *betsHandler) returnPostResData(ID, lastID uint64) interface{} {

	return []betIDData{
		{
			ID,
		},
	}
}
