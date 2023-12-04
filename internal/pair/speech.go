package pair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"yyyoichi/Collo-API/pkg/stream"
)

var (
	ErrUnExpected       = errors.New("unexpected")
	ErrCongestedServer  = errors.New("congested")
	ErrBadRequest       = errors.New("bad request")
	ErrExcessCharacters = errors.New("excess characters")
	ErrNoData           = errors.New("no data")
)

const urlf = "https://kokkai.ndl.go.jp/api/speech"

type Speech struct {
	config Config

	containRecords int
}

func NewSpeech(config Config) (*Speech, error) {
	s := &Speech{
		config: config,
	}
	s.init()
	result := s.Fetch(s.GetInitURL())
	if result.err != nil {
		return nil, result.Error()
	}
	s.containRecords = result.SpeechJson.NumberOfRecords
	return s, nil
}

func (s *Speech) GetURLs() []string {
	urls := []string{}
	max := 100
	for i := 1; i <= s.containRecords; i += max {
		urls = append(urls, s.createURL(i, 100))
	}
	return urls
}

func (s *Speech) GenerateURL(ctx context.Context) <-chan string {
	return stream.Generator[string](ctx, s.GetURLs()...)
}

func (s *Speech) GetInitURL() string {
	return s.createURL(1, 1)
}

func (s *Speech) Fetch(url string) *FetchResult {
	body, err := s.config.Fetcher(url)
	if err != nil {
		return &FetchResult{err: err, url: url}
	}
	result := &FetchResult{url: url}
	if err := json.Unmarshal(body, &result.SpeechJson); err != nil {
		return &FetchResult{err: err, url: url}
	}
	// regular
	if result.SpeechJson.Message == "" {
		return result
	}
	// get error
	err = ErrUnExpected
	if strings.Contains(result.SpeechJson.Message, "19001") {
		err = ErrCongestedServer
	} else if strings.Contains(result.SpeechJson.Message, "19020") {
		err = ErrExcessCharacters
	} else if strings.Contains(result.SpeechJson.Message, "19011") {
		err = ErrBadRequest
	}
	return &FetchResult{err: err, url: url}
}

func (s *Speech) createURL(start, max int) string {
	if max == 0 {
		max = 100
	}
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

func (s *Speech) init() {
	if s.config.Fetcher == nil {
		s.config.Fetcher = func(url string) (body []byte, err error) {
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
	if s.config.Search.Any == "" {
		s.config.Search.Any = "災害"
	}
	if s.config.Search.From.IsZero() {
		s.config.Search.From = time.Now().Add(time.Duration(-time.Hour * 24 * 365))
	}
	if s.config.Search.Until.IsZero() {
		s.config.Search.From = time.Now().Add(time.Duration(-time.Hour * 24 * 364))
	}
}

type FetchResult struct {
	SpeechJson *struct {
		NumberOfRecords    int    `json:"numberOfRecords"`
		NumberOfReturn     int    `json:"numberOfReturn"`
		StartRecord        int    `json:"startRecord"`
		NextRecordPosition int    `json:"nextRecordPosition"`
		Message            string `json:"message"`
		SpeechRecord       []struct {
			Speech     string      `json:"speech"`
			OtherField interface{} `json:"-"`
		} `json:"speechRecord"`
		OtherField interface{} `json:"-"`
	}
	err error
	url string
}

func (fr *FetchResult) GetSpeechs() []string {
	speechs := []string{}
	for _, r := range fr.SpeechJson.SpeechRecord {
		speechs = append(speechs, r.Speech)
	}
	return speechs
}

func (fr *FetchResult) GenerateSpeech(ctx context.Context) <-chan string {
	return stream.Generator[string](ctx, fr.GetSpeechs()...)
}

func (fr *FetchResult) Error() error {
	if fr.err == nil {
		return nil
	}
	return fmt.Errorf("error[%s]: %v", fr.url, fr.err)
}
