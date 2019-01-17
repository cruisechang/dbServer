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

func Test_officialCMSManagerIDHandler_get(t *testing.T) {

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
		{"0", "1", CodeSuccess, 1, http.StatusOK},            //success
		{"1", "2", CodeSuccess, 1, http.StatusOK},            //success
		{"2", "99999", CodeSuccess, 0, http.StatusOK},        //not found
		{"3", "xxxx", CodePathError, 0, http.StatusNotFound}, //非int 會環傳404
		{"4", "3.3", CodePathError, 0, http.StatusNotFound},  //非int 會環傳404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSManagers/" + tc.ID)

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
		router.Handle("/officialCMSManagers/{id:[0-9]+}", NewOfficialCMSManagerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
		t.Logf("resData=%+v", resData)
	}
}

func Test_officialCMSManagerIDHandler_patch(t *testing.T) {
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
		{"0", http.StatusOK, "1", CodeSuccess, 1, officialCMSManagerPatchParam{"forUnitTest01", 1, 0}},
		{"1", http.StatusOK, "1", CodeSuccess, 1, officialCMSManagerPatchParam{"test01", 1, 1}},
		{"2", http.StatusOK, "9999", CodeRequestDataUnmarshalError, 0, errParam{"brief", "result"}},
		{"3", http.StatusOK, "2", CodeRequestDataUnmarshalError, 0, errParam{"brief", "result"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSManagers/" + tc.ID)

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
		router.Handle("/officialCMSManagers/{id:[0-9]+}", NewOfficialCMSManagerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
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

func Test_officialCMSManagerIDHandler_delete(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}

	h := NewOfficialCMSManagerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))

	sqlDB := h.db.GetSQLDB()

	queryString := "INSERT  INTO official_cms_manager (account,password,role_id ) values (? ,?,?)"

	stmt, _ := sqlDB.Prepare(queryString)
	defer stmt.Close()

	result, _ := stmt.Exec("uintDeleteTest", "pass", 1)

	//affRow, err := result.RowsAffected()

	lastID, _ := result.LastInsertId()

	type param struct{ ID uint64 }
	tt := []struct {
		name  string
		param param
		code  int
		count int
	}{
		{"0", param{uint64(lastID)}, CodeSuccess, 1},
		{"1", param{9897667556}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/officialCMSManagers/" + strconv.FormatUint(tc.param.ID, 10))

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
		router.Handle("/officialCMSManagers/{id:[0-9]+}", h).Methods("DELETE")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed  got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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
