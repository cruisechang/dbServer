package test

import (
	"testing"
	"fmt"
	"net/http"
	"strings"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"encoding/json"
	"bytes"
)

func TestHandlerHeaderTest(t *testing.T) {
	path := fmt.Sprintf("/users")

	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("API-Key", "abc123")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", HandlerHeaderTest).Methods("POST")

	router.ServeHTTP(rr, req)
}

func TestHandlerContentTypeURLEncodedParseForm(t *testing.T) {
	path := fmt.Sprintf("/users")

	req, err := http.NewRequest("POST", path, strings.NewReader("location=tp&age=3"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", HandlerContentTypeURLEncodedParseForm).Methods("POST")

	router.ServeHTTP(rr, req)
}

func TestHandlerContentTypeURLEncoded(t *testing.T) {
	path := fmt.Sprintf("/users")

	req, err := http.NewRequest("POST", path, strings.NewReader("location=tp&age=3"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", HandlerContentTypeURLEncoded).Methods("POST")

	router.ServeHTTP(rr, req)
}

type jsonStruc struct {
	ID        int
	Credit    float32
	Limit     int
	Offset    int
	PerPage   int
	beginData string
}

func TestHandlerContentTypeJSON(t *testing.T) {

	js := &jsonStruc{
		1,
		3.3,
		10,
		10,
		100,
		"201810320 09:30:22",
	}

	b, _ := json.Marshal(js)

	path := fmt.Sprintf("/users")
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", HandlerContentTypeJSON).Methods("POST")

	router.ServeHTTP(rr, req)
}

func TestQueryHandler(t *testing.T) {

	tt := []struct {
		id         int
		limit      int
		perPage    int
		offset     int
		shouldPass bool
	}{
		{1, 10, 10, 10, true},
		{2, 20, 20, 20, true},
		{3, 30, 30, 30, true},
		{4, 40, 40, 40, true},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/users/%d?limit=%d&offset=%d&perPage=%d&beginDate='20181030 12:00:01'&endDate='20181030 12:10:00'", tc.id, tc.limit, tc.offset, tc.perPage)
		req, err := http.NewRequest("POST", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		//req.Header.Set("Content-Type","application/x-www-form-urlencoded; param=value")
		//req.Header.Set("Content-Type","application/json")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/users/{id:[0-9]+}", QueryHandler).Queries("limit", "{limit:[0-9]+}", "perPage", "{perPage:[0-9]+}", "offset", "{offset:[0-9]+}", "beginDate", "{beginDate}", "endDate", "{endDate}").Methods("POST")
		router.HandleFunc("/users/{id:[0-9]+}", QueryHandler).Methods("POST")
		router.HandleFunc("/users", QueryHandler).Methods("POST")

		//sub := router.Host("localhost").Subrouter()
		//sub.Path("/users/{id:[0-9]+}").Queries("limit","{limit:[0-9]+}","offset","{offset:[0-9]+}","beginDate","{beginDate}","endDate","{endDate}").HandlerFunc(UserHandler).Methods("POST").Name("users")
		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on routeVariable %d: got %v want %v",
				tc.id, rr.Code, http.StatusOK)
		}
	}
}

func TestMethod(t *testing.T) {

	//post
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/test", MethodGetHandler).Methods("GET")
	router.HandleFunc("/test", MethodPostHandler).Methods("POST")

	router.ServeHTTP(rr, req)

	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if rr.Code == http.StatusOK {
		if rr.Body.String() != `{"MethodGetHandler": true}` {

			t.Errorf("MethodGetHandler error ")
		}
	}

	//get
	req2, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}


	router.ServeHTTP(rr, req2)


	// In this case, our MetricsHandler returns a non-200 response
	// for a route variable it doesn't know about.
	if rr.Code == http.StatusOK  {
		if rr.Body.String()!=`{"MethodPostHandler": true}`{

			t.Errorf("MethodPostHandler error ")
		}
	}

}
