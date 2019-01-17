package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"strconv"
	"strings"
	"database/sql"
	"encoding/json"
	"errors"
)

func NewPartnerIDHandler(base baseHandler) *partnerIDHandler {
	return &partnerIDHandler{
		baseHandler: base,
	}
}

type partnerIDHandler struct {
	baseHandler
}

func (h *partnerIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "partnerIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	gotVar, err := h.getVariable(r, "id")
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id error",logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	id, err := strconv.ParseUint(gotVar, 10, 64)
	if err != nil {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id to uint64 error id=%s", logPrefix,gotVar))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if id <= 0 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get id ==0 ",logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "aesKey") {

			queryString := "SELECT  aes_key from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "aes_key", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "aes_key", queryString, h.returnTargetColumnResDataCount)
			return
		}
		if strings.Contains(r.URL.Path, "active") {

			queryString := "SELECT  active from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "active", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "active", queryString, h.returnTargetColumnResDataCount)
			return
		}
		if strings.Contains(r.URL.Path, "login") {

			queryString := "SELECT  login from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "login", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "login", queryString, h.returnTargetColumnResDataCount)
			return
		}
		if strings.Contains(r.URL.Path, "apiBindIP") {

			queryString := "SELECT  api_bind_ip from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "api_bind_ip", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "api_bind_ip", queryString, h.returnTargetColumnResDataCount)
			return
		}
		if strings.Contains(r.URL.Path, "cmsBinkIP") {

			queryString := "SELECT  cms_bind_ip from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "cms_bind_ip", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "cms_bind_ip", queryString, h.returnTargetColumnResDataCount)
			return
		}

		if strings.Contains(r.URL.Path, "accessToken") {

			queryString := "SELECT  access_token from partner WHERE partner_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, id, "access_token", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
			//h.getTargetColumnValue(w, r, "partner", id, "access_token", queryString, h.returnTargetColumnResDataCount)
			return
		}

		queryString := "SELECT partner_id,account,name,level,category,active,api_bind_ip,cms_bind_ip,create_date from partner WHERE partner_id = ? LIMIT 1"
		//h.getTargetRow(w, r, "partner", id, queryString, h.returnResDataFunc)
		h.dbQuery(w, r, logPrefix, id, "access_token", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		//check body
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, "partnerIDHandler patchTargetColumn "+errMsg)
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		column := ""
		queryString := ""

		if strings.Contains(r.URL.Path, "login") {
			column = "login"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "active") {
			column = "active"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "aesKey") {
			column = "aes_key"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "apiBindIP") {
			column = "api_bind_ip"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "cmsBindIP") {
			column = "cms_bind_ip"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else if strings.Contains(r.URL.Path, "accessToken") {
			column = "access_token"
			queryString = "UPDATE partner set " + column + "  = ?  WHERE partner_id = ? LIMIT 1"
		} else {
			//error
			//h.logger.Log(dbex.LevelError, fmt.Sprintf("%s  patch path error :%s", logPrefix, r.URL.Path))
			//h.writeError(w, http.StatusNotFound, CodePathError, "")
			//return

			queryString = "UPDATE partner set password=?, name = ?, level=?, category=? , active=? , aes_key=?, access_token=?, api_bind_ip =? , cms_bind_ip =?  WHERE partner_id = ? LIMIT 1"
		}


		//unmarshal request body
		patchData, err := h.getPatchData(column, body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix,err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, err.Error())
			return
		}
		h.patch(w, r, logPrefix, id, queryString, patchData, h.patchExec, h.returnIDResData)
		return
	}

	h.writeError(w, http.StatusNotFound, CodeMethodError, "")
}
func (h *partnerIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}

func (h *partnerIDHandler) returnTargetColumnResponseData() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "aes_key":
			resData := []aesKeyData{}
			count := 0
			var aesKey string
			for rows.Next() {
				err := rows.Scan(&aesKey)
				if err == nil {
					count += 1
					resData = append(resData,
						aesKeyData{
							aesKey,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "active":
			resData := []activeData{}
			count := 0
			var active uint
			for rows.Next() {
				err := rows.Scan(&active)
				if err == nil {
					count += 1
					resData = append(resData,
						activeData{
							active,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "login":
			resData := []loginData{}
			count := 0
			var login uint
			for rows.Next() {
				err := rows.Scan(&login)
				if err == nil {
					count += 1
					resData = append(resData,
						loginData{
							login,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "api_bind_ip":
			resData := []apiBindIPData{}
			count := 0
			var ip string
			for rows.Next() {
				err := rows.Scan(&ip)
				if err == nil {
					count += 1
					resData = append(resData,
						apiBindIPData{
							ip,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "cms_bind_ip":
			resData := []cmsBindIPData{}
			count := 0
			var ip string
			for rows.Next() {
				err := rows.Scan(&ip)
				if err == nil {
					count += 1
					resData = append(resData,
						cmsBindIPData{
							ip,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "access_token":
			resData := []accessTokenData{}
			count := 0
			var token string
			for rows.Next() {
				err := rows.Scan(&token)
				if err == nil {
					count += 1
					resData = append(resData,
						accessTokenData{
							token,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		default:
			return &responseData{}
		}
	}
}
func (h *partnerIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := partnerDB{}
		resData := []partnerData{}

		for rows.Next() {
			err := rows.Scan(&ud.partner_id, &ud.account, &ud.name, &ud.level, &ud.category, &ud.active, &ud.api_bind_ip, &ud.cms_bind_ip, &ud.create_date)
			if err == nil {
				count ++
				resData = append(resData,
					partnerData{
						ud.partner_id,
						ud.account,
						ud.name,
						ud.level,
						ud.category,
						ud.active,
						ud.api_bind_ip,
						ud.cms_bind_ip,
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
/*
func (h *partnerIDHandler) returnResDataFunc() (func(rows *sql.Rows) (interface{}, int)) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		ud := partnerDB{}
		resData := []partnerData{}

		for rows.Next() {
			err := rows.Scan(&ud.partner_id, &ud.account, &ud.name, &ud.level, &ud.category, &ud.active, &ud.api_bind_ip, &ud.cms_bind_ip, &ud.create_date)
			if err == nil {
				count += 1
				resData = append(resData,
					partnerData{
						ud.partner_id,
						ud.account,
						ud.name,
						ud.level,
						ud.category,
						ud.active,
						ud.api_bind_ip,
						ud.cms_bind_ip,
						ud.create_date,})
			}
		}
		return resData, count
	}
}

func (h *partnerIDHandler) returnTargetColumnResDataCount(column string, rows *sql.Rows) (interface{}, int) {

	switch column {
	case "aes_key":
		resData := []aesKeyData{}
		count := 0
		var aesKey string
		for rows.Next() {
			err := rows.Scan(&aesKey)
			if err == nil {
				count += 1
				resData = append(resData,
					aesKeyData{
						aesKey,
					})
			}
		}
		return resData, count
	case "active":
		resData := []activeData{}
		count := 0
		var active uint
		for rows.Next() {
			err := rows.Scan(&active)
			if err == nil {
				count += 1
				resData = append(resData,
					activeData{
						active,
					})
			}
		}
		return resData, count
	case "login":
		resData := []loginData{}
		count := 0
		var login uint
		for rows.Next() {
			err := rows.Scan(&login)
			if err == nil {
				count += 1
				resData = append(resData,
					loginData{
						login,
					})
			}
		}
		return resData, count
	case "api_bind_ip":
		resData := []apiBindIPData{}
		count := 0
		var ip string
		for rows.Next() {
			err := rows.Scan(&ip)
			if err == nil {
				count += 1
				resData = append(resData,
					apiBindIPData{
						ip,
					})
			}
		}
		return resData, count
	case "cms_bind_ip":
		resData := []cmsBindIPData{}
		count := 0
		var ip string
		for rows.Next() {
			err := rows.Scan(&ip)
			if err == nil {
				count += 1
				resData = append(resData,
					cmsBindIPData{
						ip,
					})
			}
		}
		return resData, count
	case "access_token":
		resData := []accessTokenData{}
		count := 0
		var token string
		for rows.Next() {
			err := rows.Scan(&token)
			if err == nil {
				count += 1
				resData = append(resData,
					accessTokenData{
						token,
					})
			}
		}
		return resData, count
	default:
		return "[{}]", 0
	}
}
*/

func (h *partnerIDHandler) getPatchData(column string, body []byte) (interface{}, error) {
	switch column {

	case "login":
		ug := &loginData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "active":
		ug := &activeData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "aes_key":
		ug := &aesKeyData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "api_bind_ip":
		ug := &apiBindIPData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "cms_bind_ip":
		ug := &cmsBindIPData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "access_token":
		ug := &accessTokenData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
	case "":
		d:=&partnerPatchParam{}
		err := json.Unmarshal(body, d)
		if err != nil {
			return nil, err
		}
		if len(d.Password)<3 || len(d.Name)<3 || len(d.AccessToken)<3 || len(d.AESKey)<3 {
			return nil, errors.New("")
		}
		return d, nil
	default:
		return nil, errors.New("column error")
	}
}

func (h *partnerIDHandler) patchExec(stmt *sql.Stmt, ID uint64, param interface{}) (sql.Result, error) {

	//檢查參數是否合法

	if p, ok := param.(*loginData); ok {
		return stmt.Exec(p.Login, ID)
	}
	if p, ok := param.(*activeData); ok {
		return stmt.Exec(p.Active, ID)
	}
	if p, ok := param.(*aesKeyData); ok {
		return stmt.Exec(p.AESKey, ID)
	}
	if p, ok := param.(*apiBindIPData); ok {
		return stmt.Exec(p.APIBindIP, ID)
	}
	if p, ok := param.(*cmsBindIPData); ok {
		return stmt.Exec(p.CMSBindIP, ID)
	}
	if p, ok := param.(*accessTokenData); ok {
		return stmt.Exec(p.AccessToken, ID)
	}

	if p,ok:=param.(*partnerPatchParam);ok{
		return stmt.Exec(p.Password,p.Name,p.Level,p.Category,p.Active,p.AESKey,p.AccessToken,p.APIBindIP,p.CMSBindIP, ID)
	}

	return nil, errors.New("parsing param error")
}

func (h *partnerIDHandler) returnIDResData(ID uint64) interface{} {
	return []partnerIDData{
		{
			ID,
		},
	}
}
