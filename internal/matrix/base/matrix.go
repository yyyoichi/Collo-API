package base

import (
	"iter"
	"maps"
)

// 文書-単語カウント行列
type Matrix struct {
	// 単語識別子が配列dataのいずれの列Xに位置するか
	xindex map[int]int
	// 文書識別子が配列dataのいずれの行Yに位置するか
	yindex map[string]int

	// カウント行列
	data []float64
}

func New(numDoc, numWord int, docIds iter.Seq[string], wordIds iter.Seq[int]) Matrix {
	var m = Matrix{
		xindex: make(map[int]int, numDoc),
		yindex: make(map[string]int, numWord),
		data:   make([]float64, numDoc*numWord),
	}
	var yi int
	for id := range docIds {
		m.yindex[id] = yi
		yi++
	}
	var xi int
	for id := range wordIds {
		m.xindex[id] = xi
		xi++
	}
	return m
}

func (m *Matrix) Set(docId string, wordId int, n float64) {
	*m.N(docId, wordId) = n
}

// 文書数
func (m *Matrix) CountDoc() float64 {
	return float64(len(m.yindex))
}

// 単語数
func (m *Matrix) CountWord() float64 {
	return float64(len(m.xindex))
}

func (m *Matrix) IterDocId() iter.Seq[string] {
	return maps.Keys(m.yindex)
}

func (m *Matrix) IterWordId() iter.Seq[int] {
	return maps.Keys(m.xindex)
}

// 単語idの指数をすべて文書中から加算する。
func (m *Matrix) SumDocWithWord(id int) float64 {
	var count float64
	for n := range m.nInWord(id) {
		count += n
	}
	return count
}

// 単語idの指数が0を超える文書数
func (m *Matrix) CountDocWithWord(id int) float64 {
	var count float64
	for n := range m.nInWord(id) {
		if n > 0 {
			count++
		}
	}
	return count
}

// 文書idに含まれる指数をすべて加算する。
func (m *Matrix) SumWordInDoc(id string) float64 {
	var count float64
	for _, n := range m.nIndoc(id) {
		count += n
	}
	return count
}

// 文書idに含まれる指数0を超えるの単語数
func (m *Matrix) CountWordInDoc(id string) float64 {
	var count float64
	for _, n := range m.nIndoc(id) {
		if n > 0 {
			count++
		}
	}
	return count
}

// 文書idの指数をスライスで返す。
func (m *Matrix) nIndoc(id string) []float64 {
	yi := m.yindex[id] * len(m.xindex)
	return m.data[yi : yi+len(m.xindex)]
}

// 単語idの指標をiterで返す。
func (m *Matrix) nInWord(id int) iter.Seq[float64] {
	return func(yield func(float64) bool) {
		for i := m.xindex[id]; i < len(m.data); i += len(m.xindex) {
			if ok := yield(m.data[i]); !ok {
				return
			}
		}
	}
}

// 文書における単語指数
func (m *Matrix) N(docId string, wordId int) *float64 {
	yi := m.yindex[docId]
	xi := m.xindex[wordId]
	return &m.data[yi*len(m.xindex)+xi]
}

// 同一メトリクスをデータ空で作成する。
func (m *Matrix) Copy() Matrix {
	return New(len(m.yindex), len(m.xindex), m.IterDocId(), m.IterWordId())
}
