package tfidf

// 文書におけるtf計算機
type tfdoc struct {
	nt map[int]int // 単語t出現回数。キーは単語識別子
	T  int         // 単語数count
}

// 単語[word]の出現回数
func (d tfdoc) n(wordId int) int {
	return d.nt[wordId]
}

// 単語[word]の出現頻度
func (d tfdoc) tf(wordId int) float64 {
	n, found := d.nt[wordId]
	if !found {
		return 0.0
	}
	// 出現回数n / 総単語数T
	return float64(n) / float64(d.T)
}
