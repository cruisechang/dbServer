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

func NewUsersHandler(base baseHandler) *usersHandler {
	return &usersHandler{
		baseHandler: base,
	}
}

type usersHandler struct {
	baseHandler
}

func (h *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logPrefix := "usersHandler"

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
		ug := &userGetParam{}
		err := json.Unmarshal(body, ug)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler get data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if ug.OrderBy != "userID" && ug.OrderBy != "" {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  get data orderBy error, orderBy=%s", logPrefix, ug.OrderBy))
			h.writeError(w, http.StatusOK, CodeRequestDataError, "")
			return
		}
		if ug.Order != "asc" && ug.Order != "desc" && ug.Order != "" {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  get data order error, order=%s", logPrefix, ug.Order))
			h.writeError(w, http.StatusOK, CodeRequestDataError, "")
			return
		}

		queryString, queryArgs := h.getQueryStringArgs(ug)
		h.dbQuery(w, r, logPrefix, 0, "", queryString, queryArgs, h.sqlQuery, h.returnResponseDataFunc)
		return
	}
	if r.Method == "POST" {
		param := &userPostParam{}
		err := json.Unmarshal(body, param)
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  post data unmarshal error=%s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestDataUnmarshalError, "")
			return
		}

		if len(param.Account) < 4 || len(param.Password) < 4 || len(param.Name) < 4 || param.PartnerID < 0 {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  post data illegal=%+v", logPrefix, param))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, "post data illegal")
			return
		}
		ID, err := util.GetUniqueID()
		if err != nil {
			h.logger.LogFile(dbex.LevelError, fmt.Sprintf("%s handler  post get unique userID error %s", logPrefix, err.Error()))
			h.writeError(w, http.StatusOK, CodeRequestPostDataIllegal, fmt.Sprintf("%s handler post get unique userID error %s", logPrefix, err.Error()))
		}

		queryString := "INSERT  INTO user (user_id,partner_id,account,password,name,ip,platform) values (? ,? ,? ,?, ? ,?, ?)"
		//h.post(w, r, "users", userID, queryString, param, h.sqlExec, h.returnPostResData)
		h.dbExec(w, r, logPrefix, ID, "", queryString, param, h.sqlPost, h.returnPostResponseData)
		return
	}
	h.writeError(w, http.StatusOK, CodeMethodError, "")
}

//query
func (h *usersHandler) getQueryStringArgs(ug *userGetParam) (queryString string, queryArgs []interface{}) {

	queryString = "select user_id,partner_id,account,name,credit,level,category,active,ip,platform,login,create_date from user "

	//partnerID,status
	if ug.PartnerID > -1 && ug.Active > -1 {
		queryString += "WHERE partner_id = ? AND active = ? "
		queryArgs = append(queryArgs, ug.PartnerID, ug.Active)
	} else if ug.PartnerID > -1 {
		queryString += "WHERE partner_id = ? "
		queryArgs = append(queryArgs, ug.PartnerID)
	} else if ug.Active > -1 {
		queryString += "WHERE active = ? "
		queryArgs = append(queryArgs, ug.Active)
	}

	//orderBy, order
	if ug.OrderBy != "" {
		queryString += " ORDER BY ? "
		queryArgs = append(queryArgs, "user_id")

		if ug.Order == "asc" {
			queryString += " asc "
		} else if ug.Order == "desc" {
			queryString += " desc "
		}
	}
	//limit, offset
	if ug.Limit > 0 {
		queryString += " LIMIT ? "
		queryArgs = append(queryArgs, ug.Limit)
		if ug.Offset > 0 {
			queryString += " OFFSET ?"
			queryArgs = append(queryArgs, ug.Offset)

		}
	}
	return
}

func (h *usersHandler) sqlQuery(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (*sql.Rows, error) {

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
func (h *usersHandler) returnResponseDataFunc() func(IDOrAccount interface{}, targetColumn string, rows *sql.Rows) *responseData {

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

//post
func (h *usersHandler) sqlPost(stmt *sql.Stmt, IDOrAccount interface{}, param interface{}) (sql.Result, error) {

	if p, ok := param.(*userPostParam); ok {

		return stmt.Exec(IDOrAccount, p.PartnerID, p.Account, p.Password, p.Name, p.IP, p.Platform)
	}
	return nil, errors.New("parsing param error")

}

//id預先產生
func (h *usersHandler) returnPostResponseData(IDOrAccount interface{}, column string, result sql.Result) (*responseData) {

	affRow, err := result.RowsAffected()
	if err != nil {
		return &responseData{
			Code:    CodeDBExecResultError,
			Count:   0,
			Message: "",
			Data:    []*userIDData{{}},
		}

	}

	if id, ok := IDOrAccount.(uint64); ok {
		return &responseData{
			Code:    CodeSuccess,
			Count:   int(affRow),
			Message: "",
			Data: []*userIDData{
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
		Data:    []*userIDData{{}},
	}
}

