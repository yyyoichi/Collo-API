package matrix

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatrix(t *testing.T) {
	t.Run("Builder add", func(t *testing.T) {
		tbuilder := NewMatrixBuilder()
		tbuilder.AppendDoc([]string{"hoge", "huga"})
		tbuilder.AppendDoc([]string{"hoge", "hoge", "foo"})

		require.Equal(t, tbuilder.indexByWord, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
		})

		matrix := tbuilder.BuildDocMatrix().docs
		require.EqualValues(t, matrix[0].row, []float64{1, 1, 0})
		require.EqualValues(t, matrix[1].row, []float64{2, 0, 1})

		tbuilder.AppendDoc([]string{"huga", "bar"})
		require.Equal(t, tbuilder.indexByWord, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
			"bar":  3,
		})
		matrix = tbuilder.BuildDocMatrix().docs
		require.Equal(t, matrix[0].row, []float64{1, 1, 0, 0})
		require.Equal(t, matrix[1].row, []float64{2, 0, 1, 0})
		require.Equal(t, matrix[2].row, []float64{0, 1, 0, 1})

		l := 100000
		str := make([]string, l)
		for i := 0; i < l; i++ {
			str[i] = "hoge"
		}
		tbuilder.AppendDoc(str)
		matrix = tbuilder.BuildDocMatrix().docs
		require.EqualValues(t, matrix[3].row[0], l)
		require.EqualValues(t, matrix[3].wordsCount, l)
	})

	t.Run("Builder sync add", func(t *testing.T) {
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
		matrix := tbuilder.BuildDocMatrix().docs

		ihoge := tbuilder.indexByWord["hoge"]
		ihuga := tbuilder.indexByWord["huga"]
		if matrix[0].row[ihoge] == float64(2) {
			require.EqualValues(t, matrix[0].row[ihuga], 0)
			require.EqualValues(t, matrix[1].row[ihoge], 1)
			require.EqualValues(t, matrix[1].row[ihuga], 1)

		} else if matrix[0].row[ihoge] == float64(1) {
			require.EqualValues(t, matrix[0].row[ihuga], 1)
			require.EqualValues(t, matrix[1].row[ihoge], 2)
			require.EqualValues(t, matrix[1].row[ihuga], 0)

		} else {
			t.Errorf("hoge is '%v'", matrix[0].row[ihoge])
		}
	})

	t.Run("DocMatrix IDF", func(t *testing.T) {
		tmatrix := NewDocMatrix(
			map[string]int{
				"hoge": 0,
				"huga": 1,
				"foo":  2,
				"bar":  3,
			},
			[][]string{
				{"hoge", "hoge"},
				{"hoge", "foo"},
				{"huga", "foo"},
				{"foo", "foo"},
			},
		)
		totalDocs := 4.0
		require.EqualValues(t, math.Log(totalDocs/2), tmatrix.getIDFAt(0))
		require.EqualValues(t, math.Log(totalDocs/1), tmatrix.getIDFAt(1))
		require.EqualValues(t, math.Log(totalDocs/3), tmatrix.getIDFAt(2))
		require.EqualValues(t, 0, tmatrix.getIDFAt(3))
	})
	t.Run("DocMatrix TF", func(t *testing.T) {
		tmatrix := NewDocMatrix(
			map[string]int{
				"hoge": 0,
				"huga": 1,
			},
			[][]string{
				{"hoge", "hoge", "huga"},
			},
		)
		require.EqualValues(t, 2.0/3.0, tmatrix.docs[0].tfAt(0))
		require.EqualValues(t, 1.0/3.0, tmatrix.docs[0].tfAt(1))
	})

	t.Run("DocMatrix", func(t *testing.T) {
		// exp tfs = {1,2}
		doc1 := &doc{
			row:        []float64{1, 2},
			wordsCount: 1,
		}
		// exp tfs = {2,3}
		doc2 := &doc{
			row:        []float64{2, 3},
			wordsCount: 1,
		}
		tmatrix := &DocMatrix{
			words: []string{
				"hoge",
				"huga",
			},
			docs:     []*doc{doc1, doc2},
			idfStore: []float64{2, 1},
		}

		exp := [][]float64{
			{1 * 2, 2 * 1}, {2 * 2, 3 * 1},
		}

		tmatrix.replaceWeight()
		require.Equal(t, exp[0], tmatrix.docs[0].row)
		require.Equal(t, exp[1], tmatrix.docs[1].row)
	})
}
