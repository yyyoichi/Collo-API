package matrix

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"yyyoichi/Collo-API/internal/analyzer"
	"yyyoichi/Collo-API/internal/ndl"
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
	t.Run("CoMatrix", func(t *testing.T) {
		t.Parallel()
		b := NewBuilder()
		for _, doc := range docs {
			b.AppendDocument(doc)
		}
		n, comCh := NewCoMatrixesFromBuilder(context.Background(), b, Config{})
		require.Equal(t, 1, n)
		var m *CoMatrix
		for com := range comCh {
			m = com
		}
		for p := range m.progress {
			if p == ProgressDone || p == ErrDone {
				break
			}
		}

		require.NoError(t, m.err)
		require.NotNil(t, m.meta)
		require.NotEmpty(t, m.meta.Key)
		require.NotEmpty(t, m.meta.Name)
		require.NotEmpty(t, m.meta.Description)
		require.NotNil(t, m.meta.At)
		require.Equal(t, docs[0].Key, m.meta.Key)
		require.Equal(t, docs[0].Name, m.meta.Name)
		require.Equal(t, len(docs), len(strings.Split(m.meta.Description, "- "))-1) // 各メタ情報の先頭に-1が付く
		require.Equal(t, docs[0].At.Format("2006-01-02"), m.meta.At.Format("2006-01-02"))
		n0 := m.NodeRank(0)
		n1 := m.NodeRank(len(m.words) - 1)
		require.EqualValues(t, 1, n0.Rate)
		require.EqualValues(t, 0, n1.Rate)
	})
	t.Run("Multi CoMatrix", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		b := NewBuilder()
		for _, doc := range docs {
			b.AppendDocument(doc)
		}
		var config Config
		config.PickDocGroupID = func(d *Document) string {
			return d.Key
		}
		n, _ := NewCoMatrixesFromBuilder(ctx, b, config)
		require.Equal(t, len(docs), n)
	})
}

func generateDocs() []*Document {
	ctx := context.Background()
	m := ndl.NewMeeting(ndl.CreateMeetingConfigMock(ndl.Config{}, ""))
	// 会議APIから結果取得
	meetingResultCh := m.GenerateMeeting(ctx)
	// 会議ごとの発言
	meetingCh := stream.Demulti[*ndl.MeetingResult, *ndl.MeetingRecode](ctx, meetingResultCh, func(mr *ndl.MeetingResult) []*ndl.MeetingRecode {
		if mr.Error() != nil {
			panic(mr.Error())
		}
		return ndl.NewMeetingRecodes(mr)
	})
	// 会議-単語
	docsCh := stream.FunIO[*ndl.MeetingRecode, *Document](ctx, meetingCh, func(meeting *ndl.MeetingRecode) *Document {
		ar := analyzer.Analysis(meeting.Speeches)
		if ar.Error() != nil {
			panic(ar.Error())
		}
		doc := &Document{}
		doc.Words = ar.Get(analyzer.Config{
			Includes: []analyzer.PartOfSpeechType{
				analyzer.Noun,
			},
		})
		doc.Key = meeting.IssueID
		doc.Name = fmt.Sprintf("%s %s", meeting.NameOfHouse, meeting.NameOfMeeting)
		doc.At = meeting.Date
		doc.Description = fmt.Sprintf("%s %s", meeting.NameOfHouse, meeting.NameOfMeeting)
		return doc
	})

	docs := []*Document{}
	for doc := range docsCh {
		docs = append(docs, doc)
	}
	return docs
}
