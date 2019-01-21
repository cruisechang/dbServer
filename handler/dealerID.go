package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"database/sql"
)

func NewDealerIDHandler(base baseHandler) *dealerIDHandler {
	return &dealerIDHandler{
		baseHandler: base,
	}
}

type dealerIDHandler struct {
	baseHandler
}

func (h *dealerIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "dealerIDHandler"

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
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", mid, logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if ID == 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if r.Method == "GET" || r.Method == "get" {

		queryString := "SELECT dealer_id,name,account,active, portrait_url,create_date FROM dealer where dealer_id = ? LIMIT 1"
		//h.getTargetRow(w, r, "dealer", ID, queryString, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM dealer  where dealer_id = ? LIMIT 1"
		//h.delete(w, r, "dealer", ID, queryString, h.returnIDResData)
		h.dbExec(w, r, logPrefix, ID, "", queryString, nil, h.sqlDelete, h.returnExecResponseData)

		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s patch checkBody error=%s", errMsg, logPrefix))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		queryString := "UPDATE dealer set name= ? , password =? , active =? , portrait_url =?  WHERE dealer_id = ? LIMIT 1"

		//unmarshal request body
		param, err := h.getPatchData(body)

		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		//h.patch(w, r, logPrefix, ID, queryString, patchData, h.patchExec, h.returnIDResData)
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

func (h *dealerIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *dealerIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := dealerDB{}
		resData := []dealerData{}

		for rows.Next() {
			err := rows.Scan(&ud.dealer_id, &ud.name, &ud.account, &ud.active, &ud.portrait_url, &ud.create_date)
			if err == nil {
				count ++
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

/*
func (h *dealerIDHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
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
		return resData, count
	}
}
*/


//delete
func (h *dealerIDHandler) sqlDelete(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	return stmt.Exec(IDOrAccount)
}

//patch
func (h *dealerIDHandler) getPatchData(body []byte) (interface{}, error) {
	p := &dealerPatchParam{}
	p.Active = -1 //for test
	err := json.Unmarshal(body, p)
	if err != nil {
		return nil, err
	}
	if len(p.Name) < 3 || p.Active == -1 || len(p.Password) < 3 || len(p.PortraitURL) < 3 {
		return nil, errors.New("patch data marshal error")
	}
	return p, nil
}
func (h *dealerIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	//檢查參數是否合法
	if p, ok := param.(*dealerPatchParam); ok {
		return stmt.Exec(p.Name, p.Password, p.Active, p.PortraitURL, IDOrAccount)
	}
	return nil, errors.New("parsing param error")
}

func (h *dealerIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*dealerIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint64)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*dealerIDData{
			{
				uint(ID),
			},
		},
	}
}


//func (h *dealerIDHandler) patchExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {
//
//	//檢查參數是否合法
//	if p, ok := param.(*dealerPatchParam); ok {
//
//		return stmt.Exec(p.Name, p.Password, p.Active, p.PortraitURL, ID)
//	}
//	return nil, errors.New("parsing param error")
//
//}
//
//func (h *dealerIDHandler) returnIDResData(ID uint64) interface{} {
//	return []dealerIDData{
//		{
//			uint(ID),
//		},
//	}
//}
