package matrix

import (
	"sync"
)

// 文書全体に出現した重要な単語の共起関係から共起ネットワークを生成する
// step.1 各文書における出現単語を行列化する
// step.2 TF-IDFで単語の重要度を計算する
// step.3 単語の重要度で重みを付けた文書全体の共起行列を作成する
// step.4 共起行列の固有値から各単語の固有ベクトルを計算する
// step.5 各単語の固有ベクトルを0~1スケールした固有ベクトル´を計算する
// step.6 固有ベクトル´をしきい値で単語の足切りを行い重要共起単語を計算する
// step.7 step.3の共起行列から重要共起単語のみを用いて共起ネットワークを生成する

type MatrixBuilder struct {
	indexByWord map[string]int
	docs        [][]string

	mu sync.Mutex
}

func NewMatrixBuilder() *MatrixBuilder {
	return &MatrixBuilder{
		indexByWord: map[string]int{},
		docs:        [][]string{},
	}
}

func (b *MatrixBuilder) AppendDoc(words []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, word := range words {
		if _, found := b.indexByWord[word]; !found {
			b.indexByWord[word] = len(b.indexByWord)
		}
	}
	b.docs = append(b.docs, words)
}

func (b *MatrixBuilder) Build() [][]uint {
	return b.toMatrix()
}

func (b *MatrixBuilder) toMatrix() [][]uint {
	matrix := make([][]uint, len(b.docs))
	cols := len(b.indexByWord)
	for i, doc := range b.docs {
		row := make([]uint, cols)
		for _, word := range doc {
			i, found := b.indexByWord[word]
			if !found {
				continue
			}
			if row[i] > 0 {
				row[i] += 1
			} else {
				row[i] = 1
			}
		}
		matrix[i] = row
	}
	return matrix
}
