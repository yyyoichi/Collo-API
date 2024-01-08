package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/api/v3/apiv3connect"
)

func TestOptionReq(t *testing.T) {
	mux := getHandler()
	server := httptest.NewUnstartedServer(mux)
	server.Start()
	t.Cleanup(server.Close)

	test := []string{
		apiv2connect.ColloRateWebServiceColloRateWebStreamProcedure,
		apiv3connect.MintGreenServiceNetworkStreamProcedure,
		apiv3connect.MintGreenServiceNodeRateStreamProcedure,
	}
	for _, tt := range test {
		client := server.Client()
		url := fmt.Sprintf("%s/rpc%s", server.URL, tt)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	}
}
