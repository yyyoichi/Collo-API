package matrix

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

var tFruit = map[string]int{"apple": 1, "orange": 2, "banana": 3, "cherry": 4, "grape": 5}

type tTfidfDoc struct {
	words []string
}

func (d *tTfidfDoc) WordIDS() []int {
	var resp []int
	for _, w := range d.words {
		resp = append(resp, tFruit[w])
	}
	return resp
}

func TestTFIDF(t *testing.T) {
	t.Run("builder", func(t *testing.T) {
		docA1 := tTfidfDoc{
			words: []string{"apple", "orange"},
		}
		docA2 := tTfidfDoc{
			words: []string{"orange", "banana"},
		}
		docB1 := tTfidfDoc{
			words: []string{"banana", "orange"},
		}
		docB2 := tTfidfDoc{
			words: []string{"cherry", "cherry", "grape"},
		}
		docC := tTfidfDoc{
			words: []string{"banana", "grape", "grape"},
		}
		var b BuildTFIDF
		b.AddChunkedDocs("a", &docA1, &docA2)
		b.AddChunkedDocs("b", &docB1, &docB2)
		b.AddChunkedDocs("c", &docC)

		// tf 値
		require.EqualValues(t, b.TF("a", tFruit["apple"]), 1.0/4.0)
		require.EqualValues(t, b.TF("a", tFruit["orange"]), 2.0/4.0)
		require.EqualValues(t, b.TF("a", tFruit["banana"]), 1.0/4.0)
		require.EqualValues(t, b.TF("a", tFruit["cherry"]), 0)
		require.EqualValues(t, b.TF("a", tFruit["grape"]), 0)

		require.EqualValues(t, b.TF("b", tFruit["apple"]), 0)
		require.EqualValues(t, b.TF("b", tFruit["orange"]), 1.0/5.0)
		require.EqualValues(t, b.TF("b", tFruit["banana"]), 1.0/5.0)
		require.EqualValues(t, b.TF("b", tFruit["cherry"]), 2.0/5.0)
		require.EqualValues(t, b.TF("b", tFruit["grape"]), 1.0/5.0)

		require.EqualValues(t, b.TF("c", tFruit["apple"]), 0)
		require.EqualValues(t, b.TF("c", tFruit["orange"]), 0)
		require.EqualValues(t, b.TF("c", tFruit["banana"]), 1.0/3.0)
		require.EqualValues(t, b.TF("c", tFruit["cherry"]), 0)
		require.EqualValues(t, b.TF("c", tFruit["grape"]), 2.0/3.0)

		require.EqualValues(t, b.IDF(tFruit["apple"]), math.Log(3.0/float64(1+1)))
		require.EqualValues(t, b.IDF(tFruit["orange"]), math.Log(3.0/float64(2+1)))
		require.EqualValues(t, b.IDF(tFruit["banana"]), math.Log(3.0/float64(3+1)))
		require.EqualValues(t, b.IDF(tFruit["cherry"]), math.Log(3.0/float64(1+1)))
		require.EqualValues(t, b.IDF(tFruit["grape"]), math.Log(3.0/float64(2+1)))

		tfidf := b.Build()
		for i := 0; i < len(tfidf.indices)-1; i++ {
			current := tfidf.maxTFIDF[tfidf.indices[i]]
			next := tfidf.maxTFIDF[tfidf.indices[i+1]]
			require.True(t, current >= next)
		}
		require.EqualValues(t, []int{5, 4, 2, 1, 3}, tfidf.indices)
	})

	t.Run("tfidf cap", func(t *testing.T) {
		test := []struct {
			num    int     // 所持数
			th     float64 // 上位
			min    int     // 最小
			expcap int     // 期待値
		}{
			{10, 0.1, 3, 3},   // 最小数を保証してほしい
			{10, 0.1, 15, 10}, // 最小数は保証できないので、すべて
			{10, 0.5, 3, 5},   // 50%
			{10, 0.4, 3, 4},   // 40%
		}
		for _, tt := range test {
			m := TFIDF{
				maxTFIDF: make(map[int]float64, tt.num),
			}
			for i := 0; i < tt.num; i++ {
				m.maxTFIDF[i] = 1.0
			}
			cap := m.cap(tt.th, tt.min)
			require.Equal(t, tt.expcap, cap)
		}
	})

	t.Run("tfidf top", func(t *testing.T) {
		m := TFIDF{
			maxTFIDF: map[int]float64{
				1: 1.1,
				2: 2.2,
				3: 2.2,
				4: 2.2,
				5: 5.5,
				6: 6.6,
			},
			indices: []int{6, 5, 4, 3, 2, 1},
		}
		test := []struct {
			cap    int
			explen int
		}{
			{0, 0}, // 追加しない
			{1, 1},
			{3, 5}, // 6.6, 5.5, 2.2, 2.2, 2.2, 2.2
			{6, 6},
			{7, 6},
		}
		for _, tt := range test {
			words := m.getTopWords(tt.cap)
			require.Equal(t, tt.explen, len(words))
		}
	})
}
