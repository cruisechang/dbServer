package handler

import (
	"net/http"
	"testing"

	"github.com/cruisechang/dbServer/util"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strconv"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_hallHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	type param struct{ ID uint64 }
	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "100", CodeSuccess, 1, http.StatusOK},
		{"2", "200", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls/" + tc.ID)

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
		router.Handle("/halls/{id:[0-9]+}", NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("GET")

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
		t.Logf("%+v", resData)
	}
}

func Test_hallHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	tt := []struct {
		name       string
		ID         string
		param      interface{}
		code       int
		count      int
		httpStatus int
	}{
		{"0", "100", hallPatchParam{100, "unittest", 1}, CodeSuccess, 1, http.StatusOK},
		{"1", "100", hallPatchParam{100, "unittest", 1}, CodeSuccess, 0, http.StatusOK},
		{"2", "100", hallPatchParam{100, "皇家廳", 1}, CodeSuccess, 1, http.StatusOK},
		{"3", "9999", hallPatchParam{100, "test", 1}, CodeSuccess, 0, http.StatusOK},              //id not found
		{"4", "100", struct{ X string }{"test"}, CodeRequestDataUnmarshalError, 0, http.StatusOK}, //param error
		{"5", "100", 999, CodeRequestDataUnmarshalError, 0, http.StatusOK},                        //param error
		{"6", "xxx", hallPatchParam{100, "test", 1}, CodeSuccess, 0, http.StatusNotFound},         //mux parsing route error, return 404
		{"7", "33.3", hallPatchParam{100, "test", 1}, CodeSuccess, 0, http.StatusNotFound},        //mux parsing route error, return 404
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls/" + tc.ID)

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
		router.Handle("/halls/{id:[0-9]+}", NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("PATCH")

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
		t.Logf("%+v", resData)
	}
}

/*
func Test_hallHandler_patchName(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct {
		Name string `json:"name"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 100, param{"test999"}, CodeSuccess, 1},
		{"1", 100, param{"test999"}, CodeSuccess, 0},   //update 內容相同時，count=0
		{"2", 100, param{"皇家廳"}, CodeSuccess, 1},       //改回去
		{"3", 99999, param{"test999"}, CodeSuccess, 0}, //id not found
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls/" + strconv.FormatUint(tc.ID, 10) + "/name")

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
		router.Handle("/halls/{id:[0-9]+}/name", NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
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
	}
}

func Test_hallHandler_patchActive(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	type param struct {
		Active uint `json:"active"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 200, param{0}, CodeSuccess, 1},
		{"1", 200, param{0}, CodeSuccess, 0}, //update 內容相同時，count=0
		{"2", 200, param{1}, CodeSuccess, 1},
		{"3", 99999, param{0}, CodeSuccess, 0}, //id not found
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls/" + strconv.FormatUint(tc.ID, 10) + "/active")

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
		router.Handle("/halls/{id:[0-9]+}/active", NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
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
	}
}
*/

func Test_hallHandler_delete(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)

	uniqueIDProvider, _ := util.CreateUniqueIDProvider()

	//insert first
	h := NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))
	sqlDB := h.db.GetSQLDB()
	insertID := 8888
	queryString := "INSERT  INTO hall (hall_id,name) values (? ,?)"

	stmt, _ := sqlDB.Prepare(queryString)
	defer stmt.Close()

	stmt.Exec(insertID, "test")

	type param struct{ ID uint64 }
	tt := []struct {
		name  string
		param param
		code  int
		count int
	}{
		{"0", param{uint64(insertID)}, CodeSuccess, 1},
		{"1", param{9897667556}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/halls/" + strconv.FormatUint(tc.param.ID, 10))

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
		router.Handle("/halls/{id:[0-9]+}", NewHallIDHandler(NewBaseHandler(dbx.DB, dbx.Logger, uniqueIDProvider))).Methods("DELETE")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed  got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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
