package network

import (
	"context"
	"sync"
	"yyyoichi/Collo-API/pkg/stream"
)

type NodeID uint
type EdgeID uint
type NodeWord string

type Network struct {
	nodes map[NodeWord]*Node
	edges map[EdgeID]*Edge
	mu    sync.RWMutex
}

func NewNetwork() *Network {
	return &Network{nodes: map[NodeWord]*Node{}}
}
func (nw *Network) AddNodes(ctx context.Context, words ...string) {
	nw.mu.Lock()
	defer nw.mu.Unlock()
	nodeCh := stream.GeneratorWithFn[string, *Node](
		ctx,
		func(w string) *Node {
			word := NodeWord(w)
			if node, found := nw.nodes[word]; found {
				return node
			}
			return nw.addNode(word)
		},
		words...,
	)

	nodes := []*Node{}
	for nodeA := range nodeCh {
		for _, nodeB := range nodes {
			nw.addEdge(nodeA, nodeB)
		}
		nodes = append(nodes, nodeA)
	}
}

func (nw *Network) addNode(word NodeWord) *Node {
	nw.mu.Lock()
	defer nw.mu.Unlock()
	node := &Node{
		nodeID: NodeID(len(nw.nodes)),
		word:   word,
		edges:  map[NodeID]*Edge{},
	}
	nw.nodes[word] = node
	return node
}

func (nw *Network) addEdge(nodeA, nodeB *Node) *Edge {
	if nodeA.nodeID == nodeB.nodeID {
		return nil
	}
	if edge, found := nw.foundEdge(nodeA, nodeB); found {
		return edge.countUP()
	}

	var nodeID1, nodeID2 NodeID
	if nodeA.nodeID < nodeB.nodeID {
		nodeID1 = nodeA.nodeID
		nodeID2 = nodeB.nodeID
	} else {
		nodeID1 = nodeB.nodeID
		nodeID2 = nodeA.nodeID
	}
	edge := &Edge{
		edgeID:  EdgeID(len(nw.edges)),
		nodeID1: nodeID1,
		nodeID2: nodeID2,
		count:   1,

		mu: sync.RWMutex{},
	}
	nw.mu.Lock()
	defer nw.mu.Unlock()
	nw.edges[edge.edgeID] = edge
	nodeA.edges[nodeB.nodeID] = edge
	nodeB.edges[nodeA.nodeID] = edge
	return edge
}

func (nw *Network) foundEdge(nodeA, nodeB *Node) (edge *Edge, found bool) {
	edge, found = nodeA.edges[nodeB.nodeID]
	return edge, found
}

type Node struct {
	nodeID NodeID
	word   NodeWord
	edges  map[NodeID]*Edge
}

type Edge struct {
	edgeID EdgeID
	// 必ずid1 < id2
	nodeID1, nodeID2 NodeID
	count            uint

	mu sync.RWMutex
}

func (e *Edge) countUP() *Edge {
	e.mu.RLock()
	defer e.mu.Unlock()
	e.count++
	return e
}
