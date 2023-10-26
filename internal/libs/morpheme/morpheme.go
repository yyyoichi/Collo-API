package morpheme

import (
	"errors"
	"strings"
	"yyyoichi/Collo-API/pkg/apperror"

	"github.com/shogo82148/go-mecab"
)

var model mecab.Model
var tagger mecab.MeCab

func init() {
	var err error
	model, err = mecab.NewModel(map[string]string{})
	if err != nil {
		panic(err)
	}
	tagger, err = model.NewMeCab()
	if err != nil {
		panic(err)
	}
}

type ParseError struct {
	error
}

var (
	ErrParse = errors.New("parser error")
)

func UseMorphologicalAnalytics() (*MorphologicalAnalytics, error) {
	stop := []string{}
	for _, w := range strings.Split(stopwords, "\n") {
		stop = append(stop, strings.Trim(w, " "))
	}
	return &MorphologicalAnalytics{stop}, nil
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
	stopwords []string
}

// Result 形態素リスト
type ParseResult struct {
	Result []string
	Err    error
}

func (ma *MorphologicalAnalytics) Parse(s string) *ParseResult {
	parseResult := &ParseResult{}
	lattice, err := mecab.NewLattice()
	if err != nil {
		parseResult.Err = ParseError{apperror.WrapError(err, err.Error())}
		return parseResult
	}
	defer lattice.Destroy()

	lattice.SetSentence(s)
	if err := tagger.ParseLattice(lattice); err != nil {
		parseResult.Err = ParseError{apperror.WrapError(err, err.Error())}
		return parseResult
	}
	parseResult.Result = strings.Split(lattice.String(), "\n")
	return parseResult
}

func NewMorpheme(s string) *Morpheme {
	ss := strings.Split(s, "\t")
	data := strings.Split(ss[1], ",")
	if len(data) < 8 {
		return &Morpheme{}
	}
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
func IsEnd(s string) bool {
	return s == "EOS"
}
