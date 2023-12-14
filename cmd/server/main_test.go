package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionReq(t *testing.T) {
	mux := getHandler()
	server := httptest.NewUnstartedServer(mux)
	server.Start()
	t.Cleanup(server.Close)

	t.Run("ColloWeb", func(t *testing.T) {
		t.Parallel()
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
	})
	t.Run("ColloRateWeb", func(t *testing.T) {
		t.Parallel()
		// Optionリクエストを送信
		client := server.Client()
		req, err := http.NewRequest(http.MethodOptions, server.URL+"/api.v2.ColloRateWebService/ColloRateWebStream", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})
}
