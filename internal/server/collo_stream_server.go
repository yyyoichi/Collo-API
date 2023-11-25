package server

import (
	"context"
	"errors"
	"fmt"
	"time"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
	"yyyoichi/Collo-API/internal/pair"

	"connectrpc.com/connect"
)

type ColloServer struct {
	pairConfig pair.Config
}

func NewColloServer() *ColloServer {
	return &ColloServer{
		pairConfig: pair.Config{},
	}
}

func (svr *ColloServer) ColloStream(ctx context.Context, req *connect.Request[apiv1.ColloStreamRequest], str *connect.ServerStream[apiv1.ColloStreamResponse]) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	config := svr.pairConfig
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.Any = req.Msg.Keyword
	config.Search.From = req.Msg.From.AsTime().In(l)
	config.Search.Until = req.Msg.Until.AsTime().In(l)
	handler := pair.Handler{
		Err: func(err error) {
			cancel(err)
		},
		Resp: func(resp *apiv1.ColloStreamResponse) {
			if err := str.Send(resp); err != nil {
				cancel(err)
			}
		},
		Done: func() {
			cancel(nil)
		},
	}
	if ps, err := pair.NewPairStore(config, handler); err != nil {
		handler.Err(err)
	} else {
		go ps.Stream(ctx)
	}
	<-ctx.Done()
	err := context.Cause(ctx)
	if err == context.Canceled {
		return nil
	}
	switch err.(type) {
	case pair.TimeoutError:
		return connect.NewError(
			connect.CodeDeadlineExceeded,
			errors.New("タイムアウトしました。期間を短くしてください。;"),
		)
	case pair.FetchError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録データの取得に失敗しました。; %s", err.Error()),
		)
	case pair.ParseError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録を形態素解析結果中にエラーが発生しました。; %s", err.Error()),
		)
	default:
		return connect.NewError(
			connect.CodeUnknown,
			errors.New("予期せぬエラーが発生しました。"),
		)
	}
}
