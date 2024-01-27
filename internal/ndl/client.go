package ndl

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"yyyoichi/Collo-API/pkg/stream"
)

type Client struct {
	config          Config
	numberOfRecords *int   // 検索結果数
	urlf            string // 検索用
	initUrlf        string // 初期件数検索用
	maxFetch        int    // 最大取得数
	newResult       resultInitializerInterface
	DoGet           func(url string) (body []byte, err error)
}

func NewClient(config Config) Client {
	var c Client
	config.init()
	c.config = config
	switch c.config.NDLAPI {
	case MeetingAPI:
		c.urlf = "https://kokkai.ndl.go.jp/api/meeting"          // 会議単位出力AP
		c.initUrlf = "https://kokkai.ndl.go.jp/api/meeting_list" // 会議単位簡易出力API
		c.maxFetch = 10
		c.newResult = MeetingResult{}
	case SpeechAPI:
		c.urlf = "https://kokkai.ndl.go.jp/api/speech"     // 発言出力API
		c.initUrlf = "https://kokkai.ndl.go.jp/api/speech" // 発言出力API
		c.maxFetch = 100
		c.newResult = SpeechResult{}
	}
	var doget cachedDoGet = cachedDoGet{
		useCache:    config.UseCache,
		createCache: config.CreateCache,
		dir:         config.CacheDir,
	}
	c.DoGet = doget.DoGet
	return c
}

// 検索APIを叩く数とその結果を返す。
func (c *Client) GenerateResult(ctx context.Context) (int, <-chan ResultInterface) {
	urls := c.getURLs()
	if len(urls) == 0 {
		r := c.newResult.error("", errors.New("not found"))
		return 1, stream.Generator[ResultInterface](ctx, r)
	}
	urlCh := stream.Generator[string](ctx, urls...)
	return len(urls), stream.Line[string, ResultInterface](ctx, urlCh, c.search)
}

// 期待されるNDLRecord数とNDLRecordチャネルを返す。（NDLRecord数と送信チャネル数は必ずしも一致しない）
func (c *Client) GenerateNDLResultWithErrorHook(ctx context.Context, errHook stream.ErrorHook) (int, <-chan NDLRecode) {
	_, resultCh := c.GenerateResult(ctx)
	return c.GetNumberOfRecords(), stream.DemultiWithErrorHook[ResultInterface, NDLRecode](
		ctx,
		errHook,
		resultCh,
		func(r ResultInterface) []NDLRecode {
			return r.NewNDLRecodes()
		},
	)

}

func (c *Client) GetNumberOfRecords() int {
	if c.numberOfRecords != nil {
		return *c.numberOfRecords
	}
	r := c.search(c.createURL(1, 1, c.initUrlf))
	if r.Error() != nil {
		return 0
	}
	numRecords := r.numberOfRecords()
	c.numberOfRecords = &numRecords
	return *c.numberOfRecords
}

func (c *Client) search(url string) ResultInterface {
	body, err := c.DoGet(url)
	if err != nil {
		return c.newResult.error("", err)
	}
	if result := c.newResult.unmarshl(url, body); result.Error() != nil {
		// has http error
		return result
	} else if result.message() != "" {
		// has request error
		return c.newResult.error(url, errors.New(result.message()))
	} else {
		// ok
		return result
	}
}

// [start] から[numberOfRecords]個のデータを取得するためのURLを取得する
func (c *Client) getURLs() []string {
	urls := []string{}
	for i := 1; i <= c.GetNumberOfRecords(); i += c.maxFetch {
		urls = append(urls, c.createURL(i, c.maxFetch, c.urlf))
	}
	return urls
}

func (c *Client) createURL(start, max int, urlf string) string {
	params := url.Values{}
	params.Add("maximumRecords", strconv.Itoa(max))
	params.Add("any", c.config.Search.Any)
	params.Add("from", c.config.Search.From.Format("2006-01-02"))
	params.Add("until", c.config.Search.Until.Format("2006-01-02"))
	params.Add("startRecord", strconv.Itoa(start))
	params.Add("recordPacking", "json")
	url := fmt.Sprintf("%s?%s", urlf, params.Encode())
	return url
}
