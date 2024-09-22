package datalake

import "time"

type Config struct {
	Search struct {
		From  time.Time // 開始日
		Until time.Time // 終了日(含)
		Any   string    // キーワード
	}
	NoCache  bool   // キャッシュを利用しない
	CacheDir string // キャッシュするディレクトリ
}
