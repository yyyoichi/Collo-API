package ndl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMeeting(t *testing.T) {
	t.Run("mock", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test-meeting-mock")
		require.NoError(t, err)
		foundfile := func(url string) {
			file := filepath.Join(fmt.Sprintf("%s/%s.json", dir, md5Hash(url)))
			b, err := os.ReadFile(file)
			require.NoError(t, err)
			var mr *MeetingResult
			err = json.Unmarshal(b, &mr)
			require.NoError(t, err)
		}

		config := Config{}
		config.Search.Any = "少子高齢化"
		tconfig := CreateMeetingConfigMock(config, dir)

		tmeeting := NewMeeting(tconfig)
		initURL := tmeeting.getInitURL()
		foundfile(initURL)
		for _, url := range tmeeting.getURLs() {
			foundfile(url)
		}
		require.Equal(t, "少子高齢化", tmeeting.config.Search.Any)
		require.NoError(t, os.RemoveAll(dir))
	})

	t.Run("speech", func(t *testing.T) {
		config := Config{}
		config.Search.Any = "科学"
		l, _ := time.LoadLocation("Asia/Tokyo")
		config.Search.From = time.Date(2023, 11, 30, 0, 0, 0, 0, l)
		config.Search.Until = time.Date(2023, 12, 9, 0, 0, 0, 0, l)
		m := NewMeeting(config)

		results := []*MeetingResult{}
		for mr := range m.GenerateMeeting(context.Background()) {
			results = append(results, mr)
		}
		require.Equal(t, 2, len(results))
		speechs := results[0].GetSpeechsPerMeeting()
		require.Equal(t, 10, len(speechs))
		require.Truef(t, strings.HasPrefix(speechs[0], "これより会議を開きます。日程第一国立大学法人法の"), "got '%s'", speechs[0][:20])

		mrs := NewMeetingRecodes(results[0])
		require.Equal(t, 10, len(mrs))
		for _, result := range results {
			mrs := NewMeetingRecodes(result)
			for _, mr := range mrs {
				require.NotEmpty(t, mr.Issue)
				require.NotEmpty(t, mr.IssueID)
				require.NotEmpty(t, mr.NameOfHouse)
				require.NotEmpty(t, mr.NameOfMeeting)
				require.NotEmpty(t, mr.Session)
				require.NotEmpty(t, mr.Speeches)
				require.NotNil(t, mr.Date)
			}
		}
	})
}
