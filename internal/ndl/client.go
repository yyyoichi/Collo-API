package ndl

import (
	"context"
	"encoding/json"
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
}

type resultInitializerInterface interface {
	unmarshl(url string, body []byte) ResultInterface
	error(url string, err error) ResultInterface
}

func NewClient(config Config) *Client {
	c := &Client{config: config}
	c.config.init()
	switch c.config.NDLAPI {
	case MeetingAPI:
		c.urlf = "https://kokkai.ndl.go.jp/api/meeting"          // 会議単位出力AP
		c.initUrlf = "https://kokkai.ndl.go.jp/api/meeting_list" // 会議単位簡易出力API
		c.maxFetch = 10
		c.newResult = &meetingInitializer{}
	case SpeechAPI:
		c.urlf = "https://kokkai.ndl.go.jp/api/speech"     // 発言出力API
		c.initUrlf = "https://kokkai.ndl.go.jp/api/speech" // 発言出力API
		c.maxFetch = 100
		c.newResult = &speechInitializer{}
	}
	return c
}

func (c *Client) GenerateResult(ctx context.Context) <-chan ResultInterface {
	urls := c.getURLs()
	if len(urls) == 0 {
		r := c.newResult.error("", errors.New("not found"))
		return stream.Generator[ResultInterface](ctx, r)
	}
	urlCh := stream.Generator[string](ctx, urls...)
	return stream.Line[string, ResultInterface](ctx, urlCh, c.search)
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
	body, err := c.config.DoGet(url)
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

///////////////////////////////////////////////////
///////////////////////////////////////////////////
//////////////// Metting API //////////////////////
///////////////////////////////////////////////////

type meetingInitializer struct{}

func (i *meetingInitializer) unmarshl(url string, body []byte) ResultInterface {
	result := &MeetingResult{url: url}
	if err := json.Unmarshal(body, &result.Result); err != nil {
		return i.error(url, err)
	}
	return result
}
func (i *meetingInitializer) error(url string, err error) ResultInterface {
	return &MeetingResult{
		err: err,
		url: url,
	}
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
//////////////// Speech API ///////////////////////
///////////////////////////////////////////////////

type speechInitializer struct{}

func (i *speechInitializer) unmarshl(url string, body []byte) ResultInterface {
	result := &SpeechResult{url: url}
	if err := json.Unmarshal(body, &result.Result); err != nil {
		return i.error(url, err)
	}
	return result
}
func (i *speechInitializer) error(url string, err error) ResultInterface {
	return &SpeechResult{
		err: err,
		url: url,
	}
}
