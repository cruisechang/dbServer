package test

import (
	"net/http"
	"time"
	"io"
)

type httpClient struct {
	targetAddr string
	client     *http.Client
	postURI    string
	postQuery  string
}

func newHTTPClient(address,port string,connectTimeout,handshakeTimeout, requestTimeout int) (*httpClient, error) {
	re := &httpClient{
		targetAddr:"http://" + address + ":" + port ,
	}

	var netTransport = &http.Transport{
		//Dial: (&net.Dialer{
		//	Timeout: connectTimeout * time.Second,
		//}).Dial,
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second,
	}

	re.client = &http.Client{
		Timeout:   time.Duration(requestTimeout) * time.Second,
		Transport: netTransport,

	}

	return re, nil
}
func (h *httpClient) SetTargetAddress(addr, port string) {
	h.targetAddr = "http://" + addr + ":" + port
}

func (h *httpClient) Do(method ,path string,body io.Reader) (*http.Response,error) {


	req, err := http.NewRequest(method, h.targetAddr+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

	return  h.client.Do(req)

}
