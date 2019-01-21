package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbServer/util"
	"github.com/cruisechang/dbex"
	"github.com/juju/errors"
)

func NewTransfersHandler(base baseHandler) *transfersHandler {
	return &transfersHandler{
		baseHandler: base,
	}
}

type transfersHandler struct {
	baseHandler
}

func (h *transfersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "transfersHandler"

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
		param := &transferGetParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString, queryArgs := h.getQueryStringArgs(param)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "POST" {
		param := &transferPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.PartnerTransferID) < 5 || param.Credit < 0 || param.PartnerID < 0 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		transferID, err := util.GetUniqueID()
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, fmt.Sprintf("%s handler post get unique ID error %s", logPrefix, err.Error()))
		}
		queryString := "INSERT  INTO transfer (transfer_id,partner_transfer_id,partner_id,user_id, category,transfer_credit,credit,status) values (? ,? ,? ,?, ? ,?, ?,?)"
		h.dbExec(w, r, logPrefix, transferID, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
//query
func (h *transfersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

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
func (h *transfersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := transferDB{}
		resData := []transferData{}

		for rows.Next() {
			err := rows.Scan(&ud.transfer_id, &ud.partner_transfer_id, &ud.partner_id, &ud.user_id, &ud.category, &ud.transfer_credit, &ud.credit, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count += 1
				resData = append(resData,
					transferData{
						ud.transfer_id,
						ud.partner_transfer_id,
						ud.partner_id,
						ud.user_id,
						ud.category,
						ud.transfer_credit,
						ud.credit,
						ud.status,
						ud.create_date,
						ud.account,
						ud.name})
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

func (h *transfersHandler) getQueryStringArgs(param *transferGetParam) (queryString string, queryArgs []interface{}) {

	queryString = "select transfer.transfer_id, transfer.partner_transfer_id, transfer.partner_id, transfer.user_id, transfer.category, transfer.transfer_credit, transfer.credit, transfer.status, transfer.create_date, user.account, user.name from transfer LEFT JOIN user on transfer.user_id=user.user_id  WHERE "
	hasFilter := true
	if param.PartnerID > -1 && param.UserID > -1 {
		queryString += " transfer.partner_id = ?  AND transfer.user_id = ?"
		queryArgs = append(queryArgs, param.PartnerID, param.UserID)
	} else if param.PartnerID > -1 {
		queryString += " transfer.partner_id = ?  "
		queryArgs = append(queryArgs, param.PartnerID)
	} else if param.UserID > -1 {
		queryString += " transfer.user_id = ? "
		queryArgs = append(queryArgs, param.UserID)
	} else {
		//全-1
		hasFilter = false
	}

	if param.Category > -1 {
		if hasFilter {
			queryString += " AND transfer.category = ? "
		} else {
			queryString += " transfer.category = ? "
		}
		queryArgs = append(queryArgs, param.Category)
		hasFilter = true
	}

	if param.Status > -1 {
		if hasFilter {
			queryString += " AND transfer.status = ? "

		} else {
			queryString += " transfer.status = ? "
		}
		queryArgs = append(queryArgs, param.Status)
		hasFilter = true
	}

	if hasFilter {
		queryString += " AND ( transfer.create_date BETWEEN ? AND ? ) "
	} else {
		queryString += "  ( transfer.create_date BETWEEN ? AND ? ) "
	}
	queryArgs = append(queryArgs, param.BeginDate, param.EndDate)

	return
}

//func (h *transfersHandler) sqlExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {
//
//	if p, ok := param.(*transferPostParam); ok {
//
//		return stmt.Exec(ID, p.PartnerTransferID, p.PartnerID, p.UserID, p.Category, p.TransferCredit, p.Credit, p.Status)
//	}
//	return nil, errors.New("parsing param error")
//
//}
//func (h *transfersHandler) returnPostResData(ID, lastID uint64) interface{} {
//
//	return []transferIDData{
//		{
//			ID,
//		},
//	}
//}

//post
func (h *transfersHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*transferPostParam); ok {

		return stmt.Exec(IDOrAccount, p.PartnerTransferID, p.PartnerID, p.UserID, p.Category, p.TransferCredit, p.Credit, p.Status)
	}
	return nil, errors.New("parsing param error")

}

//id預先產生
func (h *transfersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*transferIDData{{}},
		}

	}

	if id, ok := IDOrAccount.(uint64); ok {
		return &responseData{
			Code:    CodeSuccess,
			Count:   int(affRow),
			Message: "",
			Data: []*transferIDData{
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
		Data:    []*transferIDData{{}},
	}
}
