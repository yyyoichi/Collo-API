package tfidf

import (
	"iter"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type tdoc []int

func (d tdoc) IterWordId() iter.Seq[int] {
	return slices.Values(d)
}

func TestTFIDF(t *testing.T) {
	docA1 := tdoc{10, 20}
	docA2 := tdoc{20, 30}
	docB1 := tdoc{30, 20}
	docB2 := tdoc{40, 40, 50}
	docC := tdoc{30, 50, 50}

	var b Builder
	b.Append("a", &docA1, &docA2)
	b.Append("b", &docB1, &docB2)
	b.Append("c", &docC)

	// tf 値
	require.EqualValues(t, b.docs["a"].tf(10), 1.0/4.0)
	require.EqualValues(t, b.docs["a"].tf(20), 2.0/4.0)
	require.EqualValues(t, b.docs["a"].tf(30), 1.0/4.0)
	require.EqualValues(t, b.docs["a"].tf(40), 0)
	require.EqualValues(t, b.docs["a"].tf(50), 0)

	require.EqualValues(t, b.docs["b"].tf(10), 0)
	require.EqualValues(t, b.docs["b"].tf(20), 1.0/5.0)
	require.EqualValues(t, b.docs["b"].tf(30), 1.0/5.0)
	require.EqualValues(t, b.docs["b"].tf(40), 2.0/5.0)
	require.EqualValues(t, b.docs["b"].tf(50), 1.0/5.0)

	require.EqualValues(t, b.docs["c"].tf(10), 0)
	require.EqualValues(t, b.docs["c"].tf(20), 0)
	require.EqualValues(t, b.docs["c"].tf(30), 1.0/3.0)
	require.EqualValues(t, b.docs["c"].tf(40), 0)
	require.EqualValues(t, b.docs["c"].tf(50), 2.0/3.0)

	require.EqualValues(t, b.idf(10), math.Log(3.0/float64(1+1)))
	require.EqualValues(t, b.idf(20), math.Log(3.0/float64(2+1)))
	require.EqualValues(t, b.idf(30), math.Log(3.0/float64(3+1)))
	require.EqualValues(t, b.idf(40), math.Log(3.0/float64(1+1)))
	require.EqualValues(t, b.idf(50), math.Log(3.0/float64(2+1)))

	result := b.Build()
	require.Equal(t, 3, len(result.yindex)) // a,b,c
	require.Equal(t, 5, len(result.xindex)) // 10,20,30,40,50
	require.Equal(t, 3, len(result.data))
	for _, d := range result.data {
		require.Equal(t, 5, len(d))
	}

	tf10a := 1.0 / 4.0
	idf10 := math.Log(3.0 / float64(1+1))
	yi := result.yindex["a"]
	xi := result.xindex[10]
	require.EqualValues(t, calcTFIDF(tf10a, idf10), result.data[yi][xi])
}
