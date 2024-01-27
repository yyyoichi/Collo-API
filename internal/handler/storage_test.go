package handler

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"yyyoichi/Collo-API/internal/matrix"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	var handler ProcessHandler
	var config Config
	newMock := func() Storage {
		var storage Storage
		storage.new = func(ctx context.Context, ph ProcessHandler, c Config) CoMatrixes {
			var cm matrix.CoMatrix
			cm.Words = []string{"hoge"}
			cm.Indices = []int{0}
			cm.Matrix = []float64{1}
			cm.Priority = []float64{1}
			cm.Meta = matrix.MultiDocMeta{
				GroupID: "ID",
				Metas: []matrix.DocMeta{
					{Key: "Key"},
				},
			}
			return CoMatrixes{cm}
		}
		storage.getFilename = func(s string) string {
			panic(errors.New(""))
		}
		return storage
	}
	newStoragedMock := func(filename string) (Storage, error) {
		storage := newMock()
		storage.getFilename = func(s string) string { return filename }
		storage.CoMatrixes = storage.new(context.Background(), handler, config)
		storage.Words = []string{"hoge"}

		// save storage
		err := storage.saveCoMatrixes(context.Background(), config)
		return storage, err
	}

	t.Run("read and write", func(t *testing.T) {
		t.Parallel()
		var filename = "/tmp/storaged.json"
		storage, err := newStoragedMock(filename)
		require.NoError(t, err)
		defer os.Remove(filename)

		// reset in memory
		storage.CoMatrixes = CoMatrixes{}
		storage.Words = []string{}
		require.Empty(t, storage.CoMatrixes)
		require.Empty(t, storage.Words)

		found := storage.readCoMatrixes(context.Background(), handler, config)
		require.True(t, found)

		require.Equal(t, 1, len(storage.Words))
		require.Equal(t, "hoge", storage.Words[0])
		require.Equal(t, 1, len(storage.CoMatrixes))
		cm := storage.CoMatrixes[0]
		require.EqualValues(t, 0, cm.Indices[0])
		require.EqualValues(t, 1, cm.Priority[0])
		require.EqualValues(t, 1, cm.Matrix[0])
		require.Equal(t, "ID", cm.Meta.GroupID)
		require.Equal(t, 1, len(cm.Meta.Metas))
		require.Equal(t, "Key", cm.Meta.Metas[0].Key)
	})

	t.Run("not yet storage", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		storage := newMock()
		storage.getFilename = func(s string) string {
			return "/notfoundfile.json"
		}
		require.Empty(t, storage.CoMatrixes)
		require.Empty(t, storage.Words)
		found := storage.readCoMatrixes(ctx, handler, config)
		require.False(t, found)
		require.NotEmpty(t, storage.CoMatrixes)
		require.Empty(t, storage.Words)
		t.Run("", func(t *testing.T) {
		})
	})

	t.Run("new func", func(t *testing.T) {
		t.Parallel()
		var filename = "/tmp/newfunc.json"
		storage, err := newStoragedMock(filename)
		require.NoError(t, err)
		defer os.Remove(filename)
		cm := storage.NewCoMatrixes(context.Background(), handler, storagePermission{
			useStorage:  true,
			saveStorage: true,
			Config:      config,
		})
		require.Equal(t, 1, len(cm))
		require.Equal(t, 1, len(storage.Words))
	})
	t.Run("new func", func(t *testing.T) {
		t.Parallel()
		var filename = "/tmp/newfunc.json"
		storage := newMock()
		storage.getFilename = func(s string) string { return filename }
		defer os.Remove(filename)
		cm := storage.NewCoMatrixes(context.Background(), handler, storagePermission{
			useStorage:  false, // !
			saveStorage: true,
			Config:      config,
		})

		// created
		_, err := os.ReadFile(filename)
		if err != nil {
			require.Equal(t, io.EOF, err)
		}
		require.Equal(t, 1, len(cm))
		require.Equal(t, 1, len(storage.Words))
	})
}
