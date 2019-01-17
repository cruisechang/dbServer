package handler

import (
	"net/http"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_userAccountHandler_getPassword(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		account    string
		code       int
		count      int
		httpStatus int
		param      userAccountGetParam
	}{
		{"0", "account100", CodeSuccess, 1, http.StatusOK, userAccountGetParam{100, "pass"}},
		{"1", "99999", CodeSuccess, 0, http.StatusOK, userAccountGetParam{100, "pass"}},               //account not found
		{"2", "account100", CodeSuccess, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}}, //partnerID not found
		{"3", "xxx", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},      //account 太短
		{"4", "0", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},        //account 太短
		{"5", "-1", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},       //account 太短
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently, userAccountGetParam{987656776553, "pass"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.account + "/password")

		//t.Logf("path %s", path)
		//t.Logf("tc %+v", tc)
		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("handerl unmarshal param error, param=%+v", tc)
		}

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/users/{account}/password", NewUserAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
		t.Logf("resData=%+v", resData)
	}
}

func Test_userAccountHandler_getUserID(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		account    string
		code       int
		count      int
		httpStatus int
		param      userAccountGetParam
	}{
		{"0", "account100", CodeSuccess, 1, http.StatusOK, userAccountGetParam{100, "pass"}},
		{"1", "99999", CodeSuccess, 0, http.StatusOK, userAccountGetParam{100, "pass"}},               //account not found
		{"2", "account100", CodeSuccess, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}}, //partnerID not found
		{"3", "xxx", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},      //account 太短
		{"4", "0", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},        //account 太短
		{"5", "-1", CodePathError, 0, http.StatusOK, userAccountGetParam{987656776553, "pass"}},       //account 太短
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently, userAccountGetParam{987656776553, "pass"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.account + "/id")

		//t.Logf("path %s", path)
		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("handerl unmarshal param error, param=%+v", tc)
		}

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/users/{account}/id", NewUserAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
		t.Logf("resData=%+v", resData)
	}
}

func Test_userAccountHandler_getAccessToken(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name    string
		account string
		code    int
		count   int
		httpStatus int
		param   userAccountGetParam
	}{
		{"0", "account1", CodeSuccess, 1, http.StatusOK,userAccountGetParam{100, "pass"}},
		{"1", "99999", CodeSuccess, 0, http.StatusOK,userAccountGetParam{100, "pass"}}, //active 0
		{"2", "account1", CodeSuccess, 0, http.StatusOK,userAccountGetParam{0, "pass"}},     //partnerID not found
		{"3", "xxx", CodePathError, 0, http.StatusOK,userAccountGetParam{0, "pass"}},     //account 太短
		{"4", "0", CodePathError, 0, http.StatusOK,userAccountGetParam{100, "pass"}},     //account 太短
		{"5", "-1", CodePathError, 0, http.StatusOK,userAccountGetParam{9999, "pass"}},   //account 太短
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently,userAccountGetParam{9999, "pass"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.account + "/accesstoken")

		//t.Logf("path %s", path)
		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("handerl unmarshal param error, param=%+v", tc)
		}

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/users/{account}/accesstoken", NewUserAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
		t.Logf("resData=%+v", resData)
	}
}
