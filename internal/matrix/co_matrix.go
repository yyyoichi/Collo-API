package matrix

import (
	"context"
	"errors"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

var (
	ErrSymEigendeCompFailed = errors.New("symmetric eigendecomposition failed")
	ErrInvalidVectorsLangth = errors.New("vectors do not have the same length as the words")
	ErrInvalidAlgorithm     = errors.New("invalid algorithm")
)

type CoMatrixProgress int

const (
	DwMStart              CoMatrixProgress = 10
	DwMReduceCol          CoMatrixProgress = 11 // Reducing coloumns of doc-word-matrix
	CoMStart              CoMatrixProgress = 20
	CoMCreateMatrix       CoMatrixProgress = 21 // creating co-occurrencet matrix
	CoMCalcNodeImportance CoMatrixProgress = 22 // calcing node importance
	ErrDone               CoMatrixProgress = 88 // error is occuered
	ProgressDone          CoMatrixProgress = 99 // done initialization

	// Reteは小数第5位未満を四捨五入する
	RateFixed = 5
)

// 共起関係の解釈に責任を持つ
type CoMatrix struct {
	config Config
	// 共起行列
	matrix []float64
	// priority 順。単語インデックスを持つ。
	indices []int
	// 単語の重要度。位置はDocMatrixのwordsの位置に対応する
	priority []float64
	// 単語。位置はwindex
	words []string
	// 起動進捗
	progress chan CoMatrixProgress
	// 進捗が終了しているか
	done bool
	// errors
	err error
}

func NewCoMatrixFromBuilder(builder *Builder, config Config) *CoMatrix {
	m := &CoMatrix{
		config:   config,
		progress: make(chan CoMatrixProgress),
	}
	m.init()
	go func() {
		m.setProgress(DwMStart)

		dwm, tfidf := builder.Build()

		m.setProgress(DwMReduceCol)
		col := tfidf.TopPercentageWIndexes(m.config.ReduceThreshold, m.config.MinNodes)
		// 列数削減
		col.Reduce(dwm)

		// 単語数
		n := len(dwm.words)
		m.matrix = make([]float64, n*n)
		m.indices = make([]int, n)
		m.priority = make([]float64, n)
		m.words = dwm.words
		for i := range m.indices {
			m.indices[i] = i
		}

		m.setProgress(CoMStart)
		m.setup(dwm)
	}()

	return m
}

func NewCoMatrixFromDocWordM(dwm *DocWordMatrix, config Config) *CoMatrix {
	n := len(dwm.words)
	m := &CoMatrix{
		config:   config,
		matrix:   make([]float64, n*n),
		indices:  make([]int, n),
		priority: make([]float64, n),
		words:    dwm.words,
		progress: make(chan CoMatrixProgress),
	}
	m.init()
	m.indices = make([]int, len(m.words))
	for i := range m.indices {
		m.indices[i] = i
	}
	go func() {
		m.setProgress(CoMStart)
		m.setup(dwm)
	}()
	return m
}

func (m *CoMatrix) MostImportantNode() *Node {
	return m.Node(0)
}

func (m *CoMatrix) Node(i int) *Node {
	if i >= len(m.words) {
		return nil
	}
	id := m.indices[i]
	return &Node{
		ID:   uint(id),
		Word: m.words[id],
		Rate: m.priority[id],
	}
}

func (m *CoMatrix) ConsumeProgress() <-chan CoMatrixProgress {
	return m.progress
}

func (m *CoMatrix) Error() error {
	return m.err
}

// exp called go routine
func (m *CoMatrix) setup(dwm *DocWordMatrix) {
	m.setProgress(CoMCreateMatrix)
	switch m.config.CoOccurrencetNormalization {
	case Dice:
		m.matrixByDice(dwm)
	}
	m.setProgress(CoMCalcNodeImportance)

	var err error
	switch m.config.NodeRatingAlgorithm {
	case VectorCentrality:
		err = m.useVectorCentrality()
	}
	if err != nil {
		m.doneProgressWithError(err)
	}
	m.doneProgress()
}

// 共起回数の正規化にDice係数を利用して共起行列を作成する
func (m *CoMatrix) matrixByDice(dwm *DocWordMatrix) {
	// create matrix by dice
	occuerences := make([]DocWordOccurances, len(m.words))
	for windex := range m.words {
		occuerences[windex] = dwm.Occurances(windex)
	}

	ctx := context.Background()
	frequencyCh := dwm.GenerateCoOccurrencetFrequency(ctx)
	for f := range frequencyCh {
		// Dice(Wi,Wj) = 2 x 共起回数Wi,Wj / (出現回数Wi + 出現回数Wj )
		d := float64(occuerences[f.Windex1].Occurances + occuerences[f.Windex2].Occurances)
		value := float64(2*f.Frequency) / d
		m.syncSet(f.Windex1, f.Windex2, value)
	}
}

// 共起行列に共起回数をセットする
func (m *CoMatrix) syncSet(windex1, windex2 int, value float64) {
	n := len(m.words) // 単語数
	var i int
	i = windex1*n + windex2
	m.matrix[i] = value
	i = windex1 + windex2*n
	m.matrix[i] = value
}

// どれほど共起関係の中心にあるかで単語の重要度を決定する。
// （固有ベクトル中心性を単語の重要度に使用する。）
func (m *CoMatrix) useVectorCentrality() error {
	// 単語数
	n := len(m.words)
	// 対称行列化
	dence := mat.NewSymDense(n, m.matrix)
	// 固有値分解
	var eigsym mat.EigenSym
	if ok := eigsym.Factorize(dence, true); !ok {
		return ErrSymEigendeCompFailed
	}

	// 最大固有値
	maxEigenvalue := math.Inf(-1)
	maxEigenvalueIndex := 0
	for i, v := range eigsym.Values(nil) {
		if maxEigenvalue < v {
			maxEigenvalue = v
			maxEigenvalueIndex = i
		}
	}

	// 固有ベクトル行列
	var ev mat.Dense
	eigsym.VectorsTo(&ev)

	rows, _ := ev.Dims()
	if rows != n {
		return ErrInvalidVectorsLangth
	}

	// 中心性を単語の重要度とする
	for i := 0; i < rows; i++ {
		// 各要素の二乗を足す(内積)
		m.priority[i] = ev.RawRowView(i)[maxEigenvalueIndex]
	}

	// 新しいpriorityをスケーリングする
	m.scalingPriority()
	return nil
}

// 優先度を0-1に標準化する
func (m *CoMatrix) scalingPriority() {

	// 重要度に基づいて単語のインデックスを降順ソート
	sort.Slice(m.indices, func(i, j int) bool {
		return m.priority[m.indices[i]] > m.priority[m.indices[j]]
	})

	// 重要度最小値
	minVal := m.priority[m.indices[len(m.indices)-1]]
	// 重要度最大値
	maxVal := m.priority[m.indices[0]]

	// 小数第RateFixed位未満四捨五入
	s := math.Pow10(RateFixed)
	for i, val := range m.priority {
		p := (val - minVal) / (maxVal - minVal)
		m.priority[i] = math.Round(p*s) / s
	}
}

func (m *CoMatrix) setProgress(p CoMatrixProgress) {
	if !m.done {
		m.progress <- p
	}
}

func (m *CoMatrix) doneProgress() {
	defer close(m.progress)
	m.setProgress(ProgressDone)
	m.done = true
}

func (m *CoMatrix) doneProgressWithError(err error) {
	defer close(m.progress)
	m.setProgress(ErrDone)
	m.done = true
	m.err = err
}

func (m *CoMatrix) init() {
	if m.config.ReduceThreshold <= 0 || 1 < m.config.ReduceThreshold {
		m.config.ReduceThreshold = 0.1
	}
	if m.config.MinNodes == 0 {
		m.config.MinNodes = 300
	}
	if m.config.CoOccurrencetNormalization == 0 {
		m.config.CoOccurrencetNormalization = Dice
	}
	if m.config.NodeRatingAlgorithm == 0 {
		m.config.NodeRatingAlgorithm = VectorCentrality
	}
}

type Node struct {
	ID   uint
	Word string
	Rate float64
}
