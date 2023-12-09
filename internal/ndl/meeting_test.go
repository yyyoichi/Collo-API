package ndl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
}
