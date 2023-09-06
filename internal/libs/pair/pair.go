package pair

import (
	"fmt"
	"sync"
)

func NewPair() *Pair {
	return &Pair{make(map[string]string), sync.Mutex{}}
}

type Pair struct {
	idByWord map[string]string
	mu       sync.Mutex
}

func (c *Pair) getWordID(s string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id, found := c.idByWord[s]
	if found {
		return id, true
	}
	id = fmt.Sprint(len(c.idByWord) + 1)
	c.idByWord[s] = id
	return id, false
}

func NewPairResult() *PairResult {
	return &PairResult{
		WordByID: make(map[string]string),
		Pairs:    []string{},
	}
}

type PairResult struct {
	WordByID map[string]string // 新しく追加された単語とそのID
	Pairs    []string          // id,idの共起ペア
	Err      error
}

// ペアリストを結合する
func (p *PairResult) Concat(pr *PairResult) {
	for id, word := range pr.WordByID {
		p.WordByID[id] = word
	}
	// 2. concat Pairs
	p.Pairs = append(p.Pairs, pr.Pairs...)
}

// 名詞からペアリストを返す。インスタンス中、新出名詞はWordByIDとして返す。
func (c *Pair) Get(nouns []string) *PairResult {
	result := NewPairResult()
	for i := 0; i < len(nouns); i++ {
		id1, found := c.getWordID(nouns[i])
		if !found {
			result.WordByID[id1] = nouns[i]
		}
		for j := i + 1; j < len(nouns); j++ {
			id2, found := c.getWordID(nouns[j])
			if !found {
				result.WordByID[id2] = nouns[j]
			}
			// pair order
			if id1 > id2 {
				tmp := id1
				id1 = id2
				id2 = tmp
			}
			pair := fmt.Sprintf("%s,%s", id1, id2)
			result.Pairs = append(result.Pairs, pair)
		}
	}

	return result
}
