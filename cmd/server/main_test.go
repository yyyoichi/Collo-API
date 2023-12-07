package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestOptionReq(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Optionリクエストを送信
	client := server.Client()
	req, err := http.NewRequest(http.MethodOptions, server.URL+"/api.v2.ColloWebService/ColloWebStream", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func createTestServer() *httptest.Server {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	server := httptest.NewUnstartedServer(getHandler())
	server.StartTLS()
	return server
}
