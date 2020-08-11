package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"
)

func TestBasicAuthentication_Return401IfNoCredentialProvided(t *testing.T) {
	server := createTestServer()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		_ = server.ListenAndServe()
	}()

	defer cancel()
	defer server.Shutdown(ctx)

	config := ProxyConfig{
		isBasicAuthenticationEnabled: true,
		basicAuthenticationCredentials: map[string]string{
			base64.StdEncoding.EncodeToString([]byte("test:password12")): "",
		},
	}

	proxy := CreateProxy("localhost:8080", "http://localhost:9090", config)
	go func() {
		_ = proxy.ListenAndServe()
	}()
	defer proxy.Shutdown(ctx)

	time.Sleep(10 * time.Millisecond)

	reqHeaders := make(map[string][]string)
	reqHeaders["Content-Type"] = []string{"application/json"}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 401 {
		t.Error("Expected status code 401 but got", resp.StatusCode)
	}
}

//
//type TestServerHandler struct{}
//
//type TestServerResponse struct {
//	Method  string              `json:"method"`
//	Body    string              `json:"body"`
//	Headers map[string][]string `json:"headers"`
//}
//
//func (s *TestServerHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
//
//	body, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		wr.WriteHeader(500)
//		wr.Write([]byte(err.Error()))
//	}
//
//	data := TestServerResponse{
//		Method:  req.Method,
//		Body:    string(body),
//		Headers: req.Header,
//	}
//
//	body, err = json.Marshal(data)
//	wr.Write(body)
//}
//
//func createTestServer() http.Server {
//	address := "localhost:9090"
//	handler := &TestServerHandler{}
//	fmt.Println("HTTP server is listening on", address)
//	server := http.Server{Addr: address, Handler: handler}
//	return server
//}
