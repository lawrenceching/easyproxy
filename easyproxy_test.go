package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestProxySupportHttpGet(t *testing.T) {
	go startTestServer()
	go StartProxy("localhost:8080", "http://localhost:9090")

	resp, err := httpGet("http://localhost:8080")

	if err != nil {
		t.Error()
	}

	expected := "GET, \"/\""
	if expected != *resp {
		t.Error("Expected", expected, "but got", *resp)
	}

}

func TestProxySupportHttpPost(t *testing.T) {
	go startTestServer()
	go StartProxy("localhost:8080", "http://localhost:9090")

	resp, err := httpPost("http://localhost:8080")

	if err != nil {
		t.Error()
	}

	expected := "POST, \"/\""
	if expected != *resp {
		t.Error("Expected", expected, "but got", *resp)
	}
}

func TestProxySupportHttpPut(t *testing.T) {
	go startTestServer()
	go StartProxy("localhost:8080", "http://localhost:9090")

	resp, err := httpPut("http://localhost:8080")

	if err != nil {
		t.Error()
	}

	expected := "PUT, \"/\""
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

func sendHttpRequest(address string, method string) (content *string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, address, strings.NewReader(""))
	if err != nil {
		return nil, err
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

type TestServer struct{}

func (s *TestServer) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	method := req.Method
	fmt.Fprintf(wr, "%s, %q", method, html.EscapeString(req.URL.Path))
}

func startTestServer() {
	address := "localhost:9090"
	handler := &TestServer{}
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	fmt.Println("HTTP server is listening on", address)
}
