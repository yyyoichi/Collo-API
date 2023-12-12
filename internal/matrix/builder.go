package matrix

import (
	"math"
	"sync"
)

type Builder struct {
	indexByWord map[string]int
	docs        [][]string

	mu sync.Mutex
}

func NewBuilder() *Builder {
	return &Builder{
		indexByWord: map[string]int{},
		docs:        [][]string{},
	}
}

func (b *Builder) AppendDoc(words []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, word := range words {
		if _, found := b.indexByWord[word]; !found {
			b.indexByWord[word] = len(b.indexByWord)
		}
	}
	b.docs = append(b.docs, words)
}

// 文書ごとの単語出現回数の行列を生成する
func (b *Builder) build() (*DocWordMatrix, *TFIDFMatrix) {
	lenWords := len(b.indexByWord)
	matrix := make([][]float64, len(b.docs))
	for dindex, doc := range b.docs {
		matrix[dindex] = make([]float64, lenWords)
		for _, word := range doc {
			windex := b.indexByWord[word]
			matrix[dindex][windex] += 1.0
		}
	}
	return &DocWordMatrix{matrix: matrix}, b.buildTFIDF(matrix)
}

func (b *Builder) buildTFIDF(matrix [][]float64) *TFIDFMatrix {
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
func (b *Builder) createTF(matrix [][]float64) [][]float64 {
	lenWords := len(matrix[0])

	tfmatrix := make([][]float64, len(matrix))
	for dindex, d := range matrix {
		tfmatrix[dindex] = make([]float64, lenWords)
		// t = 文書d内の合計単語数
		T := float64(len(b.docs[dindex]))
		// n = 文書d内の単語tの出現回数
		for windex, n := range d {
			// TFd,t = 文書dにおける単語tの出現頻度
			tfmatrix[dindex][windex] = n / T
		}
	}
	return tfmatrix
}

func (b *Builder) createIDF(matrix [][]float64) []float64 {
	// N = 文書数
	N := float64(len(matrix))
	idfmatrix := make([]float64, len(matrix))
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
