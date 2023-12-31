package ndl

import (
	"regexp"
	"strings"
	"time"
)

// 会議情報
type NDLRecode struct {
	IssueID       string    // 会議録ID
	Session       uint16    // 国会回次(ex:第212臨時国会)
	NameOfHouse   string    // 院名(ex:衆議院)
	NameOfMeeting string    // 会議名(ex:本会議)
	Issue         string    // 号数(ex:第x号)
	Date          time.Time // 日付
	Speeches      string    // 会議中のすべての発言
}

type NdlError struct{ error }

// result initializer

type ResultInterface interface {
	Error() error
	URL() string
	NewNDLRecodes() []*NDLRecode
	message() string
	numberOfRecords() int
}

var re, _ = regexp.Compile(`^○.*?　`)
var replacer = strings.NewReplacer(
	"\u3000", "",
	"\r\n", "",
	"\r", "",
	"\n", "",
)

///////////////////////////////////////////////////
///////////////////////////////////////////////////
//////////////// Metting API //////////////////////
///////////////////////////////////////////////////

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

func (r *MeetingResult) Error() error {
	if r.err != nil {
		return NdlError{r.err}
	}
	return nil
}
func (r *MeetingResult) URL() string          { return r.url }
func (r *MeetingResult) message() string      { return r.Result.Message }
func (r *MeetingResult) numberOfRecords() int { return r.Result.NumberOfRecords }

// MeetingAPI取得結果から会議情報を作成する
func (r *MeetingResult) NewNDLRecodes() []*NDLRecode {
	records := make([]*NDLRecode, len(r.Result.MeetingRecord))
	for i, meeting := range r.Result.MeetingRecord {
		records[i] = &NDLRecode{}
		for _, speech := range meeting.SpeechRecord {
			if speech.Speaker == "会議録情報" {
				continue
			}
			s := re.ReplaceAllLiteralString(speech.Speech, "")
			records[i].Speeches += replacer.Replace(s)
		}
		records[i].IssueID = meeting.IssueID
		records[i].Session = meeting.Session
		records[i].NameOfHouse = meeting.NameOfHouse
		records[i].NameOfMeeting = meeting.NameOfMeeting
		records[i].Issue = meeting.Issue
		records[i].Date, _ = time.Parse("2006-01-02 MST", meeting.Date+" JST")
	}
	return records
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
//////////////// Speech API ///////////////////////
///////////////////////////////////////////////////

// API取得結果
type SpeechResult struct {
	url    string
	err    error
	Result struct {
		Message         string `json:"message"`
		NumberOfRecords int    `json:"numberOfRecords"`
		SpeechRecord    []struct {
			SpeechID      string `json:"speechID"`
			IssueID       string `json:"issueID"`       // 会議録ID
			Session       uint16 `json:"session"`       // 国会回次(ex:第212臨時国会)
			NameOfHouse   string `json:"nameOfHouse"`   // 院名(ex:衆議院)
			NameOfMeeting string `json:"nameOfMeeting"` // 会議名(ex:本会議)
			Issue         string `json:"issue"`         // 号数(ex:第x号)
			Date          string `json:"date"`          // 日付
			Speaker       string `json:"speaker"`
			Speech        string `json:"speech"`
		} `json:"speechRecord"`
	}
}

func (r *SpeechResult) Error() error {
	if r.err != nil {
		return NdlError{r.err}
	}
	return nil
}
func (r *SpeechResult) URL() string          { return r.url }
func (r *SpeechResult) numberOfRecords() int { return r.Result.NumberOfRecords }
func (r *SpeechResult) message() string      { return r.Result.Message }

// SpeechAPI取得結果から会議情報を作成する
func (r *SpeechResult) NewNDLRecodes() []*NDLRecode {
	records := make([]*NDLRecode, len(r.Result.SpeechRecord))
	for i, speech := range r.Result.SpeechRecord {
		if speech.Speaker == "会議録情報" {
			continue
		}
		s := re.ReplaceAllLiteralString(speech.Speech, "")
		records[i].Speeches += replacer.Replace(s)
		records[i].IssueID = speech.IssueID
		records[i].Session = speech.Session
		records[i].NameOfHouse = speech.NameOfHouse
		records[i].NameOfMeeting = speech.NameOfMeeting
		records[i].Issue = speech.Issue
		records[i].Date, _ = time.Parse("2006-01-02 MST", speech.Date+" JST")
	}
	return records
}
