package handler

import (
	"fmt"
	"github.com/cruisechang/dbServer/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func TestNotFoundHandler(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)

	uniqueIDProvider,_:=util.CreateUniqueIDProvider()

	tt := []struct {
		name string
	}{

		{"0"},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/")

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
		router.Handle("/", NewNotFoundHandler(NewBaseHandler(dbx.DB, dbx.Logger,uniqueIDProvider)))

		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusNotFound {
			t.Errorf("handler   failed http.status got %v want %v, name=%s", rr.Code, http.StatusNotFound, tc.name)
		}
	}
}
