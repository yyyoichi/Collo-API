package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMorpheme(t *testing.T) {
	t.Run("run analysis", func(t *testing.T) {
		ar := Analysis("それに凍り付け！特別な時間。")
		m := ar.GetAt(0)
		require.Equal(t, "其れ", m.Lemma)

		words := ar.Get(Config{
			Includes: []PartOfSpeechType{
				Verb,
				AdjectiveVerb,
			},
		})
		require.Equal(t, 2, len(words))
		require.Equal(t, []string{"凍り付く", "特別"}, words)
	})
	t.Run("run empty", func(t *testing.T) {
		ar := Analysis("")
		m := ar.GetAt(0)
		require.True(t, m.isEnd())
		words := ar.Get(Config{
			Includes: []PartOfSpeechType{
				Verb,
				AdjectiveVerb,
			},
		})
		require.Equal(t, 0, len(words))
	})
}
