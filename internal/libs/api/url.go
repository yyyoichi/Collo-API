package api

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type URLOptions struct {
	StartRecord    int       // 取得開始位置
	MaximumRecords int       // 最大取得数/url
	From           time.Time // 開始日
	Until          time.Time // 終了日(含)
	Any            string    // キーワード
}

func CreateURL(opt URLOptions) string {
	urlf := "https://kokkai.ndl.go.jp/api/speech"
	params := url.Values{}
	params.Add("maximumRecords", strconv.Itoa(opt.MaximumRecords))
	params.Add("any", opt.Any)
	params.Add("from", opt.From.Format("YYYY-MM-DD"))
	params.Add("until", opt.Until.Format("YYYY-MM-DD"))
	params.Add("startRecord", strconv.Itoa(opt.StartRecord))
	params.Add("recordPacking", "json")
	url := fmt.Sprintf("%s?%s", urlf, params.Encode())
	return url
}

// 取得開始位置から取得数[endRecord](含)までのURLを作成する
func CreateURLs(opt URLOptions, endRecord int) []string {
	urls := []string{}
	for start := opt.StartRecord; start <= endRecord; start += opt.MaximumRecords {
		opt.StartRecord = start
		urls = append(urls, CreateURL(opt))
	}
	return urls
}
