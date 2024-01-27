package analyzer

import (
	"fmt"
	"strings"
)

type (
	Config struct {
		Includes  []PartOfSpeechType // 対象品詞
		StopWords []string           // 追加するストップワード
	}

	// 参照: https://repository.ninjal.ac.jp/records/2872
	PartOfSpeechType uint16
)

const (
	Noun       PartOfSpeechType = 101 // 普通名詞
	PersonName PartOfSpeechType = 111 // 固有名詞人名
	PlaceName  PartOfSpeechType = 112 // 固有名詞地名
	Number     PartOfSpeechType = 121 // 数詞

	Adjective     PartOfSpeechType = 201 // 形容詞
	AdjectiveVerb PartOfSpeechType = 301 // 形容動詞 (形状詞)

	Verb PartOfSpeechType = 401 // 動詞

	strIncludes  = "i!:"
	strStopWords = "s!:"
)

func (c *Config) ToString() string {
	c.init()

	var buf strings.Builder
	buf.WriteString(strIncludes)
	// buf.WriteString(strings.Join(c.Includes, "-"))
	for i, pos := range c.Includes {
		if i > 0 {
			buf.WriteByte('-')
		}
		buf.WriteString(fmt.Sprint(pos))
	}
	buf.WriteString(strStopWords)
	buf.WriteString(strings.Join(c.StopWords, "-"))
	return buf.String()
}

func (c *Config) init() {
	if len(c.Includes) == 0 {
		c.Includes = []PartOfSpeechType{
			Noun, PlaceName,
			Adjective,
			AdjectiveVerb,
			Verb,
		}
	}
}
