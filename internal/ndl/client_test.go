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
		t.Parallel()
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
		config.UseCache = true
		config.CreateCache = true
		config.CacheDir = dir

		client := NewClient(config)
		_, resultCh := client.GenerateResult(context.Background())
		for range resultCh {
		}
		initURL := client.createURL(1, 1, client.initUrlf)
		foundfile(initURL)
		for _, url := range client.getURLs() {
			foundfile(url)
		}
		require.Equal(t, "少子高齢化", client.config.Search.Any)
		require.NoError(t, os.RemoveAll(dir))
	})

	t.Run("meeting", func(t *testing.T) {
		t.Parallel()
		config := Config{}
		config.Search.Any = "科学"
		l, _ := time.LoadLocation("Asia/Tokyo")
		config.Search.From = time.Date(2023, 11, 30, 0, 0, 0, 0, l)
		config.Search.Until = time.Date(2023, 12, 9, 0, 0, 0, 0, l)
		m := NewClient(config)

		results := []ResultInterface{}
		numDoGet, resultCh := m.GenerateResult(context.Background())
		require.Equal(t, 4, numDoGet)
		for r := range resultCh {
			results = append(results, r)
		}
		require.Equal(t, 4, len(results))
		records := results[0].NewNDLRecodes()
		require.Equal(t, 10, len(records))
		require.Truef(t, strings.HasPrefix(
			records[0].Speeches,
			"これより会議を開きます。"),
			"got '%s'",
			records[0].Speeches[:20],
		)

		for _, result := range results {
			mrs := result.NewNDLRecodes()
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
	t.Run("speech", func(t *testing.T) {
		t.Parallel()
		config := Config{}
		config.Search.Any = "科学"
		config.NDLAPI = SpeechAPI
		l, _ := time.LoadLocation("Asia/Tokyo")
		config.Search.From = time.Date(2023, 11, 30, 0, 0, 0, 0, l)
		config.Search.Until = time.Date(2023, 12, 9, 0, 0, 0, 0, l)
		m := NewClient(config)

		results := []ResultInterface{}
		numDoGet, resultCh := m.GenerateResult(context.Background())
		require.Equal(t, numDoGet, 2)
		for r := range resultCh {
			results = append(results, r)
		}
		require.Equal(t, 2, len(results))
		records := results[0].NewNDLRecodes()
		require.Equal(t, 21, len(records))

		for _, result := range results {
			mrs := result.NewNDLRecodes()
			for _, mr := range mrs {
				require.NotEmpty(t, mr.Issue)
				require.NotEmpty(t, mr.IssueID)
				require.NotEmpty(t, mr.NameOfHouse)
				require.NotEmpty(t, mr.NameOfMeeting)
				require.NotEmpty(t, mr.Session)
				require.NotNil(t, mr.Date)
			}
		}
	})
}
