package app

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/internal/libs/pair"
)

type CollocationServiceOptions struct {
	Any   string
	From  time.Time
	Until time.Time
}

func NewCollocationService(opt CollocationServiceOptions) (*CollocationService, error) {
	ma, err := morpheme.UseMorphologicalAnalytics()
	if err != nil {
		return nil, err
	}
	pair := pair.NewPair()

	url := api.CreateURL(api.URLOptions{StartRecord: 1, MaximumRecords: 1, From: opt.From, Until: opt.Until, Any: opt.Any})
	result := api.Fetch(url)
	if result.Err != nil {
		return nil, result.Err
	}
	log.Printf("Length: %d\n", result.SpeechJson.NumberOfRecords)
	return &CollocationService{ma, pair, result.SpeechJson.NumberOfRecords, opt}, nil
}

type CollocationService struct {
	*morpheme.MorphologicalAnalytics
	pair       *pair.Pair
	numRecords int
	options    CollocationServiceOptions
}

// 重複しない共起ペアデータをチャネルで返します。
func (cs *CollocationService) Stream(cxt context.Context) <-chan *pair.PairResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	log.Printf("Start Stream")
	// start pipeline
	// 1. urlがパイプされます。
	// 2. urlをfetchしmecabで取得された発言をすべて形態素解析します。
	// 3. 形態素解析結果を発言ごとに共起ペアを生成して結果を返します。
	url := generateURL(cxt, sourceURLs)
	fetchResult := pipeURL2Fetch(cxt, url)
	return pipeFetch2Pair(cxt, fetchResult, func(fr *api.FetchResult) *pair.PairResult {
		result := pair.NewPairResult()
		if fr.Err != nil {
			result.Err = fr.Err
			return result
		}
		speech := generateSpeech(cxt, fr.GetSpeechs())
		parseResult := pipeSpeech2Parse(cxt, speech, cs.speech2Parse)
		outPairs := useFunOutParse(cxt, func() <-chan *pair.PairResult {
			return pipeParse2Pair(cxt, parseResult, cs.parse2Pair)
		})

		for p := range useFunInPair(cxt, outPairs) {
			if p.Err != nil {
				result.Err = p.Err
			}
			result.Concat(p)
		}
		return result
	})
}

// validation and parse
func (cs *CollocationService) speech2Parse(s string) *morpheme.ParseResult {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return cs.Parse(s)
}

type nouns struct {
	d  []string
	mu sync.Mutex
}

func (n *nouns) add(lexeme string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.d = append(n.d, lexeme)
}

func (cs *CollocationService) parse2Pair(parseResult *morpheme.ParseResult) *pair.PairResult {
	if parseResult.Err != nil {
		return &pair.PairResult{Err: parseResult.Err}
	}
	nouns := nouns{d: []string{}, mu: sync.Mutex{}}
	for _, line := range parseResult.Result {
		if morpheme.IsEnd(line) {
			break
		}
		m := morpheme.NewMorpheme(line)
		isTarget := m.IsNoun() && !m.IsAsterisk() && !cs.IsStopword(m.Lexeme)
		if isTarget {
			nouns.add(m.Lexeme)
		}
	}
	return cs.pair.Get(nouns.d)
}
