package matrix

import (
	"context"
	"fmt"
	"sync"

	"gonum.org/v1/gonum/mat"
)

// 共起座標
type CoIndex [2]int

func (i *CoIndex) I1() int { return i[0] }
func (i *CoIndex) I2() int { return i[1] }

type CoMatrix struct {
	// 行列の長さ
	n int
	// 行列
	matrix [][]float64

	mu sync.Mutex
}

// 共起行列を生成する
// [totalWord]出現する単語数
func NewCoMatrix(totalWrod int) *CoMatrix {
	cm := &CoMatrix{
		matrix: make([][]float64, totalWrod),
		n:      totalWrod,
	}
	for i := range cm.matrix {
		cm.matrix[i] = make([]float64, totalWrod)
	}
	return cm
}

// 共起行列の共起ペアindexを返す
func (cm *CoMatrix) GenerateCoIndex(ctx context.Context) <-chan CoIndex {
	n := cm.n
	ch := make(chan CoIndex)
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- CoIndex{i, j}
				}
			}
		}
	}()

	return ch
}

// 共起頻度をセットする
func (cm *CoMatrix) Set(coIndex CoIndex, frequency float64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	i1 := coIndex.I1()
	i2 := coIndex.I2()
	cm.matrix[i1][i2] = frequency
	cm.matrix[i2][i1] = frequency
}

func (cm *CoMatrix) U() {
	data := []float64{}
	for _, row := range cm.matrix {
		data = append(data, row...)
	}
	dence := mat.NewDense(cm.n, cm.n, data)

	var svd mat.SVD
	svd.Factorize(dence, mat.SVDFull)

	var leftSingular mat.Dense
	svd.UTo(&leftSingular)

	fmt.Printf("左特異ベクトル:\n%v\n", mat.Formatted(&leftSingular, mat.Prefix(""), mat.Squeeze()))
}
