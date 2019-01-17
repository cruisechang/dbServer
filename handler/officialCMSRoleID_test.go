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

func Test_officialCMSRoleIDHandler_get(t *testing.T) {

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
		{"0", "1", CodeSuccess, 1, http.StatusOK},
		{"1", "2", CodeSuccess, 1, http.StatusOK},
		{"2", "99999", CodeSuccess, 0, http.StatusOK},
		{"3", "xxx", CodePathError, 0, http.StatusNotFound}, //非int mux會回傳404
		{"4", "3.3", CodePathError, 0, http.StatusNotFound}, //非int mux 會回傳404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSRoles/" + tc.ID)

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
		router.Handle("/officialCMSRoles/{id:[0-9]+}", NewOfficialCMSRoleIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
		t.Logf("resData=%+v",resData)
	}
}

func Test_officialCMSRoleIDHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type errParam struct {
		A string `json:"a"`
		B string `json"b"`
	}

	tt := []struct {
		name       string
		httpStatus int
		ID         string
		code       int
		count      int
		param      interface{}
	}{
		{"0", http.StatusOK, "1", CodeSuccess, 1, officialCMSRolePatchParam{"[999,999]"}},
		{"1", http.StatusOK, "1", CodeSuccess, 0, officialCMSRolePatchParam{"[999,999]"}},               //相同內容, count=1
		{"2", http.StatusOK, "1", CodeSuccess, 1, officialCMSRolePatchParam{"[1,2,3,4]"}},               //改回去
		{"2", http.StatusOK, "9999", CodeSuccess, 0, officialCMSRolePatchParam{"[1,2,3,4]"}},            //id not found
		{"3", http.StatusOK, "1", CodeRequestDataUnmarshalError, 0, struct{ X string }{"xxxx"}},         //param error
		{"4", http.StatusOK, "1", CodeRequestDataUnmarshalError, 0, struct{ X int }{3}},                 //param error
		{"5", http.StatusNotFound, "X", CodeRequestDataUnmarshalError, 0, struct{ X string }{"xxxx"}},   //mux route error return 404
		{"6", http.StatusNotFound, "3.3", CodeRequestDataUnmarshalError, 0, struct{ X string }{"xxxx"}}, //mux route error return 404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSRoles/" + tc.ID)

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
		router.Handle("/officialCMSRoles/{id:[0-9]+}", NewOfficialCMSRoleIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
			return
		}
		if rr.Code == http.StatusNotFound {
			//mux parsing route error
			return
		}

		if rr.Code != http.StatusOK {
			//no body
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
