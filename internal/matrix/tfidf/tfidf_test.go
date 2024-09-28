package tfidf

import (
	"iter"
	"math"
	"slices"
	"testing"
	"yyyoichi/Collo-API/internal/matrix/base"

	"github.com/stretchr/testify/require"
)

type tdoc []int

func (d tdoc) IterWordId() iter.Seq[int] {
	return slices.Values(d)
}

func TestNewMatric(t *testing.T) {

	docs := map[string][]int{
		"a": {10, 20, 20, 30},
		"b": {30, 20, 40, 40, 50},
		"c": {30, 40, 50},
	}
	// words := map[string]int{
	// 	"orange": 10,
	// 	"banana": 20,
	// 	"cherry": 30,
	// 	"apple":  40,
	// 	"grape":  50,
	// }
	var builder base.BuildCounterMatrix
	for docId, wordIds := range docs {
		builder.Append(docId, tdoc(wordIds))
	}
	var countMatrix = builder.Build()
	var m = New(countMatrix)

	expTF := map[string]map[int]float64{
		"a": {
			10: 1.0 / 4.0,
			20: 2.0 / 4.0,
			30: 1.0 / 4.0,
			40: 0,
			50: 0,
		},
		"b": {
			10: 0,
			20: 1.0 / 5.0,
			30: 1.0 / 5.0,
			40: 2.0 / 5.0,
			50: 1.0 / 5.0,
		},
		"c": {
			10: 0,
			20: 0,
			30: 1.0 / 3.0,
			40: 0,
			50: 2.0 / 3.0,
		},
	}
	expIDF := map[int]float64{
		10: math.Log(3.0 / float64(1+1)),
		20: math.Log(3.0 / float64(2+1)),
		30: math.Log(3.0 / float64(3+1)),
		40: math.Log(3.0 / float64(2+1)),
		50: math.Log(3.0 / float64(2+1)),
	}

	for docId := range m.IterDocId() {
		for wordId := range m.IterWordId() {
			tfidf := m.N(docId, wordId)
			exp := expTF[docId][wordId] * expIDF[wordId]
			require.Equal(t, exp, *tfidf)
		}
	}
}
