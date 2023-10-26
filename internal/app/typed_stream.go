package app

import (
	"context"
	"yyyoichi/Collo-API/internal/libs/api"
	"yyyoichi/Collo-API/internal/libs/morpheme"
	"yyyoichi/Collo-API/internal/libs/pair"
	"yyyoichi/Collo-API/pkg/stream/fun"
	"yyyoichi/Collo-API/pkg/stream/pipe"
)

func generateURL(cxt context.Context, souceURLs []string) <-chan string {
	return pipe.Generator[string](cxt, souceURLs...)
}

func pipeURL2Fetch(cxt context.Context, url <-chan string) <-chan *api.FetchResult {
	return pipe.Line[string, *api.FetchResult](cxt, url, api.Fetch)
}

func pipeFetch2Pair(cxt context.Context, fetchResult <-chan *api.FetchResult, fn func(*api.FetchResult) *pair.PairResult) <-chan *pair.PairResult {
	return pipe.Line[*api.FetchResult, *pair.PairResult](cxt, fetchResult, fn)
}

func demultiFetch2Parse(cxt context.Context, fetchResult <-chan *api.FetchResult, fn func(*api.FetchResult, func(*morpheme.ParseResult))) <-chan *morpheme.ParseResult {
	return pipe.Demulti[*api.FetchResult, *morpheme.ParseResult](cxt, fetchResult, fn)
}

func generateSpeech(cxt context.Context, speechs []string) <-chan string {
	return pipe.Generator[string](cxt, speechs...)
}

func pipeSpeech2Parse(cxt context.Context, speech <-chan string, fn func(string) *morpheme.ParseResult) <-chan *morpheme.ParseResult {
	return pipe.Line[string, *morpheme.ParseResult](cxt, speech, fn)
}

func pipeParse2Pair(cxt context.Context, parseResult <-chan *morpheme.ParseResult, fn func(*morpheme.ParseResult) *pair.PairResult) <-chan *pair.PairResult {
	return pipe.Line[*morpheme.ParseResult, *pair.PairResult](cxt, parseResult, fn)
}

func useParseFun(cxt context.Context, fn func() <-chan *morpheme.ParseResult) <-chan *morpheme.ParseResult {
	out := fun.Out[*morpheme.ParseResult](cxt, fn)
	return fun.In[*morpheme.ParseResult](cxt, out...)
}

func useFun(cxt context.Context, fn func() <-chan *pair.PairResult) <-chan *pair.PairResult {
	out := fun.Out[*pair.PairResult](cxt, fn)
	return fun.In[*pair.PairResult](cxt, out...)
}

func useFunOutParse(cxt context.Context, fn func() <-chan *pair.PairResult) []<-chan *pair.PairResult {
	return fun.Out[*pair.PairResult](cxt, fn)
}

func useFunInPair(cxt context.Context, chanels []<-chan *pair.PairResult) <-chan *pair.PairResult {
	return fun.In[*pair.PairResult](cxt, chanels...)
}
