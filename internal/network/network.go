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
	Nodes       map[NodeID]*Node `json:"nodes"`
	Edges       map[EdgeID]*Edge `json:"edges"`
	mu          sync.RWMutex
}

func NewNetwork() *Network {
	return &Network{
		nodesByWord: map[NodeWord]*Node{},
		Nodes:       map[NodeID]*Node{},
		Edges:       map[EdgeID]*Edge{},
		mu:          sync.RWMutex{},
	}
}

func (nw *Network) refreshMap() {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	for _, node := range nw.Nodes {
		node.edges = map[NodeID]*Edge{}
		nw.nodesByWord[node.Word] = node
	}

	ctx := context.Background()
	edgeCh := nw.generateEdge(ctx)
	doneCh := stream.Line[*Edge, interface{}](ctx, edgeCh, func(edge *Edge) interface{} {
		// 各nodeにedgeを追加する
		set := func(i, o NodeID) {
			if node, found := nw.Nodes[i]; !found {
				return
			} else {
				node.edges[o] = edge
			}
		}
		set(edge.NodeID1, edge.NodeID2)
		set(edge.NodeID2, edge.NodeID1)
		return struct{}{}
	})
	for range doneCh {
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
	node, found := nw.Nodes[NodeID(nodeID)]
	if !found {
		return nil, nil
	}

	nodes = make([]*Node, len(node.edges))
	edges = make([]*Edge, len(node.edges))
	i := 0
	for nodeID, edge := range node.edges {
		nodes[i] = nw.Nodes[nodeID]
		edges[i] = edge
		i++
	}

	return nodes, edges
}

// 共起語の種類が最も多いノードのIDを返す
func (nw *Network) GetCenterNodeID() NodeID {
	var nodeID NodeID
	var max int
	for id, node := range nw.Nodes {
		if max < len(node.edges) {
			max = len(node.edges)
			nodeID = id
		}
	}
	return nodeID
}

// 単語からノードIDを返す
func (nw *Network) GetByWord(word string) (NodeID, bool) {
	if node, found := nw.nodesByWord[NodeWord(word)]; found {
		return node.NodeID, true
	}
	return 0, false
}

func (nw *Network) addNode(word NodeWord) *Node {
	nw.mu.Lock()
	defer nw.mu.Unlock()
	if node, found := nw.nodesByWord[word]; found {
		return node
	}
	node := &Node{
		NodeID: NodeID(len(nw.Nodes) + 1),
		Word:   word,
		edges:  map[NodeID]*Edge{},
	}
	nw.nodesByWord[word] = node
	nw.Nodes[node.NodeID] = node
	return node
}

func (nw *Network) addEdge(nodeA, nodeB *Node) *Edge {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	if nodeA.NodeID == nodeB.NodeID {
		return nil
	}
	if edge, found := nodeA.edges[nodeB.NodeID]; found {
		return edge.countUP()
	}

	var nodeID1, nodeID2 NodeID
	if nodeA.NodeID < nodeB.NodeID {
		nodeID1 = nodeA.NodeID
		nodeID2 = nodeB.NodeID
	} else {
		nodeID1 = nodeB.NodeID
		nodeID2 = nodeA.NodeID
	}
	edge := &Edge{
		EdgeID:  EdgeID(len(nw.Edges) + 1),
		NodeID1: nodeID1,
		NodeID2: nodeID2,
		Count:   1,

		mu: sync.RWMutex{},
	}
	nw.Edges[edge.EdgeID] = edge
	nodeA.edges[nodeB.NodeID] = edge
	nodeB.edges[nodeA.NodeID] = edge
	return edge
}

type Node struct {
	NodeID NodeID   `json:"id"`
	Word   NodeWord `json:"word"`
	edges  map[NodeID]*Edge
}

type Edge struct {
	EdgeID EdgeID `json:"id"`
	// 必ずid1 < id2
	NodeID1 NodeID `json:"node_id1"`
	NodeID2 NodeID `json:"node_id2"`
	Count   uint   `json:"count"`

	mu sync.RWMutex
}

func (e *Edge) countUP() *Edge {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Count++
	return e
}

func (nw *Network) generateEdge(cxt context.Context) <-chan *Edge {
	ch := make(chan *Edge, len(nw.Edges))
	go func() {
		defer close(ch)
		for _, val := range nw.Edges {
			select {
			case <-cxt.Done():
				return
			case ch <- val:
			}
		}
	}()

	return ch
}
