package morpheme

type Config struct {
	Includes  []PartOfSpeechType // 対象品詞
	StopWords []string           // 追加するストップワード
}

func (c *Config) init() {
	if len(c.Includes) == 0 {
		c.Includes = append(c.Includes,
			Noun, PlaceName,
			Adjective,
			AdjectiveVerb,
			Verb,
		)
	}
}

// 参照: https://repository.ninjal.ac.jp/records/2872
type PartOfSpeechType int

const Noun PartOfSpeechType = 101       // 普通名詞
const PersonName PartOfSpeechType = 111 // 固有名詞人名
const PlaceName PartOfSpeechType = 112  // 固有名詞地名
const Number PartOfSpeechType = 121     // 数詞

const Adjective PartOfSpeechType = 201     // 形容詞
const AdjectiveVerb PartOfSpeechType = 301 // 形容動詞 (形状詞)

const Verb PartOfSpeechType = 401 // 動詞
