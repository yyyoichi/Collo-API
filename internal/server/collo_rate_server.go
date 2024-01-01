package server

import (
	"context"
	"fmt"
	"time"
	"yyyoichi/Collo-API/internal/analyzer"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/api/v2/apiv2connect"
	"yyyoichi/Collo-API/internal/matrix"
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

	handleErr := func(err error) {
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

	ndlConfig := svr.ndlConfig
	l, _ := time.LoadLocation("Asia/Tokyo")
	ndlConfig.Search.Any = req.Msg.Keyword
	ndlConfig.Search.From = req.Msg.From.AsTime().In(l)
	ndlConfig.Search.Until = req.Msg.Until.AsTime().In(l)
	ndlConfig.NDLAPI = ndl.SpeechAPI
	analyzerConfig := analyzer.Config{}
	analyzerConfig.Includes = make([]analyzer.PartOfSpeechType, len(req.Msg.PartOfSpeechTypes))
	for i, t := range req.Msg.PartOfSpeechTypes {
		analyzerConfig.Includes[i] = analyzer.PartOfSpeechType(t)
	}
	analyzerConfig.StopWords = req.Msg.Stopwords
	handler := provider.Handler[*apiv2.ColloRateWebStreamResponse]{}
	handler.Err = handleErr
	handler.Resp = handleResp
	handler.Done = handleDone

	var v2provider provider.V2RateProviderInterface
	switch req.Msg.Mode {
	case uint32(0):
		// シングルモード
		matrixConfig := matrix.Config{}
		matrixConfig.PickDocGroupID = func(*matrix.Document) string { return "all" }
		v2provider = provider.NewV2RateProvider(ctx, ndlConfig, analyzerConfig, matrixConfig, handler)
	case uint32(1):
		// マルチモード
		matrixConfig := matrix.Config{}
		matrixConfig.PickDocGroupID = func(d *matrix.Document) string { return d.Key }
		matrixConfig.ReduceThreshold = 0.01 // 1%の単語利用
		if req.Msg.ForcusGroupId == "" || req.Msg.ForcusGroupId == "all" {
			// all
			matrixConfig.AtGroupID = ""
			v2provider = provider.NewV2NultiRateProvider(ctx, ndlConfig, analyzerConfig, matrixConfig, handler)
		} else {
			// forcused
			matrixConfig.AtGroupID = req.Msg.ForcusGroupId
			v2provider = provider.NewV2RateProvider(ctx, ndlConfig, analyzerConfig, matrixConfig, handler)
		}
	}
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
	case ndl.NdlError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録データの取得に失敗しました。; %s", err.Error()),
		)
	case analyzer.AnalysisError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録を形態素解析結果中にエラーが発生しました。; %s", err.Error()),
		)
	case matrix.MatrixError:
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
