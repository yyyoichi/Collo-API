package matrix

import (
	"context"
	"sync"
	"yyyoichi/Collo-API/pkg/stream"
)

type (
	Builder struct {
		indexByWord       map[string]int
		appnededDocuments []AppendedDocument

		mu *sync.Mutex
	}

	// inmpelement TFIDFDocumentInterface
	AppendedDocument struct {
		parent *Builder
		DocumentMeta
		words []string
	}
)

// ドキュメントに含まれる単語を、builderにおけるindexByWordのindexで返す。
// IDが見つからない場合は、スルー
func (d *AppendedDocument) WordIDS() []int {
	var result []int
	for _, word := range d.words {
		if i, found := d.parent.indexByWord[word]; !found {
			continue
		} else {
			result = append(result, i)
		}
	}
	return result
}

func NewBuilder() Builder {
	var mu sync.Mutex
	return Builder{
		indexByWord:       map[string]int{},
		appnededDocuments: []AppendedDocument{},
		mu:                &mu,
	}
}

func (b *Builder) Append(meta DocumentMeta, words []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var doc AppendedDocument
	for _, word := range words {
		if _, found := b.indexByWord[word]; !found {
			b.indexByWord[word] = len(b.indexByWord)
		}
	}
	doc.words = words
	doc.parent = b
	doc.DocumentMeta = meta
	b.appnededDocuments = append(b.appnededDocuments, doc)
}

func (b *Builder) BuildTFIDF(t GroupingFuncType) TFIDF {
	// グループごとにドキュメントを分ける
	var groups = make(map[string][]TFIDFDocumentInterface)
	pick := b.pickGroupIDFunc(t)
	for _, doc := range b.appnededDocuments {
		id := pick(doc.DocumentMeta)
		if _, found := groups[id]; !found {
			groups[id] = []TFIDFDocumentInterface{&doc}
		} else {
			groups[id] = append(groups[id], &doc)
		}
	}

	var tb BuildTFIDF
	for groupID, docs := range groups {
		tb.AddChunkedDocs(groupID, docs...)
	}
	return tb.Build()
}

func (b *Builder) pickGroupIDFunc(t GroupingFuncType) func(DocumentMeta) string {
	switch t {
	case PickByMonth:
		return func(d DocumentMeta) string {
			return d.At.Format("2006-01")
		}
	case PickAsTotal:
		return func(DocumentMeta) string {
			return "total"
		}
	default: // PickByKey
		return func(d DocumentMeta) string { return d.Key }
	}
}

// appendされたすべての文書について[name]という名前で文書単語行列を返す。
func (b *Builder) BuildDocWordMatrix(ctx context.Context, name string) DocWordMatrix {
	return b.buildDocWordMatrix(name, b.appnededDocuments)
}

// グルーピング関数によってすべての文書を分割して文書単語行列を返す。
func (b *Builder) BuildDocWordMatrixByGroup(ctx context.Context, t GroupingFuncType) (int, <-chan DocWordMatrix) {
	groups := b.getByGroup(t)
	return len(groups), stream.GeneratorWithMapStringKey[[]AppendedDocument, DocWordMatrix](ctx, groups, b.buildDocWordMatrix)
}

// 特定のグループについてのみ文書単語行列を返す。
func (b *Builder) BuildDocWordMatrixByGroupAt(ctx context.Context, t GroupingFuncType, at string) (int, <-chan DocWordMatrix) {
	groups := b.getByGroup(t)
	if _, found := groups[at]; !found {
		return 0, nil
	}
	var omit = make(map[string][]AppendedDocument, 1)
	omit[at] = groups[at]
	return len(omit), stream.GeneratorWithMapStringKey[[]AppendedDocument, DocWordMatrix](ctx, omit, b.buildDocWordMatrix)
}

func (b *Builder) Words() []string {
	words := make([]string, len(b.indexByWord))
	for rawWord, word := range b.indexByWord {
		words[word] = rawWord
	}
	return words
}

func (b *Builder) getByGroup(t GroupingFuncType) map[string][]AppendedDocument {
	var groups = make(map[string][]AppendedDocument)
	pick := b.pickGroupIDFunc(t)
	for _, doc := range b.appnededDocuments {
		id := pick(doc.DocumentMeta)
		if _, found := groups[id]; !found {
			groups[id] = []AppendedDocument{doc}
		} else {
			groups[id] = append(groups[id], doc)
		}
	}
	return groups
}

func (b *Builder) buildDocWordMatrix(id string, docs []AppendedDocument) DocWordMatrix {
	var m DocWordMatrix
	// 単語ごとの出現回数
	m.matrix = make([][]int, len(docs))
	var metas []DocumentMeta = []DocumentMeta{}

	for docIndex, doc := range docs {
		var mx = make([]int, len(b.indexByWord)) // 1文書における単語行列
		for _, word := range doc.WordIDS() {
			mx[word] += 1
		}
		m.matrix[docIndex] = mx
		metas = append(metas, doc.DocumentMeta)
	}
	m.meta = NewMultiDocMeta(id, metas)
	m.wordCount = len(b.indexByWord)
	return m
}

type ColumnReduction struct {
	words []int // 残したい単語
	done  bool
}

func (r *ColumnReduction) found(word int) bool {
	for _, w := range r.words {
		if w == word {
			return true
		}
	}
	return false
}

func (r *ColumnReduction) Reduce(b *Builder) {
	if r.done {
		return
	}
	// 削除した分、新しくふりなおす
	var count int
	for rawWord, word := range b.indexByWord {
		if r.found(word) {
			b.indexByWord[rawWord] = count
			count++
		} else {
			delete(b.indexByWord, rawWord)
		}
	}
	r.done = true
}
