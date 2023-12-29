package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/ndl"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRateWebServer(t *testing.T) {
	config := ndl.Config{}
	config = ndl.CreateMeetingConfigMock(config, "")
	mux := createRateHandler(config)
	server := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
	server.Start()
	t.Cleanup(server.Close)

	t.Run("Init", func(t *testing.T) {
		stream, err := rateRequest(&apiv2.ColloRateWebStreamRequest{
			Keyword:           config.Search.Any,
			From:              timestamppb.New(config.Search.From),
			Until:             timestamppb.New(config.Search.Until),
			ForcusNodeId:      0,
			PartOfSpeechTypes: []uint32{101, 401},
			Stopwords:         []string{"発展", "開発"},
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
	})

	t.Run("Forcus NodeID", func(t *testing.T) {
		stream, err := rateRequest(&apiv2.ColloRateWebStreamRequest{
			Keyword:      config.Search.Any,
			From:         timestamppb.New(config.Search.From),
			Until:        timestamppb.New(config.Search.Until),
			ForcusNodeId: 1,
		}, server.URL)
		require.NoError(t, err)

		i := 0
		for stream.Receive() {
			require.NoError(t, stream.Err())
			resp := stream.Msg()
			if resp.Dones != resp.Needs {
				continue
			}
			i++
			// done == needsの一度目は完了通知のみ
			if i == 1 {
				continue
			}
			if resp.Dones == resp.Needs {
				require.NotNil(t, resp.Nodes)
				require.NotEqual(t, 0, len(resp.Nodes))
			}
		}
		require.NotEqual(t, 0, i)
	})

}

func rateRequest(req *apiv2.ColloRateWebStreamRequest, url string) (
	*connect.ServerStreamForClient[apiv2.ColloRateWebStreamResponse],
	error,
) {
	client := apiv2connect.NewColloRateWebServiceClient(
		http.DefaultClient,
		url,
	)
	return client.ColloRateWebStream(
		context.Background(),
		connect.NewRequest(req),
	)
}

func createRateHandler(config ndl.Config) http.Handler {
	svr := &ColloRateWebServer{
		ndlConfig: config,
	}
	mux := http.NewServeMux()
	mux.Handle(apiv2connect.NewColloRateWebServiceHandler(svr))
	return mux
}
