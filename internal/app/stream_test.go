package app

import (
	"context"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/internal/libs/pair"
)

// resultをスピーチ単位にして返す場合
func (cs *CollocationService) DemultiParseStream(cxt context.Context) <-chan *pair.PairResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	url := generateURL(cxt, sourceURLs)
	fetchResult := pipeURL2Fetch(cxt, url)
	parseResult := demultiFetch2Parse(cxt, fetchResult, func(fr *api.FetchResult, send func(*morpheme.ParseResult)) {
		if fr.Err != nil {
			result := &morpheme.ParseResult{}
			result.Err = fr.Err
			send(result)
			return
		}
		speech := generateSpeech(cxt, fr.GetSpeechs())
		for p := range pipeSpeech2Parse(cxt, speech, cs.speech2Parse) {
			send(p)
		}
	})
	outPairs := useFunOutParse(cxt, func() <-chan *pair.PairResult {
		return pipeParse2Pair(cxt, parseResult, cs.parse2Pair)
	})
	return useFunInPair(cxt, outPairs)
}

// resultをスピーチ単位にしたうえで、形態素解析をファンアウトする
func (cs *CollocationService) DemultiFunParseStream(cxt context.Context) <-chan *pair.PairResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	url := generateURL(cxt, sourceURLs)
	fetchResult := pipeURL2Fetch(cxt, url)
	parseResult := demultiFetch2Parse(cxt, fetchResult, func(fr *api.FetchResult, send func(*morpheme.ParseResult)) {
		if fr.Err != nil {
			result := &morpheme.ParseResult{}
			result.Err = fr.Err
			send(result)
			return
		}
		speech := generateSpeech(cxt, fr.GetSpeechs())
		f := useParseFun(cxt, func() <-chan *morpheme.ParseResult {
			return pipeSpeech2Parse(cxt, speech, cs.speech2Parse)
		})
		for p := range f {
			send(p)
		}
	})
	outPairs := useFunOutParse(cxt, func() <-chan *pair.PairResult {
		return pipeParse2Pair(cxt, parseResult, cs.parse2Pair)
	})
	return useFunInPair(cxt, outPairs)
}

// resultをスピーチ単位にしたうえで、フェッチand形態素解析をファンアウトする
func (cs *CollocationService) DemultiFunFetchParseStream(cxt context.Context) <-chan *pair.PairResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	url := generateURL(cxt, sourceURLs)
	parseResult := useParseFun(cxt, func() <-chan *morpheme.ParseResult {
		fetchResult := pipeURL2Fetch(cxt, url)
		return demultiFetch2Parse(cxt, fetchResult, func(fr *api.FetchResult, send func(*morpheme.ParseResult)) {
			if fr.Err != nil {
				result := &morpheme.ParseResult{}
				result.Err = fr.Err
				send(result)
				return
			}
			speech := generateSpeech(cxt, fr.GetSpeechs())
			f := useParseFun(cxt, func() <-chan *morpheme.ParseResult {
				return pipeSpeech2Parse(cxt, speech, cs.speech2Parse)
			})
			for p := range f {
				send(p)
			}
		})
	})
	outPairs := useFunOutParse(cxt, func() <-chan *pair.PairResult {
		return pipeParse2Pair(cxt, parseResult, cs.parse2Pair)
	})
	return useFunInPair(cxt, outPairs)
}

// resultをスピーチ単位にしたうえで、フェッチand形態素解析ペア作成までをファンアウトする
func (cs *CollocationService) DemultiFunStream(cxt context.Context) <-chan *pair.PairResult {
	opt := cs.options
	sourceURLs := api.CreateURLs(api.URLOptions{StartRecord: 1, MaximumRecords: 100, From: opt.From, Until: opt.Until, Any: opt.Any}, cs.numRecords)

	url := generateURL(cxt, sourceURLs)
	outPairs := useFunOutParse(cxt, func() <-chan *pair.PairResult {
		fetchResult := pipeURL2Fetch(cxt, url)
		parseResult := demultiFetch2Parse(cxt, fetchResult, func(fr *api.FetchResult, send func(*morpheme.ParseResult)) {
			if fr.Err != nil {
				result := &morpheme.ParseResult{}
				result.Err = fr.Err
				send(result)
				return
			}
			speech := generateSpeech(cxt, fr.GetSpeechs())
			f := useParseFun(cxt, func() <-chan *morpheme.ParseResult {
				return pipeSpeech2Parse(cxt, speech, cs.speech2Parse)
			})
			for p := range f {
				send(p)
			}
		})
		return pipeParse2Pair(cxt, parseResult, cs.parse2Pair)
	})
	return useFunInPair(cxt, outPairs)
}
