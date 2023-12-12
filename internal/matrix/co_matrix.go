package matrix

import (
	"context"
	"errors"
	"sort"
	"yyyoichi/Collo-API/pkg/stream"

	"gonum.org/v1/gonum/mat"
)

var (
	ErrSymEigendeCompFailed = errors.New("symmetric eigendecomposition failed")
	ErrInvalidVectorsLangth = errors.New("vectors do not have the same length as the words")
	ErrInvalidAlgorithm     = errors.New("invalid algorithm")
)

// 共起関係の解釈に責任を持つ
type CoMatrix struct {
	DocMatrixInterface
	config Config
	// priority 順。単語インデックスを持つ。
	indices []int
	// 単語の重要度。位置はDocMatrixのwordsの位置に対応する
	priority []float64
}

// 共起行列を生成する
// [totalWord]出現する単語数
func NewCoMatrix(docMatrix DocMatrixInterface, config Config) (*CoMatrix, error) {
	// 単語のインデックスを作成

	m := &CoMatrix{
		DocMatrixInterface: docMatrix,
		config:             config,
	}
	m.init()

	var err error
	switch m.config.WordRatingAlgorithm {
	case VectorCentrality:
		err = m.useVectorCentrality()
	default:
		err = ErrInvalidAlgorithm
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// 共起頻度の低い単語を重要語としたい場合、
// 左特異ベクトルを求める。

// どれほど共起関係の中心にあるかで単語の重要度を決定する。
// （固有ベクトル中心性を単語の重要度に使用する。）
func (m *CoMatrix) useVectorCentrality() error {
	n := m.LenWords()

	// 共起行列
	dence := mat.NewSymDense(n, make([]float64, n*n))

	ctx := context.Background()
	coIndexCh := m.generateCoIndex(ctx)
	doneCh := stream.FunIO[[2]int, interface{}](ctx, coIndexCh, func(co [2]int) interface{} {
		frequency := m.CoOccurrencetFrequency(co)
		dence.SetSym(co[0], co[1], frequency)
		return struct{}{}
	})
	for range doneCh {
	}

	// 固有値保持
	var eigsym mat.EigenSym
	if ok := eigsym.Factorize(dence, true); !ok {
		return ErrSymEigendeCompFailed
	}

	// 固有ベクトル行列
	var ev mat.Dense
	eigsym.VectorsTo(&ev)

	rows, _ := ev.Dims()
	if rows != m.LenWords() {
		return ErrInvalidVectorsLangth
	}

	// 各行（単語）のノルムを計算し、中心性とする
	// 中心性を単語の重要度とする
	for i := 0; i < rows; i++ {
		var row mat.VecDense
		row.CloneFromVec(ev.RowView(i))

		m.priority[i] = mat.Norm(&row, 2)
	}

	// 新しいpriorityをスケーリングする
	m.scalingPriority()
	return nil
}

func (m *CoMatrix) scalingPriority() {

	// 重要度に基づいて単語のインデックスを降順ソート
	sort.Slice(m.indices, func(i, j int) bool {
		return m.priority[m.indices[i]] > m.priority[m.indices[j]]
	})

	// 重要度最小値
	minVal := m.priority[m.indices[len(m.indices)-1]]
	// 重要度最大値
	maxVal := m.priority[0]

	minSc := m.config.MinScale
	maxSc := m.config.MaxScale

	scaledPriority := make([]float64, len(m.priority))
	for i, val := range m.priority {
		scaledPriority[i] = ((val-minVal)/(maxVal-minVal))*(maxSc-minSc) + minSc
	}

	m.priority = scaledPriority
}

// 共起行列の共起ペアindexを返す
func (m *CoMatrix) generateCoIndex(ctx context.Context) <-chan [2]int {
	n := m.LenWords()
	ch := make(chan [2]int)
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- [2]int{i, j}
				}
			}
		}
	}()

	return ch
}

func (m *CoMatrix) init() {
	if 0 < m.config.Threshold && m.config.Threshold <= 1 {
		m.config.Threshold = 0.1
	}
	if m.config.MaxScale <= m.config.MinScale {
		m.config.MaxScale += m.config.MinScale + 1
	}
	if m.config.WordRatingAlgorithm == 0 {
		m.config.WordRatingAlgorithm = VectorCentrality
	}
	if m.priority == nil {
		m.priority = make([]float64, m.LenWords())
	}
	if m.indices == nil {
		m.indices = make([]int, m.LenWords())
		for i := range m.indices {
			m.indices[i] = i
		}
	}
}
