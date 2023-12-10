package matrix

import (
	"math"
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

func (b *MatrixBuilder) BuildCountDocMatrix() DocMatrixInterface {
	return NewCountDocMatrix(b.indexByWord, b.docs)
}

func (b *MatrixBuilder) BuildTFIDFDocMatrix() DocMatrixInterface {
	m := NewCountDocMatrix(b.indexByWord, b.docs)
	return NewTFIDFDocMatrix(m)
}

type DocMatrixInterface interface {
	// 単語数を取得する
	LenWords() int
	// 文書数を取得する
	LenDocs() int
	// 文書[dindex]における単語[windex]の値を取得する
	GetValueAt(dindex, windex int) float64
	// Docs全体における共起単語[coIndex]の共起頻度を取得する
	CoOccurrencetFrequency(coIndex CoIndex) float64
	// 文書[dindex]における共起単語[coIndex]の共起頻度を取得する
	CoOccurrencetFrequencyAt(dindex int, coIndex CoIndex) float64
}

// implement DocMatrixInterface
type docMatrixBase[T any] struct {
	// 出現単語。位置はwindexとしてidfstoreやdocsのrowに対応付けられる
	words []string
	// any型の文書のスライス。
	docs []T
}

func newBaseDocMatrix[T any](words []string, docs []T) *docMatrixBase[T] {
	m := &docMatrixBase[T]{
		words: words,
		docs:  docs,
	}
	return m
}

func (m *docMatrixBase[T]) GetValueAt(dindex, windex int) float64 {
	defer panic("metrix.DocMatrix.GetValueAt is not implemented")
	return 0.0
}

func (m *docMatrixBase[T]) LenDocs() int {
	return len(m.docs)
}
func (m *docMatrixBase[T]) LenWords() int {
	return len(m.words)
}
func (m *docMatrixBase[T]) CoOccurrencetFrequency(coIndex CoIndex) float64 {
	// 共起頻度
	frequency := 0.0
	// 共起関係が存在した文書の数
	count := 0
	for dindex := 0; dindex < m.LenDocs(); dindex++ {
		f := m.CoOccurrencetFrequencyAt(dindex, coIndex)
		if f > 0 {
			count++
		}
		// 共起頻度を足す
		frequency += f
	}
	if count == 0 {
		return 0.0
	}
	// 文書の数で割り正規化する
	// 共起頻度の平均
	return frequency / float64(count)
}
func (m *docMatrixBase[T]) CoOccurrencetFrequencyAt(dindex int, coIndex CoIndex) float64 {
	return m.GetValueAt(dindex, coIndex.I1()) * m.GetValueAt(dindex, coIndex.I2())
}

// TF-IDFで重みを付けた出現単語行列
// implement DocMatrixInterface
type TFIDFDocMatrix struct {
	*docMatrixBase[[]float64]
}

func NewTFIDFDocMatrix(docMatrix *CountDocMatrix) *TFIDFDocMatrix {
	matrix := make([][]float64, docMatrix.LenDocs())
	for dindex := range docMatrix.docs {
		matrix[dindex] = make([]float64, docMatrix.LenWords())
		for windex := 0; windex < docMatrix.LenWords(); windex++ {
			matrix[dindex][windex] = docMatrix.getTFIDFAt(dindex, windex)
		}
	}
	m := &TFIDFDocMatrix{
		docMatrixBase: newBaseDocMatrix(docMatrix.words, matrix),
	}
	return m
}
func (m *TFIDFDocMatrix) GetValueAt(dindex, windex int) float64 {
	return m.docs[dindex][windex]
}

// 単語の出現回数を値として持つ出現単語行列
// implement DocMatrixInterface
type CountDocMatrix struct {
	*docMatrixBase[*countDoc]
	// windexに対応したIDF
	idfStore []float64
}

func NewCountDocMatrix(
	indexByWord map[string]int,
	docs [][]string,
) *CountDocMatrix {
	words := make([]string, len(indexByWord))
	for word, windex := range indexByWord {
		words[windex] = word
	}

	d := make([]*countDoc, len(docs))

	m := &CountDocMatrix{
		docMatrixBase: newBaseDocMatrix(words, d),
		idfStore:      make([]float64, len(indexByWord)),
	}
	m.setupDocs(indexByWord, docs)
	m.setupIDF()
	return m
}
func (m *CountDocMatrix) GetValueAt(dindex, windex int) float64 {
	return float64(m.docs[dindex].getAt(windex))
}

// 文書[dindex]における単語[windex]のTF-IDFを返す
func (m *CountDocMatrix) getTFIDFAt(dindex, windex int) float64 {
	return m.docs[dindex].tfAt(windex) * m.getIDFAt(windex)
}

// [windex]のIDFを取得する
func (m *CountDocMatrix) getIDFAt(windex int) float64 {
	return m.idfStore[windex]
}

// 各文書における出現単語を行列化する
func (m *CountDocMatrix) setupDocs(indexByWord map[string]int, docs [][]string) {
	totalWord := len(indexByWord)
	// 文書ごとの単語の出現回数を保持する
	// 各文書の単語を行列化
	for i, words := range docs {
		doc := newDoc(totalWord)
		for _, word := range words {
			windex := indexByWord[word]
			doc.addAt(windex)
		}
		m.docs[i] = doc
	}
}

func (m *CountDocMatrix) setupIDF() {
	totalDocs := float64(len(m.docs))

	idf := make([]float64, m.LenWords())
	for windex := 0; windex < m.LenWords(); windex++ {
		// 単語ごとにループ
		// 単語の出現回数
		count := 0.0
		for _, doc := range m.docs {
			if doc.hasAt(windex) {
				count += 1.0
			}
		}
		if count > 0 {
			// ある単語[windex]のidf
			idf[windex] = math.Log(totalDocs / count)
		}
	}
	m.idfStore = idf
}

// 1つの文書を表現する
type countDoc struct {
	// matrixbuilderのwordByIDに位置が対応した出現単語数
	row []int
	// 文書内の単語数
	wordsCount int
}

// [l]コの単語列をもつドキュメントを作成する
func newDoc(lenWords int) *countDoc {
	return &countDoc{row: make([]int, lenWords)}
}

func (d *countDoc) getAt(windex int) int {
	return d.row[windex]
}

// [windex]に位置する単語が出現しているか
func (d *countDoc) hasAt(windex int) bool {
	return d.row[windex] > 0
}

func (d *countDoc) addAt(windex int) {
	// 出現回数をカウントアップ
	d.row[windex] += 1
	// 総単語数をカウントアップ
	d.wordsCount++
}

func (d *countDoc) tfAt(windex int) float64 {
	return float64(d.row[windex]) / float64(d.wordsCount)
}
