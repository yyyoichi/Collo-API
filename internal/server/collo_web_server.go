package server

import (
	"context"
	"fmt"
	"time"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/network"
	"yyyoichi/Collo-API/internal/pair"

	"connectrpc.com/connect"
)

type ColloWebServer struct {
	kokkaiRequestConfig pair.Config
	apiv2connect.ColloWebServiceHandler
}

func NewColloWebServer() *ColloWebServer {
	return &ColloWebServer{
		kokkaiRequestConfig: pair.Config{},
	}
}

func (svr *ColloWebServer) ColloWebStream(
	ctx context.Context,
	req *connect.Request[apiv2.ColloWebStreamRequest],
	stream *connect.ServerStream[apiv2.ColloWebStreamResponse],
) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	handlerErr := func(err error) {
		select {
		case <-ctx.Done():
			return
		default:
			cancel(err)
		}
	}
	handleResp := func(resp *apiv2.ColloWebStreamResponse) {
		select {
		case <-ctx.Done():
			return
		default:
			if err := stream.Send(resp); err != nil {
				cancel(err)
			}
		}
	}
	handleDone := func() {
		select {
		case <-ctx.Done():
			return
		default:
			cancel(nil)
		}
	}

	config := svr.kokkaiRequestConfig
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.Any = req.Msg.Keyword
	config.Search.From = req.Msg.From.AsTime().In(l)
	config.Search.Until = req.Msg.Until.AsTime().In(l)
	handler := network.Handler{}
	handler.Err = handlerErr
	handler.Resp = handleResp
	handler.Done = handleDone

	networkProvider := network.NewNetworkProvider(ctx, config, handler)
	select {
	case <-ctx.Done():
	default:
		if req.Msg.ForcusNodeId == 0 {
			// initial request
			nodeID := networkProvider.GetByWord(req.Msg.Keyword)
			if nodeID == 0 {
				nodeID = networkProvider.GetCenterNodeID()
			}
			networkProvider.StreamNetworksWith(nodeID)
		} else {
			networkProvider.StreamNetworksAround(uint(req.Msg.ForcusNodeId))
		}
		cancel(nil)
	}

	err := context.Cause(ctx)
	if err == context.Canceled {
		return nil
	}
	switch err.(type) {
	case network.FetchError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録データの取得に失敗しました。; %s", err.Error()),
		)
	case network.ParseError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録を形態素解析結果中にエラーが発生しました。; %s", err.Error()),
		)
	default:
		return connect.NewError(
			connect.CodeUnknown,
			fmt.Errorf("予期せぬエラーが発生しました。; %s", err.Error()),
		)
	}
}
