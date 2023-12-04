package network

import (
	"context"
	"testing"
	"yyyoichi/Collo-API/internal/pair"

	"github.com/stretchr/testify/require"
)

func TestNetwork(t *testing.T) {
	t.Run("AddNode", func(t *testing.T) {
		tnetwork := NewNetwork()
		node := tnetwork.addNode("hoge")
		require.Equal(t, node.Word, NodeWord("hoge"))
		require.Equal(t, node.NodeID, NodeID(1))
		require.Equal(t, len(node.edges), 0)
		require.NotNil(t, tnetwork.Nodes[1])
		require.NotNil(t, tnetwork.nodesByWord["hoge"])

		tnetwork.addNode("fuga")
		require.NotNil(t, tnetwork.Nodes[2])
		require.NotNil(t, tnetwork.nodesByWord["fuga"])

		tnetwork.addNode("hoge")
		require.Nil(t, tnetwork.Nodes[3])
	})
	t.Run("AddEdge", func(t *testing.T) {
		tnetwork := NewNetwork()
		hoge := tnetwork.addNode("hoge")
		fuga := tnetwork.addNode("fuga")
		edge := tnetwork.addEdge(hoge, fuga)
		require.Equal(t, edge.EdgeID, EdgeID(1))
		require.Equal(t, edge.Count, uint(1))
		require.Equal(t, edge.NodeID1, NodeID(1))
		require.Equal(t, edge.NodeID2, NodeID(2))
		require.Equal(t, hoge.edges[fuga.NodeID], edge)
		require.Equal(t, fuga.edges[hoge.NodeID], edge)
		require.Equal(t, tnetwork.Edges[edge.EdgeID], edge)

		tnetwork.addEdge(hoge, hoge)
		require.Nil(t, tnetwork.Edges[EdgeID(2)])
		tnetwork.addEdge(fuga, hoge)
		require.Equal(t, edge.Count, uint(2))

		tnetwork.addEdge(hoge, tnetwork.addNode("foo"))
		require.NotNil(t, tnetwork.Edges[EdgeID(2)])
	})
	t.Run("AddNetwork", func(t *testing.T) {
		tnetwork := NewNetwork()
		tnetwork.AddNetwork(context.Background(), "foo", "bar", "baz", "foo")
		require.Equal(t, len(tnetwork.Nodes), 3)
		require.Equal(t, len(tnetwork.Edges), 3)
	})
	t.Run("GetNetworkAround", func(t *testing.T) {
		tnetwork := NewNetwork()
		tnetwork.AddNetwork(context.Background(), "foo", "bar", "baz", "foo")
		foo := tnetwork.nodesByWord["foo"]
		nodes, edges := tnetwork.GetNetworkAround(uint(foo.NodeID))
		require.Equal(t, len(nodes), 2)
		require.Equal(t, len(edges), 2)

		nodes, edges = tnetwork.GetNetworkAround(uint(99))
		require.Nil(t, nodes)
		require.Nil(t, edges)
	})
	t.Run("GetCenterNode", func(t *testing.T) {
		tnetwork := NewNetwork()
		node := tnetwork.addNode("foo")
		tnetwork.AddNetwork(context.Background(), "foo", "bar")
		tnetwork.AddNetwork(context.Background(), "foo", "baz")

		// exp foo
		require.Equal(t, node.NodeID, tnetwork.GetCenterNodeID())
	})

	t.Run("GetByWord", func(t *testing.T) {
		tnetwork := NewNetwork()
		node := tnetwork.addNode("foo")
		nodeID, found := tnetwork.GetByWord("foo")
		require.True(t, found)
		require.Equal(t, node.NodeID, nodeID)

		fr := pair.MAnalytics.Parse("学と學")
		require.NoError(t, fr.Error())
		nouns := fr.GetNouns()
		require.Equal(t, len(nouns), 2)
		require.Equal(t, nouns[0], nouns[1])

		node = tnetwork.addNode(NodeWord(nouns[0]))
		np := NetworkProvider{network: tnetwork}
		nodeID = np.GetByWord("學")
		require.Equal(t, node.NodeID, nodeID)
	})
}
