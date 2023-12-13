package matrix

import (
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
}
