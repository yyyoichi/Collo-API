package network

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Run("GetSet", func(t *testing.T) {
		dir, _ := os.MkdirTemp("", "test-network-manager")
		defer os.RemoveAll(dir)

		NManager.dir = dir
		// memory保存期間なし
		NManager.ttl = time.Duration(-1 * time.Second)
		tnetwork := NewNetwork()
		tnetwork.addNode("hoge")
		NManager.Set("hoge", tnetwork)
		// get method
		actNetwork, found := NManager.Get("hoge")
		require.True(t, found)
		require.Equal(t, tnetwork, actNetwork)
		// found on memory
		d, found := NManager.data[md5Hash("hoge")]
		require.True(t, found)
		require.NotNil(t, d)
		require.Equal(t, tnetwork, d.Network)
		// found file
		f, err := os.OpenFile(path.Join(fmt.Sprintf("%s/%s.json", dir, md5Hash("hoge"))), os.O_RDWR, 0600)
		require.NoError(t, err)
		require.NotNil(t, f)

		// run delete
		NManager.tick = time.Duration(time.Microsecond)
		NManager.StartCleanup()
		time.Sleep(time.Millisecond * 50)

		// not found on memory
		_, found = NManager.data[md5Hash("hoge")]
		require.False(t, found)
		// found file
		f, err = os.OpenFile(path.Join(fmt.Sprintf("%s/%s.json", dir, md5Hash("hoge"))), os.O_RDWR, 0600)
		require.NoError(t, err)
		require.NotNil(t, f)

		NManager.StopCleanup()
		n, _ := NManager.Get("hoge")
		require.Equal(t, len(n.Nodes), 1)
		require.Equal(t, n.Nodes[NodeID(0)].Word, NodeWord("hoge"))
	})
}
