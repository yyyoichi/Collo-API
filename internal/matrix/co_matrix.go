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
	ErrInvalidVectorsLangth = errors.New("vectors do not have the same length as the Words")
	ErrInvalidAlgorithm     = errors.New("invalid algorithm")
)

type CoMatrixProgress int

const (
	DwMStart              CoMatrixProgress = 10
	DwMReduceCol          CoMatrixProgress = 11 // Reducing coloumns of doc-word-Matrix
	CoMStart              CoMatrixProgress = 20
	CoMCreateMatrix       CoMatrixProgress = 21 // creating co-occurrencet Matrix
	CoMCalcNodeImportance CoMatrixProgress = 22 // calcing node importance
	ErrDone               CoMatrixProgress = 88 // error is occuered
	ProgressDone          CoMatrixProgress = 99 // done initialization

	// Reteは小数第5位未満を四捨五入する
	RateFixed = 5
)

type (
	MatrixError struct{ error }

	CoMatrixes struct {
		Data  []CoMatrix `json:"d"`
		Words []string   `json:"w"`
	}
	Handler struct {
		Err func(error)
	}
	CoMatrix struct {
		// 共起行列の正規化アルゴリズム
		coOccurrencetNormalization CoOccurrencetNormalization
		// ノードの中心性を求めるアルゴリズム
		nodeRatingAlgorithm NodeRatingAlgorithm
		// 共起行列に含まれる文書群のメタ情報
		Meta MultiDocMeta `json:"m"`
		// 共起行列
		Matrix []float64 `json:"x"`
		// Priority 順。単語インデックスを持つ。
		Indices []int `json:"i"`
		// 単語の重要度。位置はDocMatrixのwordsの位置に対応する
		Priority []float64 `json:"p"`
		// 単語。位置はwindex
		PtrWords *[]string
		// 起動進捗
		progress chan CoMatrixProgress
		// 進捗が終了しているか
		done bool
		// errors
		err error
	}
)

// Builderから複数の共起行列を返す。CoMatrixes.Data の各CoMatrixはConsumeProgress()して更新を開始してください。
func NewCoMatrixesFromBuilder(ctx context.Context, builder Builder, config Config) CoMatrixes {
	config.init()
	// 文書単語行列からTF-IDFを計算し列削除を準備する
	tfidf := builder.BuildTFIDF(config.GroupingFuncType)
	col := tfidf.GetColumnReduction(config.ReduceThreshold, config.MinNodes, config.MaxNodes)
	// 列数削減
	col.Reduce(&builder)

	var cos CoMatrixes
	cos.Words = builder.Words()
	newCoMatrix := func(dwm DocWordMatrix) CoMatrix {
		return NewCoMatrixFromDocWordMatrix(&cos.Words, dwm, config.CoOccurrencetNormalization, config.NodeRatingAlgorithm)
	}

	// set cos.Data
	switch {
	case config.AtGroupID == "total":
		alldwm := builder.BuildDocWordMatrix(ctx, "total")
		cos.Data = make([]CoMatrix, 1)
		cos.Data[0] = newCoMatrix(alldwm)
	case config.AtGroupID != "":
		n, dwmCh := builder.BuildDocWordMatrixByGroupAt(ctx, config.GroupingFuncType, config.AtGroupID)
		if n != 1 {
			// nodata
			return cos
		}
		cos.Data = make([]CoMatrix, 1)
		var dwm DocWordMatrix
		for dwm = range dwmCh {
		}
		cos.Data[0] = newCoMatrix(dwm)
	default:
		// AtGroupID is empty
		alldwm := builder.BuildDocWordMatrix(ctx, "total")
		n, dwmCh := builder.BuildDocWordMatrixByGroup(ctx, config.GroupingFuncType)
		if n == 0 {
			cos.Data = make([]CoMatrix, 1)
			cos.Data[0] = newCoMatrix(alldwm)
		} else {
			cos.Data = make([]CoMatrix, 1, 1+n)
			cos.Data[0] = newCoMatrix(alldwm)
			for dwm := range dwmCh {
				cos.Data = append(cos.Data, newCoMatrix(dwm))
			}
		}
	}
	return cos
}

func NewCoMatrixFromDocWordMatrix(
	PtrWords *[]string,
	dwm DocWordMatrix,
	coOccurrencetNormalization CoOccurrencetNormalization,
	nodeRatingAlgorithm NodeRatingAlgorithm,
) CoMatrix {
	m := CoMatrix{
		coOccurrencetNormalization: coOccurrencetNormalization,
		nodeRatingAlgorithm:        nodeRatingAlgorithm,
		Meta:                       dwm.meta,
		progress:                   make(chan CoMatrixProgress),
		PtrWords:                   PtrWords,
	}
	m.init()

	go func() {
		m.setProgress(CoMStart)
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
	}()
	return m
}

func (m *CoMatrix) ValidNodeID(nodeID uint) bool {
	return nodeID != 0 && int(nodeID) <= len(*m.PtrWords)
}

func (m *CoMatrix) MostImportantNode() *Node {
	return m.NodeRank(0)
}

// 重要度順位[rank](0~)のNodeを返す
func (m *CoMatrix) NodeRank(rank int) *Node {
	if rank >= len(*m.PtrWords) {
		return nil
	}
	id := m.Indices[rank]
	return m.NodeID(uint(id) + 1)
}

// Node[nodeID](1~)のNodeを返す
func (m *CoMatrix) NodeID(nodeID uint) *Node {
	if !m.ValidNodeID(nodeID) {
		return nil
	}
	// NodeIDは1から始まる。クライアントのためのインデックス。windexとは別個。
	var numEdges uint = 0
	// matrix中の行位置
	startIndex := (int(nodeID) - 1) * len(*m.PtrWords)
	for i := 0; i < len(*m.PtrWords); i++ {
		rate := m.Matrix[startIndex+i]
		if rate > 0 {
			numEdges++
		}
	}
	return &Node{
		ID:       uint(nodeID),
		Word:     (*m.PtrWords)[nodeID-1],
		Rate:     m.Priority[nodeID-1],
		NumEdges: numEdges,
	}
}

func (m *CoMatrix) Edge(nodeID1, nodeID2 uint) *Edge {
	if !m.ValidNodeID(nodeID1) || !m.ValidNodeID(nodeID2) {
		return nil
	}
	edge := &Edge{}
	edge.Node1ID = nodeID1
	edge.Node2ID = nodeID2

	// setID
	// 数字の小さいIDを行にして、1次元スライス上の共起行列の位置をEdgeのIDとして利用する
	n := len(*m.PtrWords)
	wi1 := int(nodeID1) - 1
	wi2 := int(nodeID2) - 1
	if wi1 <= wi2 {
		row := wi1
		edge.ID = uint(row*n + wi2)
	} else {
		row := wi2
		edge.ID = uint(row*n + wi1)
	}

	edge.Rate = m.Matrix[edge.ID]
	return edge
}

// NodeIDの共起関係にあるNodeとそのEdgeを返す
func (m *CoMatrix) CoOccurrenceRelation(nodeID uint) (nodes []*Node, edges []*Edge) {
	nodes = []*Node{}
	edges = []*Edge{}

	if !m.ValidNodeID(nodeID) {
		return nodes, edges
	}

	subjectNodeID := nodeID
	for objectWIndex := range *m.PtrWords {
		objectNodeID := uint(objectWIndex + 1)
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
		if !m.ValidNodeID(nodeID) {
			return struct{}{}
		}
		ns, es := m.CoOccurrenceRelation(nodeID)
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

// 共起行列のIDを取得する
func (m *CoMatrix) ID() string {
	return m.Meta.GroupID
}

func (m *CoMatrix) As(metaGroupID string) {
	m.Meta.GroupID = metaGroupID
}

func (m *CoMatrix) LenNodes() int {
	return len(*m.PtrWords)
}

func (m *CoMatrix) Error() error {
	if m.err != nil {
		return MatrixError{m.err}
	}
	return nil
}

// 共起回数の正規化にDice係数を利用して共起行列を作成する
func (m *CoMatrix) matrixByDice(dwm DocWordMatrix) {
	// create Matrix by dice
	occuerences := make([]DocWordOccurances, len(*m.PtrWords))
	for windex := range *m.PtrWords {
		occuerences[windex] = dwm.Occurances(windex)
	}

	ctx := context.Background()
	frequencyCh := dwm.GenerateCoOccurrencetFrequency(ctx)
	for f := range frequencyCh {
		// Dice(Wi,Wj) = 2 x 共起回数Wi,Wj / (出現回数Wi + 出現回数Wj )
		d := float64(occuerences[f.Windex1].Occurances + occuerences[f.Windex2].Occurances)
		if d == 0 {
			m.syncSet(f.Windex1, f.Windex2, 0)
		} else {
			value := float64(2*f.Frequency) / d
			m.syncSet(f.Windex1, f.Windex2, value)
		}
	}
}

// 共起行列に共起回数をセットする
func (m *CoMatrix) syncSet(windex1, windex2 int, value float64) {
	n := len(*m.PtrWords) // 単語数
	var i int
	i = windex1*n + windex2
	m.Matrix[i] = value
	i = windex1 + windex2*n
	m.Matrix[i] = value
}

// どれほど共起関係の中心にあるかで単語の重要度を決定する。
// （固有ベクトル中心性を単語の重要度に使用する。）
func (m *CoMatrix) useVectorCentrality() error {
	// 単語数
	n := len(*m.PtrWords)
	// 対称行列化
	dence := mat.NewSymDense(n, m.Matrix)
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
		m.Priority[i] = ev.RawRowView(i)[maxEigenvalueIndex]
	}

	// 新しいpriorityをスケーリングする
	m.scalingPriority()
	return nil
}

// 優先度を0-1に標準化する
func (m *CoMatrix) scalingPriority() {

	// 重要度に基づいて単語のインデックスを降順ソート
	sort.Slice(m.Indices, func(i, j int) bool {
		return m.Priority[m.Indices[i]] > m.Priority[m.Indices[j]]
	})

	// 重要度最小値
	minVal := m.Priority[m.Indices[len(m.Indices)-1]]
	// 重要度最大値
	maxVal := m.Priority[m.Indices[0]]

	// 小数第RateFixed位未満四捨五入
	s := math.Pow10(RateFixed)
	for i, val := range m.Priority {
		p := (val - minVal) / (maxVal - minVal)
		m.Priority[i] = math.Round(p*s) / s
	}
}

func (m *CoMatrix) setProgress(p CoMatrixProgress) {
	if !m.done {
		m.progress <- p
	}
}

func (m *CoMatrix) doneProgress() {
	if !m.done {
		defer close(m.progress)
		m.setProgress(ProgressDone)
		m.done = true
	}
}

func (m *CoMatrix) doneProgressWithError(err error) {
	if !m.done {
		defer close(m.progress)
		m.setProgress(ErrDone)
		m.done = true
		m.err = err
	}
}

func (m *CoMatrix) init() {
	if m.coOccurrencetNormalization == 0 {
		m.coOccurrencetNormalization = Dice
	}
	if m.nodeRatingAlgorithm == 0 {
		m.nodeRatingAlgorithm = VectorCentrality
	}

	n := len(*m.PtrWords)
	if m.Matrix == nil {
		m.Matrix = make([]float64, n*n)
	}
	if m.Indices == nil {
		m.Indices = make([]int, n)
	}
	if m.Priority == nil {
		m.Priority = make([]float64, n)
		for i := range m.Indices {
			m.Indices[i] = i
		}
	}
}

type Node struct {
	ID       uint
	Word     string
	Rate     float64
	NumEdges uint
}

type Edge struct {
	ID, Node1ID, Node2ID uint
	Rate                 float64
}
