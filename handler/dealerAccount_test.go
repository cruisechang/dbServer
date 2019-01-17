package handler

import (
	"net/http"
	"testing"

	"github.com/cruisechang/dbex"
	"fmt"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

func Test_dealerAccountHandler_getPassword(t *testing.T) {
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
	}{
		{"0", "test", 0, 1},
		{"1", "99999", 0, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.account + "/password")

		t.Logf("path %s", path)

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/dealers/{account}/password", NewDealerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed account=%s, got %v want %v,name=%s", tc.account, rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{
		}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

func Test_dealerAccountHandler_getDealerID(t *testing.T) {
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
	}{
		{"0", "test", 0, 1},
		{"1", "99999", 0, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.account + "/id")

		t.Logf("path %s", path)

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/dealers/{account}/id", NewDealerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed account=%s, got %v want %v,name=%s", tc.account, rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{
		}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

func Test_dealerAccountHandler_getLogin(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name    string
		code    int
		count   int
		account string
		param   interface{}
	}{
		{"0", CodeSuccess, 1, "test", dealerAccountGetParam{"1234"}},    //okd
		{"1", CodeSuccess, 0, "test999", dealerAccountGetParam{"1234"}}, //account not found
		{"2", CodeSuccess, 0, "test", dealerAccountGetParam{"xxxx"}},    //pass not found
		{"3", CodeSuccess, 0, "test1", dealerAccountGetParam{"1234"}},   //active !=1
		{"4", CodeSuccess, 0, "9999", dealerAccountGetParam{"1234"}},    //account not found
		{"5", CodeSuccess, 0, "33.3", dealerAccountGetParam{"1234"}},    //account not found
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.account + "/login")

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
		router.Handle("/dealers/{account}/login", NewDealerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed account=%s, got %v want %v,name=%s, path=%s, param%+v", tc.account, rr.Code, http.StatusOK, tc.name, path, tc.param)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{
		}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s, path=%s, param=%+v", err.Error(), path, tc.param)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s, param=%+v", resData.Code, tc.code, tc.name, path, tc.param)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s, path=%s, param=%+v", resData.Count, tc.count, tc.name, path, tc.param)

		}
	}
}
