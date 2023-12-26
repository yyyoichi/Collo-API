package server

import (
	"context"
	"fmt"
	"time"
	"yyyoichi/Collo-API/internal/analyzer"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/internal/provider"

	"connectrpc.com/connect"
)

type ColloRateWebServer struct {
	ndlConfig ndl.Config
	apiv2connect.ColloRateWebServiceHandler
}

func NewColloRateWebServer() *ColloRateWebServer {
	return &ColloRateWebServer{
		ndlConfig: ndl.Config{},
	}
}

func (svr *ColloRateWebServer) ColloRateWebStream(
	ctx context.Context,
	req *connect.Request[apiv2.ColloRateWebStreamRequest],
	stream *connect.ServerStream[apiv2.ColloRateWebStreamResponse],
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
	handleResp := func(resp *apiv2.ColloRateWebStreamResponse) {
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

	config := svr.ndlConfig
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.Any = req.Msg.Keyword
	config.Search.From = req.Msg.From.AsTime().In(l)
	config.Search.Until = req.Msg.Until.AsTime().In(l)
	handler := provider.Handler[*apiv2.ColloRateWebStreamResponse]{}
	handler.Err = handlerErr
	handler.Resp = handleResp
	handler.Done = handleDone

	v2provider := provider.NewV2RateProvider(ctx, config, analyzer.Config{}, handler)
	select {
	case <-ctx.Done():
	default:
		if req.Msg.ForcusNodeId == 0 {
			// initial request
			v2provider.StreamTop3NodesCoOccurrence()
		} else {
			v2provider.StreamForcusNodeIDCoOccurrence(uint(req.Msg.ForcusNodeId))
		}
		cancel(nil)
	}

	err := context.Cause(ctx)
	if err == context.Canceled {
		return nil
	}
	switch err.(type) {
	case provider.NdlError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録データの取得に失敗しました。; %s", err.Error()),
		)
	case provider.AnalysisError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録を形態素解析結果中にエラーが発生しました。; %s", err.Error()),
		)
	case provider.MatrixError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("共起関係の計算に失敗しました。; %s", err.Error()),
		)
	default:
		return connect.NewError(
			connect.CodeUnknown,
			fmt.Errorf("予期せぬエラーが発生しました。; %s", err.Error()),
		)
	}
}
