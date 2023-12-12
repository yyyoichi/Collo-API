package matrix

import (
	"math"
	"sort"
)

// 単語文書行列
type DocWordMatrix struct {
	matrix [][]float64
	// words  []string
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

// 上位[threshold]%の*重要度以上*の単語位置[windex]を返す
func (m *TFIDFMatrix) TopPercentageWIndexes(threshold float64) []int {
	// 単語種数
	n := float64(m.lenWrods)
	// 返却個数
	cap := int(math.Ceil(n * threshold))
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
	return wids
}
