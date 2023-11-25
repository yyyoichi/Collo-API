package pair

import (
	"context"
	"fmt"
	"sync"
	"yyyoichi/Collo-API/pkg/stream"
)

type PairStore struct {
	speech *Speech

	idByWord map[string]string
	mu       sync.Mutex
}

func NewPairStore(config Config) (*PairStore, error) {
	speech, err := NewSpeech(config)
	if err != nil {
		return nil, err
	}
	ps := &PairStore{
		speech:   speech,
		idByWord: map[string]string{},
		mu:       sync.Mutex{},
	}
	return ps, nil
}

type FetchError struct {
	error
}
type ParseError struct {
	error
}

// ストリームなし
func (ps *PairStore) stream_case0(ctx context.Context, cancel context.CancelCauseFunc) <-chan *PairChunk {
	ch := make(chan *PairChunk)
	go func(ps *PairStore) {
		defer close(ch)
		for url := range ps.speech.generateURL(ctx) {
			fetchResult := ps.speech.fetch(url)
			if fetchResult.err != nil {
				cancel(fetchResult.Error())
				break
			}
			chunk := ps.newPairChunk()
			for _, speech := range fetchResult.getSpeechs() {
				parseResult := ma.parse(speech)
				if parseResult.err != nil {
					cancel(parseResult.Error())
					break
				}
				nouns := parseResult.getNouns()
				chunk.set(nouns)
			}
			select {
			case <-ctx.Done():
				return
			default:
				ch <- chunk
			}
		}
	}(ps)
	return ch
}

// 全てを順にパイプ
func (ps *PairStore) stream_case1(ctx context.Context, cancel context.CancelCauseFunc) <-chan *PairChunk {
	urlCh := ps.speech.generateURL(ctx)
	fetchResultCh := stream.Line[string, *fetchResult](ctx, urlCh, ps.speech.fetch)
	speechCh := stream.Demulti[*fetchResult, string](ctx, fetchResultCh, func(fr *fetchResult) []string {
		if fr.err != nil {
			cancel(fr.Error())
		}
		return fr.getSpeechs()
	})
	nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
		pr := ma.parse(s)
		if pr.err != nil {
			cancel(pr.Error())
		}
		return pr.getNouns()
	})
	return stream.Line[[]string, *PairChunk](ctx, nounsCh, func(s []string) *PairChunk {
		c := ps.newPairChunk()
		c.set(s)
		return c
	})
}

// fetchから丸々funアウトする
func (ps *PairStore) stream_case2(ctx context.Context, cancel context.CancelCauseFunc) <-chan *PairChunk {
	urlCh := ps.speech.generateURL(ctx)
	return stream.FunIO[string, *PairChunk](ctx, urlCh, func(url string) *PairChunk {
		fetchResult := ps.speech.fetch(url)
		speechCh := fetchResult.generateSpeech(ctx)
		nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				cancel(pr.Error())
			}
			return pr.getNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c
	})
}

// 形態素解析からfunアウトする
func (ps *PairStore) stream_case3(ctx context.Context, cancel context.CancelCauseFunc) <-chan *PairChunk {
	urlCh := ps.speech.generateURL(ctx)
	fetchResultCh := stream.Line[string, *fetchResult](ctx, urlCh, ps.speech.fetch)
	return stream.FunIO[*fetchResult, *PairChunk](ctx, fetchResultCh, func(fr *fetchResult) *PairChunk {
		speechCh := fr.generateSpeech(ctx)
		nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				cancel(pr.Error())
			}
			return pr.getNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c
	})
}

// fetchから丸々funアウト, 形態素解析前にもfunアウトする
func (ps *PairStore) stream_case4(ctx context.Context, cancel context.CancelCauseFunc) <-chan *PairChunk {
	urlCh := ps.speech.generateURL(ctx)
	return stream.FunIO[string, *PairChunk](ctx, urlCh, func(url string) *PairChunk {
		fetchResult := ps.speech.fetch(url)
		speechCh := fetchResult.generateSpeech(ctx)
		nounsCh := stream.FunIO[string, []string](ctx, speechCh, func(s string) []string {
			pr := ma.parse(s)
			if pr.err != nil {
				cancel(pr.Error())
			}
			return pr.getNouns()
		})
		c := ps.newPairChunk()
		for nouns := range nounsCh {
			c.set(nouns)
		}
		return c
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
