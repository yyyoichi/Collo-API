package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	test := []struct {
		config Config
		exp    string
	}{
		{
			Config{
				Includes:  []PartOfSpeechType{Noun},
				StopWords: []string{},
			},
			"i!:101s!:",
		},
		{
			Config{
				Includes:  []PartOfSpeechType{Noun, Adjective},
				StopWords: []string{"hoge", "fuga"},
			},
			"i!:101-201s!:hoge-fuga",
		},
		{
			Config{
				Includes:  []PartOfSpeechType{},
				StopWords: []string{},
			},
			// depend on config.init
			"i!:101-112-201-301-401s!:",
		},
	}
	for _, tt := range test {
		require.Equal(t, tt.exp, tt.config.ToString())
	}
}
