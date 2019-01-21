package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_dealerHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct{ ID uint64 }
	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK},
		{"1", "9999", CodeSuccess, 0, http.StatusOK},
		{"2", "33.3", CodeSuccess, 0, http.StatusNotFound},
		{"3", "xxx", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.ID)

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
		router.Handle("/dealers/{id:[0-9]+}", NewDealerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
		//t.Logf("%+v",resData)
	}
}

func Test_dealerHandler_delete(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	//insert first
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	queryString := "INSERT  INTO dealer (name,account,password,portrait_url ) values (? ,?,?,?)"

	stmt, _ := sqlDB.Prepare(queryString)
	defer stmt.Close()

	result, _ := stmt.Exec("deleteTest", "accounttest", "aaaa", "url")
	lastID, _ := result.LastInsertId()

	type param struct{ ID uint64 }
	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", strconv.FormatInt(lastID, 10), CodeSuccess, 1, http.StatusOK},
		{"1", "99999", CodeSuccess, 0, http.StatusOK},     //id not found
		{"2", "xxx", CodeSuccess, 0, http.StatusNotFound}, //mux parsing route error, return 404
		{"3", "3x", CodeSuccess, 0, http.StatusNotFound},  //mux parsing route error, return 404
		{"4", "3.3", CodeSuccess, 0, http.StatusNotFound}, //mux parsing route error, return 404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.ID)

		t.Logf("path %s", path)

		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/dealers/{id:[0-9]+}", NewDealerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("DELETE")

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
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}
		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

func Test_dealerHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct {
		Active int `json:"active"`
	}
	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK, dealerPatchParam{"forUint", "1234", "url", 0}},
		{"1", "1", CodeSuccess, 1, http.StatusOK, dealerPatchParam{"forUnit", "1234", "url", 0}},       //update 多項時，內容相同，也會update count=1
		{"2", "1", CodeSuccess, 1, http.StatusOK, dealerPatchParam{"test1", "1234", "url", 1}},         //改回去原始的
		{"3", "9999", CodeSuccess, 0, http.StatusOK, dealerPatchParam{"test1", "1234", "url", 1}},      //mux parsing route error, return 404
		{"4", "abc", CodeSuccess, 0, http.StatusNotFound, dealerPatchParam{"test1", "1234", "url", 1}}, //mux parsing route error, return 404
		{"5", "0.3", CodeSuccess, 0, http.StatusNotFound, dealerPatchParam{"test1", "1234", "url", 1}}, //mux parsing route error, return 404
		{"6", "10001", CodeRequestDataUnmarshalError, 0, http.StatusOK, struct{ Room int }{1}},         //參數錯誤
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/dealers/" + tc.ID)

		t.Logf("path %s", path)

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
		router.Handle("/dealers/{id:[0-9]+}", NewDealerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

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
			t.Fatalf("handler unmarshal responseData error=%s", err.Error())
		}
		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}
