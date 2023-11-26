package pair

import (
	"log"
	"strings"

	"github.com/shogo82148/go-mecab"
)

var MAnalytics *MorphologicalAnalytics = useMorphologicalAnalytics()

func useMorphologicalAnalytics() *MorphologicalAnalytics {
	stops := []string{}
	for _, w := range strings.Split(stopwords, "\n") {
		stops = append(stops, strings.Trim(w, " "))
	}
	ma := &MorphologicalAnalytics{stopwords: stops}
	var err error
	ma.model, err = mecab.NewModel(map[string]string{})
	if err != nil {
		panic(err)
	}

	ma.tagger, err = ma.model.NewMeCab()
	if err != nil {
		panic(err)
	}

	return ma
}

type MorphologicalAnalytics struct {
	stopwords []string
	model     mecab.Model
	tagger    mecab.MeCab
}

func (a *MorphologicalAnalytics) isStopword(lexeme string) bool {
	for _, w := range a.stopwords {
		if w == lexeme {
			return true
		}
	}
	return false
}

func (a *MorphologicalAnalytics) Parse(speech string) *ParseResult {
	lattice, err := mecab.NewLattice()
	if err != nil {
		return &ParseResult{err: err}
	}
	defer lattice.Destroy()

	s := strings.ReplaceAll(speech, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	lattice.SetSentence(s)
	if err := a.tagger.ParseLattice(lattice); err != nil {
		return &ParseResult{err: err}
	}
	return &ParseResult{Result: strings.Split(lattice.String(), "\n")}
}

// Result 形態素リスト
type ParseResult struct {
	Result []string
	err    error
}

func (pr *ParseResult) GetNouns() []string {
	nouns := []string{}
	for _, line := range pr.Result {
		m := NewMorpheme(line)
		if m.IsTraget() && !MAnalytics.isStopword(m.Lexeme) {
			nouns = append(nouns, m.Lexeme)
		}
	}
	return nouns
}

func (pr *ParseResult) Error() error {
	return pr.err
}

func NewMorpheme(s string) *Morpheme {
	if s == "" || s == "EOS" {
		return &Morpheme{EOS: true}
	}
	ss := strings.Split(s, "\t")
	if len(ss) < 2 {
		log.Printf("warn: '%s' has no '\\t'", s)
		return &Morpheme{EOS: true}
	}
	data := strings.Split(ss[1], ",")
	if len(data) < 8 {
		if len(data) != 6 {
			log.Printf("warn: '%s' has no 8 ',', got='%d'", s, len(data))
			return &Morpheme{EOS: true}
		}
		// data長合わせ
		data = append(data, data[4], data[5])
	}
	return &Morpheme{
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
	EOS                  bool   // 終了
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

func (m *Morpheme) IsTraget() bool {
	return !m.isEnd() && m.isNoun() && !m.isAsterisk()
}

func (m *Morpheme) isNoun() bool {
	return m.PartOfSpeechDetails1 == "普通名詞" || m.PartOfSpeechDetails1 == "固有名詞"
}
func (m *Morpheme) isAsterisk() bool {
	return m.Lexeme == "*"
}
func (m *Morpheme) isEnd() bool {
	return m.EOS
}
