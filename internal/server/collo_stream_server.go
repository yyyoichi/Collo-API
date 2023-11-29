package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/network"
	"yyyoichi/Collo-API/internal/pair"

	"connectrpc.com/connect"
)

type ColloNetworkServer struct {
	kokkaiRequestConfig pair.Config

	apiv2connect.UnimplementedColloNetworkServiceHandler
}

func NewColloNetworkServer() *ColloNetworkServer {
	return &ColloNetworkServer{
		kokkaiRequestConfig: pair.Config{},
	}
}

func (svr *ColloNetworkServer) ColloNetworkStream(
	ctx context.Context,
	stream *connect.BidiStream[apiv2.ColloNetworkStreamRequest, apiv2.ColloNetworkStreamResponse],
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
	handleResp := func(resp *apiv2.ColloNetworkStreamResponse) {
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

	init := func(req *apiv2.ColloNetworkStreamRequest) *network.NetworkProvider {
		config := svr.kokkaiRequestConfig
		l, _ := time.LoadLocation("Asia/Tokyo")
		config.Search.Any = req.Keyword
		config.Search.From = req.From.AsTime().In(l)
		config.Search.Until = req.Until.AsTime().In(l)
		handler := network.Handler{}
		handler.Err = handlerErr
		handler.Resp = handleResp
		handler.Done = handleDone
		return network.NewNetworkProvider(ctx, config, handler)
	}
	go func() {
		var networkpv *network.NetworkProvider
		for {
			req, err := stream.Receive()
			if errors.Is(err, io.EOF) {
				cancel(nil)
			}
			if err != nil {
				cancel(err)
			}
			select {
			case <-ctx.Done():
				return
			default:
				if networkpv == nil {
					networkpv = init(req)
				} else {
					go networkpv.StreamNetworksAround(uint(req.ForcusNodeId))
				}
			}
		}
	}()
	<-ctx.Done()
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
