package base

import (
	"errors"
	"iter"
	"maps"
)

type (
	// 文書-単語カウンター
	BuildCounterMatrix struct {
		// 文書識別子とそのカウンター
		docs map[string]doc
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
func (b *BuildCounterMatrix) Append(id string, docs ...Document) error {
	b.init()

	if _, found := b.docs[id]; found {
		return errors.New("id is found in this")
	}
	var doc = doc{
		nt: make(map[int]int),
		T:  0,
	}
	for _, src := range docs {
		for id := range src.IterWordId() {
			doc.T++
			doc.nt[id] += 1
			b.uniqWordIds[id] = struct{}{}
		}
	}

	b.docs[id] = doc
	return nil
}

// Appendされたすべての文書から文書単語カウント行列を作成する。
func (b *BuildCounterMatrix) Build() Matrix {
	var m = New(len(b.docs), len(b.uniqWordIds), maps.Keys(b.docs), maps.Keys(b.uniqWordIds))

	for docId, doc := range b.docs {
		for wordId := range doc.nt {
			m.Set(docId, wordId, float64(doc.n(wordId)))
		}
	}
	return m
}

func (b *BuildCounterMatrix) init() {
	if b.docs == nil {
		b.docs = make(map[string]doc)
	}
	if b.uniqWordIds == nil {
		b.uniqWordIds = make(map[int]struct{})
	}
}

// 文書における単語カウンター
type doc struct {
	nt map[int]int // 単語t出現回数。キーは単語識別子
	T  int         // 単語数count
}

// 単語[wordId]の出現回数
func (d doc) n(wordId int) int {
	return d.nt[wordId]
}
