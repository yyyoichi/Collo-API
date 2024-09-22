package datalake

import "time"

// 1会議情報
type Meeting struct {
	IssueID       string    // 会議録ID
	Session       int       // 国会回次(ex:第212臨時国会)
	NameOfHouse   string    // 院名(ex:衆議院)
	NameOfMeeting string    // 会議名(ex:本会議)
	Issue         string    // 号数(ex:第x号)
	Date          time.Time // 日付
	Speeches      string    // 発言
}
