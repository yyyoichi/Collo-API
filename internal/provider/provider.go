package provider

import (
	"context"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/internal/pair"
	"yyyoichi/Collo-API/pkg/stream"
)

type MatrixProvider struct {
	handler Handler[any]
}

func NewV3Provider(ctx context.Context, ndlConfig ndl.Config, handler Handler[any]) *MatrixProvider {
	p := &MatrixProvider{
		handler: handler,
	}
	_ = p.getWightDocMatrix(ctx, ndlConfig)
	return p
}

func (p *MatrixProvider) getWightDocMatrix(ctx context.Context, ndlConfig ndl.Config) matrix.DocMatrixInterface {
	m := ndl.NewMeeting(ndlConfig)
	// 会議APIから結果取得
	meetingResultCh := m.GenerateMeeting(ctx)
	// 会議ごとの発言
	meetingCh := stream.Demulti[*ndl.MeetingResult, string](ctx, meetingResultCh, func(mr *ndl.MeetingResult) []string {
		if mr.Error() != nil {
			p.handler.Err(mr.Error())
		}
		return mr.GetSpeechsPerMeeting()
	})
	// 会議ごとの単語
	wordsCh := stream.FunIO[string, []string](ctx, meetingCh, func(meeting string) []string {
		pr := pair.MAnalytics.Parse(meeting)
		if pr.Error() != nil {
			p.handler.Err(pr.Error())
		}
		return pr.GetNouns()
	})

	builder := matrix.NewMatrixBuilder()
	for words := range wordsCh {
		builder.AppendDoc(words)
	}
	return builder.BuildTFIDFDocMatrix()
}
