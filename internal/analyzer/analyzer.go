package analyzer

import (
	"bufio"
	"context"
	"iter"
	"slices"
	"strings"
	"yyyoichi/Collo-API/pkg/stream"

	"github.com/shogo82148/go-mecab"
)

var (
	// ストップワード initialize in init()
	stopwords []string

	// 解析器 initialize in init()
	tagger mecab.MeCab

	// 並列処理に保持が必要。使用しない。
	_model mecab.Model
)

type Analyzer struct {
	// 語彙素がストップワードか
	isStopword   func(lemma string) bool
	isTargetType func(v PartOfSpeechType) bool
}

func New(config Config) Analyzer {
	config.init()
	var stops = make([]string, len(stopwords)+len(config.StopWords))
	_ = copy(stops, stopwords)
	_ = copy(stops[len(stopwords):], config.StopWords)

	var a Analyzer
	a.isStopword = func(lemma string) bool {
		return slices.Contains(stops, lemma)
	}
	a.isTargetType = func(v PartOfSpeechType) bool {
		return slices.Contains(config.Includes, v)
	}
	return a
}

// 語彙素を返す。
func (a *Analyzer) Parse(s string) (iter.Seq[string], error) {
	lattice, err := mecab.NewLattice()
	if err != nil {
		return nil, err
	}
	defer lattice.Destroy()

	lattice.SetSentence(replacer.Replace(s))
	if err := tagger.ParseLattice(lattice); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(lattice.String()))
	return func(yield func(string) bool) {
		for scanner.Scan() {
			m := newMorpheme(scanner.Text())
			if m.isEnd() {
				break
			}
			if m.isAsterisk() || a.isStopword(m.Lemma) || !a.isTargetType(m.Type()) {
				continue
			}
			if ok := yield(m.Lemma); !ok {
				return
			}
		}
	}, nil
}

var replacer = strings.NewReplacer(
	"\u3000", "",
	"\r\n", "",
	"\r", "",
	"\n", "",
	"\t", "",
	" ", "",
)

// TODO 以下削除

func Analysis(sentence string) AnalysisResult {
	lattice, err := mecab.NewLattice()
	if err != nil {
		return AnalysisResult{err: err}
	}
	defer lattice.Destroy()

	// 正規化
	s := strings.ReplaceAll(sentence, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")

	lattice.SetSentence(s)
	if err := tagger.ParseLattice(lattice); err != nil {
		return AnalysisResult{err: err}
	}
	return AnalysisResult{
		result: strings.Split(lattice.String(), "\n"),
	}
}

type AnalysisError struct{ error }

type AnalysisResult struct {
	result []string
	err    error
}

func (ar *AnalysisResult) Error() error {
	if ar.err != nil {
		return AnalysisError{ar.err}
	}
	return nil
}
func (ar *AnalysisResult) Get(config Config) []string {
	config.init()

	sws := append(stopwords, config.StopWords...)
	isStopword := func(m *Morpheme) bool {
		for _, w := range sws {
			if w == m.Lemma {
				return true
			}
		}
		return false
	}

	ctx := context.Background()
	morphemeCh := ar.GenerateMorpheme(ctx)
	lemmaCh := stream.Line[Morpheme, string](ctx, morphemeCh, func(m Morpheme) string {
		if m.isEnd() || m.isAsterisk() {
			return ""
		}
		if isStopword(&m) {
			return ""
		}

		// config対象の単語が調査
		for _, t := range config.Includes {
			if m.TypeIs(t) {
				// 一致すれば語彙素を返す
				return m.Lemma
			}
		}
		return ""
	})

	result := []string{}
	for lemma := range lemmaCh {
		if lemma != "" {
			result = append(result, lemma)
		}
	}
	return result
}

func (ar *AnalysisResult) GetAt(i int) Morpheme {
	if len(ar.result) <= i {
		return Morpheme{EOS: true}
	}
	return newMorpheme(ar.result[i])
}

// 解析結果の形態素をストリームで返す
func (ar *AnalysisResult) GenerateMorpheme(ctx context.Context) <-chan Morpheme {
	resultCh := stream.Generator[string](ctx, ar.result...)
	return stream.Line[string, Morpheme](ctx, resultCh, newMorpheme)
}

func init() {
	// init stopwords
	words := strings.Split(getStopwords(), "\n")
	stopwords = make([]string, len(words))
	for i, w := range words {
		stopwords[i] = strings.Trim(w, " ")
	}

	// init mecab
	var err error
	_model, err = mecab.NewModel(map[string]string{})
	if err != nil {
		panic(err)
	}
	tagger, err = _model.NewMeCab()
	if err != nil {
		panic(err)
	}
}
