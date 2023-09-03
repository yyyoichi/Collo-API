package morpheme

import (
	"errors"
	"strings"
	"yyyoichi/Collo-API/pkg/apperror"

	"github.com/shogo82148/go-mecab"
)

type ParseError struct {
	error
}

var (
	ErrParse = errors.New("parser error")
)

func UseMorphologicalAnalytics() (*MorphologicalAnalytics, error) {
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		return nil, err
	}
	stop := []string{}
	for _, w := range strings.Split(stopwords, "\n") {
		stop = append(stop, strings.Trim(w, " "))
	}
	return &MorphologicalAnalytics{tagger, stop}, err
}

func (ma *MorphologicalAnalytics) IsStopword(lexeme string) bool {
	for _, w := range ma.stopwords {
		if w == lexeme {
			return true
		}
	}
	return false
}

type MorphologicalAnalytics struct {
	tagger    mecab.MeCab
	stopwords []string
}

func (ma *MorphologicalAnalytics) Destory() { ma.tagger.Destroy() }

// Result 形態素リスト
type ParseResult struct {
	Result []string
	Err    error
}

func (ma *MorphologicalAnalytics) Parse(s string) *ParseResult {
	parseResult := &ParseResult{}
	p, err := ma.tagger.Parse(s)
	if err != nil {
		parseResult.Err = ParseError{apperror.WrapError(err, err.Error())}
		return parseResult
	}
	parseResult.Result = strings.Split(p, "\n")
	return parseResult
}

func NewMorpheme(s string) *Morpheme {
	ss := strings.Split(s, "\t")
	data := strings.Split(ss[1], ",")
	return &Morpheme{
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
	Surface              string // 表層形
	PartOfSpeech         string // 品詞
	PartOfSpeechDetails1 string // 品詞細分類1
	PartOfSpeechDetails2 string // 品詞細分類2
	PartOfSpeechDetails3 string // 品詞細分類3
	Inflection           string // 活用型 (一段など)
	InflectedForm        string // 活用形 (基本形など)
	Pronunciation        string // 語彙素読み
	Lexeme               string // 語彙素
}

func (m *Morpheme) IsNoun() bool {
	return m.PartOfSpeechDetails1 == "普通名詞" || m.PartOfSpeechDetails1 == "固有名詞"
}
func (m *Morpheme) IsAsterisk() bool {
	return m.Lexeme == "*"
}
func (m *Morpheme) IsPipe() bool {
	return m.PartOfSpeech == "補助記号" && m.Surface == "|"
}
func (m *Morpheme) IsEnd() bool {
	return m.Surface == "EOS"
}
