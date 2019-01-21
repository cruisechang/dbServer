package handler

import (
	"github.com/cruisechang/dbServer/util"
	"net/http"
	"testing"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func Test_transferPartnerTransferIDHandler_get(t *testing.T) {

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
		//id 是字串
		{"0", "test100", CodeSuccess, 1, http.StatusOK}, //success
		{"1", "99999", CodeSuccess, 0, http.StatusOK},   //not found
		{"2", "0", CodePathError, 0, http.StatusOK},     //not found
		{"3", "-1", CodePathError, 0, http.StatusOK},    //not found
		{"4", "xxx", CodePathError, 0, http.StatusOK},
		{"5", "1.1", CodePathError, 0, http.StatusOK},
		{"6", "", CodePathError, 0, http.StatusNotFound},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/transfers/ptID/" + tc.ID)

		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.Handle("/transfers/ptID/{partnerTransferID}", NewTransferPartnerTransferIDHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

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
		t.Logf("resData=%+v", resData)
	}
}
