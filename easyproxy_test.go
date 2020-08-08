package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestProxySupportHttpGet(t *testing.T) {
	server := createTestServer()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		_ = server.ListenAndServe()
	}()

	defer cancel()
	defer server.Shutdown(ctx)

	proxy := CreateProxy("localhost:8080", "http://localhost:9090")
	go proxy.ListenAndServe()
	defer proxy.Shutdown(ctx)

	reqHeaders := make(map[string][]string)
	reqHeaders["Content-Type"] = []string{"application/json"}
	resp, err := sendHttpRequest("http://localhost:8080", "GET", reqHeaders)

	if err != nil {
		t.Error()
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(*resp), &data)

	method := data["method"]
	if method != "GET" {
		t.Error("Expected GET but got", method)
	}

	body := data["body"]
	if body != "" {
		t.Error("Expected \"\" but got", body)
	}

	headers := data["headers"].(map[string]interface{})

	contentType := headers["Content-Type"].([]interface{})
	if contentType[0].(string) != "application/json" {
		t.Error("Expected \"application/json\" but got", contentType)
	}

}

func TestProxySupportHttpPost(t *testing.T) {
	server := createTestServer()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		_ = server.ListenAndServe()
	}()

	defer cancel()

	proxy := CreateProxy("localhost:8080", "http://localhost:9090")
	go proxy.ListenAndServe()
	defer proxy.Shutdown(ctx)

	resp, err := httpPost("http://localhost:8080")

	if err != nil {
		t.Error()
	}

	expected := "{\"method\":\"POST\",\"body\":\"id=123\\u0026key=Value\",\"headers\":{\"Accept-Encoding\":[\"gzip\"],\"Content-Type\":[\"application/x-www-form-urlencoded\"],\"User-Agent\":[\"Go-http-client/1.1\"]}}"
	if expected != *resp {
		t.Error("Expected", expected, "but got", *resp)
	}

	server.Shutdown(ctx)
}

func TestProxySupportHttpPut(t *testing.T) {
	server := createTestServer()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		_ = server.ListenAndServe()
	}()

	defer cancel()
	defer server.Shutdown(ctx)

	proxy := CreateProxy("localhost:8080", "http://localhost:9090")
	go proxy.ListenAndServe()
	defer proxy.Shutdown(ctx)

	resp, err := httpPut("http://localhost:8080")

	if err != nil {
		t.Error()
	}

	expected := "{\"method\":\"PUT\",\"body\":\"\",\"headers\":{\"Accept-Encoding\":[\"gzip\"],\"Content-Length\":[\"0\"],\"Content-Type\":[\"application/json; charset=utf-8\"],\"User-Agent\":[\"Go-http-client/1.1\"]}}"
	if expected != *resp {
		t.Error("Expected", expected, "but got", *resp)
	}
}

func httpGet(url string) (content *string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	c := string(body)
	return &c, nil
}

func httpPost(address string) (content *string, err error) {
	resp, err := http.PostForm(address, url.Values{"key": {"Value"}, "id": {"123"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	c := string(body)
	return &c, nil
}

func httpPut(address string) (content *string, err error) {
	req, err := http.NewRequest("PUT", address, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	c := string(body)
	return &c, nil
}

func sendHttpRequest(address string, method string, headers map[string][]string) (content *string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, address, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	c := string(body)
	return &c, nil
}

type TestServerHandler struct{}

type TestServerResponse struct {
	Method  string              `json:"method"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

func (s *TestServerHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		wr.WriteHeader(500)
		wr.Write([]byte(err.Error()))
	}

	data := TestServerResponse{
		Method:  req.Method,
		Body:    string(body),
		Headers: req.Header,
	}

	body, err = json.Marshal(data)
	wr.Write(body)
}

func createTestServer() http.Server {
	address := "localhost:9090"
	handler := &TestServerHandler{}
	fmt.Println("HTTP server is listening on", address)
	server := http.Server{Addr: address, Handler: handler}
	return server
}
