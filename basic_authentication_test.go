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
			base64.StdEncoding.EncodeToString([]byte("test:password")): "",
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

func TestBasicAuthentication_Return401IfWrongCredentialProvided(t *testing.T) {
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
			base64.StdEncoding.EncodeToString([]byte("test:password")): "",
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

	credential := base64.StdEncoding.EncodeToString([]byte("test:wrongpassword"))
	req.Header.Set("Authorization", "Basic "+credential)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 401 {
		t.Error("Expected status code 401 but got", resp.StatusCode)
	}
}

func TestBasicAuthentication_Return200IfCorrectCredentialProvided(t *testing.T) {
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
			base64.StdEncoding.EncodeToString([]byte("test:password")): "",
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

	credential := base64.StdEncoding.EncodeToString([]byte("test:password"))
	req.Header.Set("Authorization", "Basic "+credential)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Expected status code 200 but got", resp.StatusCode)
	}
}
