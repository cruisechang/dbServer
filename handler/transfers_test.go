package handler

import (
	"net/http"
	"testing"

	"github.com/cruisechang/dbex"
	"fmt"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"github.com/cruisechang/dbServer/util"
	"strconv"
)

func Test_transfersHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		code  int
		count int
		param transferGetParam
	}{
		{"0", CodeSuccess, 6, transferGetParam{100, -1, -1, -1, "2018-11-28 10:00:00", "2018-12-21 23:59:59"}},
		{"1", CodeSuccess, 6, transferGetParam{100, 100, -1, -1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"2", CodeSuccess, 6, transferGetParam{101, 101, -1, -1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"3", CodeSuccess, 3, transferGetParam{101, 101, 0, -1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"4", CodeSuccess, 1, transferGetParam{101, 101, 0, 1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"5", CodeSuccess, 0, transferGetParam{101, 99999, -1, -1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"6", CodeSuccess, 0, transferGetParam{99999, 101, -1, -1, "2018-11-20 01:01:01", "2018-12-21 23:59:59"}},
		{"7", CodeSuccess, 12, transferGetParam{-1, -1, -1, -1, "2018-11-28 10:00:00", "2018-12-21 23:59:59"}},
		{"8", CodeSuccess, 11, transferGetParam{-1, -1, -1, -1, "2018-11-28 10:00:00", "2018-11-28 15:00:59"}},
		{"9", CodeSuccess, 1, transferGetParam{-1, -1, -1, -1, "2018-11-28 15:00:00", "2018-12-21 23:59:59"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/transfers")

		b, _ := json.Marshal(tc.param)

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
		router.Handle("/transfers", NewTransfersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed ID=%d, got %v want %v,name=%s, path=%s ,param=%+v", tc.param.PartnerID, rr.Code, http.StatusOK, tc.name, path, tc.param)
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

func TestTransfersHandlerPost(t *testing.T) {

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
		param transferPostParam
	}{

		{"0", CodeSuccess, 1, transferPostParam{strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10), 100, 100, 1, 100, 33.3, 2}},
		{"1", CodeSuccess, 1, transferPostParam{strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10), 100, 100, 1, 100, 33.3, 2}},
		{"2", CodeSuccess, 1, transferPostParam{strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10), 100, 100, 1, 100, 33.3, 2}},
		{"3", CodeRequestPostDataIllegal, 0, transferPostParam{strconv.FormatInt(int64(util.RandomInt(1, 9999999)), 10), 100, 100, 1, 100, -1, 2}},
		{"4", CodeSuccess, 1, transferPostParam{"99999", 100, 100, 1, 100, 100.3, 2}},
		{"5", CodeDBExecError, 0, transferPostParam{"99999", 100, 100, 1, 100, 100.3, 2}}, //partner transfer id 重複
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/transfers")
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
		router.Handle(path, NewTransfersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed ID=%d, got %v want %v,name=%s, path=%s ,param=%+v", tc.param.PartnerID, rr.Code, http.StatusOK, tc.name, path, tc.param)
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
