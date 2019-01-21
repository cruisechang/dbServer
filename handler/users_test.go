package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cruisechang/dbServer/util"
	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func TestUsersHandlerGet(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)
	tt := []struct {
		name  string
		code  int
		count int
		param userGetParam

		//PartnerID uint64 `json:"partnerID"`
		//Limit     int    `json:"limit"`
		//Offset    int    `json:"offset"`
		//Status    int    `json:"status"`
		//OrderBy   string `json:"orderBy"`
		//Order     string `json:"order"`
	}{

		{"0", 0, 50, userGetParam{100, -1, -1, -1, "", ""}},
		{"1", 0, 2, userGetParam{100, 1, -1, -1, "", ""}},
		{"2", 0, 50, userGetParam{101, -1, -1, -1, "", ""}},
		{"3", 0, 1, userGetParam{102, -1, 0, -1, "", ""}},
		{"4", 0, 0, userGetParam{999, 1, 9, -1, "", ""}},
		{"5", 0, 2, userGetParam{-1, 1, 19, -1, "", ""}},
		{"6", 0, 0, userGetParam{101, 2, 29, -1, "userID", "asc"}},
		{"7", 0, 0, userGetParam{100, 2, 29, -1, "userID", "desc"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users")
		b, _ := json.Marshal(tc.param)

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		//req.Header.Set("Content-Type","application/x-www-form-urlencoded; param=value")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/users", NewUsersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}
		t.Logf("resData count=%d", resData.Count)

	}
}

func TestUsersHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	//db
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	var ids []uint64 //放ids，刪掉用

	accounts := []string{}
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))

	tt := []struct {
		name  string
		code  int
		count int
		param interface{}
	}{
		{"0", CodeSuccess, 1, userPostParam{0, accounts[0], "pass", accounts[0], "111.111.111.111", 0}},
		{"1", CodeDBExecError, 0, userPostParam{0, accounts[0], "passd", accounts[0], "111.111.111.111", 0}}, //partner + account 跟第一組重複
		{"2", CodeRequestDataUnmarshalError, 0, 0},
		{"3", CodeRequestDataUnmarshalError, 0, ""},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users")
		b, _ := json.Marshal(tc.param)

		req, err := http.NewRequest("POST", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/users", NewUsersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler should have failed on  httpStatus got %v want %v", rr.Code, http.StatusOK)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*userIDData
		}{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}

		//insert success
		if resData.Count == 1 {
			t.Logf("ID=%d ", resData.Data[0].UserID)
			ids = append(ids, resData.Data[0].UserID)
		}
	}

	if len(ids) > 0 {
		queryString := "DELETE FROM user  where user_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
