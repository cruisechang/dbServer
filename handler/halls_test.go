package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHallsHandlerGet(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)
	tt := []struct {
		name string
	}{

		{"0"},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls")

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
		router.Handle("/halls", NewHallsHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler   failed http.status got %v want %v, name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}
		t.Logf("resData+%+v", resData)
	}
}

func TestHallsHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}

	//db
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	var ids []uint //放ids，刪掉用

	type param struct {
		HallID uint   `json:"hallID"`
		Name   string `json:"name"`
	}
	tt := []struct {
		name  string
		param param
		code  int
	}{

		{"0", param{90000, "測試90000"}, CodeSuccess},
		{"1", param{90000, "測試90001"}, CodeDBExecError}, //id duplicate
		{"2", param{90001, "測試90000"}, CodeDBExecError}, //name  duplicate
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls")
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
		router.Handle("/halls", NewHallsHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusOK {
			t.Errorf("TestHallsHandlerPost failed  http.Statuc got %v want %v, name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*hallData
		}{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("TestHallsHandlerPost unmarshal responseData error=%s", err.Error())
		}

		if resData.Code != tc.code {
			t.Errorf("TestHallsHandlerPost resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		//insert success
		if resData.Count == 1 {
			t.Logf("id %v ", resData.Data[0].HallID)
			ids = append(ids, resData.Data[0].HallID)
		}
	}

	if len(ids) > 0 {
		queryString := "DELETE FROM hall  where hall_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
