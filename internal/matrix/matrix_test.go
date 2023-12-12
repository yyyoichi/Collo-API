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
		var windexes []int

		// cap ciel(4 * 0.1) -> 1
		windexes = m.TopPercentageWIndexes(0.1)
		sort.Ints(windexes)
		require.EqualValues(t, []int{0, 2}, windexes)
		// cap ciel(4 * 0.6) -> 3
		windexes = m.TopPercentageWIndexes(0.6)
		sort.Ints(windexes)
		require.EqualValues(t, []int{0, 2, 3}, windexes)
	})
}
