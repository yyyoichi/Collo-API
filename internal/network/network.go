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
	nodesByWord map[NodeWord]*Node
	nodes       map[NodeID]*Node
	edges       map[EdgeID]*Edge
	mu          sync.RWMutex
}

func NewNetwork() *Network {
	return &Network{
		nodesByWord: map[NodeWord]*Node{},
		nodes:       map[NodeID]*Node{},
		edges:       map[EdgeID]*Edge{},
		mu:          sync.RWMutex{},
	}
}

func (nw *Network) AddNetwork(ctx context.Context, words ...string) {
	nodeCh := stream.GeneratorWithFn[string, *Node](
		ctx,
		func(word string) *Node {
			return nw.addNode(NodeWord(word))
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

// [nodeID]に関連するNodeとEdgeを返す
func (nw *Network) GetNetworkAround(nodeID uint) (nodes []*Node, edges []*Edge) {
	node, found := nw.nodes[NodeID(nodeID)]
	if !found {
		return nil, nil
	}

	nodes = make([]*Node, len(node.edges))
	edges = make([]*Edge, len(node.edges))
	i := 0
	for nodeID, edge := range node.edges {
		nodes[i] = nw.nodes[nodeID]
		edges[i] = edge
		i++
	}

	return nodes, edges
}

func (nw *Network) addNode(word NodeWord) *Node {
	nw.mu.Lock()
	defer nw.mu.Unlock()
	if node, found := nw.nodesByWord[word]; found {
		return node
	}
	node := &Node{
		nodeID: NodeID(len(nw.nodes)),
		word:   word,
		edges:  map[NodeID]*Edge{},
	}
	nw.nodesByWord[word] = node
	nw.nodes[node.nodeID] = node
	return node
}

func (nw *Network) addEdge(nodeA, nodeB *Node) *Edge {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	if nodeA.nodeID == nodeB.nodeID {
		return nil
	}
	if edge, found := nodeA.edges[nodeB.nodeID]; found {
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
	nw.edges[edge.edgeID] = edge
	nodeA.edges[nodeB.nodeID] = edge
	nodeB.edges[nodeA.nodeID] = edge
	return edge
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
	e.mu.Lock()
	defer e.mu.Unlock()
	e.count++
	return e
}
