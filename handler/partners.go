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

//NewPartnersHandler returns PartnersHandler structure
func NewPartnersHandler(base baseHandler) *PartnersHandler {
	return &PartnersHandler{
		baseHandler: base,
	}
}

//PartnersHandler select and insert new data
type PartnersHandler struct {
	baseHandler
}

func (h *PartnersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPrefix := "partnersHandler"

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
		param := &partnerGetParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler get data unmarshal error=%s", err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if param.OrderBy != "partnerID" && param.OrderBy != "" {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler get data orderBy error, orderBy=%s", param.OrderBy))
			h.writeError(w, http.StatusOK, CodeRequestDataError, "")
			return
		}
		if param.Order != "asc" && param.Order != "desc" && param.Order != "" {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler get data order error, order=%s", param.Order))
			h.writeError(w, http.StatusOK, CodeRequestDataError, "")
			return
		}

		queryString, queryArgs := h.getQueryStringArgs(param)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		return
	}
	if r.Method == "POST" {
		param := &partnerPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler post data unmarshal error=%s", err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.Account) < 5 || len(param.Password) < 5 || len(param.Name) < 5 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler post data illegal=%+v", param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		partnerID, err := util.GetUniqueID()
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("partnersHandler post get unique hallID error %s", err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, fmt.Sprintf("partnersHandler post get unique hallID error %s", err.Error()))
		}
		queryString := "INSERT  INTO partner (partner_id,account,password,name,level,category,aes_key,access_token,api_bind_ip,cms_bind_ip) values (? ,? ,? , ?, ? ,?, ?,?,?,?)"
		h.dbExec(w, r, logPrefix, partnerID, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *PartnersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

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
func (h *PartnersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		dc := partnerDB{} //db container
		resData := []partnerData{}

		for rows.Next() {
			err := rows.Scan(&dc.partner_id, &dc.account, &dc.name, &dc.level, &dc.category, &dc.active, &dc.api_bind_ip, &dc.cms_bind_ip, &dc.create_date)
			if err == nil {
				count++
				resData = append(resData,
					partnerData{
						dc.partner_id,
						dc.account,
						dc.name,
						dc.level,
						dc.category,
						dc.active,
						dc.api_bind_ip,
						dc.cms_bind_ip,
						dc.create_date})
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

func (h *PartnersHandler) getQueryStringArgs(param *partnerGetParam) (queryString string, queryArgs []interface{}) {

	queryString = "select partner_id,account,name,level,category,active,api_bind_ip,cms_bind_ip,create_date from partner "

	//status
	if param.Active > -1 {
		queryString += "WHERE  active = ? "
		queryArgs = append(queryArgs, param.Active)
	}

	//orderBy, order
	if param.OrderBy != "" {
		queryString += " ORDER BY ? "
		queryArgs = append(queryArgs, "partner_id")

		if param.Order == "asc" {
			queryString += " asc "
		} else if param.Order == "desc" {
			queryString += " desc "
		}
	}
	//limit, offset
	if param.Limit > 0 {
		queryString += " LIMIT ? "
		queryArgs = append(queryArgs, param.Limit)
		if param.Offset > 0 {
			queryString += " OFFSET ?"
			queryArgs = append(queryArgs, param.Offset)

		}
	}
	return
}
func (h *PartnersHandler) returnResDataFunc() func(rows *sql.Rows) (interface{}, int) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		dc := partnerDB{} //db container
		resData := []partnerData{}

		for rows.Next() {
			err := rows.Scan(&dc.partner_id, &dc.account, &dc.name, &dc.level, &dc.category, &dc.active, &dc.api_bind_ip, &dc.cms_bind_ip, &dc.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					partnerData{
						dc.partner_id,
						dc.account,
						dc.name,
						dc.level,
						dc.category,
						dc.active,
						dc.api_bind_ip,
						dc.cms_bind_ip,
						dc.create_date})
			}
		}
		return resData, count
	}
}

//post
func (h *PartnersHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*partnerPostParam); ok {

		return stmt.Exec(IDOrAccount, p.Account, p.Password, p.Name, p.Level, p.Category, p.AESKey, p.AccessToken, p.APIBindIP, p.CMSBindIP)
	}
	return nil, errors.New("parsing param error")

}

//id 預先產生
func (h *PartnersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*partnerIDData{{}},
		}

	}

	if id, ok := IDOrAccount.(uint64); ok {
		return &responseData{
			Code:    CodeSuccess,
			Count:   int(affRow),
			Message: "",
			Data: []*partnerIDData{
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
		Data:    []*partnerIDData{{}},
	}
}
