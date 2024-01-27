package matrix

import (
	"context"
	"fmt"
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
		Matrix: []float64{
			0, 1, 1, 1, 1, 0,
			1, 0, 0, 0, 0, 0,
			1, 0, 0, 1, 0, 0,
			1, 0, 1, 0, 1, 1,
			1, 0, 0, 1, 0, 1,
			0, 0, 0, 1, 1, 0},
		Words: words,
	}
	m.init()
	require.NoError(t, m.useVectorCentrality())
	require.Equal(t, []int{3, 0, 4, 2, 5, 1}, m.Indices)
	for i, p := range m.Priority {
		t.Logf("'word%d' priority: %v\n", i, p)
	}

	topNode := m.MostImportantNode()
	require.EqualValues(t, 4, topNode.ID)
	node1 := m.NodeID(4)
	require.EqualValues(t, topNode.ID, node1.ID)
	require.EqualValues(t, 4, node1.NumEdges)
	require.EqualValues(t, 1, m.NodeID(2).NumEdges)

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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		b := NewBuilder()
		for _, doc := range docs {
			b.AppendDocument(doc)
		}
		n, m, _ := NewMultiCoMatrixFromBuilder(ctx, b, Config{
			ReduceThreshold:  0.001,
			GroupingFuncType: PickAsTotal,
		})
		require.Equal(t, 1, n)
		for p := range m.progress {
			if p == ProgressDone || p == ErrDone {
				break
			}
		}

		require.NoError(t, m.err)
		require.NotNil(t, m.Meta)
		require.NotNil(t, m.Meta.From)
		require.NotNil(t, m.Meta.Until)
		require.Equal(t, len(docs), len(m.Meta.Metas)) // 各メタ情報の先頭に-1が付く
		n0 := m.NodeRank(0)
		n1 := m.NodeRank(len(m.Words) - 1)
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

		t.Run("Test Get All", func(t *testing.T) {
			t.Parallel()
			var config Config
			config.GroupingFuncType = PickByKey
			config.ReduceThreshold = 0.001
			n, _, _ := NewMultiCoMatrixFromBuilder(ctx, b, config)
			require.Equal(t, len(docs), n)
		})
		t.Run("Test Get Group", func(t *testing.T) {
			t.Parallel()
			var config Config
			config.GroupingFuncType = PickByKey
			config.ReduceThreshold = 0.001
			config.AtGroupID = docs[0].Key
			n, _, _ := NewMultiCoMatrixFromBuilder(ctx, b, config)
			require.Equal(t, 1, n)
		})
	})
}

func generateDocs() []Document {
	ctx := context.Background()
	ndlConfig := ndl.Config{
		UseCache:    true,
		CreateCache: true,
	}
	c := ndl.NewClient(ndlConfig)
	// 会議APIから結果取得
	_, recordCh := c.GenerateNDLResultWithErrorHook(ctx, func(err error) {
		panic(err)
	})
	// 会議-単語
	docsCh := stream.FunIO[ndl.NDLRecode, Document](ctx, recordCh, func(record ndl.NDLRecode) Document {
		ar := analyzer.Analysis(record.Speeches)
		if ar.Error() != nil {
			panic(ar.Error())
		}
		var doc Document
		doc.Words = ar.Get(analyzer.Config{
			Includes: []analyzer.PartOfSpeechType{
				analyzer.Noun,
			},
		})
		doc.Key = record.IssueID
		doc.Name = fmt.Sprintf("%s %s", record.NameOfHouse, record.NameOfMeeting)
		doc.At = record.Date
		doc.Description = fmt.Sprintf("%s %s", record.NameOfHouse, record.NameOfMeeting)
		return doc
	})

	docs := []Document{}
	for doc := range docsCh {
		if len(doc.Words) > 0 {
			docs = append(docs, doc)
		}
	}
	return docs
}
