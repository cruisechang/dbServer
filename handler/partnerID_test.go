package handler

import (
	"fmt"
	"strconv"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"testing"
	"github.com/cruisechang/dbex"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func Test_partnerHandler_get(t *testing.T) {

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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}", NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_getAESKey(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/aesKey"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" +tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_getAvtive(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/active"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_getLogin(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/login"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_getAPIBindIP(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/apiBindIP"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
func Test_partnerHandler_getCMSBindIP(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/cmsBindIP"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_getAccessToken(t *testing.T) {
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
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"1", "101", CodeSuccess, 1, http.StatusOK},
		{"2", "102", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusMovedPermanently},
	}

	targetPath := "/accessToken"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID + targetPath)

		t.Logf("path %s", path)

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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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

func Test_partnerHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	ip := []string{"xxx.xxx.xxx", "xxx.xxx.xxx"}
	ipb, _ := json.Marshal(ip)

	tt := []struct {
		name  string
		ID    string
		code  int
		count int
		param interface{}
	}{
		{"0", "100", CodeSuccess, 1, partnerPatchParam{"test", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}},   //修改
		{"1", "100", CodeSuccess, 0, partnerPatchParam{"test", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}},   //相同
		{"2", "100", CodeSuccess, 1, partnerPatchParam{"pass", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}},   //改回去
		{"3", "99999", CodeSuccess, 0, partnerPatchParam{"pass", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}}, //id not found
		{"4", "xxx", CodeSuccess, 0, partnerPatchParam{"pass", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}},   //mux parsing route error
		{"5", "3.3", CodeSuccess, 0, partnerPatchParam{"pass", "account100", 0, 0, "12345678", "12345678", string(ipb), string(ipb), 0}},   //mux parsing route error
		{"6", "100", CodeRequestDataUnmarshalError, 0, struct{ Login int }{1}},                                                                               //參數錯誤
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + tc.ID)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}", NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

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

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}

func Test_partnerHandler_patchLogin(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param struct{ Login int }
	}{
		{"0", 100, CodeSuccess, 1, struct{ Login int }{1}},     //修改
		{"1", 100, CodeSuccess, 0, struct{ Login int }{1}},     //改相同值,count=0
		{"2", 100, CodeSuccess, 1, struct{ Login int }{0}},     //改回去
		{"3", 9999999, CodeSuccess, 0, struct{ Login int }{0}}, //not found
	}

	targetPath := "/login"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

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

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}

func Test_partnerHandler_patchActive(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param struct{ Active int }
	}{
		{"0", 100, CodeSuccess, 1, struct{ Active int }{9}},
		{"1", 100, CodeSuccess, 0, struct{ Active int }{9}},
		{"2", 100, CodeSuccess, 1, struct{ Active int }{1}},
		{"3", 9999999, CodeSuccess, 0, struct{ Active int }{0}},
	}

	targetPath := "/active"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

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

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}

func Test_partnerHandler_patchAESKey(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param struct{ AESKey string }
	}{
		{"0", 100, CodeSuccess, 1, struct{ AESKey string }{"testaes1"}},
		{"1", 100, CodeSuccess, 0, struct{ AESKey string }{"testaes1"}},
		{"2", 100, CodeSuccess, 1, struct{ AESKey string }{"aeskey"}},
		{"3", 9999999, CodeSuccess, 0, struct{ AESKey string }{"testaes"}},
	}

	targetPath := "/aesKey"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

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

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
	}
}

func Test_partnerHandler_patchAPIBindIP(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	ip := []string{"xxx.xxx.xxx"}
	ip2 := []string{"xxx.xxx.xxx", "xxx.xxx.xxx"}

	ipb, _ := json.Marshal(ip)
	ipb2, _ := json.Marshal(ip2)

	apis := apiBindIPData{
		APIBindIP: string(ipb),
	}
	apis2 := apiBindIPData{
		APIBindIP: string(ipb2),
	}

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param apiBindIPData
	}{
		{"0", 100, CodeSuccess, 1, apis},
		{"1", 100, CodeSuccess, 0, apis},
		{"2", 100, CodeSuccess, 1, apis2},
		{"3", 9999999, CodeSuccess, 0, apis},
	}

	targetPath := "/apiBindIP"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed ID=%d, got %v want %v,name=%s", tc.ID, rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

func Test_partnerHandler_patchCMSBindIP(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	ip := []string{"xxx.xxx.xxx"}
	ip2 := []string{"xxx.xxx.xxx", "xxx.xxx.xxx"}

	ipb, _ := json.Marshal(ip)
	ipb2, _ := json.Marshal(ip2)

	apis := cmsBindIPData{
		CMSBindIP: string(ipb),
	}
	apis2 := cmsBindIPData{
		CMSBindIP: string(ipb2),
	}

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param cmsBindIPData
	}{
		{"0", 100, CodeSuccess, 1, apis},
		{"1", 100, CodeSuccess, 0, apis},
		{"2", 100, CodeSuccess, 1, apis2},
		{"3", 9999999, CodeSuccess, 0, apis},
	}

	targetPath := "/cmsBindIP"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed ID=%d, got %v want %v,name=%s", tc.ID, rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

//update 與db相同資料，會回傳affecterd row =0
func Test_partnerHandler_patchAccessToken(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	apis := accessTokenData{
		AccessToken: "00000000",
	}
	apis2 := accessTokenData{
		AccessToken: "12345678",
	}

	tt := []struct {
		name  string
		ID    uint64
		code  int
		count int
		param accessTokenData
	}{
		{"0", 100, CodeSuccess, 1, apis},
		{"1", 100, CodeSuccess, 0, apis},
		{"2", 100, CodeSuccess, 1, apis2},
		{"3", 9999999, CodeSuccess, 0, apis},
	}

	targetPath := "/accessToken"

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.ID, 10) + targetPath)

		t.Logf("path %s", path)
		t.Logf("tc %+v", tc)
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
		router.Handle("/partners/{id:[0-9]+}"+targetPath, NewPartnerIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed ID=%d, got %v want %v,name=%s", tc.ID, rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)

		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}
