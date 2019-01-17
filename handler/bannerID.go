package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
)

//NewBroadcastIDHandler returns handler for banner
func NewBannerIDHandler(base baseHandler) *BannerIDHandler {
	return &BannerIDHandler{
		baseHandler: base,
	}
}

//BannerIDHandler presents structure of handling banner
type BannerIDHandler struct {
	baseHandler
}

func (h *BannerIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "BroadcastIDHandler"

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

	if r.Method == "DELETE" || r.Method == "delete" {
		queryString := "DELETE FROM banner  where banner_id = ? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, nil, h.sqlDelete, h.returnExecResponseData)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		//check body
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s  patch %s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		param := &bannerPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		queryString := "Update banner set pic_url=? ,link_url=? ,description=?, platform=?, active=? WHERE banner_id =? LIMIT 1"
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}
	h.writeError(w, http.StatusNotFound, CodeMethodError, "")
}

//delete
func (h *BannerIDHandler) sqlDelete(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	return stmt.Exec(IDOrAccount)
}

//patch

func (h *BannerIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	ID, _ := IDOrAccount.(uint64)

	if p, ok := param.(*bannerPostParam); ok {
		return stmt.Exec(p.PicURL, p.LinkURL, p.Description, p.Platform, p.Active, ID)
	}

	return nil, errors.New("parsing param error")
}

func (h *BannerIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*bannerIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint64)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*bannerIDData{
			{
				ID,
			},
		},
	}
}
