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

func Test_officialCMSManagersHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		code  int
		count int
	}{
		{"0", CodeSuccess, 6},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialcmsmanagers")

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
		router.Handle(path, NewOfficialCMSManagersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed  got %v want %v,name=%s, path=%s ", rr.Code, http.StatusOK, tc.name, path)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{
		}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s, path=%s", err.Error(), path)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s", resData.Code, tc.code, tc.name, path)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s, path=%s ", resData.Count, tc.count, tc.name, path)

		}
		t.Logf("resData+%+v", resData)
	}
}

func Test_officialCMSManagersHandler_post(t *testing.T) {

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
		param officialCMSManagerPostParam
	}{

		{"0", CodeSuccess, 1, officialCMSManagerPostParam{"forUnitTest01", "forUnitTest01", 1}},
		{"1", CodeRequestPostDataIllegal, 0, officialCMSManagerPostParam{"err", "err", 1}},
		{"2", CodeDBExecError, 0, officialCMSManagerPostParam{"forUnitTest01", "forUnitTest01", 1}}, //帳號重複
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSManagers")
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
		router.Handle(path, NewOfficialCMSManagersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed ,got %v want %v,name=%s, path=%s ,param=%+v", rr.Code, http.StatusOK, tc.name, path, tc.param)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*managerIDData
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
			t.Logf("ID=%d", resData.Data[0].ManagerID)
			ids = append(ids, resData.Data[0].ManagerID)
		}
	}

	if len(ids) > 0 {
		queryString := "DELETE FROM dealer  where dealer_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
