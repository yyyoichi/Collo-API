package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/pair"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWebServer(t *testing.T) {
	config := pair.Config{}
	config = pair.CreateMockConfig(config)
	mux := createHandler(config)
	server := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
	server.Start()
	t.Cleanup(server.Close)

	t.Run("init", func(t *testing.T) {
		stream, err := request(&apiv2.ColloWebStreamRequest{
			Keyword: config.Search.Any,
			From:    timestamppb.New(config.Search.From),
			Until:   timestamppb.New(config.Search.Until),
		}, server.URL)
		require.NoError(t, err)
		i := 0
		needs := 0
		for stream.Receive() {
			require.NoError(t, stream.Err())
			resp := stream.Msg()
			if i == 0 {
				needs = int(resp.Needs)
			} else if i <= needs {
				require.Equal(t, i, int(resp.Dones))
			} else {
				require.NotEqual(t, 0, len(resp.Nodes))
			}
			i++
		}
		require.NotEqual(t, 0, i)

		stream, err = request(&apiv2.ColloWebStreamRequest{
			Keyword:      config.Search.Any,
			From:         timestamppb.New(config.Search.From),
			Until:        timestamppb.New(config.Search.Until),
			ForcusNodeId: 1,
		}, server.URL)
		require.NoError(t, err)

		i = 0
		for stream.Receive() {
			require.NoError(t, stream.Err())
			resp := stream.Msg()
			require.NotNil(t, resp.Nodes)
			require.NotEqual(t, 0, len(resp.Nodes))
			i++
		}
		require.NotEqual(t, 0, i)
	})

}

func request(req *apiv2.ColloWebStreamRequest, url string) (
	*connect.ServerStreamForClient[apiv2.ColloWebStreamResponse],
	error,
) {
	client := apiv2connect.NewColloWebServiceClient(
		http.DefaultClient,
		url,
	)
	return client.ColloWebStream(
		context.Background(),
		connect.NewRequest(req),
	)
}

func createHandler(config pair.Config) http.Handler {
	svr := &ColloWebServer{
		kokkaiRequestConfig: config,
	}
	mux := http.NewServeMux()
	mux.Handle(apiv2connect.NewColloWebServiceHandler(svr))
	return mux
}
