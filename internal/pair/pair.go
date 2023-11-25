package pair

import (
	"context"
	"fmt"
	"sync"
	"time"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
	"yyyoichi/Collo-API/pkg/stream"
)

type FetchError struct{ error }
type ParseError struct{ error }
type TimeoutError struct{ error }

type PairStore struct {
	speech  *Speech
	handler Handler

	idByWord map[string]string
	mu       sync.Mutex
}

func NewPairStore(config Config, handler Handler) (*PairStore, error) {
	var err error
	ps := &PairStore{
		handler:  handler,
		idByWord: map[string]string{},
		mu:       sync.Mutex{},
	}

	ps.speech, err = NewSpeech(config)
	if err != nil {
		return nil, FetchError{err}
	}
	return ps, err
}

func (ps *PairStore) Stream(ctx context.Context) {
	done := make(chan struct{}, 1)
	defer close(done)

	go func() {
		chunkCh := ps.stream_case3(ctx)
		stream.Line[*apiv1.ColloStreamResponse, interface{}](ctx, chunkCh, func(resp *apiv1.ColloStreamResponse) interface{} {
			ps.handler.Resp(resp)
			return nil
		})
		done <- struct{}{}
	}()

	timelimit := time.Second * 60
	select {
	case <-time.After(timelimit):
		ps.handleError(TimeoutError{})
	case <-done:
		ps.handler.Done()
	case <-ctx.Done():
	}
}

func (ps *PairStore) handleError(err error) {
	ps.handler.Err(err)
}

type Handler struct {
	Resp func(resp *apiv1.ColloStreamResponse)
	Err  func(err error)
	Done func()
}

// ストリームなし
func (ps *PairStore) stream_case0(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	ch := make(chan *apiv1.ColloStreamResponse)
	go func(ps *PairStore) {
		defer close(ch)
		for url := range ps.speech.generateURL(ctx) {
			fetchResult := ps.speech.fetch(url)
			if fetchResult.err != nil {
				ps.handleError(fetchResult.Error())
				break
			}
			chunk := ps.newPairChunk()
			for _, speech := range fetchResult.getSpeechs() {
				parseResult := ma.parse(speech)
				if parseResult.err != nil {
					ps.handleError(parseResult.Error())
					break
				}
				nouns := parseResult.getNouns()
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
	urlCh := ps.speech.generateURL(ctx)
	fetchResultCh := stream.Line[string, *fetchResult](ctx, urlCh, ps.speech.fetch)
	speechCh := stream.Demulti[*fetchResult, string](ctx, fetchResultCh, func(fr *fetchResult) []string {
		if fr.err != nil {
			ps.handleError(fr.Error())
		}
		return fr.getSpeechs()
	})
	nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
		pr := ma.parse(s)
		if pr.err != nil {
			ps.handleError(pr.Error())
		}
		return pr.getNouns()
	})
	return stream.Line[[]string, *apiv1.ColloStreamResponse](ctx, nounsCh, func(s []string) *apiv1.ColloStreamResponse {
		c := ps.newPairChunk()
		c.set(s)
		return c.ConvResp()
	})
}

// fetchから丸々funアウトする
func (ps *PairStore) stream_case2(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	urlCh := ps.speech.generateURL(ctx)
	return stream.FunIO[string, *apiv1.ColloStreamResponse](ctx, urlCh, func(url string) *apiv1.ColloStreamResponse {
		fetchResult := ps.speech.fetch(url)
		speechCh := fetchResult.generateSpeech(ctx)
		nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				ps.handleError(pr.Error())
			}
			return pr.getNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c.ConvResp()
	})
}

// 形態素解析からfunアウトする
func (ps *PairStore) stream_case3(ctx context.Context) <-chan *apiv1.ColloStreamResponse {
	urlCh := ps.speech.generateURL(ctx)
	fetchResultCh := stream.Line[string, *fetchResult](ctx, urlCh, ps.speech.fetch)
	return stream.FunIO[*fetchResult, *apiv1.ColloStreamResponse](ctx, fetchResultCh, func(fr *fetchResult) *apiv1.ColloStreamResponse {
		speechCh := fr.generateSpeech(ctx)
		nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				ps.handleError(pr.Error())
			}
			return pr.getNouns()
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
	urlCh := ps.speech.generateURL(ctx)
	return stream.FunIO[string, *apiv1.ColloStreamResponse](ctx, urlCh, func(url string) *apiv1.ColloStreamResponse {
		fetchResult := ps.speech.fetch(url)
		speechCh := fetchResult.generateSpeech(ctx)
		nounsCh := stream.FunIO[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				ps.handleError(pr.Error())
			}
			return pr.getNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c.ConvResp()
	})
}

func (ps *PairStore) append(word string) (string, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if id, found := ps.idByWord[word]; found {
		return id, true
	}
	id := fmt.Sprint(len(ps.idByWord) + 1)
	ps.idByWord[word] = id
	return id, false
}

type PairChunk struct {
	ps       *PairStore
	WordByID map[string]string
	Pairs    []string
}

func (pc *PairChunk) ConvResp() *apiv1.ColloStreamResponse {
	resp := &apiv1.ColloStreamResponse{}
	resp.Words = pc.WordByID
	resp.Pairs = pc.Pairs
	return resp
}

func (ps *PairStore) newPairChunk() *PairChunk {
	return &PairChunk{
		ps:       ps,
		WordByID: make(map[string]string),
		Pairs:    []string{},
	}
}

func (pc *PairChunk) set(nouns []string) {
	for i := 0; i < len(nouns); i++ {
		id1, found := pc.ps.append(nouns[i])
		if !found {
			pc.WordByID[id1] = nouns[i]
		}
		for j := i + 1; j < len(nouns); j++ {
			if nouns[i] == nouns[j] {
				continue
			}
			id2, found := pc.ps.append(nouns[j])
			if !found {
				pc.WordByID[id2] = nouns[j]
			}
			// pair order
			if id1 > id2 {
				tmp := id1
				id1 = id2
				id2 = tmp
			}
			pair := fmt.Sprintf("%s,%s", id1, id2)
			pc.Pairs = append(pc.Pairs, pair)
		}
	}
}
