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

type Search struct {
	config          Config
	numberOfRecords *int   // 検索結果数
	urlf            string // 検索用
	initUrlf        string // 初期件数検索用
	maxFetch        int    // 最大取得数
	newResult       resultInitializerInterface
}

func NewSearch(config Config) *Search {
	s := &Search{config: config}
	s.config.init()
	switch s.config.SearchAPI {
	case MeetingAPI:
		s.urlf = "https://kokkai.ndl.go.jp/api/meeting"          // 会議単位出力AP
		s.initUrlf = "https://kokkai.ndl.go.jp/api/meeting_list" //I会議単位簡易出力API
		s.maxFetch = 10
		s.newResult = &meetingInitializer{fetcher: config.Fetcher}
	case SpeechAPI:
		s.urlf = "https://kokkai.ndl.go.jp/api/speech"     //発言出力API
		s.initUrlf = "https://kokkai.ndl.go.jp/api/speech" //発言出力API
		s.maxFetch = 100
		s.newResult = &speechInitializer{fetcher: config.Fetcher}
	}
	return s
}

func (s *Search) GenerateResult(ctx context.Context) <-chan ResultInterface {
	urls := s.getURLs()
	if len(urls) == 0 {
		r := s.newResult.error("", errors.New("not found"))
		return stream.Generator[ResultInterface](ctx, r)
	}
	urlCh := stream.Generator[string](ctx, urls...)
	return stream.Line[string, ResultInterface](ctx, urlCh, s.search)
}

func (s *Search) GetNumberOfRecords() int {
	if s.numberOfRecords != nil {
		return *s.numberOfRecords
	}
	r := s.search(s.createURL(1, 1, s.initUrlf))
	if r.Error() != nil {
		return 0
	}
	numRecords := r.numberOfRecords()
	s.numberOfRecords = &numRecords
	return *s.numberOfRecords
}

func (s *Search) search(url string) ResultInterface {
	body, err := s.config.Fetcher(url)
	if err != nil {
		return s.newResult.error("", err)
	}
	result := s.newResult.unmarshl(url, body)
	// has error
	if result.Error() != nil {
		return result
	}
	// regular
	if result.message() == "" {
		return result
	}
	// has error
	return s.newResult.error(url, errors.New(result.message()))
}

// [start] から[numberOfRecords]個のデータを取得するためのURLを取得する
func (s *Search) getURLs() []string {
	urls := []string{}
	for i := 1; i <= s.GetNumberOfRecords(); i += s.maxFetch {
		urls = append(urls, s.createURL(i, s.maxFetch, s.urlf))
	}
	return urls
}

func (s *Search) createURL(start, max int, urlf string) string {
	params := url.Values{}
	params.Add("maximumRecords", strconv.Itoa(max))
	params.Add("any", s.config.Search.Any)
	params.Add("from", s.config.Search.From.Format("2006-01-02"))
	params.Add("until", s.config.Search.Until.Format("2006-01-02"))
	params.Add("startRecord", strconv.Itoa(start))
	params.Add("recordPacking", "json")
	url := fmt.Sprintf("%s?%s", urlf, params.Encode())
	return url
}

type resultInitializerInterface interface {
	unmarshl(url string, body []byte) ResultInterface
	error(url string, err error) ResultInterface
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
//////////////// Metting API //////////////////////
///////////////////////////////////////////////////

type meetingInitializer struct {
	fetcher Fetcher
}

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

type speechInitializer struct {
	fetcher Fetcher
}

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
