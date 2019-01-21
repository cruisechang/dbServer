package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cruisechang/dbServer/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_broadcastsHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	tt := []struct {
		name       string
		code       int
		count      int
		httpStatus int
	}{
		{"0", CodeSuccess, 4, http.StatusOK},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/broadcasts")

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
		router.Handle(path, NewBroadcastsHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed  got %d want %d, name=%s", rr.Code, tc.httpStatus, tc.name)
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
	}
}

func TestBroadcastsHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	tt := []struct {
		name       string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{

		{"0", CodeSuccess, 1, http.StatusOK, broadcastPostParam{"for test...", 2, 999, 1}},
		{"1", CodeRequestPostDataIllegal, 0, http.StatusOK, broadcastPostParam{"", 0, 999, 1}}, //internal <1, repeat times < 1
		{"2", CodeRequestPostDataIllegal, 0, http.StatusOK, broadcastPostParam{"", 2, 0, 1}},   //internal <1, repeat times < 1
		{"3", CodeRequestDataUnmarshalError, 0, http.StatusOK, ""},
		{"4", CodeRequestDataUnmarshalError, 0, http.StatusOK, 1},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/broadcasts")
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
		router.Handle(path, NewBroadcastsHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("http status got %d want %d, name=%s", rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Fatalf("resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}
