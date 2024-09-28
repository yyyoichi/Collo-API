package tfidf

import (
	"math"
	"yyyoichi/Collo-API/internal/matrix/base"
)

type Matrix struct {
	base.Matrix
}

func New(countMatrix base.Matrix) Matrix {
	// 総文書数
	n := countMatrix.CountDoc()
	wordId_Idf := make(map[int]float64)
	for wordId := range countMatrix.IterWordId() {
		// 単語が出現する文書数
		df := countMatrix.CountDocWithWord(wordId)
		// 0除算のためプラス1
		wordId_Idf[wordId] = math.Log(n / (df + 1))
	}

	var m = countMatrix.Copy()
	for docId := range countMatrix.IterDocId() {
		// 文書における出現単語数
		t := countMatrix.SumWordInDoc(docId)
		for wordId := range countMatrix.IterWordId() {
			// 出現回数n / 総単語数T
			tf := *countMatrix.N(docId, wordId) / t
			idf := wordId_Idf[wordId]
			m.Set(docId, wordId, tf*idf)
		}
	}
	return Matrix{Matrix: m}
}
