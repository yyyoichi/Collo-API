package pair

import (
	"context"
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
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case0(context.Background()) {
		}
	}
}

func BenchmarkCase1(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case1(context.Background()) {
		}
	}
}

func BenchmarkCase2(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case2(context.Background()) {
		}
	}
}

func BenchmarkCase3(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case3(context.Background()) {
		}
	}
}

func BenchmarkCase4(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case4(context.Background()) {
		}
	}
}

// ストリームなし
func (ps *PairStore) stream_case0(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	ch := make(chan *apiv1.ColloStreamResponse)
	go func(ps *PairStore) {
		defer close(ch)
		for url := range ps.speech.GenerateURL(ctx) {
			FetchResult := ps.speech.Fetch(url)
			if FetchResult.err != nil {
				ps.handleError(FetchResult.Error())
				break
			}
			chunk := ps.newPairChunk()
			for _, speech := range FetchResult.GetSpeechs() {
				parseResult := MAnalytics.Parse(speech)
				if parseResult.err != nil {
					ps.handleError(parseResult.Error())
					break
				}
				nouns := parseResult.GetNouns()
				chunk.set(nouns)
			}
			select {
			case <-ctx.Done():
				return
			default:
				ch <- chunk.ConvResp()
			}
		}
	}(ps)
	return ch
}

// 全てを順にパイプ
func (ps *PairStore) stream_case1(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	urlCh := ps.speech.GenerateURL(ctx)
	fetchResultCh := stream.Line[string, *FetchResult](ctx, urlCh, ps.speech.Fetch)
	speechCh := stream.Demulti[*FetchResult, string](ctx, fetchResultCh, func(fr *FetchResult) []string {
		if fr.err != nil {
			ps.handleError(fr.Error())
		}
		return fr.GetSpeechs()
	})
	nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
		pr := MAnalytics.Parse(s)
		if pr.err != nil {
			ps.handleError(pr.Error())
		}
		return pr.GetNouns()
	})
	return stream.Line[[]string, *apiv1.ColloStreamResponse](ctx, nounsCh, func(s []string) *apiv1.ColloStreamResponse {
		c := ps.newPairChunk()
		c.set(s)
		return c.ConvResp()
	})
}

// fetchから丸々funアウトする
func (ps *PairStore) stream_case2(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	urlCh := ps.speech.GenerateURL(ctx)
	return stream.FunIO[string, *apiv1.ColloStreamResponse](ctx, urlCh, func(url string) *apiv1.ColloStreamResponse {
		FetchResult := ps.speech.Fetch(url)
		speechCh := FetchResult.GenerateSpeech(ctx)
		nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
			pr := MAnalytics.Parse(s)
			if pr.err != nil {
				ps.handleError(pr.Error())
			}
			return pr.GetNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c.ConvResp()
	})
}

// fetchから丸々funアウト, 形態素解析前にもfunアウトする
func (ps *PairStore) stream_case4(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	urlCh := ps.speech.GenerateURL(ctx)
	return stream.FunIO[string, *apiv1.ColloStreamResponse](ctx, urlCh, func(url string) *apiv1.ColloStreamResponse {
		FetchResult := ps.speech.Fetch(url)
		speechCh := FetchResult.GenerateSpeech(ctx)
		nounsCh := stream.FunIO[string, []string](ctx, speechCh, func(s string) []string {
			pr := MAnalytics.Parse(s)
			if pr.err != nil {
				ps.handleError(pr.Error())
			}
			return pr.GetNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c.ConvResp()
	})
}

func initConfigMock() {
	config := Config{}
	config.Search.Any = "自動車"
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.From = time.Date(2022, 3, 1, 0, 0, 0, 0, l)
	config.Search.Until = time.Date(2022, 5, 1, 0, 0, 0, 0, l)
	tconfig = CreateMockConfig(config)
	fmt.Println("init config")
}
