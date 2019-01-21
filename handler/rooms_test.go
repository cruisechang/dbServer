package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func TestRoomsHandlerGet(t *testing.T) {

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
		path := fmt.Sprintf("/rooms")

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
		router.Handle("/rooms", NewRoomsHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusOK {
			t.Errorf("TestHallsHandlerGet   failed http.status got %v want %v, name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}
		t.Logf("resData=%+v", resData)
	}
}

func TestRoomsHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}

	//db
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	var ids []uint //放ids，刪掉用

	tt := []struct {
		name  string
		code  int
		count int
		param interface{}
	}{

		{"0", CodeSuccess, 1, roomPostParam{900, 100, "測試90000", 0, 20, "url", 1}},
		{"1", CodeSuccess, 1, roomPostParam{901, 100, "測試90001", 0, 20, "url", 1}},
		{"2", CodeDBExecError, 0, roomPostParam{900, 100, "測試90002", 0, 20, "url", 1}}, //duplicate id
		{"3", CodeDBExecError, 0, roomPostParam{902, 100, "測試90000", 0, 20, "url", 1}}, //duplicate name
		{"4", CodeRequestPostDataIllegal, 0, struct{ RoomID int }{904}},                //wrong param
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms")
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
		router.Handle("/rooms", NewRoomsHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusOK {
			t.Errorf("TestRoomsHandlerPost failed  http.Statuc got %v want %v, name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*roomIDData
		}{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}

		//insert success
		if resData.Count == 1 {
			t.Logf("ID=%d ", resData.Data[0].RoomID)
			ids = append(ids, resData.Data[0].RoomID)
		}

	}

	if len(ids) > 0 {
		queryString := "DELETE FROM room  where room_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
