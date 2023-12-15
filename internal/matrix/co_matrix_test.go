package matrix

import (
	"context"
	"testing"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/internal/pair"
	"yyyoichi/Collo-API/pkg/stream"

	"github.com/stretchr/testify/require"
)

func TestCoMatrixExample(t *testing.T) {
	words := []string{"word1", "word2", "word3", "word4", "word5", "word6"}
	// testcase-> https://qiita.com/igenki/items/a673140ecbfda4ee7dba
	m := &CoMatrix{
		matrix: []float64{
			0, 1, 1, 1, 1, 0,
			1, 0, 0, 0, 0, 0,
			1, 0, 0, 1, 0, 0,
			1, 0, 1, 0, 1, 1,
			1, 0, 0, 1, 0, 1,
			0, 0, 0, 1, 1, 0},
		words: words,
	}
	m.init()
	require.NoError(t, m.useVectorCentrality())
	require.Equal(t, []int{3, 0, 4, 2, 5, 1}, m.indices)
	for i, p := range m.priority {
		t.Logf("'word%d' priority: %v\n", i, p)
	}

	topNode := m.MostImportantNode()
	require.EqualValues(t, 4, topNode.ID)
	node1 := m.NodeID(4)
	require.EqualValues(t, topNode.ID, node1.ID)

	edge := m.Edge(3, 1)
	require.EqualValues(t, 2, edge.ID)
	require.EqualValues(t, 1, edge.Rate)
	edge = m.Edge(4, 3)
	require.EqualValues(t, 15, edge.ID)
	require.EqualValues(t, 1, edge.Rate)

	nodes, edges := m.CoOccurrenceRelation(2)
	require.Equal(t, 1, len(nodes))
	require.Equal(t, 1, len(edges))
	require.EqualValues(t, 1, nodes[0].ID)
	require.EqualValues(t, 1, edges[0].ID)

	nodes, edges = m.CoOccurrenceDept(2, 2)
	require.Equal(t, 5, len(nodes))
	require.Equal(t, 4, len(edges))

	nodes, edges = m.CoOccurrenceDept(3, 6)
	require.Equal(t, 6, len(nodes))
	require.Equal(t, 8, len(edges))

	nodes, edges = m.CoOccurrences(1, 2, 3)
	require.Equal(t, 5, len(nodes))
	require.Equal(t, 5, len(edges))

}

func TestCoMatrix(t *testing.T) {
	docs := generateDocs()
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
	n0 := m.NodeRank(0)
	n1 := m.NodeRank(len(m.words) - 1)
	require.EqualValues(t, 1, n0.Rate)
	require.EqualValues(t, 0, n1.Rate)
}

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
