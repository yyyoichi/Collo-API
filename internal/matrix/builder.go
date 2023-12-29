package matrix

import (
	"math"
	"sync"
	"time"
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
	lenWords := len(b.indexByWord)
	matrix := make([][]int, len(b.documents))
	for dindex, doc := range b.documents {
		matrix[dindex] = make([]int, lenWords)
		for _, word := range doc.Words {
			windex := b.indexByWord[word]
			matrix[dindex][windex] += 1.0
		}
	}
	words := make([]string, len(b.indexByWord))
	for word, windex := range b.indexByWord {
		words[windex] = word
	}
	return &DocWordMatrix{matrix: matrix, words: words}, b.buildTFIDF(matrix)
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
