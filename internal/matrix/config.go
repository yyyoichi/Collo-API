package matrix

type WordRatingAlgorithm uint

const VectorCentrality WordRatingAlgorithm = 1

type Config struct {
	Threshold           float64             // しきい値。上位Threshold%の値を使用する
	MinScale            float64             // 重要度のスケーリング最小値
	MaxScale            float64             // 重要度のスケーリング最大値
	WordRatingAlgorithm WordRatingAlgorithm // 重要語判定アルゴリズム
}
