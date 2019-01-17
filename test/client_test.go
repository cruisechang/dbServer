package test

import (
	"testing"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/cruisechang/dbServer/handler"
	"bytes"
	"fmt"
)

type responseData struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type transferGetParam struct {
	PartnerID int64  `json:"partnerID"` //有-1
	UserID    int64  `json:"userID"`    //有-1
	Category  int    `json:"category"`  //有-1
	Status    int    `json:"status"`    //有-1
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

func Test_httpClient_Get(t *testing.T) {

	client, _ := newHTTPClient("139.162.68.65", "15000",5,5, 5)
	//client.SetTargetAddress("139.162.68.65", "15000")

	type args struct {
		path string
	}
	tests := []struct {
		name        string
		httpStatus  int
		args        args
		wantResBody []byte
		wantErr     bool
	}{
		{
			"0",
			http.StatusOK,
			args{"/transfers/100"},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Do("GET", tt.args.path, nil)

			if err != nil {
				t.Fatalf("get err=%s,name=%s", err.Error(), tt.name)

			}
			if resp.StatusCode != tt.httpStatus {
				t.Fatalf("statusCode=%d,name=%s", resp.StatusCode, tt.name)
			}

			defer resp.Body.Close()

			resBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body err=%s,name=%s", err.Error(), tt.name)
			}

			resData := &responseData{
			}
			err = json.Unmarshal(resBody, resData)

			if (err != nil) != tt.wantErr {
				t.Errorf("httpClient.Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(resData, tt.wantResBody) {
			//	t.Errorf("httpClient.Post() = %v, want %v", gotResBody, tt.wantResBody)
			//}

		})
	}
}

func Test_httpClient_RoundGet(t *testing.T) {

	client, _ := newHTTPClient("139.162.68.65", "15000",5,5, 5)
	//client.SetTargetAddress("139.162.68.65", "15000")
	//client.SetTargetAddress("localhost", "15000")

	type roundGetParam struct {
		HallID    int    `json:"round"`    //有-1
		RoomID    int    `json:"roomID"`   //有-1
		RoomType  int    `json:"roomType"` //有-1
		Status    int    `json:"status"`   //有-1
		BeginDate string `json:"beginDate"`
		EndDate   string `json:"endDate"`
	}

	tests := []struct {
		name  string
		code  int
		count int
		param roundGetParam
	}{
		{"0", handler.CodeSuccess, 11, roundGetParam{-1, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"1", handler.CodeSuccess, 6, roundGetParam{100, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"2", handler.CodeSuccess, 2, roundGetParam{200, 600, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"3", handler.CodeSuccess, 1, roundGetParam{100, 1, 0, 0, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"4", handler.CodeSuccess, 3, roundGetParam{-1, -1, 0, -1, "2018-12-18 14:57:53", "2018-12-20 23:59:59"}},
	}

	path := fmt.Sprintf("/rounds")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b, _ := json.Marshal(tt.param)

			resp, err := client.Do("GET", path, bytes.NewBuffer(b))

			if err != nil {
				t.Fatalf("get err=%s,name=%s", err.Error(), tt.name)

			}

			defer resp.Body.Close()

			resBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body err=%s,name=%s", err.Error(), tt.name)
			}

			resData := &responseData{
			}
			err = json.Unmarshal(resBody, resData)

			t.Logf("handler unmarshal responseData resData=%+v", resData)

		})
	}
}

func Test_httpClient_RoundGetID(t *testing.T) {

		client, _ := newHTTPClient("139.162.68.65", "15000",5,5, 5)
		//client.SetTargetAddress("139.162.68.65", "15000")

		path := fmt.Sprintf("/rounds/16818356")

		resp, err := client.Do("GET", path, nil)

		if err != nil {
			//t.Fatalf("get err=%s,name=%s", err.Error(), tt.name)

		}

		defer resp.Body.Close()

		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//t.Fatalf("read body err=%s,name=%s", err.Error(), tt.name)
		}

		resData := &responseData{
		}
		err = json.Unmarshal(resBody, resData)

		t.Logf("handler unmarshal responseData resData=%+v", resData)

	/*

	type roundGetParam struct {
		HallID    int    `json:"round"`    //有-1
		RoomID    int    `json:"roomID"`   //有-1
		RoomType  int    `json:"roomType"` //有-1
		Status    int    `json:"status"`   //有-1
		BeginDate string `json:"beginDate"`
		EndDate   string `json:"endDate"`
	}

	tests := []struct {
		name  string
		code  int
		count int
		param roundGetParam
	}{
		{"0", handler.CodeSuccess, 11, roundGetParam{-1, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"1", handler.CodeSuccess, 6, roundGetParam{100, -1, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"2", handler.CodeSuccess, 2, roundGetParam{200, 600, -1, -1, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"3", handler.CodeSuccess, 1, roundGetParam{100, 1, 0, 0, "2018-12-10 10:00:00", "2018-12-20 23:59:59"}},
		{"4", handler.CodeSuccess, 3, roundGetParam{-1, -1, 0, -1, "2018-12-18 14:57:53", "2018-12-20 23:59:59"}},
	}


	path := fmt.Sprintf("/rounds/16818356")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {



			resp, err := client.Do("POST",path,nil)

			if err != nil {
				t.Fatalf("get err=%s,name=%s", err.Error(), tt.name)

			}

			defer resp.Body.Close()

			resBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body err=%s,name=%s", err.Error(), tt.name)
			}


			resData := &responseData{
			}
			err = json.Unmarshal(resBody, resData)

			t.Logf("handler unmarshal responseData resData=%+v",  resData)

		})
	}
	*/
}

/*
func Test_httpClient_Post(t *testing.T) {

	client, _ := newHTTPClient(5, 5)
	client.SetTargetAddress("139.162.68.65", "15000")



	b0, _ := json.Marshal(transferGetParam{100, -1, -1, -1, "2018-11-28 10:00:00", "2018-12-21 23:59:59"})

	type args struct {
		path string
		body io.Reader
	}
	tests := []struct {
		name        string
		httpStatus int
		args        args
		wantResBody []byte
		wantErr     bool
	}{
		{
			"0",
			http.StatusOK,
			args{"/transfers/100", bytes.NewBuffer(b0),},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.GET(tt.args.path)

			if err!=nil{
				t.Fatalf("get err=%s,name=%s",err.Error(),tt.name)

			}
			if resp.StatusCode!=tt.httpStatus{
				t.Fatalf("statusCode=%d,name=%s",resp.StatusCode,tt.name)
			}

			if resp.Body != nil {
				defer resp.Body.Close()

				resBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				return
			}

			if err!=nil{
				t.Logf("post err =%s",err.Error())
			}

			resData := &responseData{
			}
			err = json.Unmarshal(gotResBody, resData)

			if (err != nil) != tt.wantErr {
				t.Errorf("httpClient.Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResBody, tt.wantResBody) {
				t.Errorf("httpClient.Post() = %v, want %v", gotResBody, tt.wantResBody)
			}

		})
	}
}
*/
