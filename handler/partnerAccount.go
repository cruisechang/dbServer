package handler

import (
	"net/http"
	"github.com/cruisechang/dbex"
	"fmt"
	"github.com/gorilla/mux"
	"strings"
	"database/sql"
	"encoding/json"
)

func NewPartnerAccountHandler(base baseHandler) *PartnerAccountHandler {
	return &PartnerAccountHandler{
		baseHandler: base,
	}
}

type PartnerAccountHandler struct {
	baseHandler
}

func (h *PartnerAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "partnerAccountHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()


	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s get account not found", logPrefix))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	if len(account) < 4 {
		h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s account length =%d ", logPrefix, len(account)))
		h.writeError(w, http.StatusOK, CodePathError, "")
		return
	}

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "password") {
			queryString := "SELECT password  FROM partner where account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "password", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
			return

		} else if strings.Contains(r.URL.Path, "id") {
			queryString := "SELECT partner_id  FROM partner where account = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "partner_id", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
			return

		}else if strings.Contains(r.URL.Path, "login") {

			body, errCode, errMsg := h.checkBody(w, r)
			if errCode != CodeSuccess {
				h.logger.Log(dbex.LevelError, errMsg)
				h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
				return
			}

			param := &partnerAccountGetParam{}
			err := json.Unmarshal(body, param)
			if err != nil {
				h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s  get data unmarshal error=%s", logPrefix, err.Error()))
				h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
				return
			}

			queryString := "SELECT partner_id,account,name,level,category,active,api_bind_ip,cms_bind_ip,create_date from partner WHERE  account = ? AND password=? LIMIT 1"
			h.dbQuery(w, r, logPrefix, account, "login", queryString, param, h.sqlQuery, h.returnResponseDataFunc)
			return
		}
	}

	h.writeError(w, http.StatusOK, CodeMethodError, "")
}
func (h *PartnerAccountHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

	if p,ok:=param.(*partnerAccountGetParam);ok{
		return stmt.Query(IDOrAccount,p.Password)
	}
	return stmt.Query(IDOrAccount)
}

func (h *PartnerAccountHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "password":
			resData := []passwordData{}
			count := 0
			var password string
			for rows.Next() {
				err := rows.Scan(&password)
				if err == nil {
					count ++
					resData = append(resData,
						passwordData{
							password,
						})
				}
			}
			return &responseData{
				Code:    CodeSuccess,
				Count:   count,
				Message: "",
				Data:    resData,
			}
		case "partner_id":
			resData := []partnerIDData{}
			count := 0
			var id uint64
			for rows.Next() {
				err := rows.Scan(&id)
				if err == nil {
					count ++
					resData = append(resData,
						partnerIDData{
							id,
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
		default:
			return &responseData{}
		}
	}
}

