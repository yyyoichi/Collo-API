package matrix

type NodeRatingAlgorithm uint

const VectorCentrality NodeRatingAlgorithm = 1

type CoOccurrencetNormalization uint

const Dice CoOccurrencetNormalization = 1

type Config struct {
	ReduceThreshold            float64                    // しきい値。上位Threshold%の単語を使用する
	MinNodes                   int                        // 最小単語数
	NodeRatingAlgorithm        NodeRatingAlgorithm        // 重要語判定アルゴリズム
	CoOccurrencetNormalization CoOccurrencetNormalization // 共起行列の正規化アルゴリズム
}
