package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
	"yyyoichi/Collo-API/internal/api/v1/apiv1connect"
	"yyyoichi/Collo-API/internal/pair"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer(t *testing.T) {
	config := pair.Config{}
	config = pair.CreateMockConfig(config)
	server := createServer(config)
	defer server.Close()

	t.Run("Reguler", func(t *testing.T) {
		stream, err := request(config, server.URL)
		if err != nil {
			t.Error(err)
		}
		defer stream.Close()
		for stream.Receive() {
		}
		if stream.Err() != nil {
			t.Error(stream.Err())
		}
	})

	t.Run("ExpError", func(t *testing.T) {
		config.Search.Any = "機関車"
		stream, err := request(config, server.URL)
		if err != nil {
			t.Error(err)
		}
		defer stream.Close()
		for stream.Receive() {
		}
		if stream.Err() == nil {
			t.Error("expected error")
		}
		if connect.CodeOf(stream.Err()) != connect.CodeInternal {
			t.Errorf("expected 'connect.CodeInternal', but got='%s'", connect.CodeOf(stream.Err()))
		}
	})
}

func request(config pair.Config, url string) (
	*connect.ServerStreamForClient[apiv1.ColloStreamResponse],
	error,
) {
	client := apiv1connect.NewColloServiceClient(
		http.DefaultClient,
		url,
	)
	return client.ColloStream(
		context.Background(),
		connect.NewRequest(&apiv1.ColloStreamRequest{
			Keyword: config.Search.Any,
			From:    timestamppb.New(config.Search.From),
			Until:   timestamppb.New(config.Search.Until),
		}),
	)
}

func createServer(config pair.Config) *httptest.Server {
	svr := &ColloServer{
		pairConfig: config,
	}
	api := http.NewServeMux()
	api.Handle(apiv1connect.NewColloServiceHandler(svr))
	server := httptest.NewUnstartedServer(h2c.NewHandler(api, &http2.Server{}))
	server.Start()
	return server
}
