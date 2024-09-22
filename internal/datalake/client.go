package datalake

import (
	"context"
	"iter"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/yyyoichi/httpcache-go"
	kokkaiapi "github.com/yyyoichi/kokkai-api"
)

type Client struct {
	config  Config
	request kokkaiapi.Request[*kokkaiapi.KaigiResult]
}

func NewClient(config Config) Client {
	var client = kokkaiapi.KaigiClient
	client.Interval = 0

	// 国会議事録APIにリクエストを飛ばすときは、2回目以降、1秒の間隔を開ける。
	// NoCacheでなければ常にキャッシュを利用する。
	httpClient := &httpcache.Client{
		Client: &httpClientWithInterval{
			client:   http.DefaultClient,
			interval: time.Duration(time.Second),
		},
		Cache:   httpcache.NewStorageCache(config.CacheDir),
		Handler: httpcache.NewDefaultHandler(),
	}
	if config.NoCache {
		httpClient.Handler = httpcache.NewLatestHandler()
	}
	client.HttpClient = httpClient

	request := kokkaiapi.DefaultKaigiRequest
	request.Client = client
	return Client{
		config:  config,
		request: request,
	}
}

func (c *Client) IterMeeting(ctx context.Context) iter.Seq2[*Meeting, error] {
	ctx, cancel := context.WithCancelCause(ctx)

	resultCh := make(chan *kokkaiapi.KaigiResult)
	go func() {
		defer close(resultCh)

		param := kokkaiapi.NewParam()
		param.Any(c.config.Search.Any)
		param.From(c.config.Search.From.Format(apiDateF))
		param.Until(c.config.Search.Until.Format(apiDateF))
		for result, err := range kokkaiapi.IterResult(param, c.request) {
			if err != nil {
				cancel(err)
				return
			}
			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return func(yield func(*Meeting, error) bool) {
		defer cancel(nil)

		if err := context.Cause(ctx); err != nil {
			yield(nil, err)
			return
		}
		for r := range resultCh {
			select {
			case <-ctx.Done():
				yield(nil, context.Cause(ctx))
				return
			default:
				for m := range convertMeeting(r) {
					if ok := yield(m, nil); !ok {
						return
					}
				}
			}
		}
	}
}

var (
	re, _    = regexp.Compile(`^○.*?　`)
	replacer = strings.NewReplacer(
		"\u3000", "",
		"\r\n", "",
		"\r", "",
		"\n", "",
	)
	apiDateF = "2006-01-02"
	// parse by 2006-01-02
	parseTime = func(datestr string) time.Time {
		d, _ := time.Parse(apiDateF, datestr)
		return d
	}
)

func convertMeeting(r *kokkaiapi.KaigiResult) iter.Seq[*Meeting] {
	return func(yield func(*Meeting) bool) {
		for _, mr := range r.MeetingRecord {
			m := &Meeting{
				IssueID:       mr.IssueID,
				Session:       mr.Session,
				NameOfHouse:   mr.NameOfHouse,
				NameOfMeeting: mr.NameOfMeeting,
				Issue:         mr.Issue,
				Date:          parseTime(mr.Date),
			}
			var buf strings.Builder
			for _, ms := range mr.SpeechRecord {
				if ms.Speaker == "会議録情報" {
					continue
				}
				s := re.ReplaceAllLiteralString(ms.Speech, "")
				_, _ = buf.WriteString(replacer.Replace(s))
			}
			m.Speeches = buf.String()
			if ok := yield(m); !ok {
				return
			}
		}
	}
}

type httpClientWithInterval struct {
	client   *http.Client
	interval time.Duration
	reqCount int
}

func (c *httpClientWithInterval) Do(req *http.Request) (*http.Response, error) {
	if c.reqCount > 0 {
		<-time.After(c.interval)
	}
	resp, err := c.client.Do(req)
	c.reqCount++
	return resp, err
}
