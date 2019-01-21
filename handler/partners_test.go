package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cruisechang/dbServer/util"
	"github.com/cruisechang/dbex"
	"github.com/gorilla/mux"
)

func TestPartnersHandlerGet(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}
	fmt.Sprintf("%v", dbx)

	tt := []struct {
		name  string
		count int
		param partnerGetParam
	}{

		{"0", 3, partnerGetParam{-1, "", "", -1, -1}},
		{"0", 3, partnerGetParam{-1, "", "", -1, 0}},
		{"0", 3, partnerGetParam{-1, "", "", -1, 5}},
		{"0", 1, partnerGetParam{0, "", "", -1, 5}},
		{"0", 2, partnerGetParam{1, "partnerID", "asc", 0, -1}},
		{"0", 0, partnerGetParam{2, "partnerID", "desc", -1, -1}},
		{"0", 1, partnerGetParam{0, "", "", -1, -1}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners")
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
		router.Handle("/partners", NewPartnersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("GET")

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
	}
}

func TestPartnersHandlerPost(t *testing.T) {

	dbx, err := dbex.NewDBEX("dbexConfig.json")
	dbx.Logger.SetLevel(dbex.LevelInfo)

	if err != nil {
		t.Fatalf("dbex error %s", err.Error())
	}

	//db
	h := NewDealersHandler(NewBaseHandler(dbx.DB, dbx.Logger))
	sqlDB := h.db.GetSQLDB()
	var ids []uint64 //放ids，刪掉用

	//
	//account 不能重複
	//name 不能重複
	//prefix 不能重複
	//

	accounts := []string{}
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999999)), 10))
	accounts = append(accounts, "account"+strconv.FormatInt(int64(util.RandomInt(1, 9999999999)), 10))

	//type param struct {
	//	Account   string `json:"account"`
	//	Password  string `json:"password"`
	//	Name      string `json:"name"`
	//	Level     int    `json:"level"`
	//	Category  int    `json:"category"`
	//	AESKey    string `json:"aesKey"`
	//	AccessToken    string `json:"accessToken"`
	//	APIBindIP string `json:"apiBindIP"`
	//	CMSBindIP string `json:"cmsBindIP"`
	//}

	tt := []struct {
		name  string
		code  int
		param partnerPostParam
	}{

		{"0", CodeSuccess, partnerPostParam{accounts[0], "passd", accounts[0], 0, 0, "ssssssss", "12345678", "[]", "[]"}},
		{"1", CodeSuccess, partnerPostParam{accounts[1], "passd", accounts[1], 0, 0, "ssssssss", "12345678", "[]", "[]"}},
		{"2", CodeSuccess, partnerPostParam{accounts[2], "passd", accounts[2], 0, 0, "ssssssss", "12345678", "[]", "[]"}},
		{"3", CodeSuccess, partnerPostParam{accounts[3], "passd", accounts[3], 0, 0, "ssssssss", "12345678", "[]", "[]"}},
		{"4", CodeSuccess, partnerPostParam{accounts[4], "passd", accounts[4], 0, 0, "ssssssss", "12345678", "[]", "[]"}},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/partners")
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
		router.Handle("/partners", NewPartnersHandler(NewBaseHandler(dbx.DB, dbx.Logger))).Methods("POST")

		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code != http.StatusOK {
			t.Fatalf("handler failed  got %v want %v, name=%s", rr.Code, http.StatusOK, tc.name)
		}

		body, _ := ioutil.ReadAll(rr.Body)

		resData := &struct {
			Code    int
			Count   int
			Message string
			Data    []*partnerIDData
		}{}
		err = json.Unmarshal(body, resData)
		if err != nil {
			t.Fatalf("handler unmarshal responseData error=%s, path=%s, param=%+v", err.Error(), path, tc.param)
		}

		if resData.Code != tc.code {
			t.Fatalf("handler resData code  got %d want %d, name=%s, path=%s, param=%+v", resData.Code, tc.code, tc.name, path, tc.param)

		}

		//insert success
		if resData.Count == 1 {
			t.Logf("id %v ", resData.Data[0].PartnerID)
			ids = append(ids, resData.Data[0].PartnerID)
		}
	}

	if len(ids) > 0 {
		queryString := "DELETE FROM partner  where partner_id = ? LIMIT 1"
		stmt, _ := sqlDB.Prepare(queryString)
		defer stmt.Close()

		for _, v := range ids {
			stmt.Exec(v)
		}
	}
}
