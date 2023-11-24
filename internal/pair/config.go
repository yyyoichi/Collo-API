package pair

import "time"

type Fetcher func(url string) (body []byte, err error)

type Config struct {
	Search struct {
		From  time.Time // 開始日
		Until time.Time // 終了日(含)
		Any   string    // キーワード
	}
	Fetcher Fetcher
}
