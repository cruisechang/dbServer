package handler

import (
	"testing"
	"fmt"
	"github.com/cruisechang/dbex"
	"net/http"
	"bytes"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/cruisechang/dbServer/util"
	"strconv"
	"io/ioutil"
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

		{"0", 0, 50, userGetParam{100, -1, -1, -1, "", ""},},
		{"1", 0, 2, userGetParam{100, 1, -1, -1, "", ""},},
		{"2", 0, 50, userGetParam{101, -1, -1, -1, "", ""},},
		{"3", 0, 1, userGetParam{102, -1, 0, -1, "", ""},},
		{"4", 0, 0, userGetParam{999, 1, 9, -1, "", ""},},
		{"5", 0, 2, userGetParam{-1, 1, 19, -1, "", ""},},
		{"6", 0, 0, userGetParam{101, 2, 29, -1, "userID", "asc"},},
		{"7", 0, 0, userGetParam{100, 2, 29, -1, "userID", "desc"},},
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
	fmt.Sprintf("%v", dbx)

	accounts := []string{}
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10))

	tt := []struct {
		PartnerID uint64 `json:"partnerID"`
		Account   string `json:"account"`
		Password  string `json:"password"`
		Name      string `json:"name"`
		IP        string `json:"ip"`
		Platform  int    `json:"platform"`
	}{

		{0, accounts[0], "passd", accounts[0], "111.111.111.111", 0},
		{0, accounts[1], "passd", accounts[0], "111.111.111.111", 0},
		{0, accounts[2], "passd", accounts[0], "111.111.111.111", 0},
		{0, accounts[3], "passd", accounts[0], "111.111.111.111", 0},
		{0, accounts[4], "passd", accounts[0], "111.111.111.111", 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users")
		b, _ := json.Marshal(tc)

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

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusOK {
			t.Errorf("handler should have failed on  partnerID=%d, got %v want %v", tc.PartnerID, rr.Code, http.StatusOK)
		}
	}
}
