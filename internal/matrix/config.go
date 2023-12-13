package matrix

type WordRatingAlgorithm uint

const VectorCentrality WordRatingAlgorithm = 1

type CoOccurrencetNormalization uint

const Dice CoOccurrencetNormalization = 1

type Config struct {
	ReduceThreshold            float64                    // しきい値。上位Threshold%の単語を使用する
	MinNodes                   int                        // 最小単語数
	MinNodeImportanceScale     float64                    // ノード重要度のスケーリング最小値
	MaxNodeImportanceScale     float64                    // ノード重要度のスケーリング最大値
	NodeRatingAlgorithm        WordRatingAlgorithm        // 重要語判定アルゴリズム
	CoOccurrencetNormalization CoOccurrencetNormalization // 共起行列の正規化アルゴリズム
}
