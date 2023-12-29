package matrix

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"yyyoichi/Collo-API/pkg/stream"
)

type Builder struct {
	indexByWord map[string]int
	documents   []*Document

	mu sync.Mutex
}

func NewBuilder() *Builder {
	return &Builder{
		indexByWord: map[string]int{},
		documents:   []*Document{},
	}
}

func (b *Builder) AppendDocument(doc *Document) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, word := range doc.Words {
		if _, found := b.indexByWord[word]; !found {
			b.indexByWord[word] = len(b.indexByWord)
		}
	}
	b.documents = append(b.documents, doc)
}

// 文書ごとの単語出現回数の行列を生成する
func (b *Builder) Build() (*DocWordMatrix, *TFIDFMatrix) {
	matrix := b.build(b.documents)
	return b.build(b.documents), b.buildTFIDF(matrix.matrix)
}

// グループ数とグルーピングした単語文書行列を返す。
func (b *Builder) BuildByGroup(ctx context.Context, pickID PickDocGroupID) (int, <-chan *DocWordMatrix) {
	// groupIDをキーにしたドキュメントリスト
	group := map[string][]*Document{}
	for _, doc := range b.documents {
		id := pickID(doc)
		if _, found := group[id]; !found {
			group[id] = []*Document{doc}
		} else {
			group[id] = append(group[id], doc)
		}
	}
	docsCh := stream.GeneratorWithMapStringKey[[]*Document, []*Document](ctx, group, func(_ string, v []*Document) []*Document { return v })
	matrixCh := stream.FunIO[[]*Document, *DocWordMatrix](ctx, docsCh, b.build)
	return len(group), matrixCh
}

func (b *Builder) build(documents []*Document) *DocWordMatrix {
	lenWords := len(b.indexByWord)
	// matrix文書単語行列
	matrix := make([][]int, len(documents))
	// meta情報列
	metas := make([]*DocMeta, len(documents))
	for dindex, doc := range documents {
		// matrix
		matrix[dindex] = make([]int, lenWords)
		for _, word := range doc.Words {
			windex := b.indexByWord[word]
			matrix[dindex][windex] += 1.0
		}
		// meta
		metas[dindex] = doc.pickMeta()
	}
	words := make([]string, len(b.indexByWord))
	for word, windex := range b.indexByWord {
		words[windex] = word
	}
	return &DocWordMatrix{
		matrix: matrix,
		words:  words,
		metas:  metas,
	}
}

func (b *Builder) buildTFIDF(matrix [][]int) *TFIDFMatrix {
	tf := b.createTF(matrix)
	idf := b.createIDF(matrix)
	TF := func(dindex, windex int) float64 {
		return tf[dindex][windex]
	}
	IDF := func(windex int) float64 {
		return idf[windex]
	}

	lenWords := len(matrix[0])
	tfidfMatrix := make([][]float64, len(matrix))
	for dindex := range matrix {
		tfidfMatrix[dindex] = make([]float64, lenWords)
		for windex := range matrix[0] {
			tfidfMatrix[dindex][windex] = TF(dindex, windex) * (IDF(windex) + 1)
		}
	}
	return NewTFIDFMatrix(tfidfMatrix)
}

// b.buildされた行列[matrix]を受け取り、対応するTFを計算する
func (b *Builder) createTF(matrix [][]int) [][]float64 {
	lenWords := len(matrix[0])

	tfmatrix := make([][]float64, len(matrix))
	for dindex, d := range matrix {
		tfmatrix[dindex] = make([]float64, lenWords)
		// t = 文書d内の合計単語数
		T := float64(len(b.documents[dindex].Words))
		// n = 文書d内の単語tの出現回数
		for windex, n := range d {
			// TFd,t = 文書dにおける単語tの出現頻度
			tfmatrix[dindex][windex] = float64(n) / T
		}
	}
	return tfmatrix
}

func (b *Builder) createIDF(matrix [][]int) []float64 {
	// N = 文書数
	N := float64(len(matrix))
	idfmatrix := make([]float64, len(matrix[0]))
	// t = 単語
	for _, t := range b.indexByWord {
		// DFt = 単語tが出現する文書数
		df := 0.0
		// d = 文書
		for _, d := range matrix {
			// n = 文書d内の単語tの出現回数
			n := d[t]
			if n > 0 {
				df += 1.0
			}
		}
		idfmatrix[t] = math.Log((N + 1) / (df + 1))
	}
	return idfmatrix
}

// メタ情報を含む一文書
type Document struct {
	Key         string    // 識別子
	Name        string    // 任意の名前
	At          time.Time // 日付
	Description string    // 説明
	Words       []string  // 単語
}

func (d *Document) pickMeta() *DocMeta {
	meta := &DocMeta{}
	meta.Key = d.Key
	meta.Name = d.Name
	meta.At = d.At
	meta.Description = d.Description
	return meta
}

// 文書のメタ情報
type DocMeta struct {
	Key         string    // 識別子
	Name        string    // 任意の名前
	At          time.Time // 日付
	Description string    // 説明
}

// 複数のメタ情報を連結する
func joinDocMeta(metas []*DocMeta) *DocMeta {
	if len(metas) == 0 {
		return nil
	}
	meta := &DocMeta{}
	meta.Key = metas[0].Key
	meta.Name = metas[0].Name
	meta.At = metas[0].At
	if len(metas) == 1 {
		meta.Description = metas[0].Description
		return meta
	}
	descriptions := make([]string, len(metas))
	for i, m := range metas {
		d := m.At.Format("2006-01-02")
		// - Name(Key)
		// At
		// Description ~~
		// ~~~~~~~~~~~~~~
		descriptions[i] = fmt.Sprintf("- %s(%s)\n%s\n%s", m.Name, m.Key, d, m.Description)
	}
	meta.Description = strings.Join(descriptions, "\n")
	return meta
}
