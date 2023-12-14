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
	t.Run("CoMatrix", func(t *testing.T) {
		words := []string{"word1", "word2", "word3", "word4", "word5", "word6"}
		// testcase-> https://qiita.com/igenki/items/a673140ecbfda4ee7dba
		m := &CoMatrix{
			matrix: []float64{0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0},
			words:  words,
		}
		m.init()
		require.NoError(t, m.useVectorCentrality())
		require.Equal(t, []int{3, 0, 4, 2, 5, 1}, m.indices)
		for i, p := range m.priority {
			t.Logf("'word%d' priority: %v\n", i, p)
		}
	})
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
		n0 := m.Node(0)
		n1 := m.Node(len(m.words) - 1)
		require.EqualValues(t, 1, n0.Rate)
		require.EqualValues(t, 0, n1.Rate)
		t.Logf("the most important node is %s(%v), the bottom node is %s(%v)", n0.Word, n0.Rate, n1.Word, n1.Rate)
	})
}
