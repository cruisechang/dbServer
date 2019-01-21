package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_officialCMSManagerAccountHandler_getCheckLoginData(t *testing.T) {
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
		{"0", "test01", CodeSuccess, 1, http.StatusOK},
		{"1", "99999", CodeSuccess, 0, http.StatusOK},
		{"2", "3.3", CodePathError, 0, http.StatusOK},            //account 太短
		{"3", "0", CodePathError, 0, http.StatusOK},              //account 太短
		{"4", "", CodePathError, 0, http.StatusMovedPermanently}, //account 太短
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSManagers/" + tc.account + "/login")

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
		router.Handle("/officialCMSManagers/{account}/login", NewOfficialCMSManagerAccountHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed account=%s, http status got %d want %d, name=%s", tc.account, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code == http.StatusNotFound || rr.Code == http.StatusMovedPermanently {
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
