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

		require.Equal(t, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
		}, tbuilder.indexByWord)

		matrix := tbuilder.BuildCountDocMatrix().(*CountDocMatrix).docs
		require.EqualValues(t, []int{1, 1, 0}, matrix[0].row)
		require.EqualValues(t, []int{2, 0, 1}, matrix[1].row)

		tbuilder.AppendDoc([]string{"huga", "bar"})
		require.Equal(t, map[string]int{
			"hoge": 0,
			"huga": 1,
			"foo":  2,
			"bar":  3,
		}, tbuilder.indexByWord)
		matrix = tbuilder.BuildCountDocMatrix().(*CountDocMatrix).docs
		require.Equal(t, []int{1, 1, 0, 0}, matrix[0].row)
		require.Equal(t, []int{2, 0, 1, 0}, matrix[1].row)
		require.Equal(t, []int{0, 1, 0, 1}, matrix[2].row)

		l := 100000
		str := make([]string, l)
		for i := 0; i < l; i++ {
			str[i] = "hoge"
		}
		tbuilder.AppendDoc(str)
		matrix = tbuilder.BuildCountDocMatrix().(*CountDocMatrix).docs
		require.EqualValues(t, l, matrix[3].row[0])
		require.EqualValues(t, l, matrix[3].wordsCount)
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
		matrix := tbuilder.BuildCountDocMatrix().(*CountDocMatrix).docs

		ihoge := tbuilder.indexByWord["hoge"]
		ihuga := tbuilder.indexByWord["huga"]
		if matrix[0].row[ihoge] == 2 {
			require.EqualValues(t, 0, matrix[0].row[ihuga])
			require.EqualValues(t, 1, matrix[1].row[ihoge])
			require.EqualValues(t, 1, matrix[1].row[ihuga])

		} else if matrix[0].row[ihoge] == 1 {
			require.EqualValues(t, 1, matrix[0].row[ihuga])
			require.EqualValues(t, 2, matrix[1].row[ihoge])
			require.EqualValues(t, 0, matrix[1].row[ihuga])

		} else {
			t.Errorf("hoge is '%v'", matrix[0].row[ihoge])
		}
	})

	t.Run("DocMatrix IDF", func(t *testing.T) {
		tmatrix := NewCountDocMatrix(
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
		tmatrix := NewCountDocMatrix(
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
		doc1 := &countDoc{
			row:        []int{1, 2},
			wordsCount: 1,
		}
		// exp tfs = {2,3}
		doc2 := &countDoc{
			row:        []int{2, 3},
			wordsCount: 1,
		}
		tmatrix := &CountDocMatrix{
			docMatrixBase: &docMatrixBase[*countDoc]{
				words: []string{
					"hoge",
					"huga",
				},
				docs: []*countDoc{doc1, doc2},
			},
			idfStore: []float64{2, 1},
		}

		exp := [][]float64{
			{1 * 2, 2 * 1}, {2 * 2, 3 * 1},
		}

		ttfidfmatrix := NewTFIDFDocMatrix(tmatrix)
		require.Equal(t, exp[0], ttfidfmatrix.docs[0])
		require.Equal(t, exp[1], ttfidfmatrix.docs[1])
	})
}
