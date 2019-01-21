package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
	"github.com/sony/sonyflake"
	"io"
	"io/ioutil"
	"net/http"
)

type sqlExecFunc func(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error)
type returnExecResponseData func(IDOrAccount interface{}, column string, result sql.Result) *responseData

type sqlQueryFunc func(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error)
type returnResponseDataStructFunc func() func(IDOrAccount interface{}, column string, rows *sql.Rows) *responseData

//NewBaseHandler returns structure of base handler
func NewBaseHandler(db *dbex.DB, logger *dbex.Logger, provider *sonyflake.Sonyflake) baseHandler {
	return baseHandler{
		db:               db,
		logger:           logger,
		uniqueIDProvider: provider,
	}
}

//BaseHandler 是所有handler的基礎，負責實際db操作及log行為
type baseHandler struct {
	db               *dbex.DB
	logger           *dbex.Logger
	uniqueIDProvider *sonyflake.Sonyflake
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

func (h *baseHandler) getUniqueID() (uint64, error) {
	if id, err := h.uniqueIDProvider.NextID(); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

//post creates a new record
//ID 可能是userID 可能是hallID, 或是沒用
//
/*
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
*/

/*
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
*/
/*
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
*/

func (h *baseHandler) dbExec(w http.ResponseWriter, r *http.Request, logPrefix string, IDOrAccount interface{}, targetColumn, queryString string, param interface{}, sqlExec sqlExecFunc, returnResFunc returnExecResponseData) {
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
		h.writeError(w, http.StatusOK, CodeDBExecError, fmt.Sprintf("%s dbExec sqlDB exec error %s", logPrefix, err.Error()))
		return
	}

	resData := returnResFunc(IDOrAccount, targetColumn, result)

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
		h.writeError(w, http.StatusOK, CodeDBExecError, fmt.Sprintf("%s get query error %s", logPrefix, err.Error()))
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s dbQuery sqlDB rows error=%s", logPrefix, err.Error()))
		h.writeError(w, http.StatusOK, CodeDBRowsError, fmt.Sprintf("%s dbQuery sqlDB rows error=%s", logPrefix, err.Error()))
		return
	}

	resFunc := returnResFunc()
	ID, _ := IDOrAccount.(uint64)
	resData := resFunc(ID, targetColumn, rows)

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
