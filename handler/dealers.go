package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cruisechang/dbex"
)

//NewDealersHandler returns DealersHandler structure
func NewDealersHandler(base baseHandler) *DealersHandler {
	return &DealersHandler{
		baseHandler: base,
	}
}

//DealersHandler select dealers and insert new dealer
type DealersHandler struct {
	baseHandler
}

func (h *DealersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "dealersHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT dealer_id,name,account,active, portrait_url,create_date FROM dealer "
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

		param := &dealerPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("dealersHandler post data unmarshal error=%s", err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.Name) < 3 || len(param.Account) < 3 || len(param.Password) < 3 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("dealersHandler post data illegal=%+v", param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		queryString := "INSERT  INTO dealer (name,account,password,portrait_url ) values (? ,?,?,?)"
		h.dbExec(w, r, logPrefix, 0, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//query
func (h *DealersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *DealersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := dealerDB{}
		resData := []dealerData{}

		for rows.Next() {
			err := rows.Scan(&ud.dealer_id, &ud.name, &ud.account, &ud.active, &ud.portrait_url, &ud.create_date)
			if err == nil {
				count += 1
				resData = append(resData,
					dealerData{
						ud.dealer_id,
						ud.name,
						ud.account,
						ud.active,
						ud.portrait_url,
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
func (h *DealersHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*dealerPostParam); ok {

		return stmt.Exec(p.Name, p.Account, p.Password, p.PortraitURL)
	}
	return nil, errors.New("")

}

//id自動產生
func (h *DealersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*dealerIDData{{}},
		}

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecLastIDError,
			Count:   0,
			Message: "",
			Data:    []*dealerIDData{{}},
		}
	}

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*dealerIDData{
			{
				uint(lastID),
			},
		},
	}
}
