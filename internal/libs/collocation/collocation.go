package collocation

import "fmt"

func NewCollocation() *Collocation {
	return &Collocation{make(map[string]string)}
}

type Collocation struct {
	idByWord map[string]string
}

func (c *Collocation) getWordID(s string) (string, bool) {
	id, found := c.idByWord[s]
	if found {
		return id, true
	}
	id = fmt.Sprint(len(c.idByWord) + 1)
	c.idByWord[s] = id
	return id, false
}

type CollocationResult struct {
	WordByID map[string]string // 新しく追加された単語とそのID
	Pairs    []string          // id,idの共起ペア
	Err      error
}

// 名詞からペアリストを返す。インスタンス中、新出名詞はWordByIDとして返す。
func (c *Collocation) Get(nouns []string) *CollocationResult {
	result := &CollocationResult{}

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
			pair := fmt.Sprintf("%s,%s", id1, id2)
			result.Pairs = append(result.Pairs, pair)
		}
	}

	return result
}
