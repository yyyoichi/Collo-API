package matrix

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTFIDFMatrix(t *testing.T) {
	t.Run("TFIDF TopPercentage", func(t *testing.T) {
		matrix := [][]float64{
			{1, 2, 4, 0},
			{1, 3, 5, 0},
			{1, 3, 6, 0},
			{6, 2, 4, 5},
		}
		m := NewTFIDFMatrix(matrix)
		// cap ciel(4 * 0.1) -> 1
		cap := m.cap(0.1, 3)
		require.Equal(t, 3, cap)
		// cap ciel(4 * 0.7) -> 3
		cap = m.cap(0.7, 2)
		require.Equal(t, 3, cap)
		// cap cail(4 * 0.1) -> 1
		cap = m.cap(0.1, 4)
		require.Equal(t, 4, cap)

		var windexes []int

		// cap ciel(4 * 0.1) -> 1
		windexes = m.TopPercentageWIndexes(0.1, 1).windexes
		sort.Ints(windexes)
		require.EqualValues(t, []int{0, 2}, windexes)
		// cap ciel(4 * 0.6) -> 3
		windexes = m.TopPercentageWIndexes(0.6, 1).windexes
		sort.Ints(windexes)
		require.EqualValues(t, []int{0, 2, 3}, windexes)
	})

	t.Run("Reduce DocWordMatrix", func(t *testing.T) {
		docmatrix := DocWordMatrix{
			words: []string{"hoge", "fuga", "foo"},
			matrix: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
		}
		c := ColumnReduction{
			windexes: []int{2, 0},
		}
		c.Reduce(&docmatrix)

		require.Equal(t, []string{"foo", "hoge"}, docmatrix.words)
		require.Equal(t, []int{2, 0}, docmatrix.matrix[0])
		require.Equal(t, []int{5, 3}, docmatrix.matrix[1])
		require.Equal(t, []int{8, 6}, docmatrix.matrix[2])
	})

	t.Run("CoOccurrencetFrequency", func(t *testing.T) {
		docmatrix := &DocWordMatrix{
			words: []string{"hoge", "fuga", "foo"},
			matrix: [][]int{
				{0, 1, 1},
				{1, 1, 2},
				{2, 0, 3},
			},
		}
		f := docmatrix.CoOccurrencetFrequency(0, 1)
		require.Equal(t, 1, f.Frequency)
		require.Equal(t, 1, f.Count)

		f = docmatrix.CoOccurrencetFrequency(0, 2)
		require.Equal(t, 8, f.Frequency)
		require.Equal(t, 2, f.Count)

		o := docmatrix.Occurances(2)
		require.Equal(t, 6, o.Occurances)
		require.Equal(t, 3, o.Count)

		count := 0
		for range docmatrix.generateCoIndex(context.Background()) {
			count++
		}
		require.Equal(t, 3, count)
	})
}
