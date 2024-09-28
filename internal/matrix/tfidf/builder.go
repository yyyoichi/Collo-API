package tfidf

import (
	"errors"
	"iter"
	"math"
)

// TFIDF計算機
type (
	Builder struct {
		// 文書識別子とそのtf
		docs map[string]tfdoc
		// 出現したユニークな単語識別子
		uniqWordIds map[int]struct{}
	}
	Document interface {
		// 文書に出現した単語の識別子を返す関数。
		// 出現毎に識別子を渡されることが期待される。
		IterWordId() iter.Seq[int]
	}
)

// 複数の文書[docs]を一つの文書とみなして、TF計算に追加します。
// [id]にはインスタンス中一意の文字列を指定してください。
func (b *Builder) Append(id string, docs ...Document) error {
	b.init()

	if _, found := b.docs[id]; found {
		return errors.New("id is found in this")
	}
	var tfdoc = tfdoc{
		nt: make(map[int]int),
		T:  0,
	}
	for _, doc := range docs {
		for id := range doc.IterWordId() {
			tfdoc.T++
			tfdoc.nt[id] += 1
			b.uniqWordIds[id] = struct{}{}
		}
	}

	b.docs[id] = tfdoc
	return nil
}

// Appendされたすべての文書からTFIDFを計算する。
func (b *Builder) Build() TFIDF {
	numX := len(b.uniqWordIds)
	numY := len(b.docs)
	var tfidf = TFIDF{
		xindex: make(map[int]int, numX),
		yindex: make(map[string]int, numY),
		data:   make([]float64, numY*numX),
	}
	for wordId := range b.uniqWordIds {
		xi := len(tfidf.xindex)
		tfidf.xindex[wordId] = xi
	}
	for docId := range b.docs {
		yi := len(tfidf.yindex)
		tfidf.yindex[docId] = yi
	}

	for docId, doc := range b.docs {
		yi := tfidf.yindex[docId]
		for wordId := range doc.nt {
			xi := tfidf.xindex[wordId]
			tfidf.data[yi*numY+xi] = calcTFIDF(doc.tf(wordId), b.idf(wordId))
		}
	}
	return tfidf
}

// 単語[wordId]のIDFを返します。
func (b *Builder) idf(wordId int) float64 {
	var N = len(b.docs)   // 文書数
	var DF = b.df(wordId) // 単語[wordId]の出現する単語数
	// 0除算のためプラス1
	return math.Log(float64(N) / float64(DF+1))
}

// 単語[wordId]が出現する文書数
func (b *Builder) df(wordId int) int {
	var c int
	for _, doc := range b.docs {
		if doc.n(wordId) > 0 {
			c++
		}
	}
	return c
}

func (b *Builder) init() {
	if b.docs == nil {
		b.docs = make(map[string]tfdoc)
	}
	if b.uniqWordIds == nil {
		b.uniqWordIds = make(map[int]struct{})
	}
}

// TFIDFを返す。
func calcTFIDF(tf, idf float64) float64 {
	// 0除算のためプラス1
	return tf * (idf + 1)
}
