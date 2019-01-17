package handler

import (
	"github.com/cruisechang/dbex"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"io/ioutil"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
)

//type returnTargetColumnResDataCount func(column string, rows *sql.Rows) (interface{}, int)
//type returnResDataFunc func() (func(rows *sql.Rows) (interface{}, int))

type sqlExecFunc func(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error)
type returnPostResDataFunc func(ID uint64, lastID uint64) interface{}
type returnIDResDataFunc func(ID uint64) interface{}

type sqlExecFunc2 func(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error)
type returnExecResponseData func(IDOrAccount interface{},column string,result sql.Result) (*responseData)


type sqlQueryFunc func(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error)
type returnResponseDataStructFunc func() (func(IDOrAccount interface{},column string,rows *sql.Rows) (*responseData))

func NewBaseHandler(db *dbex.DB, logger *dbex.Logger) baseHandler {
	return baseHandler{
		db:     db,
		logger: logger,
	}

}

type baseHandler struct {
	db     *dbex.DB
	logger *dbex.Logger
}

func (h *baseHandler) checkHead(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("API-Key") != "qwerASDFzxcv!@#$" {
		return false
	}
	return true
}
func (h *baseHandler) checkBody(w http.ResponseWriter, r *http.Request) (body []byte, errorCode int, errorMessage string) {
	if r.Body == nil {
		return nil, CodeBodyNil, fmt.Sprintf("handler get body=nil")
	}
	//ioBody, err := r.GetBody()
	//if err != nil {
	//	return nil, CodeBodyError0, fmt.Sprintf("handler get GetBody error=%s", err.Error())
	//}
	//
	//if ioBody == nil {
	//	return nil, CodeBodyError1, fmt.Sprintf("handler get body=nil")
	//}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, CodeBodyReadError, fmt.Sprintf("handler get read body error=%s", err.Error())
	}

	if body == nil {
		return nil, CodeBodyReadNil, fmt.Sprintf("handler get body=nil")
	}
	return body, CodeSuccess, ""
}

func (h *baseHandler) getVariable(r *http.Request, variable string) (string, error) {
	vars := mux.Vars(r)
	v, ok := vars[variable]
	if !ok {
		return "", errors.New("get varialbe error")
	}
	return v, nil
}

//post creates a new record
//ID 可能是userID 可能是hallID, 或是沒用
//
func (h *baseHandler) post(w http.ResponseWriter, r *http.Request, logPrefix string, ID uint64, queryString string, param interface{}, sqlExec sqlExecFunc, returnPostResData returnPostResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("handler post param =%+v", param))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	result, err := sqlExec(stmt, ID, param)

	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecError, "insert error , maybe insert duplicate value into unique column")
		return
	}

	affRow, err := result.RowsAffected()
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecResultError, "")
		return
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post sqlDB exec error =%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecLastIDError, "insert get dealer id error")
		return
	}

	resData := returnPostResData(ID, uint64(lastID))

	rd := responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler post res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler post response data=%s", logPrefix, resStr))
}

func (h *baseHandler) patch(w http.ResponseWriter, r *http.Request, logPrefix string, ID uint64, queryString string, param interface{}, sqlExec sqlExecFunc, returnExecResData returnIDResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("handler post param =%+v", param))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler exec sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	result, err := sqlExec(stmt, ID, param)

	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler exec sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecError, "insert error , maybe insert duplicate value into unique column")
		return
	}

	affRow, err := result.RowsAffected()
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler exec sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecResultError, "")
		return
	}

	resData := returnExecResData(ID)

	rd := responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler exec res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler exec response data=%s", logPrefix, resStr))
}

func (h *baseHandler) delete(w http.ResponseWriter, r *http.Request, logPrefix string, ID uint64, queryString string, returnExecResData returnIDResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler delete sqlDB queryString=%s,ID=%d", logPrefix, queryString, ID))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler delete sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(ID)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler delete sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecError, "")
		return
	}

	affRow, err := result.RowsAffected()
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler delete sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecResultError, "delete result error")
		return
	}

	resData := returnExecResData(ID)

	rd := responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler delete res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler delete response data=%s", logPrefix, resStr))
}

func (h *baseHandler) dbExec(w http.ResponseWriter, r *http.Request, logPrefix string, IDOrAccount interface{}, targetColumn, queryString string, param interface{}, sqlExec sqlExecFunc2, returnResFunc returnExecResponseData) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s dbExec sqlDB queryString=%s", logPrefix, queryString))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbExec sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	result, err := sqlExec(stmt, IDOrAccount, param)

	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbExec sqlDB exec error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecError, fmt.Sprintf("%s dbExec sqlDB exec error %s",logPrefix,err.Error()))
		return
	}

	resData := returnResFunc(IDOrAccount,targetColumn,result)

	js, err := json.Marshal(resData)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbExec res marshal error=%s, resData=%+v", logPrefix, err.Error(), resData))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, fmt.Sprintf("%s dbExec get res marshal error=%s, resData=%+v", logPrefix, err.Error(), resData))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s dbExec response data=%s", logPrefix, resStr))

}

func (h *baseHandler) dbQuery(w http.ResponseWriter, r *http.Request, logPrefix string, IDOrAccount interface{}, targetColumn, queryString string, param interface{}, sqlQuery sqlQueryFunc, returnResFunc returnResponseDataStructFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s dbQuery sqlDB queryString=%s", logPrefix, queryString))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbQuery sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	rows, err := sqlQuery(stmt, IDOrAccount, param)

	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbQuery sqlDB quert error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBExecError, fmt.Sprintf("%s get query error %s",logPrefix,err.Error()))
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbQuery sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s dbQuery sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	resFunc := returnResFunc()
	ID,_:=IDOrAccount.(uint64)
	resData := resFunc(ID,targetColumn,rows)

	js, err := json.Marshal(resData)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbQuery res marshal error=%s, resData=%+v", logPrefix, err.Error(), resData))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, fmt.Sprintf("%s dbQuery get res marshal error=%s, resData=%+v", logPrefix, err.Error(), resData))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s dbQuery response data=%s", logPrefix, resStr))

}

//get select rows without filter
/*
func (h *baseHandler) get(w http.ResponseWriter, r *http.Request, logPrefix, queryString string, returnResFunc returnResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s", logPrefix, queryString))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, "")
		return
	}
	defer rows.Close()

	resFunc := returnResFunc()
	resData, count := resFunc(rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s handler get sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, fmt.Sprintf("%s handler get res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get response data=%s", logPrefix, resStr))
}
*/


//getByFilter
/*
func (h *baseHandler) getByFilter(w http.ResponseWriter, r *http.Request, logPrefix, queryString string, queryArgs []interface{}, returnResFunc returnResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s, queryArgs=%+v", logPrefix, queryString, queryArgs))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	rows, err := h.sqlQuery(stmt, queryArgs)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, "")
		return
	}
	defer rows.Close()

	//res
	resFunc := returnResFunc()
	resData, count := resFunc(rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, "")
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get response data=%s", logPrefix, resStr))

}
*/

/*
func (h *baseHandler) getTargetRow(w http.ResponseWriter, r *http.Request, logPrefix string, ID uint64, queryString string, returnResFunc returnResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s, ID=%d", logPrefix, queryString, ID))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, "")
		return
	}
	defer rows.Close()

	resFunc := returnResFunc()
	resData, count := resFunc(rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, "")
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get response data=%s", logPrefix, resStr))

}
*/
/*
func (h *baseHandler) getTargetRowByStringData(w http.ResponseWriter, r *http.Request, logPrefix string, targetData string, queryString string, returnResFunc returnResDataFunc) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s, target data=%s", logPrefix, queryString, targetData))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, "")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(targetData)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, "")
		return
	}
	defer rows.Close()

	resFunc := returnResFunc()
	resData, count := resFunc(rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, "")
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get response data=%s", logPrefix, resStr))

}
*/
/*
func (h *baseHandler) getTargetColumnValueByUserAccount(w http.ResponseWriter, r *http.Request, logPrefix string, param *userAccountGetParam, column string, queryString string, returnResFunc returnTargetColumnResDataCount) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s, param=%+v, column=%s", logPrefix, queryString, param, column))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(param.PartnerID, param.Account)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		return
	}
	defer rows.Close()

	//res
	//resFun := returnResFunc(column)
	//resData, count := resFun(rows)

	resData, count := returnResFunc(column, rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler getTargetColumn response data=%s", logPrefix, resStr))
}
*/

/*
func (h *baseHandler) getTargetColumnValueByAccount(w http.ResponseWriter, r *http.Request, logPrefix string, account string, column string, queryString string, returnResFunc returnTargetColumnResDataCount) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s, account=%s, column=%s", logPrefix, queryString, account, column))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(account)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		return
	}
	defer rows.Close()

	resData, count := returnResFunc(column, rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler getTargetColumn response data=%s", logPrefix, resStr))
}

func (h *baseHandler) getTargetColumnValue(w http.ResponseWriter, r *http.Request, logPrefix string, ID uint64, column string, queryString string, returnResFunc returnTargetColumnResDataCount) {
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler get sqlDB queryString=%s,ID=%d", logPrefix, queryString, ID))

	sqlDB := h.db.GetSQLDB()

	stmt, err := sqlDB.Prepare(queryString)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBPrepareError, fmt.Sprintf("%s handler getTargetColumn sqlDB prepare error=%s", logPrefix, err.Error()))
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBQueryError, fmt.Sprintf("%s handler getTargetColumn sqlDB query error=%s", logPrefix, err.Error()))
		return
	}
	defer rows.Close()

	resData, count := returnResFunc(column, rows)

	//
	if err = rows.Err(); err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s handler getTargetColumn sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	rd := responseData{
		Code:    CodeSuccess,
		Count:   count,
		Message: "",
		Data:    resData,
	}

	js, err := json.Marshal(rd)
	if err != nil {
		h.writeError(w, http.StatusOK, CodeResponseDataMarshalError, "")
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler getTargetColumn res marshal error=%s, resData=%+v", logPrefix, err.Error(), rd))
		return
	}
	resStr := string(js)

	h.writeSuccess(w, resStr)
	h.logger.LogFile(dbex.LevelInfo, fmt.Sprintf("%s handler getTargetColumn response data=%s", logPrefix, resStr))
}
*/

func (h *baseHandler) writeError(w http.ResponseWriter, httpStatusCode int, errorCode int, errorMsg string) {

	rd := responseData{
		Code:    errorCode,
		Count:   0,
		Message: errorMsg,
		Data:    []struct{}{},
		//Data: "[{}]",
	}
	b, _ := json.Marshal(rd)

	w.WriteHeader(httpStatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies.
	io.WriteString(w, string(b))
}

func (h *baseHandler) writeSuccess(w http.ResponseWriter, resStr string) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies.

	io.WriteString(w, resStr)

}

//func (h *baseHandler) sqlQuery(stmt *sql.Stmt, args []interface{}) (*sql.Rows, error) {
//	switch len(args) {
//	case 0:
//		return stmt.Query()
//	case 1:
//		return stmt.Query(args[0])
//	case 2:
//		return stmt.Query(args[0], args[1])
//	case 3:
//		return stmt.Query(args[0], args[1], args[2])
//	case 4:
//		return stmt.Query(args[0], args[1], args[2], args[3])
//	case 5:
//		return stmt.Query(args[0], args[1], args[2], args[3], args[4])
//	case 6:
//		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5])
//	case 7:
//		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5], args[6])
//	case 8:
//		return stmt.Query(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7])
//
//	}
//	return nil, errors.New("args error")
//
//}
