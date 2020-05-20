package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
)

//NewUserIDHandler returns *UserIDHandler
func NewUserIDHandler(base baseHandler) *UserIDHandler {
	return &UserIDHandler{
		baseHandler: base,
	}
}

//UserIDHandler presents structure of user id
type UserIDHandler struct {
	baseHandler
}

func (h *UserIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "userIDHandler"

	defer func() {
		if r := recover(); r != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s panic=%v", logPrefix, r))
			h.writeError(w, http.StatusOK, CodePanic, fmt.Sprintf("%s panic %v", logPrefix, r))
		}
	}()

	vars := mux.Vars(r)
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

	//get
	if r.Method == "GET" || r.Method == "get" {
		if strings.Contains(r.URL.Path, "credit") {
			queryString := "SELECT  credit from user WHERE user_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, ID, "credit", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
		} else if strings.Contains(r.URL.Path, "active") {
			queryString := "SELECT  active from user WHERE user_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, ID, "active", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
		} else if strings.Contains(r.URL.Path, "login") {
			queryString := "SELECT  login from user WHERE user_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, ID, "login", queryString, nil, h.sqlQuery, h.returnTargetColumnResponseData)
		} else {
			queryString := "SELECT user_id,partner_id,account,name,credit,level,category,active,ip,platform,login,create_date from user WHERE user_id = ? LIMIT 1"
			h.dbQuery(w, r, logPrefix, ID, "", queryString, nil, h.sqlQuery, h.returnResponseDataFunc)
		}
		return
	}

	if r.Method == "PATCH" || r.Method == "patch" {
		//check body
		body, errCode, errMsg := h.checkBody(w, r)
		if errCode != CodeSuccess {
			h.logger.Log(dbex.LevelError, fmt.Sprintf("%s  patchTargetColumn %s", logPrefix, errMsg))
			h.writeError(w, http.StatusOK, CodeBodyError0, errMsg)
			return
		}

		column := ""

		if strings.Contains(r.URL.Path, "credit") {
			column = "credit"
		} else if strings.Contains(r.URL.Path, "login") {
			column = "login"
		} else if strings.Contains(r.URL.Path, "active") {
			column = "active"
		} else {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patch path error :%s", logPrefix, r.URL.Path))
			h.writeError(w, http.StatusOK, CodeRequestPathError, "")
			return
		}
		queryString := "UPDATE user set " + column + "  = ?  WHERE user_id = ? LIMIT 1"

		//unmarshal request body
		param, err := h.getPatchData(column, body)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s patchTargetColumn data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		//h.patch(w, r, logPrefix, id, queryString, patchData, h.patchExec, h.returnIDResData)
		h.dbExec(w, r, logPrefix, ID, column, queryString, param, h.sqlPatch, h.returnExecResponseData)
		return
	}
	h.writeError(w, http.StatusNotFound, CodeMethodError, "")
}

func (h *UserIDHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {
	return stmt.Query(IDOrAccount)
}
func (h *UserIDHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		count := 0
		ud := userDB{}
		resData := []userData{}

		for rows.Next() {
			err := rows.Scan(&ud.user_id, &ud.partner_id, &ud.account, &ud.name, &ud.credit, &ud.level, &ud.category, &ud.active, &ud.ip, &ud.platform, &ud.login, &ud.create_date)
			if err == nil {
				count++
				resData = append(resData,
					userData{
						ud.user_id,
						ud.partner_id,
						ud.account,
						ud.name,
						ud.credit,
						ud.level,
						ud.category,
						ud.active,
						ud.ip,
						ud.platform,
						ud.login,
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

func (h *UserIDHandler) returnTargetColumnResponseData() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

	return func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {
		switch targetColumn {
		case "credit":
			resData := []creditData{}
			count := 0
			var credit float32
			for rows.Next() {
				err := rows.Scan(&credit)
				if err == nil {
					count++
					resData = append(resData,
						creditData{
							credit,
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
					count++
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
					count++
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
		default:
			return &responseData{}
		}
	}
}

//////
func (h *UserIDHandler) returnResDataFunc() func(rows *sql.Rows) (interface{}, int) {

	return func(rows *sql.Rows) (interface{}, int) {
		count := 0
		ud := userDB{}
		resData := []userData{}

		for rows.Next() {
			err := rows.Scan(&ud.user_id, &ud.partner_id, &ud.account, &ud.name, &ud.credit, &ud.level, &ud.category, &ud.active, &ud.ip, &ud.platform, &ud.login, &ud.create_date)
			if err == nil {
				count++
				resData = append(resData,
					userData{
						ud.user_id,
						ud.partner_id,
						ud.account,
						ud.name,
						ud.credit,
						ud.level,
						ud.category,
						ud.active,
						ud.ip,
						ud.platform,
						ud.login,
						ud.create_date})
			}
		}
		return resData, count
	}
}

//patch
func (h *UserIDHandler) getPatchData(column string, body []byte) (interface{}, error) {
	switch column {
	case "credit":
		ug := &creditData{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			return nil, err
		}
		return ug, nil
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
	default:
		return nil, errors.New("column error")
	}
}
func (h *UserIDHandler) sqlPatch(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*creditData); ok {
		return stmt.Exec(p.Credit, IDOrAccount)
	}
	if p, ok := param.(*loginData); ok {
		return stmt.Exec(p.Login, IDOrAccount)
	}
	if p, ok := param.(*activeData); ok {

		return stmt.Exec(p.Active, IDOrAccount)
	}

	return nil, errors.New("parsing param error")
}

func (h *UserIDHandler) returnExecResponseData(IDOrAccount interface{}, column string, result sql.Result) *responseData {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*userIDData{{}},
		}
	}

	ID, _ := IDOrAccount.(uint64)

	return &responseData{
		Code:    CodeSuccess,
		Count:   int(affRow),
		Message: "",
		Data: []*userIDData{
			{
				ID,
			},
		},
	}
}
