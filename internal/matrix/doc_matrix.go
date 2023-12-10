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

func (b *MatrixBuilder) BuildDocMatrix() *DocMatrix {
	return NewDocMatrix(b.indexByWord, b.docs)
}

func (b *MatrixBuilder) BuildWeightDocMatrix() *DocMatrix {
	m := NewDocMatrix(b.indexByWord, b.docs)
	m.replaceWeight()
	return m
}

// 各文書の単語出現回数を保持する
type DocMatrix struct {
	// 出現単語。位置はwindexとしてidfstoreやdocsのrowに対応付けられる
	words []string
	// docs 行列
	docs []*doc
	// windexに対応したIDF
	idfStore []float64
}

func NewDocMatrix(
	indexByWord map[string]int,
	docs [][]string,
) *DocMatrix {
	dm := &DocMatrix{
		words:    make([]string, len(indexByWord)),
		docs:     make([]*doc, len(docs)),
		idfStore: make([]float64, len(indexByWord)),
	}
	for word, windex := range indexByWord {
		dm.words[windex] = word
	}

	dm.setupDocs(indexByWord, docs)
	dm.setupIDF()
	return dm
}

func (dm *DocMatrix) GetWordLen() int {
	return len(dm.words)
}

// [windex1]と[windex2]の共起頻度を返す
func (dm *DocMatrix) CoOccurrencetFrequency(windex1, windex2 int) float64 {
	// 共起頻度
	frequency := 0.0
	// 共起関係が存在した文書の数
	count := 0
	for dindex := range dm.docs {
		f := dm.CoOccurrencetFrequencyAt(dindex, windex1, windex2)
		if f > 0 {
			count++
		}
		// 共起頻度を足す
		frequency += f
	}
	// 文書の数で割り正規化する
	// 共起頻度の平均
	return frequency / float64(count)
}

// [dindex]の[windex1]と[windex2]の共起頻度を返す
func (dm *DocMatrix) CoOccurrencetFrequencyAt(dindex, windex1, windex2 int) float64 {
	return dm.GetAt(dindex, windex1) * dm.GetAt(dindex, windex2)
}

func (dm *DocMatrix) GetAt(dindex, windex int) float64 {
	return dm.docs[dindex].getAt(windex)
}

// TFIDFで重み付けされた共起行列を返す
func (dm *DocMatrix) replaceWeight() {
	// 重みづけされた文書の単語出現回数行列
	for dindex, doc := range dm.docs {
		for windex := range dm.words {
			tfidf := doc.tfAt(windex) * dm.getIDFAt(windex)
			dm.docs[dindex].row[windex] = tfidf
		}
	}
}

// [windex]のIDFを取得する
func (dm *DocMatrix) getIDFAt(windex int) float64 {
	return dm.idfStore[windex]
}

// 各文書における出現単語を行列化する
func (dm *DocMatrix) setupDocs(indexByWord map[string]int, docs [][]string) {
	totalWord := len(indexByWord)
	// 文書ごとの単語の出現回数を保持する
	// 各文書の単語を行列化
	for i, words := range docs {
		doc := newDoc(totalWord)
		for _, word := range words {
			windex := indexByWord[word]
			doc.addAt(windex)
		}
		dm.docs[i] = doc
	}
}

func (dm *DocMatrix) setupIDF() {
	totalDocs := float64(len(dm.docs))

	idf := make([]float64, len(dm.words))
	for windex := range dm.words {
		// 単語ごとにループ
		// 単語の出現回数
		count := 0.0
		for _, doc := range dm.docs {
			if doc.hasAt(windex) {
				count += 1.0
			}
		}
		if count > 0 {
			// ある単語[windex]のidf
			idf[windex] = math.Log(totalDocs / count)
		}
	}
	dm.idfStore = idf
}

// 1つの文書を表現する
type doc struct {
	// matrixbuilderのwordByIDに位置が対応した出現単語数
	row []float64
	// 文書内の単語数
	wordsCount int
}

// [l]コの単語列をもつドキュメントを作成する
func newDoc(l int) *doc {
	return &doc{row: make([]float64, l)}
}

func (d *doc) getAt(i int) float64 {
	return d.row[i]
}

// [i]に位置する単語が出現しているか
func (d *doc) hasAt(i int) bool {
	return d.row[i] > 0
}

func (d *doc) addAt(i int) {
	// 出現回数をカウントアップ
	d.row[i] += 1.0
	// 総単語数をカウントアップ
	d.wordsCount++
}

func (d *doc) tfAt(i int) float64 {
	return d.row[i] / float64(d.wordsCount)
}