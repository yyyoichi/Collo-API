package matrix

import (
	"context"
	"testing"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/internal/pair"
	"yyyoichi/Collo-API/pkg/stream"

	"github.com/stretchr/testify/require"
)

func generateDocs() [][]string {
	ctx := context.Background()
	m := ndl.NewMeeting(ndl.CreateMeetingConfigMock(ndl.Config{}, ""))
	// 会議APIから結果取得
	meetingResultCh := m.GenerateMeeting(ctx)
	// 会議ごとの発言
	meetingCh := stream.Demulti[*ndl.MeetingResult, string](ctx, meetingResultCh, func(mr *ndl.MeetingResult) []string {
		if mr.Error() != nil {
			panic(mr.Error())
		}
		return mr.GetSpeechsPerMeeting()
	})
	// 会議ごとの単語
	wordsCh := stream.FunIO[string, []string](ctx, meetingCh, func(meeting string) []string {
		pr := pair.MAnalytics.Parse(meeting)
		if pr.Error() != nil {
			panic(pr.Error())
		}
		return pr.GetNouns()
	})

	docs := [][]string{}
	for doc := range wordsCh {
		docs = append(docs, doc)
	}
	return docs
}

func TestCoMatrix(t *testing.T) {
	docs := generateDocs()
	t.Run("Create CoMatrix", func(t *testing.T) {
		b := NewBuilder()
		for _, doc := range docs {
			b.AppendDoc(doc)
		}
		m := NewCoMatrixFromBuilder(b, Config{})
		for p := range m.progress {
			t.Log(p)
			if p == ProgressDone || p == ErrDone {
				break
			}
		}

		require.NoError(t, m.err)
	})
}
