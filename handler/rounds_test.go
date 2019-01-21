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

func Test_roundsHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	tt := []struct {
		name  string
		code  int
		count int
		param roundGetParam
	}{
		{"0", CodeSuccess, 11, roundGetParam{-1, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"1", CodeSuccess, 6, roundGetParam{100, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"2", CodeSuccess, 2, roundGetParam{200, 600, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"3", CodeSuccess, 1, roundGetParam{100, 1, 0, 0, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"4", CodeSuccess, 3, roundGetParam{-1, -1, 0, -1, "2018-12-18 14:57:53", "2018-12-20 23:59:59"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rounds")

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
		router.Handle(path, NewRoundsHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed  got %v want %v,name=%s, path=%s ,param=%+v", rr.Code, http.StatusOK, tc.name, path, tc.param)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
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

func Test_roundsHandler_post(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	//db
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))
	sqlDB := h.db.GetSQLDB()
	var ids []uint64 //放ids，刪掉用

	tt := []struct {
		name  string
		code  int
		count int
		param roundPostParam
	}{

		{"0", CodeSuccess, 1, roundPostParam{100, 1, 0, "banker win", "{}", 1}},
		{"1", CodeSuccess, 1, roundPostParam{100, 1, 0, "banker win", "{}", 0}},
		{"2", CodeSuccess, 1, roundPostParam{200, 700, 7, "3點", "{}", 2}},
		{"3", CodeSuccess, 1, roundPostParam{200, 700, 7, "3點", "{}", 1}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rounds")
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
		router.Handle(path, NewRoundsHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed  got %v want %v,name=%s, path=%s ,param=%+v", rr.Code, http.StatusOK, tc.name, path, tc.param)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*roundIDData
		}{}
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
		//insert success
		if resData.Count == 1 {
			t.Logf("ID=%d ", resData.Data[0].Round)
			ids = append(ids, resData.Data[0].Round)
		}
	}

	if len(ids) > 0 {
		queryString := "DELETE FROM round  where round_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
