package matrix

import (
	"context"
	"yyyoichi/Collo-API/pkg/stream"
)

// 単語文書行列
type DocWordMatrix struct {
	matrix [][]int
	meta   MultiDocMeta
	words  []string
}

// すべての共起ペアについて共起回数を返す
func (m *DocWordMatrix) GenerateCoOccurrencetFrequency(ctx context.Context) <-chan DocWordFrequency {
	coIndexCh := m.generateCoIndex(ctx)
	return stream.Line[[2]int, DocWordFrequency](ctx, coIndexCh, func(coIndex [2]int) DocWordFrequency {
		return m.CoOccurrencetFrequency(coIndex[0], coIndex[1])
	})
}

// すべての文書における2単語[windex1][windex2]の共起回数[f]と共起文書数[c]を返す
func (m *DocWordMatrix) CoOccurrencetFrequency(windex1, windex2 int) DocWordFrequency {
	f := DocWordFrequency{
		Windex1: windex1,
		Windex2: windex2,
	}
	if len(m.words) <= windex1 || len(m.words) <= windex2 {
		return f
	}
	for _, doc := range m.matrix {
		// 共起頻度
		f.Add(doc[windex1] * doc[windex2])
	}
	return f
}

// すべての文書内での単語[windex]の出現回数[o]と出現文書数[c]を返す
func (m *DocWordMatrix) Occurances(windex int) DocWordOccurances {
	o := DocWordOccurances{
		Windex: windex,
	}
	if len(m.words) <= windex {
		return o
	}
	for _, doc := range m.matrix {
		o.Add(doc[windex])
	}
	return o
}

// 共起ペアをループする
func (m *DocWordMatrix) generateCoIndex(ctx context.Context) <-chan [2]int {
	n := len(m.words)
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

// 文書-単語行列における共起回数
type DocWordFrequency struct {
	// 単語位置
	Windex1, Windex2 int
	// 共起回数
	Frequency int
	// 出現文書数
	Count int
}

func (f *DocWordFrequency) Add(frequency int) {
	f.Frequency += frequency
	if frequency > 0 {
		f.Count++
	}
}

// 文書-単語行列におけるある単語の出現回数
type DocWordOccurances struct {
	// 単語位置
	Windex int
	// 出現回数
	Occurances int
	// 出現文書数
	Count int
}

func (o *DocWordOccurances) Add(occurances int) {
	o.Occurances += occurances
	if occurances > 0 {
		o.Count++
	}
}
