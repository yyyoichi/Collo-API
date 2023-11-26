package network

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetwork(t *testing.T) {
	t.Run("AddNode", func(t *testing.T) {
		tnetwork := NewNetwork()
		node := tnetwork.addNode("hoge")
		require.Equal(t, node.word, NodeWord("hoge"))
		require.Equal(t, node.nodeID, NodeID(0))
		require.Equal(t, len(node.edges), 0)
		require.NotNil(t, tnetwork.nodes[0])
		require.NotNil(t, tnetwork.nodesByWord["hoge"])

		tnetwork.addNode("fuga")
		require.NotNil(t, tnetwork.nodes[1])
		require.NotNil(t, tnetwork.nodesByWord["fuga"])

		tnetwork.addNode("hoge")
		require.Nil(t, tnetwork.nodes[2])
	})
	t.Run("AddEdge", func(t *testing.T) {
		tnetwork := NewNetwork()
		hoge := tnetwork.addNode("hoge")
		fuga := tnetwork.addNode("fuga")
		edge := tnetwork.addEdge(hoge, fuga)
		require.Equal(t, edge.edgeID, EdgeID(0))
		require.Equal(t, edge.count, uint(1))
		require.Equal(t, edge.nodeID1, NodeID(0))
		require.Equal(t, edge.nodeID2, NodeID(1))
		require.Equal(t, hoge.edges[fuga.nodeID], edge)
		require.Equal(t, fuga.edges[hoge.nodeID], edge)
		require.Equal(t, tnetwork.edges[edge.edgeID], edge)

		tnetwork.addEdge(hoge, hoge)
		require.Nil(t, tnetwork.edges[EdgeID(1)])
		tnetwork.addEdge(fuga, hoge)
		require.Equal(t, edge.count, uint(2))

		tnetwork.addEdge(hoge, tnetwork.addNode("foo"))
		require.NotNil(t, tnetwork.edges[EdgeID(1)])
	})
	t.Run("AddNetwork", func(t *testing.T) {
		tnetwork := NewNetwork()
		tnetwork.AddNetwork(context.Background(), "foo", "bar", "baz", "foo")
		require.Equal(t, len(tnetwork.nodes), 3)
		require.Equal(t, len(tnetwork.edges), 3)
	})
	t.Run("GetNetworkAround", func(t *testing.T) {
		tnetwork := NewNetwork()
		tnetwork.AddNetwork(context.Background(), "foo", "bar", "baz", "foo")
		foo := tnetwork.nodesByWord["foo"]
		nodes, edges := tnetwork.GetNetworkAround(uint(foo.nodeID))
		require.Equal(t, len(nodes), 2)
		require.Equal(t, len(edges), 2)

		nodes, edges = tnetwork.GetNetworkAround(uint(99))
		require.Nil(t, nodes)
		require.Nil(t, edges)
	})
}
