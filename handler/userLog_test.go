package handler

import (
	"net/http"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strconv"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_userLogHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct {
		BeginDate string
		EndDate   string
	}
	tt := []struct {
		name   string
		userID uint64
		count  int
		param  param
	}{
		{"0", 100, 2, param{"2018-11-20 01:01:01", "2018-11-21 23:59:59"}},
		{"1", 101, 1, param{"2018-11-20 01:01:01", "2018-11-21 23:59:59"}},
		{"2", 102, 1, param{"2018-11-20 01:01:01", "2018-11-21 23:59:59"}},
		{"3", 103, 6, param{"2018-11-20 01:01:01", "2018-11-21 23:59:59"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + strconv.FormatUint(tc.userID, 10) + "/log")

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
		router.Handle("/users/{id:[0-9]+}/log", NewUserLogHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}
