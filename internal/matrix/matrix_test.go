package matrix

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatrix(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		tbuilder := NewMatrixBuilder()
		tbuilder.AppendDoc([]string{"hoge", "huga"})
		tbuilder.AppendDoc([]string{"hoge", "hoge", "foo"})

		require.Equal(t, tbuilder.indexByWord, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
		})

		matrix := tbuilder.toMatrix()
		require.Equal(t, matrix[0], []uint{1, 1, 0})
		require.Equal(t, matrix[1], []uint{2, 0, 1})

		tbuilder.AppendDoc([]string{"huga", "bar"})
		require.Equal(t, tbuilder.indexByWord, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
			"bar":  3,
		})
		matrix = tbuilder.toMatrix()
		require.Equal(t, matrix[0], []uint{1, 1, 0, 0})
		require.Equal(t, matrix[1], []uint{2, 0, 1, 0})
		require.Equal(t, matrix[2], []uint{0, 1, 0, 1})

		l := 100000
		str := make([]string, l)
		for i := 0; i < l; i++ {
			str[i] = "hoge"
		}
		tbuilder.AppendDoc(str)
		matrix = tbuilder.toMatrix()
		require.Equal(t, matrix[3][0], uint(l))
	})

	t.Run("sync add", func(t *testing.T) {
		tbuilder := NewMatrixBuilder()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			tbuilder.AppendDoc([]string{"hoge", "huga"})
		}()
		go func() {
			defer wg.Done()
			tbuilder.AppendDoc([]string{"hoge", "hoge"})
		}()
		wg.Wait()

		require.Equal(t, len(tbuilder.indexByWord), 2)
		matrix := tbuilder.toMatrix()

		ihoge := tbuilder.indexByWord["hoge"]
		ihuga := tbuilder.indexByWord["huga"]
		if matrix[0][ihoge] == uint(2) {
			require.EqualValues(t, matrix[0][ihuga], 0)
			require.EqualValues(t, matrix[1][ihoge], 1)
			require.EqualValues(t, matrix[1][ihuga], 1)

		} else if matrix[0][ihoge] == (1) {
			require.EqualValues(t, matrix[0][ihuga], 1)
			require.EqualValues(t, matrix[1][ihoge], 2)
			require.EqualValues(t, matrix[1][ihuga], 0)

		} else {
			t.Errorf("hoge is '%v'", matrix[0][ihoge])
		}
	})
}
