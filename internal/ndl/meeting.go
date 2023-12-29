package ndl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yyyoichi/Collo-API/pkg/stream"
)

const meetingListUrlf = "https://kokkai.ndl.go.jp/api/meeting_list" // 会議単位簡易出力API
const meetingUrlf = "https://kokkai.ndl.go.jp/api/meeting"          //会議単位出力API
const maximumRecords = 10                                           // 一回の最大取得件数

type Meeting struct {
	config          Config
	numberOfRecords *int
}

func NewMeeting(config Config) *Meeting {
	m := &Meeting{config: config}
	m.init()
	return m
}

func (m *Meeting) GenerateMeeting(ctx context.Context) <-chan *MeetingResult {
	urls := m.getURLs()
	if len(urls) == 0 {
		mr := &MeetingResult{
			err: errors.New("not found data"),
		}
		return stream.Generator[*MeetingResult](ctx, mr)
	}
	urlCh := stream.Generator[string](ctx, urls...)
	return stream.Line[string, *MeetingResult](ctx, urlCh, m.fetch)
}

func (m *Meeting) GetNumberOfRecords() int {
	if m.numberOfRecords == nil {
		mr := m.fetch(m.getInitURL())
		if mr.err != nil {
			return 0
		}
		m.numberOfRecords = &mr.Result.NumberOfRecords
	}
	return *m.numberOfRecords
}

func (m *Meeting) fetch(url string) *MeetingResult {
	body, err := m.config.Fetcher(url)
	if err != nil {
		return &MeetingResult{err: err, url: url}
	}
	result := &MeetingResult{url: url}
	if err := json.Unmarshal(body, &result.Result); err != nil {
		return &MeetingResult{err: err, url: url}
	}
	// regular
	if result.Result.Message == "" {
		return result
	}
	// has error
	return &MeetingResult{err: errors.New(result.Result.Message), url: url}
}

func (m *Meeting) getURLs() []string {
	numberOfRecords := m.GetNumberOfRecords()
	urls := []string{}
	for i := 1; i <= numberOfRecords; i += maximumRecords {
		urls = append(urls, m.createURL(i, maximumRecords, meetingUrlf))
	}
	return urls
}

func (m *Meeting) getInitURL() string {
	return m.createURL(1, 1, meetingListUrlf)
}

func (m *Meeting) createURL(start, max int, urlf string) string {
	if max == 0 {
		max = maximumRecords
	}
	params := url.Values{}
	params.Add("maximumRecords", strconv.Itoa(max))
	params.Add("any", m.config.Search.Any)
	params.Add("from", m.config.Search.From.Format("2006-01-02"))
	params.Add("until", m.config.Search.Until.Format("2006-01-02"))
	params.Add("startRecord", strconv.Itoa(start))
	params.Add("recordPacking", "json")
	url := fmt.Sprintf("%s?%s", urlf, params.Encode())
	return url
}

func (m *Meeting) init() {
	if m.config.Fetcher == nil {
		m.config.Fetcher = func(url string) (body []byte, err error) {
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
	if m.config.Search.Any == "" {
		m.config.Search.Any = "科学技術"
	}
	l, _ := time.LoadLocation("Asia/Tokyo")
	if m.config.Search.From.IsZero() {
		m.config.Search.From = time.Date(2023, 11, 1, 0, 0, 0, 0, l)
	}
	if m.config.Search.Until.IsZero() {
		m.config.Search.Until = time.Date(2023, 11, 15, 0, 0, 0, 0, l)
	}
}

type NdlError struct{ error }

// API取得結果
type MeetingResult struct {
	url    string
	err    error
	Result struct {
		Message         string `json:"message"`
		NumberOfRecords int    `json:"numberOfRecords"`
		MeetingRecord   []struct {
			IssueID       string `json:"issueID"`       // 会議録ID
			Session       uint16 `json:"session"`       // 国会回次(ex:第212臨時国会)
			NameOfHouse   string `json:"nameOfHouse"`   // 院名(ex:衆議院)
			NameOfMeeting string `json:"nameOfMeeting"` // 会議名(ex:本会議)
			Issue         string `json:"issue"`         // 号数(ex:第x号)
			Date          string `json:"date"`          // 日付
			SpeechRecord  []struct {
				Speaker string `json:"speaker"`
				Speech  string `json:"speech"`
			} `json:"speechRecord"`
		} `json:"meetingRecord"`
	}
}

var re, _ = regexp.Compile(`^○.*?　`)
var replacer = strings.NewReplacer(
	"\u3000", "",
	"\r\n", "",
	"\r", "",
	"\n", "",
)

func (mr *MeetingResult) Error() error {
	if mr.err != nil {
		return NdlError{mr.err}
	}
	return nil
}
func (mr *MeetingResult) URL() string { return mr.url }

// 会議ごとの発言をひとまとめにして返す。
func (mr *MeetingResult) GetSpeechsPerMeeting() []string {
	speechs := make([]string, len(mr.Result.MeetingRecord))
	for i, meeting := range mr.Result.MeetingRecord {
		for _, speech := range meeting.SpeechRecord {
			if speech.Speaker == "会議録情報" {
				continue
			}
			s := re.ReplaceAllLiteralString(speech.Speech, "")
			speechs[i] += replacer.Replace(s)
		}
	}
	return speechs
}

// MeetingAPI取得結果から会議情報を作成する
func NewMeetingRecodes(mr *MeetingResult) []*MeetingRecode {
	mrs := make([]*MeetingRecode, len(mr.Result.MeetingRecord))
	for i, meeting := range mr.Result.MeetingRecord {
		mrs[i] = &MeetingRecode{}
		for _, speech := range meeting.SpeechRecord {
			if speech.Speaker == "会議録情報" {
				continue
			}
			s := re.ReplaceAllLiteralString(speech.Speech, "")
			mrs[i].Speeches += replacer.Replace(s)
		}
		mrs[i].IssueID = meeting.Issue
		mrs[i].Session = meeting.Session
		mrs[i].NameOfHouse = meeting.NameOfHouse
		mrs[i].NameOfMeeting = meeting.NameOfMeeting
		mrs[i].Issue = meeting.Issue
		mrs[i].Date, _ = time.Parse("2006-01-02 MST", meeting.Date+" JST")
	}
	return mrs
}

// 会議情報
type MeetingRecode struct {
	IssueID       string    // 会議録ID
	Session       uint16    // 国会回次(ex:第212臨時国会)
	NameOfHouse   string    // 院名(ex:衆議院)
	NameOfMeeting string    // 会議名(ex:本会議)
	Issue         string    // 号数(ex:第x号)
	Date          time.Time // 日付
	Speeches      string    // 会議中のすべての発言
}
