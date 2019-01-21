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

func Test_userIDHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK},     //success
		{"1", "99999", CodeSuccess, 0, http.StatusOK}, //not found
		{"2", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"3", "1.1", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID)

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
		router.Handle("/users/{id:[0-9]+}", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
	}
}

func Test_userIDHandler_getCredit(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		code       int
		count      int
		httpStatus int
		ID         string
	}{
		{"0", CodeSuccess, 1, http.StatusOK, "1"},           //success
		{"1", CodeSuccess, 0, http.StatusOK, "99999"},       //not found
		{"2", CodeSuccess, 0, http.StatusNotFound, "xxxxx"}, //mux parsing error httpStatusNotFound
		{"3", CodeSuccess, 0, http.StatusNotFound, "3.3"},   //mux parsing error httpStatusNotFound
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/credit")

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
		router.Handle("/users/{id:[0-9]+}/credit", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
	}
}

func Test_userIDHandler_getLogin(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK},
		{"1", "9999", CodeSuccess, 0, http.StatusOK},
		{"2", "xxxx", CodeSuccess, 1, http.StatusNotFound},
		{"3", "3.3", CodeSuccess, 1, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/login")

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
		router.Handle("/users/{id:[0-9]+}/login", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
	}
}

func Test_userIDHandler_getActive(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK},
		{"1", "99999", CodeSuccess, 0, http.StatusOK},
		{"2", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"3", "3.3", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/active")

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
		router.Handle("/users/{id:[0-9]+}/active", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed id=%s, http status got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
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
	}
}

func Test_userIDHandler_patchCredit(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK, struct{ Credit float32 }{100}},         //success
		{"0", "1", CodeSuccess, 0, http.StatusOK, struct{ Credit float32 }{100}},         //重複，不update
		{"0", "1", CodeSuccess, 1, http.StatusOK, struct{ Credit float32 }{10}},          //還原
		{"0", "99999", CodeSuccess, 0, http.StatusOK, struct{ Credit float32 }{100}},     //not found
		{"0", "xxx", CodeSuccess, 0, http.StatusNotFound, struct{ Credit float32 }{100}}, //http status not found
		{"0", "3.3", CodeSuccess, 0, http.StatusNotFound, struct{ Credit float32 }{100}}, //http status not found
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/credit")

		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("Test_userHandler_patchCredit unmarshal param error, param=%+v", tc)
		}

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
		router.Handle("/users/{id:[0-9]+}/credit", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed id=%s, http status got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
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
	}
}

func Test_userIDHandler_patchLogin(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK, struct{ Login int }{1}},     //success
		{"1", "1", CodeSuccess, 0, http.StatusOK, struct{ Login int }{1}},     //the same
		{"2", "1", CodeSuccess, 1, http.StatusOK, struct{ Login int }{0}},     //success
		{"3", "99999", CodeSuccess, 0, http.StatusOK, struct{ Login int }{0}}, //not found
		{"4", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, struct {Login string `json:"login"`}{"x"}},
		{"5", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, ""},
		{"6", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, 9},
		{"7", "xxxx", CodeSuccess, 0, http.StatusNotFound, struct{ Login int }{0}}, //404
		{"8", "33,3", CodeSuccess, 0, http.StatusNotFound, struct{ Login int }{0}}, //404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/login")

		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("Test_userHandler_patchCredit unmarshal param error, param=%+v", tc)
		}

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
		router.Handle("/users/{id:[0-9]+}/login", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed id=%s, http status got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
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
	}
}

func Test_userIDHandler_patchActive(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK, struct{ Active int `json:"active"` }{0}}, //success
		{"1", "1", CodeSuccess, 0, http.StatusOK, struct{ Active int `json:"active"` }{0}}, //the same
		{"2", "1", CodeSuccess, 1, http.StatusOK, struct{ Active int `json:"active"` }{1}}, //success
		{"3", "99999", CodeSuccess, 0, http.StatusOK, struct{ Active int }{0}},             //not found
		{"4", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, struct{ Active string }{"x"}},
		{"5", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, ""},
		{"6", "1", CodeRequestDataUnmarshalError, 0, http.StatusOK, 9},
		{"7", "xxxxx", CodeSuccess, 0, http.StatusNotFound, struct{ Active int }{0}}, //404
		{"8", "33,3", CodeSuccess, 0, http.StatusNotFound, struct{ Active int }{0}},  //404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.ID + "/active")

		b, err := json.Marshal(tc.param)
		if err != nil {
			t.Fatalf("handerl unmarshal param error, param=%+v", tc)
		}

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
		router.Handle("/users/{id:[0-9]+}/active", NewUserIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed id=%s, http status got %d want %d, name=%s", tc.ID, rr.Code, tc.httpStatus, tc.name)
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
	}
}
