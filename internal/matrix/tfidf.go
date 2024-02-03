package matrix

import (
	"errors"
	"math"
	"sort"
)

type TFIDF struct {
	// 単語識別子とその文書全体でのTFIDF最大値
	maxTFIDF map[int]float64
	// 大きい順にソートされた単語識別子
	indices []int
}

// 上位[threshold]%の*重要度以上*の単語を持つColumnReductionを返す。返却数は[min]を単語の実数が下回らない限り保証する。
func (m *TFIDF) GetColumnReduction(threshold float64, min, max int) ColumnReduction {
	cap := m.cap(threshold, min, max)
	return ColumnReduction{
		words: m.getTopWords(cap),
	}
}

// 重要度順に[cap]コの単語を返す。重要度が等しいとき[cap]よりも多い数の単語を返す。
func (m *TFIDF) getTopWords(cap int) []int {
	var result []int
	if cap <= 0 {
		return result
	}
	if cap > len(m.maxTFIDF) {
		// 単語数のが少ない
		cap = len(m.maxTFIDF)
	}
	// 重要度順[i]の追加を判定、追加する再帰関数。
	var fn func(i int)
	fn = func(i int) {
		// しきい値までの個数を追加したか
		if len(result) < cap {
			result = append(result, m.indices[i])
			fn(i + 1)
			return
		}
		// capと同じ数だけ追加完了
		// 追加された最低重要度
		bottomScore := m.maxTFIDF[m.indices[i-1]]
		j := i
		for {
			if j >= len(m.maxTFIDF) {
				// 探索上限
				break
			}
			// 最低重要度と同じ値なら追加し続ける
			if bottomScore == m.maxTFIDF[m.indices[j]] {
				result = append(result, m.indices[j])
				j++
				continue
			}
			return
		}
	}
	fn(0)
	return result
}

// 返却個数を返す。単語数の[threshold]%切り上げか、[max]の小さいほうを返す。[min]数は保証する
func (m *TFIDF) cap(threshold float64, min, max int) int {
	if min > len(m.maxTFIDF) {
		return len(m.maxTFIDF)
	}
	// 以下、最小値は長さより大きいことは保証
	if max > len(m.maxTFIDF) {
		max = len(m.maxTFIDF) // 最大値補正
	}
	// 上位[threshold]%の単語数
	upperCount := int(math.Ceil(float64(len(m.maxTFIDF)) * threshold))
	if min > upperCount { // 最小数の保証
		return min
	}

	if max < upperCount {
		return max
	}
	return upperCount
}

// TFIDF計算機
type (
	BuildTFIDF struct {
		docs map[string]tfdoc
	}
	TFIDFDocumentInterface interface {
		WordIDS() []int // 文書に含まれる単語の識別子を返す関数
	}
)

// 複数の文書[docs]を一つの文書とみなして、TF計算に追加します。
// [id]にはインスタンス中一意の文字列を指定してください。
func (b *BuildTFIDF) AddChunkedDocs(id string, docs ...TFIDFDocumentInterface) error {
	if b.docs == nil {
		b.docs = make(map[string]tfdoc)
	}

	if _, found := b.docs[id]; found {
		return errors.New("id is found in this")
	}
	var tfdoc tfdoc
	tfdoc.nt = make(map[int]int)
	for _, d := range docs {
		words := d.WordIDS()
		tfdoc.T += len(words) // 単語数を追加
		for _, word := range words {
			// 単語の出現回数を追加
			tfdoc.nt[word] += 1
		}
	}

	b.docs[id] = tfdoc

	return nil
}

func (b *BuildTFIDF) Build() TFIDF {
	// 各単語ごとにTFIDFの最大値を保持する
	var maxTFIDF = map[int]float64{}
	for id, doc := range b.docs {
		for word := range doc.nt {
			tfidf := b.TFIDF(id, word)
			src := maxTFIDF[word]
			if src < tfidf {
				// 新しく取得したほうのが大きい
				maxTFIDF[word] = tfidf
			}
		}
	}
	// 大きい順にソートされた単語識別子
	var sortedWordIndices = make([]int, len(maxTFIDF))
	i := 0
	for word := range maxTFIDF {
		sortedWordIndices[i] = word
		i++
	}
	sort.Slice(sortedWordIndices, func(i, j int) bool {
		wordi := sortedWordIndices[i]
		wordj := sortedWordIndices[j]
		return maxTFIDF[wordi] > maxTFIDF[wordj]
	})
	return TFIDF{
		maxTFIDF: maxTFIDF,
		indices:  sortedWordIndices,
	}
}

// 文書[id]における単語[word]のTFIDFを返します。
func (b *BuildTFIDF) TFIDF(id string, word int) float64 {
	// 0除算のためプラス1
	return b.TF(id, word) * (b.IDF(word) + 1)
}

// 文書[id]における単語[word]のTFを返します。
func (b *BuildTFIDF) TF(id string, word int) float64 {
	tfdoc, found := b.docs[id]
	if !found {
		return 0.0
	}
	return tfdoc.tf(word)
}

// 単語[word]のIDFを返します。
func (b *BuildTFIDF) IDF(word int) float64 {
	var N = len(b.docs) // 文書数
	var DF = b.DF(word) // 単語[word]の出現する単語数
	// 0除算のためプラス1
	return math.Log(float64(N) / float64(DF+1))
}

// 単語[word]が出現する文書数
func (b *BuildTFIDF) DF(word int) int {
	var c int
	for _, doc := range b.docs {
		if doc.n(word) > 0 {
			c++
		}
	}
	return c
}

// 文書におけるtf計算機
type tfdoc struct {
	nt map[int]int // 単語t出現回数。キーは単語識別子
	T  int         // 単語数count
}

// 単語[word]の出現回数
func (d *tfdoc) n(word int) int {
	return d.nt[word]
}

// 単語[word]の出現頻度
func (d *tfdoc) tf(word int) float64 {
	n, found := d.nt[word]
	if !found {
		return 0.0
	}
	// 出現回数n / 総単語数T
	return float64(n) / float64(d.T)
}
