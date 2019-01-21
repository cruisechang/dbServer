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

func Test_partnerLogHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	type param struct {
		BeginDate string
		EndDate   string
	}
	tt := []struct {
		name      string
		count     int
		partnerID uint64
		param     param
	}{
		{"0", 0, 100, param{"2018-11-20 01:01:01", "2018-11-27 23:59:59"}},
		{"1", 6, 101, param{"2018-11-20 01:01:01", "2018-11-27 23:59:59"}},
		{"2", 3, 102, param{"2018-11-20 01:01:01", "2018-11-27 23:59:59"}},
		{"3", 0, 103, param{"2018-11-20 01:01:01", "2018-11-27 23:59:59"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners/" + strconv.FormatUint(tc.partnerID, 10) + "/log")

		b, _ := json.Marshal(tc.param)

		req, err := http.NewRequest("GET", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle("/partners/{id:[0-9]+}/log", NewPartnerLogHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

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
		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}
		t.Logf("resData=%+v", resData)
	}
}
