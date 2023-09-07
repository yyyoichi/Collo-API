package server

import (
	"context"
	"errors"
	"log"
	"time"
	collov1 "yyyoichi/Collo-API/gen/proto/collo/v1"
	"yyyoichi/Collo-API/internal/app"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/pkg/apperror"

	"connectrpc.com/connect"
)

var ErrTimeout = errors.New("timeout")

type TimeoutError struct {
	error
}

type ColloServer struct{}

func (*ColloServer) ColloStream(cxt context.Context, req *connect.Request[collov1.ColloStreamRequest], str *connect.ServerStream[collov1.ColloStreamResponse]) error {
	cxt, cancel := context.WithCancelCause(cxt)
	defer cancel(nil)
	defer func() {
		time.Sleep(time.Second * 30)
		cancel(TimeoutError{ErrTimeout})
	}()

	log.Printf("Get Request: %s\n", req.Header())
	service, err := app.NewCollocationService(app.CollocationServiceOptions{
		Any:   req.Msg.Keyword,
		From:  req.Msg.From.AsTime(),
		Until: req.Msg.Until.AsTime(),
	})
	if err != nil {
		cancel(err)
	}

	for c := range service.Stream(cxt) {
		if c.Err != nil {
			cancel(c.Err)
			break
		}
		var resp *collov1.ColloStreamResponse
		resp.Words = c.WordByID
		resp.Pairs = c.Pairs
		if err := str.Send(resp); err != nil {
			cancel(err)
			break
		}
	}

	cancel(nil)

	switch err := context.Cause(cxt).(type) {
	case nil:
		return nil
	case TimeoutError:
		log.Fatalf(err.Error())
		handleError(err, "タイムアウトしました。期間を短くするか、キーワードをより具体的にしてください。; "+err.Error())
		return err
	case api.FetchError:
		handleError(err, "議事録データの取得に失敗しました。; "+err.Error())
		return err
	case morpheme.ParseError:
		handleError(err, "議事録を形態素解析結果中にエラーが発生しました。; "+err.Error())
		return err
	default:
		err = apperror.WrapError(err, ("There was an unexpected issue; please report this as a bug."))
		handleError(err, "予期せぬエラーが発生しました。")
		return err
	}
}

func handleError(err error, message string) error {
	log.Printf("[Error]: %v\n", message)
	log.Printf("%v \n", err)
	return errors.New(message)
}
