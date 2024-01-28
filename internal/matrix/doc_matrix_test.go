package matrix

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocWordMatrix(t *testing.T) {
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
