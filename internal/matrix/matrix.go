package matrix

import (
	"math"
	"sort"
	"sync"
)

// 単語文書行列
type DocWordMatrix struct {
	matrix [][]int
	words  []string
}

// すべての文書における2単語[windex1][windex2]の共起回数[f]と共起文書数[c]を返す
func (m *DocWordMatrix) CoOccurrencetFrequency(windex1, windex2 int) (f int, c int) {
	if len(m.words) <= windex1 {
		return 0, 0
	}
	if len(m.words) <= windex2 {
		return 0, 0
	}
	for _, doc := range m.matrix {
		// 共起頻度
		fq := doc[windex1] * doc[windex2]
		if fq > 0 {
			c++
		}
		f += fq
	}
	return f, c
}

// すべての文書内での単語[windex]の出現回数[o]と出現文書数[c]を返す
func (m *DocWordMatrix) Occurances(windex int) (o int, c int) {
	if len(m.words) <= windex {
		return 0, 0
	}
	for _, doc := range m.matrix {
		oc := doc[windex]
		if oc > 0 {
			c++
		}
		o += oc
	}
	return o, c
}

type TFIDFMatrix struct {
	// TF-IDFの生データ。
	flatMatrix []float64
	// 重要度順のflatMatrixの位置
	indices []int
	// 単語数
	lenWrods int
}

func NewTFIDFMatrix(tfidfMatrix [][]float64) *TFIDFMatrix {
	cap := len(tfidfMatrix) * len(tfidfMatrix[0])
	m := &TFIDFMatrix{
		flatMatrix: make([]float64, 0, cap),
		indices:    make([]int, cap),
		lenWrods:   len(tfidfMatrix[0]),
	}

	for _, row := range tfidfMatrix {
		m.flatMatrix = append(m.flatMatrix, row...)
	}

	for i := range m.flatMatrix {
		m.indices[i] = i
	}

	// m.flatMatrixを降順ソート
	sort.Slice(m.indices, func(i, j int) bool {
		return m.flatMatrix[m.indices[i]] > m.flatMatrix[m.indices[j]]
	})

	return m
}

// 返却個数を返す。単語数の[threshold]%切り上げか、[minWords]の大きいほうを返す。単語数がどちらよりも小さいときすべての単語数を返す。
func (m *TFIDFMatrix) cap(threshold float64, minWords int) int {
	// 単語数
	l := m.lenWrods
	// 上位[threshold]%
	upper := int(math.Ceil(float64(l) * threshold))
	// 返却個数
	cap := l
	if upper > minWords {
		cap = upper
	} else {
		cap = minWords
	}
	// 実際の単語数が小さいとき、すべての単語数を返す。
	if cap > l {
		return l
	} else {
		return cap
	}
}

// 上位[threshold]%の*重要度以上*の単語位置[windex]を返す。返却数は[minWords]を単語の実数が下回らない限り保証する。
func (m *TFIDFMatrix) TopPercentageWIndexes(threshold float64, minWords int) ColumnReduction {
	// 単語種数
	n := float64(m.lenWrods)
	// 返却個数
	cap := m.cap(threshold, minWords)
	windexes := make(map[int]interface{}, cap)

	// 重要度順[i]番目の単語位置を追加
	add := func(i int) {
		// 重要度[i]番目のflatMatrixの位置
		findex := float64(m.indices[i])
		// 単語位置[windex]化
		windex := math.Mod(findex, n)
		windexes[int(windex)] = struct{}{}
	}

	// 重要度順[i]の追加を判定、追加する再帰関数。
	var fn func(i int) int
	fn = func(i int) int {
		// しきい値までの個数を追加したか
		if len(windexes) == cap {
			// 最低重要度
			bottom := m.flatMatrix[m.indices[i-1]]
			j := i
			for {
				// 最低重要度と同じ値なら追加し続ける
				tfidf := m.flatMatrix[m.indices[j]]
				if bottom != tfidf {
					break
				}
				add(j)
				j++
			}
			return 0
		} else {
			add(i)
			return fn(i + 1)
		}
	}
	fn(0)

	wids := make([]int, 0, cap)
	for windex := range windexes {
		wids = append(wids, windex)
	}
	return ColumnReduction{
		windexes: wids,
	}
}

// 列縮小
type ColumnReduction struct {
	windexes []int

	done bool
}

// 縮小後サイズ
func (r ColumnReduction) Len() int { return len(r.windexes) }

func (r ColumnReduction) Reduce(m *DocWordMatrix) {
	if r.done {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		matrix := make([][]int, len(m.matrix))
		for dindex := range m.matrix {
			matrix[dindex] = make([]int, r.Len())
			for newWIndex, windex := range r.windexes {
				matrix[dindex][newWIndex] = m.matrix[dindex][windex]
			}
		}
		m.matrix = matrix
	}()

	go func() {
		defer wg.Done()
		words := make([]string, r.Len())
		for newIndex, windex := range r.windexes {
			words[newIndex] = m.words[windex]
		}
		m.words = words
	}()

	wg.Wait()
	r.done = true
}

// 縮小前のwindexを取得する
func (r ColumnReduction) Pre(newIndex int) int {
	return r.windexes[newIndex]
}
