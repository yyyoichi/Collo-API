package matrix

import (
	"strconv"
	"strings"
)

type (
	Config struct {
		ReduceThreshold            float64                    // しきい値。上位Threshold%の単語を使用する
		MinNodes                   int                        // 最小単語数
		MaxNodes                   int                        // 最大単語数
		NodeRatingAlgorithm        NodeRatingAlgorithm        // ノードの中心性を導くアルゴリズム
		CoOccurrencetNormalization CoOccurrencetNormalization // 共起回数を正規化するアルゴリズム
		GroupingFuncType           GroupingFuncType
		AtGroupID                  string // 指定グループ識別子の共起行列を取得する
	}

	NodeRatingAlgorithm        uint // ノードの中心性を導くアルゴリズム
	CoOccurrencetNormalization uint // 共起回数を正規化するアルゴリズム
	// Documentからグループ識別子を取り出す関数。複数Documentをグルーピングするために使用する
	GroupingFuncType uint
)

const (
	VectorCentrality NodeRatingAlgorithm = 1

	Dice CoOccurrencetNormalization = 1

	PickByKey   GroupingFuncType = 1
	PickByMonth GroupingFuncType = 2
	PickAsTotal GroupingFuncType = 3

	strReduceThreshold            = "r!:"
	strMinNodes                   = "m!:"
	strMaxNodes                   = "x!:"
	strNodeRatingAlgorithm        = "a!:"
	strCoOccurrencetNormalization = "c!:"
	strGroupingFuncType           = "p!:"
)

func (c *Config) ToString() string {
	c.init()
	var buf strings.Builder
	buf.WriteString(strReduceThreshold)
	buf.WriteString(strconv.Itoa(int(c.ReduceThreshold * 100)))

	buf.WriteString(strMinNodes)
	buf.WriteString(strconv.Itoa(c.MinNodes))

	buf.WriteString(strMaxNodes)
	buf.WriteString(strconv.Itoa(c.MaxNodes))

	buf.WriteString(strNodeRatingAlgorithm)
	buf.WriteString(strconv.Itoa(int(c.NodeRatingAlgorithm)))

	buf.WriteString(strNodeRatingAlgorithm)
	buf.WriteString(strconv.Itoa(int(c.CoOccurrencetNormalization)))

	buf.WriteString(strGroupingFuncType)
	buf.WriteString(strconv.Itoa(int(c.GroupingFuncType)))

	return buf.String()
}

func (c *Config) init() {
	if c.ReduceThreshold <= 0 || 1 < c.ReduceThreshold {
		c.ReduceThreshold = 0.1
	}
	if c.MinNodes == 0 {
		c.MinNodes = 100
	}
	if c.MaxNodes == 0 {
		c.MaxNodes = 300
	}
	if c.MinNodes > c.MaxNodes {
		c.MinNodes = c.MaxNodes / 2
	}
	if c.NodeRatingAlgorithm == 0 {
		c.NodeRatingAlgorithm = VectorCentrality
	}
	if c.CoOccurrencetNormalization == 0 {
		c.CoOccurrencetNormalization = Dice
	}
	if c.GroupingFuncType == 0 {
		c.GroupingFuncType = PickByKey
	}
}
