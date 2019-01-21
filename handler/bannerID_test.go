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

func Test_bannerIDHandler_patch(t *testing.T) {
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
		{"0", http.StatusOK, "1", CodeSuccess, 1, bannerPostParam{"forTest", "", "", 0, 1}},
		{"1", http.StatusOK, "1", CodeSuccess, 0, bannerPostParam{"forTest", "", "", 0, 1}},
		{"2", http.StatusOK, "1", CodeSuccess, 1, bannerPostParam{"http://139.168.672.65/banner/", "", "", 0, 1}},
		{"3", http.StatusOK, "9999", CodeSuccess, 0, bannerPostParam{"http://139.168.62.65/banner/test", "", "", 0, 1}},
		{"4", http.StatusNotFound, "xxxx", CodeRequestDataUnmarshalError, 0, errParam{"brief", "result"}},
		{"5", http.StatusNotFound, "3.3", CodeRequestDataUnmarshalError, 0, errParam{"brief", "result"}},
		{"6", http.StatusOK, "1", CodeRequestDataUnmarshalError, 0, ""},
		{"7", http.StatusOK, "1", CodeRequestDataUnmarshalError, 0, 33},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/banners/" + tc.ID)

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
		router.Handle("/banners/{id:[0-9]+}", NewBannerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
		}

		if rr.Code != http.StatusOK {
			return
		}

		body, _ := ioutil.ReadAll(rr.Body)

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
		t.Logf("resData=%+v", resData)
	}
}

func Test_bannerIDHandler_delete(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}

	//insert first
	h := NewBroadcastIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	queryString := "INSERT  INTO banner (pic_url,link_url,description,platform,active ) values (? ,?,?,?,?)"

	stmt, _ := sqlDB.Prepare(queryString)
	defer stmt.Close()

	result, _ := stmt.Exec("deleteTest", "", "", 0, 1)
	lastID, _ := result.LastInsertId()

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", strconv.FormatInt(lastID, 10), CodeSuccess, 1, http.StatusOK},
		{"1", "99999", CodeSuccess, 0, http.StatusOK},
		{"2", "xxxx", CodeSuccess, 0, http.StatusNotFound}, //id not int
		{"3", "3.3", CodeSuccess, 0, http.StatusNotFound},  //id not int
		{"4", "", CodeSuccess, 0, http.StatusNotFound},     //id not int
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/banners/" + tc.ID)

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
		router.Handle("/banners/{id:[0-9]+}", NewBannerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("DELETE")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
		}

		if rr.Code != http.StatusOK {
			return
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
