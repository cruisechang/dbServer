package handler

import (
	"github.com/cruisechang/dbServer/util"
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

func Test_transferHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct{ ID uint64 }

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "100", CodeSuccess, 1, http.StatusOK},   //success
		{"1", "101", CodeSuccess, 1, http.StatusOK},   //success
		{"2", "99999", CodeSuccess, 0, http.StatusOK}, //not found
		{"3", "0", CodePathError, 0, http.StatusOK},
		{"4", "-1", CodeSuccess, 0, http.StatusNotFound},
		{"5", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"6", "1.1", CodeSuccess, 0, http.StatusNotFound},
		{"7", "", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/transfers/" + tc.ID)

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
		router.Handle("/transfers/{id:[0-9]+}", NewTransferIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed userID=%s, got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
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

func Test_transferHandler_patchStatus(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		Status int `json:"status"`
	}
	tt := []struct {
		name  string
		ID    string
		code  int
		count int
		param param
	}{
		{"0", "100", CodeSuccess, 1, param{1}},
		{"1", "100", CodeSuccess, 0, param{1}}, //update 內容相同時，count=0
		{"2", "100", CodeSuccess, 1, param{0}},
		{"3", "99999", CodeSuccess, 0, param{3}}, //無此id
		{"4", "abc", CodeSuccess, 0, param{3}},   //404
		{"5", "33.3", CodeSuccess, 0, param{3}},  //404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/transfers/" + tc.ID + "/status")

		b, _ := json.Marshal(tc.param)

		req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/transfers/{id:[0-9]+}/status", NewTransferIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			return
			//t.Errorf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s, path=%s", err.Error(), path)
		}
		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s, path=%s", resData.Code, tc.code, tc.name, path)
		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s, path=%s", resData.Count, tc.count, tc.name, path)

		}
	}
}
