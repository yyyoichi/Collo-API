package app

import (
	"context"
	"strings"
	"time"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/collocation"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/pkg/stream/fun"
	"yyyoichi/Collo-API/pkg/stream/pipe"
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
	collo := collocation.NewCollocation()

	url := api.CreateURL(api.URLOptions{StartRecord: 1, MaximumRecords: 1, From: opt.From, Until: opt.Until, Any: opt.Any})
	result := api.Fetch(url)
	if result.Err != nil {
		return nil, result.Err
	}
	return &CollocationService{ma, collo, result.SpeachJson.NumberOfRecords, opt}, nil
}

type CollocationService struct {
	*morpheme.MorphologicalAnalytics
	c          *collocation.Collocation
	numRecords int
	options    CollocationServiceOptions
}

// 重複しない共起ペアデータをチャネルで返します。
func (cs *CollocationService) Stream(cxt context.Context) <-chan *collocation.CollocationResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	// start pipeline
	// 1. urlがパイプされます。
	// 2. ファンアウトしてurlをfetchしmecabで取得された発言をすべて形態素解析します。
	// 3. ファンインした形態素解析結果を発言ごとに共起ペアを生成して結果を返します。
	url := pipe.Generator[string](cxt, sourceURLs...)
	parseResults := fun.Out[*morpheme.ParseResult](cxt, func() <-chan *morpheme.ParseResult {
		return cs.parse(cxt, cs.fetch(cxt, url))
	})
	return cs.pairs(cxt, fun.In[*morpheme.ParseResult](cxt, parseResults...))
}

func (cs *CollocationService) fetch(cxt context.Context, url <-chan string) <-chan *api.FetchResult {
	return pipe.Line[string, *api.FetchResult](cxt, url, api.Fetch)
}

// 形態素解析した結果を返す
func (cs *CollocationService) parse(cxt context.Context, fetchResult <-chan *api.FetchResult) <-chan *morpheme.ParseResult {
	return pipe.Line[*api.FetchResult, *morpheme.ParseResult](cxt, fetchResult, func(fr *api.FetchResult) *morpheme.ParseResult {
		if fr.Err != nil {
			return &morpheme.ParseResult{Err: fr.Err}
		}
		// 発言を||で区切りまとめて形態素にする
		speach := strings.Join(fr.GetSpeachs(), "||")
		return cs.Parse(speach)
	})
}

// 形態素解析結果から共起ペアを返す
func (cs *CollocationService) pairs(cxt context.Context, parseResult <-chan *morpheme.ParseResult) <-chan *collocation.CollocationResult {
	// 形態素リスト[morphemes]をスピーチ単位を1チャンクとして名詞(語彙素)リストと、探査数を返す。
	getNounsInSpeachChunk := func(morphemes []string) pipe.ChunkFnResp[[]string] {
		lexemes := []string{}
		i := 0
		for {
			m := morpheme.NewMorpheme(morphemes[i])
			if m.IsEnd() {
				// 全てを探査済みと見なす
				return pipe.ChunkFnResp[[]string]{Out: lexemes, Len: len(morphemes)}
			}
			// ||のときチャンク終了。結果を返す
			if m.IsPipe() && morpheme.NewMorpheme(morphemes[i+1]).IsPipe() {
				return pipe.ChunkFnResp[[]string]{Out: lexemes, Len: i + 2}
			}

			isTarget := m.IsNoun() && !m.IsAsterisk() && !cs.IsStopword(m.Lexeme)
			if isTarget {
				lexemes = append(lexemes, m.Lexeme)
			}
			i++
		}
	}
	getCollocationResult := func(nouns []string) *collocation.CollocationResult { return cs.c.Get(nouns) }
	return pipe.Line[*morpheme.ParseResult, *collocation.CollocationResult](cxt, parseResult, func(pr *morpheme.ParseResult) *collocation.CollocationResult {
		if pr.Err != nil {
			return &collocation.CollocationResult{Err: pr.Err}
		}
		results := &collocation.CollocationResult{}
		for result := range pipe.Line[[]string, *collocation.CollocationResult](cxt, pipe.Chunk[string, []string](cxt, getNounsInSpeachChunk, pr.Result...), getCollocationResult) {
			for id, word := range result.WordByID {
				results.WordByID[id] = word
			}
			results.Pairs = append(results.Pairs, result.Pairs...)
		}
		return results
	})
}
