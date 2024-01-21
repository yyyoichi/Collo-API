package ndl

import (
	"os"
	"time"
)

type NDLAPI int             // 検索AP種類
const MeetingAPI NDLAPI = 1 // 会議単位検索
const SpeechAPI NDLAPI = 2  // 発言単位検索

type Config struct {
	Search struct {
		From  time.Time // 開始日
		Until time.Time // 終了日(含)
		Any   string    // キーワード
	}
	NDLAPI      NDLAPI // 検索API
	UseCache    bool   // 検索時キャッシュ利用
	CreateCache bool   // 検索後キャッシュ作成
	CacheDir    string // キャッシュ利用するディレクトリ
}

func (c *Config) init() {
	if c.CacheDir == "" {
		c.CacheDir = "/tmp/ndl-cache"
	}
	if _, err := os.Stat(c.CacheDir); err != nil {
		os.Mkdir(c.CacheDir, 0700)
	}
	if c.Search.Any == "" {
		c.Search.Any = "科学技術"
	}
	l, _ := time.LoadLocation("Asia/Tokyo")
	if c.Search.From.IsZero() {
		c.Search.From = time.Date(2023, 11, 1, 0, 0, 0, 0, l)
	}
	if c.Search.Until.IsZero() {
		c.Search.Until = time.Date(2023, 11, 15, 0, 0, 0, 0, l)
	}
	if c.NDLAPI == 0 {
		c.NDLAPI = 1
	}
}
