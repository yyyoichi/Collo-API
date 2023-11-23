package server

import (
	"context"
	"errors"
	"log"
	"time"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
	"yyyoichi/Collo-API/internal/app"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/pkg/apperror"

	"connectrpc.com/connect"
)

type TimeoutError struct{ error }

type ColloServer struct{}

func (*ColloServer) ColloStream(cxt context.Context, req *connect.Request[apiv1.ColloStreamRequest], str *connect.ServerStream[apiv1.ColloStreamResponse]) error {
	done := make(chan interface{})
	cxt, cancel := context.WithCancelCause(cxt)
	defer close(done)
	defer cancel(nil)

	service, err := app.NewCollocationService(app.CollocationServiceOptions{
		Any:   req.Msg.Keyword,
		From:  req.Msg.From.AsTime(),
		Until: req.Msg.Until.AsTime(),
	})
	if err != nil {
		cancel(err)
	}

	go func() {
		for pr := range service.Stream(cxt) {
			if pr.Err != nil {
				cancel(pr.Err)
				return
			}
			resp := &apiv1.ColloStreamResponse{}
			resp.Words = pr.WordByID
			resp.Pairs = pr.Pairs
			if err := str.Send(resp); err != nil {
				cancel(err)
				return
			}
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Println("End stream")
		return nil

	case <-time.After(time.Second * 60):
		err = errors.New("timeout")
		handleError(err, "タイムアウトしました。期間を短くするか、キーワードをより具体的にしてください。; "+err.Error())
		return err

	case <-cxt.Done():
		switch err := context.Cause(cxt).(type) {
		case nil:
			return nil
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
}

func handleError(err error, message string) error {
	log.Printf("[Error]: %v\n", message)
	log.Printf("%v \n", err)
	return errors.New(message)
}
