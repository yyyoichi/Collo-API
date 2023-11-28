package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/pair"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer(t *testing.T) {
	config := pair.Config{}
	config = pair.CreateMockConfig(config)
	mux := createServer(config)
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	t.Run("Reguler", func(t *testing.T) {
		client := apiv2connect.NewColloNetworkServiceClient(
			server.Client(),
			server.URL,
		)
		stream := client.ColloNetworkStream(context.Background())
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			defer stream.CloseRequest()
			err := stream.Send(&apiv2.ColloNetworkStreamRequest{
				Keyword:      config.Search.Any,
				From:         timestamppb.New(config.Search.From),
				Until:        timestamppb.New(config.Search.Until),
				ForcusNodeId: 1,
			})
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
			err = stream.Send(&apiv2.ColloNetworkStreamRequest{
				Keyword:      config.Search.Any,
				From:         timestamppb.New(config.Search.From),
				Until:        timestamppb.New(config.Search.Until),
				ForcusNodeId: 2,
			})
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}()
		go func() {
			defer wg.Done()
			defer stream.CloseResponse()
			// needs
			needs := 0
			// node,edgeデータ送信された回数
			count := 0
			// i データ受け取り回数
			i := 0
			for {
				resp, err := stream.Receive()
				if err != nil && !errors.Is(err, io.EOF) {
					require.NoError(t, err)
				}
				if i == 0 {
					require.NotEmpty(t, resp.Needs)
					needs = int(resp.Needs)
					continue
				}
				if i <= needs {
					require.Equal(t, resp.Dones, i)
					return
				}
				if i > needs {
					count++
					require.NotNil(t, resp.Nodes)
					require.NotNil(t, resp.Edges)
					if count == 2 {
						return
					}
				}
				i++
			}
		}()
		wg.Wait()
	})
}

func createServer(config pair.Config) http.Handler {
	svr := &ColloNetworkServer{
		kokkaiRequestConfig: config,
	}
	mux := http.NewServeMux()
	mux.Handle(apiv2connect.NewColloNetworkServiceHandler(svr))
	return mux
}
