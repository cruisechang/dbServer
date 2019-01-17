package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"github.com/gorilla/mux"
	"database/sql"
)

func NewTransferPartnerTransferIDHandler(base baseHandler) *transferPartnerTransferIDHandler {
	return &transferPartnerTransferIDHandler{
		baseHandler: base,
	}
}

type transferPartnerTransferIDHandler struct {
	baseHandler
}

func (h *transferPartnerTransferIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "transferPartnerTransferID"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
	tid, ok := vars["partnerTransferID"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get partnerTransferID not found %s", logPrefix,r.RequestURI))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}


	if len(tid) <5  {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  get path error %s ", logPrefix,r.RequestURI))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {

		queryString := "select transfer.transfer_id, transfer.partner_transfer_id, transfer.partner_id, transfer.user_id, transfer.category, transfer.transfer_credit, transfer.credit, transfer.status, transfer.create_date, user.account, user.name from transfer LEFT JOIN user on transfer.user_id=user.user_id  WHERE transfer.partner_transfer_id = ? "
		h.dbQuery(w, r, logPrefix, tid, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}


func (h *transferPartnerTransferIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *transferPartnerTransferIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := transferDB{}
		resData := []transferData{}

		for rows.Next() {
			err := rows.Scan(&ud.transfer_id, &ud.partner_transfer_id, &ud.partner_id, &ud.user_id, &ud.category, &ud.transfer_credit, &ud.credit, &ud.status, &ud.create_date, &ud.account, &ud.name)
			if err == nil {
				count ++
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
						ud.name,})
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


