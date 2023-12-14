package matrix

import (
	"context"
	"errors"
	"math"
	"sort"
	"sync"
	"yyyoichi/Collo-API/pkg/stream"
	"yyyoichi/Collo-API/pkg/structuer"

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
	// 共起行列の正規化アルゴリズム
	coOccurrencetNormalization CoOccurrencetNormalization
	// ノードの中心性を求めるアルゴリズム
	nodeRatingAlgorithm NodeRatingAlgorithm
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
		coOccurrencetNormalization: config.CoOccurrencetNormalization,
		nodeRatingAlgorithm:        config.NodeRatingAlgorithm,
	}
	m.init()
	go func() {
		m.setProgress(DwMStart)
		dwm, tfidf := builder.Build()

		// 列数削減
		m.setProgress(DwMReduceCol)
		col := tfidf.TopPercentageWIndexes(config.ReduceThreshold, config.MinNodes)
		col.Reduce(dwm)

		m.setup(dwm)
	}()

	return m
}

func NewCoMatrixFromDocWordM(
	dwm *DocWordMatrix,
	coOccurrencetNormalization CoOccurrencetNormalization,
	nodeRatingAlgorithm NodeRatingAlgorithm,
) *CoMatrix {
	m := &CoMatrix{
		coOccurrencetNormalization: coOccurrencetNormalization,
		nodeRatingAlgorithm:        nodeRatingAlgorithm,
	}

	go m.setup(dwm)
	return m
}

func (m *CoMatrix) MostImportantNode() *Node {
	return m.NodeRank(0)
}

// 重要度順位[rank](0~)のNodeを返す
func (m *CoMatrix) NodeRank(rank int) *Node {
	if rank >= len(m.words) {
		return nil
	}
	id := m.indices[rank]
	return m.NodeID(id)
}

// Node[id]のNodeを返す
func (m *CoMatrix) NodeID(id int) *Node {
	return &Node{
		ID:   uint(id),
		Word: m.words[id],
		Rate: m.priority[id],
	}
}

func (m *CoMatrix) Edge(id1, id2 int) *Edge {
	edge := &Edge{}
	edge.Node1ID = uint(id1)
	edge.Node2ID = uint(id2)

	// setID
	// 数字の小さいIDを行にして、1次元スライス上の共起行列の位置をEdgeのIDとして利用する
	n := len(m.words)
	if id1 <= id2 {
		row := id1
		edge.ID = uint(row*n + id2)
	} else {
		row := id2
		edge.ID = uint(row*n + id1)
	}

	edge.Rate = m.matrix[edge.ID]
	return edge
}

// NodeIDの共起関係にあるNodeとそのEdgeを返す
func (m *CoMatrix) CoOccurrenceRelation(nodeID uint) (nodes []*Node, edges []*Edge) {
	nodes = []*Node{}
	edges = []*Edge{}

	if int(nodeID) >= len(m.words) {
		return nodes, edges
	}

	subjectNodeID := int(nodeID)
	for objectNodeID := range m.words {
		edge := m.Edge(subjectNodeID, objectNodeID)
		if edge.Rate <= 0 {
			continue
		}
		nodes = append(nodes, m.NodeID(objectNodeID))
		edges = append(edges, edge)
	}

	return nodes, edges
}

// NodeIDsと共起関係にあるNodeとそのEdgeを返す
func (m *CoMatrix) CoOccurrences(nodeIDs ...uint) (nodes []*Node, edges []*Edge) {
	nodeset := structuer.NewSet[*Node](func(n *Node) any { return n.ID })
	edgeset := structuer.NewSet[*Edge](func(e *Edge) any { return e.ID })

	ctx := context.Background()
	nodeIDCh := stream.Generator[uint](ctx, nodeIDs...)
	doneCh := stream.Line[uint, interface{}](ctx, nodeIDCh, func(nodeID uint) interface{} {
		ns, es := m.CoOccurrenceRelation(uint(nodeID))
		for _, n := range ns {
			nodeset.Add(n)
		}
		for _, e := range es {
			edgeset.Add(e)
		}
		return struct{}{}
	})
	for range doneCh {
	}

	return nodeset.ToSlice(), edgeset.ToSlice()
}

// [dept]回、NodeIDの共起関係にあるNodeを再帰的に取得する
func (m *CoMatrix) CoOccurrenceDept(dept int, nodeID uint) (nodes []*Node, edges []*Edge) {

	nodeset := structuer.NewSet[*Node](func(n *Node) any { return n.ID })
	edgeset := structuer.NewSet[*Edge](func(e *Edge) any { return e.ID })

	var fn func(d int, id uint) int
	fn = func(d int, id uint) int {
		if d == 0 {
			return 1
		}
		ns, es := m.CoOccurrenceRelation(id)
		for _, e := range es {
			edgeset.Add(e)
		}
		var nodewait sync.WaitGroup
		for _, n := range ns {
			nodewait.Add(1)
			go func(n *Node) {
				defer nodewait.Done()
				nodeset.Add(n)
				fn(d-1, n.ID)
			}(n)
		}
		nodewait.Wait()
		return 0
	}

	fn(dept, nodeID)

	return nodeset.ToSlice(), edgeset.ToSlice()
}

func (m *CoMatrix) ConsumeProgress() <-chan CoMatrixProgress {
	return m.progress
}

func (m *CoMatrix) Error() error {
	return m.err
}

// exp called go routine
func (m *CoMatrix) setup(dwm *DocWordMatrix) {
	m.setProgress(CoMStart)
	m.words = dwm.words
	m.init()
	m.setProgress(CoMCreateMatrix)
	switch m.coOccurrencetNormalization {
	case Dice:
		m.matrixByDice(dwm)
	}
	m.setProgress(CoMCalcNodeImportance)

	var err error
	switch m.nodeRatingAlgorithm {
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
	if m.progress == nil {
		m.progress = make(chan CoMatrixProgress)
	}
	if m.coOccurrencetNormalization == 0 {
		m.coOccurrencetNormalization = Dice
	}
	if m.nodeRatingAlgorithm == 0 {
		m.nodeRatingAlgorithm = VectorCentrality
	}

	if m.words == nil {
		return
	}

	n := len(m.words)
	if m.matrix == nil {
		m.matrix = make([]float64, n*n)
	}
	if m.indices == nil {
		m.indices = make([]int, n)
	}
	if m.priority == nil {
		m.priority = make([]float64, n)
		for i := range m.indices {
			m.indices[i] = i
		}
	}
}

type Node struct {
	ID   uint
	Word string
	Rate float64
}

type Edge struct {
	ID, Node1ID, Node2ID uint
	Rate                 float64
}
