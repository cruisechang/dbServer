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

func Test_betIDHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	type param struct{ ID string }

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1000001", CodeSuccess, 1, http.StatusOK},     //success
		{"1", "1000002", CodeSuccess, 1, http.StatusOK},     //success
		{"2", "99999", CodeSuccess, 0, http.StatusOK},       //not found
		{"3", "xxxxx", CodeSuccess, 0, http.StatusNotFound}, //非int字串 gorilla mux 回傳404
		{"4", "3.3", CodeSuccess, 0, http.StatusNotFound},   //非int字串 gorilla mux 404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/bets/" + tc.ID)

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
		router.Handle("/bets/{id:[0-9]+}", NewBetIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed userID=%s, got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code == http.StatusNotFound {
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
	}
}

func Test_betHandler_patchStatus(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	tt := []struct {
		name       string
		httpStatus int
		ID         string
		code       int
		count      int
		param      statusData
	}{
		{"0", http.StatusOK, "1000001", CodeSuccess, 1, statusData{0}},
		{"1", http.StatusOK, "1000001", CodeSuccess, 0, statusData{0}},
		{"2", http.StatusOK, "1000001", CodeSuccess, 1, statusData{1}},
		{"3", http.StatusOK, "99999", CodeSuccess, 0, statusData{2}},      //無此id
		{"4", http.StatusNotFound, "qwe3", CodeSuccess, 0, statusData{2}}, //非整數字串
		{"5", http.StatusNotFound, "33.3", CodeSuccess, 0, statusData{2}}, //非整數字串
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/bets/" + tc.ID + "/status")

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
		router.Handle("/bets/{id:[0-9]+}/status", NewBetIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, tc.httpStatus, tc.name, path)
			return
		}

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed userID=%s, got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code == http.StatusNotFound {
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
	}
}
