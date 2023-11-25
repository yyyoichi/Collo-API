package pair

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
	"yyyoichi/Collo-API/pkg/stream"
)

var tconfig Config
var thandler = Handler{
	Err: func(err error) {
		fmt.Println(err)
		panic(err)
	},
	Resp: func(resp *apiv1.ColloStreamResponse) {
	},
	Done: func() {},
}

func TestMain(m *testing.M) {
	initConfigMock()
	os.Exit(m.Run())
}
func BenchmarkCase0(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ps := NewPairStore(tconfig, thandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case0(ctx, cancel) {
		}
	}
}

func BenchmarkCase1(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ps := NewPairStore(tconfig, thandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case1(ctx, cancel) {
		}
	}
}

func BenchmarkCase2(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ps := NewPairStore(tconfig, thandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case2(ctx, cancel) {
		}
	}
}

func BenchmarkCase3(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ps := NewPairStore(tconfig, thandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case3(ctx, cancel) {
		}
	}
}

func BenchmarkCase4(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ps := NewPairStore(tconfig, thandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case4(ctx, cancel) {
		}
	}
}

func initConfigMock() {
	ctx := context.Background()
	config := Config{}
	config.Search.Any = "自動車"
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.From = time.Date(2022, 3, 1, 0, 0, 0, 0, l)
	config.Search.Until = time.Date(2022, 5, 1, 0, 0, 0, 0, l)
	config.Fetcher = nil
	store := map[string][]byte{}

	// 始めの件数取得fetchをモック化
	spe := &Speech{config: config}
	spe.init()
	url := spe.createURL(1, 1)
	fr := spe.fetch(url)
	if fr.err != nil {
		panic(fr.err)
	}
	body, err := json.Marshal(fr.SpeechJson)
	if err != nil {
		panic(err)
	}
	store[fr.url] = body

	// 取得件数分モック化
	ps := NewPairStore(config, thandler)
	urlCh := ps.speech.generateURL(ctx)
	for fr := range stream.FunIO[string, *fetchResult](ctx, urlCh, ps.speech.fetch) {
		if fr.err != nil {
			panic(err)
		}
		body, err := json.Marshal(fr.SpeechJson)
		if err != nil {
			panic(err)
		}
		store[fr.url] = body
	}

	config.Fetcher = func(url string) (body []byte, err error) {
		return store[url], nil
	}
	tconfig = config
	fmt.Println("init config")
}
