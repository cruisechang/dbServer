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

func Test_partnerAccountHandler_getPassword(t *testing.T) {
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
	}{
		{"0", "account100", CodeSuccess, 1, http.StatusOK},
		{"1", "account101", CodeSuccess, 1, http.StatusOK},
		{"2", "account102", CodeSuccess, 1, http.StatusOK},
		{"3", "8838987676", CodeSuccess, 0, http.StatusOK},
		{"4", "3.3", CodePathError, 0, http.StatusOK},
		{"5", "", CodePathError, 0, http.StatusMovedPermanently},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.account + "/password")

		t.Logf("path %s", path)

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.Handle("/partners/{account}/password", NewPartnerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, http status got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
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
		t.Logf("%+v", resData)
	}
}

func Test_partnerAccountHandler_getID(t *testing.T) {
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
	}{
		{"0", "account100", CodeSuccess, 1, http.StatusOK},
		{"1", "account101", CodeSuccess, 1, http.StatusOK},
		{"2", "account102", CodeSuccess, 1, http.StatusOK},
		{"3", "8838987676", CodeSuccess, 0, http.StatusOK},
		{"4", "3.3", CodePathError, 0, http.StatusOK},
		{"5", "", CodePathError, 0, http.StatusMovedPermanently},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.account + "/id")

		t.Logf("path %s", path)

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.Handle("/partners/{account}/id", NewPartnerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, http status got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
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
		t.Logf("%+v", resData)
	}
}

func Test_partnerAccountHandler_login(t *testing.T) {
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
		param      interface{}
	}{
		{"0", "account100", CodeSuccess, 1, http.StatusOK, partnerAccountGetParam{Password: "pass"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.account + "/login")

		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf(" unmarshal param error, param=%+v", tc)
		}

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.Handle("/partners/{account}/login", NewPartnerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, http status got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
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
		t.Logf("%+v", resData)
	}
}
