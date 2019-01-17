package test

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"io"
	"github.com/gorilla/mux"
)

func MethodGetHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("MethodGetHandler\n") //get variables

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"MethodGetHandler": true}`)

}

func MethodPostHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("MethodPostHandler\n") //get variables

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"MethodPostHandler": true}`)

}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	for k := range vars {
		fmt.Printf("QueryHandler %s, %s\n", k, vars[k]) //get variables
	}

	if r.Body != nil {
		r.Body.Close()
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"users": true}`)

}
func HandlerHeaderTest(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("API-Key") != "abc123" {
		fmt.Printf("HandlerContentHeaderTest error")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"users": true}`)
}
func HandlerContentTypeURLEncodedParseForm(w http.ResponseWriter, r *http.Request) {

	//body 轉成form，只能讀一次
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	v := r.Form //need parseForm
	lo := v.Get("location")
	age := v.Get("age")
	fmt.Printf("HandlerContentTypeURLEncodedParseForm location=%s, age=%s\n", lo, age)

	//body只能讀一次
	//body, _ := ioutil.ReadAll(r.Body)
	//fmt.Println("response Body:", string(body))
	r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"users": true}`)
}

func HandlerContentTypeURLEncoded(w http.ResponseWriter, r *http.Request) {

	//不需要parseForm
	lc := r.FormValue("location")
	age := r.FormValue("age")
	fmt.Printf("HandlerContentTypeURLEncoded location=%s, age=%s\n", lc, age)

	//body只能讀一次
	//body, _ := ioutil.ReadAll(r.Body)
	//fmt.Println("response Body:", string(body))
	r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "UTF-8")

	io.WriteString(w, `{"users": true}`)
}
func HandlerContentTypeJSON(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	//fmt.Println("UsersHandlerContentTypeJSON response Body:", string(body))

	js := &struct {
		ID        int
		Credit    float32
		Limit     int
		Offset    int
		PerPage   int
		beginData string
	}{}

	json.Unmarshal(body, js)
	fmt.Printf("HandlerContentTypeJSON response Body= %+v\n", js)

	r.Body.Close()
}
