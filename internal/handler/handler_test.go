package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"
	"yyyoichi/Collo-API/internal/api/v3/apiv3connect"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandlers(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(apiv3connect.NewMintGreenServiceHandler(&V3Handler{}))
	server := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
	server.Start()
	t.Cleanup(server.Close)

	l, _ := time.LoadLocation("Asia/Tokyo")
	var initV3ReqConfig = func(v3req *apiv3.RequestConfig) *apiv3.RequestConfig {
		if v3req == nil {
			v3req = &apiv3.RequestConfig{}
		}
		v3req.Keyword = "科学"
		v3req.From = timestamppb.New(time.Date(2023, 11, 1, 0, 0, 0, 0, l))
		v3req.Until = timestamppb.New(time.Date(2023, 11, 5, 0, 0, 0, 0, l))
		v3req.UseNdlCache = true
		v3req.CreateNdlCache = true
		return v3req
	}
	var createNetworkRequestStream = func(req *apiv3.NetworkStreamRequest) (
		*connect.ServerStreamForClient[apiv3.NetworkStreamResponse],
		error,
	) {
		req.Config = initV3ReqConfig(req.Config)
		client := apiv3connect.NewMintGreenServiceClient(
			http.DefaultClient,
			server.URL,
		)
		return client.NetworkStream(
			context.Background(),
			connect.NewRequest(req),
		)
	}
	var createNodeRateRequestStream = func(req *apiv3.NodeRateStreamRequest) (
		*connect.ServerStreamForClient[apiv3.NodeRateStreamResponse],
		error,
	) {
		req.Config = initV3ReqConfig(req.Config)
		client := apiv3connect.NewMintGreenServiceClient(
			http.DefaultClient,
			server.URL,
		)
		return client.NodeRateStream(
			context.Background(),
			connect.NewRequest(req),
		)
	}
	t.Run("Regular Network", func(t *testing.T) {
		t.Parallel()
		test := []*apiv3.NetworkStreamRequest{
			{},
			{ForcusNodeId: 1},
		}
		for _, req := range test {
			stream, err := createNetworkRequestStream(req)
			require.NoError(t, err)
			var hasdata bool
			for stream.Receive() {
				require.NoError(t, stream.Err())
				resp := stream.Msg()
				if resp.Nodes != nil && resp.Edges != nil {
					hasdata = true
				}
			}
			require.True(t, hasdata)
		}
	})
	t.Run("Error Network", func(t *testing.T) {
		t.Parallel()
		req := &apiv3.NetworkStreamRequest{
			Config: &apiv3.RequestConfig{
				ForcusGroupId: "not found group id",
			},
		}
		stream, err := createNetworkRequestStream(req)
		require.NoError(t, err)
		for stream.Receive() {
		}
		require.Error(t, stream.Err())
	})
	t.Run("Regular NodeRate", func(t *testing.T) {
		t.Parallel()
		req := &apiv3.NodeRateStreamRequest{}
		stream, err := createNodeRateRequestStream(req)
		require.NoError(t, err)
		var hasdata bool
		for stream.Receive() {
			require.NoError(t, stream.Err())
			resp := stream.Msg()
			if resp.Nodes != nil {
				hasdata = true
			}
		}
		require.True(t, hasdata)
	})
	t.Run("Error NodeRate", func(t *testing.T) {
		t.Parallel()
		req := &apiv3.NodeRateStreamRequest{
			Config: &apiv3.RequestConfig{
				ForcusGroupId: "not found group id",
			},
		}
		stream, err := createNodeRateRequestStream(req)
		require.NoError(t, err)
		for stream.Receive() {
		}
		require.Error(t, stream.Err())
	})
}
