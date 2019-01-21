package handler

import (
	"net/http"
	"testing"

	"github.com/cruisechang/dbex"
	"fmt"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func Test_roundHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct{ ID string }

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{

		{"0", "1", CodeSuccess, 1, http.StatusOK},     //success
		{"1", "2", CodeSuccess, 1, http.StatusOK},     //success
		{"2", "99999", CodeSuccess, 0, http.StatusOK}, //not found
		{"3", "0", CodePathError, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "1.1", CodeSuccess, 0, http.StatusNotFound},
		{"6", "-1", CodeSuccess, 0, http.StatusNotFound},
		{"7", "", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rounds/" + tc.ID)

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
		router.Handle("/rounds/{id:[0-9]+}", NewRoundIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
		t.Logf("resData=%+v", resData)
	}
}

func Test_roundHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type x struct {
		A string `json:"a"`
		B string `json"b"`
	}

	xx := x{"astr", "bstr"}

	jxx, _ := json.Marshal(xx)

	tt := []struct {
		name       string
		httpStatus int
		ID         string
		code       int
		count      int
		param      interface{}
	}{
		{"0", http.StatusOK, "1", CodeSuccess, 1, roundPatchParam{"for test", string(jxx), 1}},       //update
		{"1", http.StatusOK, "1", CodeSuccess, 0, roundPatchParam{"for test", string(jxx), 1}},       //update the same
		{"2", http.StatusOK, "1", CodeSuccess, 1, roundPatchParam{"banker win", string(jxx), 1}},     //還原
		{"3", http.StatusOK, "99999", CodeSuccess, 0, roundPatchParam{"banker win", string(jxx), 1}}, //id not found
		{"4", http.StatusOK, "1", CodeRequestDataUnmarshalError, 0, x{"brief", "result"}},            //param error
		{"5", http.StatusOK, "2", CodeRequestDataUnmarshalError, 0, x{"brief", "result"}},            //param error
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rounds/" + tc.ID + "/patch")

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
		router.Handle("/rounds/{id:[0-9]+}/patch", NewRoundIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
			return
		}

		if rr.Code != http.StatusOK {
			return
		}

		body, _ := ioutil.ReadAll(rr.Body)

		if string(body) != "" {

			resData := &responseData{}
			err = json.Unmarshal(body, resData)
			if err != nil {
				t.Fatalf("handler unmarshal responseData error=%s, path=%s", err.Error(), path)
			}
			if resData.Code != tc.code {
				t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s", resData.Code, tc.code, tc.name, path)
			}

			if resData.Count != tc.count {
				t.Fatalf("handler resData count  got %d want %d, name=%s, path=%s", resData.Count, tc.count, tc.name, path)

			}
		}
	}
}
