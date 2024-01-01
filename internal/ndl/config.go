package ndl

import (
	"io"
	"net/http"
	"time"
)

type NDLAPI int             // 検索AP種類
const MeetingAPI NDLAPI = 1 // 会議単位検索
const SpeechAPI NDLAPI = 2  // 発言単位検索

type DoGet func(url string) (body []byte, err error)

type Config struct {
	Search struct {
		From  time.Time // 開始日
		Until time.Time // 終了日(含)
		Any   string    // キーワード
	}
	DoGet  DoGet  // APIコール
	NDLAPI NDLAPI // 検索API
}

func (c *Config) init() {
	if c.DoGet == nil {
		c.DoGet = func(url string) (body []byte, err error) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, err
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			return body, err
		}
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
