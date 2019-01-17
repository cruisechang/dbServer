package handler

import (
	"net/http"
	"testing"
	"github.com/cruisechang/dbex"
	"fmt"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"database/sql"
)

func Test_bannersHandler_get(t *testing.T) {
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
	}{
		{"0", CodeSuccess, 2, http.StatusOK},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/banners")

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
		router.Handle(path, NewBannersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("handler failed  got %d want %d, name=%s", rr.Code, tc.httpStatus, tc.name)
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

func TestBannersHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)

	sqlDB := dbx.DB.GetSQLDB()
	var stmt *sql.Stmt

	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	type receiveData struct {
		Code    int            `json:"code"`
		Count   int            `json:"count"`
		Message string         `json:"message"`
		Data    []bannerIDData `json:"data"`
	}

	tt := []struct {
		name       string
		code       int
		count      int
		httpStatus int
		param      interface{}
	}{

		{"0", CodeSuccess, 1, http.StatusOK, bannerPostParam{"http://139.162.113.174/resource/lobby/banner/pc1.jpg", "", "", 0, 1}},
		{"1", CodeRequestPostDataIllegal, 0, http.StatusOK, bannerPostParam{"http://", "", "", 0, 1}}, //picurl len <10
		{"2", CodeRequestDataUnmarshalError, 0, http.StatusOK, ""},
		{"3", CodeRequestDataUnmarshalError, 0, http.StatusOK, 1},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/banners")
		b, _ := json.Marshal(tc.param)

		req, err := http.NewRequest("POST", path, bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("charset", "UTF-8")
		req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.Handle(path, NewBannersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != tc.httpStatus {
			t.Fatalf("http status got %d want %d, name=%s", rr.Code, tc.httpStatus, tc.name)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &receiveData{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("unmarshal responseData error=%s name=%s", err.Error(), tc.name)
		}

		if resData.Code != tc.code {
			t.Fatalf("resData code  got %d want %d, name=%s", resData.Code, tc.code, tc.name)
		}

		if resData.Count != tc.count {
			t.Fatalf("resData count  got %d want %d, name=%s", resData.Count, tc.count, tc.name)
		}

		//delete post data
		if len(resData.Data) > 0 {

			queryString := "DELETE FROM banner  where banner_id = ? LIMIT 1"
			stmt, err := sqlDB.Prepare(queryString)
			if err != nil {
				t.Fatalf("resData delete prepare error=%s, name=%s", err.Error(), tc.name)
			}

			_, err = stmt.Exec(resData.Data[0].BannerID)
			if err != nil {
				t.Fatalf("resData delete exec error=%s, name=%s", err.Error(), tc.name)
			}
		}

	}
}
