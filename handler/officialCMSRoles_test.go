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
)

func Test_officialCMSRolesHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		code  int
		httpStatus int
	}{
		{"0", CodeSuccess,http.StatusOK},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSRoles")

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
		router.Handle(path, NewOfficialCMSRolesHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed  got %d want %d, name=%s", rr.Code, tc.httpStatus, tc.name)
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

		t.Logf("resData=%+v",resData)
	}
}

func Test_officialCMSRolesHandler_post(t *testing.T) {

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
		param interface{}
	}{

		{"0", CodeSuccess, 1, officialCMSRolePostParam{"[999,999]"}},
		{"1", CodeRequestDataUnmarshalError, 0, officialCMSRolePostParam{"999,999"}},
		{"2", CodeRequestDataUnmarshalError, 0, struct{X int }{ 1}}, //param error
		{"3", CodeRequestDataUnmarshalError, 0, struct{Str string }{ ""}}, //param error
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSRoles")
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
		router.Handle(path, NewOfficialCMSRolesHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			//mux parsing route error
			return
		}

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed ,got %v want %v,name=%s, path=%s ,param=%+v", rr.Code, http.StatusOK, tc.name, path, tc.param)
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

