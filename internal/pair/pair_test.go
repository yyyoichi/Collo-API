package pair

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestT(t *testing.T) {
	t.Run("Morpheme", func(t *testing.T) {
		s := "明日は天気が良いと言います。乾燥しているので洗濯するとよいです。君。"
		r := ma.parse(s)
		expected := []string{
			"明日", "天気", "乾燥", "洗濯",
		}
		gots := r.getNouns()
		require.Equal(t, len(expected), len(gots))
		require.Equal(t, expected, gots)
	})
	t.Run("FetchMorpheme", func(t *testing.T) {
		speech, err := NewSpeech(tconfig)
		require.NoError(t, err)
		for url := range speech.generateURL(context.Background()) {
			fr := speech.fetch(url)
			for _, s := range fr.getSpeechs() {
				r := ma.parse(s)
				require.NoError(t, r.err)
			}
		}
	})
	t.Run("PairChunk", func(t *testing.T) {
		ps, err := NewPairStore(tconfig, thandler)
		require.NoError(t, err)
		chunk := ps.newPairChunk()
		chunk.set([]string{"1", "2", "3"})

		require.Equal(t, 3, len(ps.idByWord))
		for ky, vl := range ps.idByWord {
			require.Equal(t, ky, vl)
		}
		require.Equal(t, 3, len(chunk.ConvResp().Words))
		require.Equal(t, 3, len(chunk.WordByID))
		for ky, vl := range chunk.WordByID {
			require.Equal(t, ky, vl)
		}
		require.Equal(t, []string{"1,2", "1,3", "2,3"}, chunk.ConvResp().Pairs)
		require.Equal(t, []string{"1,2", "1,3", "2,3"}, chunk.Pairs)
	})
}
