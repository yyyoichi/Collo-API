package analyzer

import (
	"log"
	"strings"
)

func newMorpheme(s string) Morpheme {
	if s == "" || s == "EOS" {
		return Morpheme{EOS: true}
	}
	ss := strings.Split(s, "\t")
	if len(ss) < 2 {
		log.Printf("warn: '%s' has no '\\t'", s)
		return Morpheme{EOS: true}
	}
	data := strings.Split(ss[1], ",")
	if len(data) < 8 {
		if len(data) != 6 {
			log.Printf("warn: '%s' has no 8 ',', got='%d'", s, len(data))
			return Morpheme{EOS: true}
		}
		// data長合わせ
		data = append(data, data[4], data[5])
	}
	return Morpheme{
		false,
		ss[0],
		data[0],
		data[1],
		data[2],
		data[3],
		data[4],
		data[5],
		data[6],
		data[7],
	}
}

type Morpheme struct {
	EOS     bool   // 終了
	Surface string // 表層形
	Pos     string // 品詞
	Pos1    string // 品詞細分類1
	Pos2    string // 品詞細分類2
	Pos3    string // 品詞細分類3
	CType   string // 活用型 (一段など)
	CForm   string // 活用形 (基本形など)
	LForm   string // 語彙素読み
	Lemma   string // 語彙素
}

func (m *Morpheme) TypeIs(t PartOfSpeechType) bool {
	switch t {
	case Noun:
		// 副詞可能以外の普通名詞
		return m.Pos1 == "普通名詞" && m.Pos2 != "副詞可能"
	case PersonName:
		// 人名の固有名詞
		return m.Pos1 == "固有名詞" && m.Pos2 == "人名"
	case PlaceName:
		// 地名の固有名詞
		return m.Pos1 == "固有名詞" && m.Pos2 == "地名"
	case Number:
		// 数詞
		return m.Pos1 == "数詞"
	case Adjective:
		// 形容詞
		return m.Pos == "形容詞"
	case AdjectiveVerb:
		// 形状詞-助動詞語幹以外の形状詞(形容動詞)
		return m.Pos == "形状詞" && m.Pos1 != "助動詞語幹"
	case Verb:
		// 動詞
		return m.Pos == "動詞"
	default:
		return false
	}
}

func (m *Morpheme) isAsterisk() bool {
	return m.Lemma == "*"
}
func (m *Morpheme) isEnd() bool {
	return m.EOS
}
