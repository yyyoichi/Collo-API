package pair

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"yyyoichi/Collo-API/pkg/stream"
)

func CreateMockConfig(config Config) Config {
	ctx := context.Background()
	if config.Search.Any == "" {
		config.Search.Any = "自動車"
	}
	l, _ := time.LoadLocation("Asia/Tokyo")
	if config.Search.From.IsZero() {
		config.Search.From = time.Date(2022, 3, 1, 0, 0, 0, 0, l)
	}
	if config.Search.Until.IsZero() {
		config.Search.Until = time.Date(2022, 5, 1, 0, 0, 0, 0, l)
	}
	config.Fetcher = nil

	store := map[string][]byte{}
	// 始めの件数取得fetchをモック化
	spe := &Speech{config: config}
	spe.init()
	fr := spe.fetch(spe.createURL(1, 1))
	if fr.err != nil {
		panic(fr.err)
	}
	if body, err := json.Marshal(fr.SpeechJson); err != nil {
		panic(err)
	} else {
		store[fr.url] = body
		spe.containRecords = fr.SpeechJson.NumberOfRecords
	}

	// 取得件数分モック化
	ps := &PairStore{
		speech:   spe,
		handler:  Handler{},
		idByWord: map[string]string{},
		mu:       sync.Mutex{},
	}
	urlCh := ps.speech.generateURL(ctx)
	for fr := range stream.FunIO[string, *fetchResult](ctx, urlCh, ps.speech.fetch) {
		if fr.err != nil {
			panic(fr.err)
		}
		if body, err := json.Marshal(fr.SpeechJson); err != nil {
			panic(err)
		} else {
			store[fr.url] = body
		}
	}

	config.Fetcher = func(url string) (body []byte, err error) {
		return store[url], nil
	}
	return config
}