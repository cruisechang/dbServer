package handler

import (
	"testing"
	"fmt"
	"github.com/cruisechang/dbex"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	//"encoding/json"
	//"bytes"
	//"io/ioutil"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

func TestDealersHandlerGet(t *testing.T) {

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
		path := fmt.Sprintf("/dealers")

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.Handle("/dealers", NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler   failed http.status got %v want %v, name=%s", rr.Code, http.StatusOK,tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}
		t.Logf("resData+%+v",resData)
	}
}

func TestDealersHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)

	tt := []struct {
		name string
		code int
		count int
		param dealerPostParam
	}{

		{"0", CodeSuccess,1,dealerPostParam{"gg1","gg1","1234","url"}},
		{"1", CodeDBExecError,0,dealerPostParam{"gg1","gg1","1234","url"}},  //account duplicate
	}


	for _, tc := range tt {
		path := fmt.Sprintf("/dealers")
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
		router.Handle("/dealers", NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("TestRoomsHandlerPost failed  http.Statuc got %v want %v, name=%s", rr.Code, http.StatusOK,tc.name)
		}

		body,_:=ioutil.ReadAll(rr.Body)

		resData := &responseData{
		}
		err = json.Unmarshal(body,resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s",err.Error())
		}

		if resData.Code!=tc.code{
			t.Errorf("handler resData code  got %d want %d, name=%s",resData.Code,tc.code,tc.name)

		}

		if resData.Count!=tc.count{
			t.Errorf("handler resData count  got %d want %d, name=%s",resData.Count,tc.count,tc.name)

		}
	}
}





