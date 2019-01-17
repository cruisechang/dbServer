package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"encoding/json"
	"database/sql"
	"errors"
)

func NewBannersHandler(base baseHandler) *BannersHandler {
	return &BannersHandler{
		baseHandler: base,
	}
}

type BannersHandler struct {
	baseHandler
}

func (h *BannersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "BannersHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	if r.Method == "GET" || r.Method == "get" {
		queryString := "SELECT banner_id, pic_url, link_url, description, platform,active, create_date  FROM banner "
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

		param := &bannerPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.PicURL) < 10 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		queryString := "INSERT  INTO banner (pic_url,link_url,description,platform,active) values (? ,?, ?, ?,?)"
		h.dbExec(w, r, logPrefix, 0, "", queryString, param, h.sqlExec, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//get
func (h *BannersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query()
}
func (h *BannersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := bannerDB{}
		resData := []bannerData{}

		for rows.Next() {
			err := rows.Scan(&ud.banner_id, &ud.pic_url, &ud.link_url, &ud.description, &ud.platform, &ud.active, &ud.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					bannerData{
						ud.banner_id,
						ud.pic_url,
						ud.link_url,
						ud.description,
						ud.platform,
						ud.active,
						ud.create_date,})
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
func (h *BannersHandler) sqlExec(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*bannerPostParam); ok {

		return stmt.Exec(p.PicURL, p.LinkURL, p.Description, p.Platform, p.Active)
	}
	return nil, errors.New("")

}
func (h *BannersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*bannerIDData{{}},
		}

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecLastIDError,
			Count:   0,
			Message: "",
			Data:    []*bannerIDData{{}},
		}
	}

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*bannerIDData{
			{
				uint64(lastID),
			},
		},
	}
}
