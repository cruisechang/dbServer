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

func Test_userAccessTokenHandler_get(t *testing.T) {
	dbx, err := dbex.NewDBEX("dbexConfig.json")
	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	dbx.Logger.SetLevel(dbex.LevelInfo)
	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	tt := []struct {
		name  string
		code  int
		count int
		token string
	}{
		{"0", CodeSuccess, 1, "cb69b634-aeaf-41dc-a945-4a9cfc350fee"}, //success
		{"1", CodeSuccess, 0, "notFound"},                             //not found
		{"2", CodeSuccess, 0, "99999"},                                // not found
		{"3", CodePathError, 0, "9.99"},                               // path error = 太短
		{"4", CodePathError, 0, "33x"},                                // path error = 太短
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/" + tc.token + "/tokenData")

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
		router.Handle("/users/{accessToken}/tokenData", NewUserAccessTokenHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed token=%s, got %v want %v,name=%s, path=%s ", tc.token, rr.Code, http.StatusOK, tc.name, path)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &responseData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s, path=%s ", err.Error(), path)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s ", resData.Code, tc.code, tc.name, path)

		}

		if resData.Count != tc.count {
			t.Fatalf("handler resData count  got %d want %d, name=%s, path=%s ", resData.Count, tc.count, tc.name, path)

		}

		t.Logf("resData=%+v", resData)
	}
}
