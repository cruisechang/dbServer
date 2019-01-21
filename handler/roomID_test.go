package handler

import (
	"github.com/cruisechang/dbServer/util"
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

func Test_roomHandler_get(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct{ ID uint64 }
	tt := []struct {
		name       string
		ID         string
		code       int
		count      int
		httpStatus int
	}{
		{"0", "1", CodeSuccess, 1, http.StatusOK},
		{"1", "2", CodeSuccess, 1, http.StatusOK},
		{"2", "700", CodeSuccess, 1, http.StatusOK},
		{"3", "99999", CodeSuccess, 0, http.StatusOK},
		{"4", "xxx", CodeSuccess, 0, http.StatusNotFound},
		{"5", "3.3", CodeSuccess, 0, http.StatusNotFound},
		{"6", "", CodeSuccess, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + tc.ID)

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
		router.Handle("/rooms/{id:[0-9]+}", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

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

func Test_roomHandler_delete(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	//dbx.Logger.SetLevel(dblog.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	//insert first
	h := NewBroadcastIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))
	sqlDB := h.db.GetSQLDB()
	queryString := "INSERT  INTO room (room_id,hall_id,name,room_type,hls_url,bet_countdown,dealer_id,limitation_id) values (? ,?,?,?,?,?,?,?)"

	stmt, _ := sqlDB.Prepare(queryString)
	defer stmt.Close()

	insertID := 9999
	stmt.Exec(insertID, "99999", "test", 1, "url", 20, 999, 9999)

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
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.param.ID, 10))

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
		router.Handle("/rooms/{id:[0-9]+}", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("DELETE")

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
			t.Errorf("handler resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Errorf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)

		}
	}
}

func Test_roomHandler_patch(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	//type param struct {
	//	DealerID uint `json:"dealerID"`
	//}
	tt := []struct {
		name  string
		ID    uint64
		param interface{}
		code  int
		count int
	}{
		{"0", 1, roomPatchParam{1, 100, "test", 1, 1, "url", 20, 0}, CodeSuccess, 1}, //改
		{"1", 1, roomPatchParam{1, 100, "test", 1, 1, "url", 20, 0}, CodeSuccess, 0}, //相同內容
		{"2", 1, roomPatchParam{1, 100, "百家01", 1, 1, "url", 20, 0}, CodeSuccess, 1}, //改回去
		{"3", 999, struct{ X int }{3}, CodeRequestDataUnmarshalError, 0},             //error
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10))

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
		router.Handle("/rooms/{id:[0-9]+}", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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

func Test_roomHandler_patchName(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

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
		{"0", 1, param{"test99"}, CodeSuccess, 1},
		{"1", 1, param{"test99"}, CodeSuccess, 0}, //update 內容相同時，count=0
		{"2", 1, param{"百家01"}, CodeSuccess, 1},   //改回去
		{"3", 99999, param{"test99"}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/name")

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
		router.Handle("/rooms/{id:[0-9]+}/name", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s, path=%s", rr.Code, http.StatusOK, tc.name, path)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s,path=%s", err.Error(), path)
		}
		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s", resData.Code, tc.code, tc.name, path)
		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s, path=%s", resData.Count, tc.count, tc.name, path)

		}
	}
}

func Test_roomHandler_patchActive(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

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
		{"0", 1, param{0}, CodeSuccess, 1},
		{"1", 1, param{0}, CodeSuccess, 0},
		{"2", 1, param{1}, CodeSuccess, 1}, //改回去
		{"3", 99999, param{1}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/active")

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
		router.Handle("/rooms/{id:[0-9]+}/active", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

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

func Test_roomHandler_patchHLRURL(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		HLSURL string `json:"hlsURL"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{"test"}, CodeSuccess, 1},
		{"1", 1, param{"test"}, CodeSuccess, 0},
		{"2", 1, param{"url"}, CodeSuccess, 1},     //改回去
		{"3", 99999, param{"url"}, CodeSuccess, 0}, //update 內容相同時，count=0
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/hlsURL")

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
		router.Handle("/rooms/{id:[0-9]+}/hlsURL", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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
func Test_roomHandler_patchBoot(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		Boot uint `json:"boot"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{100}, CodeSuccess, 1},
		{"1", 1, param{100}, CodeSuccess, 0}, //update 內容相同時，count=0
		{"2", 1, param{10}, CodeSuccess, 1},  //改回去
		{"3", 99999, param{100}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/boot")

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
		router.Handle("/rooms/{id:[0-9]+}/boot", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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

func Test_roomHandler_patchRound(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		Round uint64 `json:"round"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{20180987}, CodeSuccess, 1},
		{"1", 1, param{20180987}, CodeSuccess, 0},
		{"2", 1, param{181211100100028}, CodeSuccess, 1}, //改回去
		{"3", 99999, param{20180987}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/round")

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
		router.Handle("/rooms/{id:[0-9]+}/round", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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

func Test_roomHandler_patchStatus(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		Status uint `json:"status"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{1}, CodeSuccess, 1},
		{"1", 1, param{1}, CodeSuccess, 0},
		{"2", 1, param{0}, CodeSuccess, 1},      //改回去
		{"3", 99999, param{11}, CodeSuccess, 0}, //update 內容相同時，count=0
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/status")

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
		router.Handle("/rooms/{id:[0-9]+}/status", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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

func Test_roomHandler_patchBetCountdown(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		BetCountdown uint `json:"betCountdown"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{11}, CodeSuccess, 1},
		{"1", 1, param{11}, CodeSuccess, 0},
		{"2", 1, param{20}, CodeSuccess, 1}, //改回去
		{"3", 99999, param{20}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/betCountdown")

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
		router.Handle("/rooms/{id:[0-9]+}/betCountdown", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
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
	}
}

func Test_roomHandler_patchBetDealerID(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		DealerID uint `json:"dealerID"`
	}
	tt := []struct {
		name  string
		ID    uint64
		param param
		code  int
		count int
	}{
		{"0", 1, param{2}, CodeSuccess, 1},
		{"1", 1, param{2}, CodeSuccess, 0},
		{"2", 1, param{1}, CodeSuccess, 1}, //改回去
		{"3", 99999, param{2}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + strconv.FormatUint(tc.ID, 10) + "/dealerID")

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
		router.Handle("/rooms/{id:[0-9]+}/dealerID", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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

func Test_roomHandler_patchNewRound(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	tt := []struct {
		name  string
		ID    string
		param roomNewRoundPatchParam
		code  int
		count int
	}{
		{"0", "1", roomNewRoundPatchParam{10, 999, 0}, CodeSuccess, 1},
		{"1", "1", roomNewRoundPatchParam{10, 999, 0}, CodeSuccess, 0},
		{"2", "1", roomNewRoundPatchParam{10, 181211100100028, 0}, CodeSuccess, 1}, //改回去
		{"3", "99999", roomNewRoundPatchParam{2, 1, 1}, CodeSuccess, 0},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/rooms/" + tc.ID + "/newRound")

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
		router.Handle("/rooms/{id:[0-9]+}/newRound", NewRoomIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("PATCH")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed http.Status   got %v want %v,name=%s", rr.Code, http.StatusOK, tc.name)
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
