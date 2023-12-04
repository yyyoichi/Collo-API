package network

import (
	"context"
	"errors"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/pair"
	"yyyoichi/Collo-API/pkg/stream"
)

type FetchError struct{ error }
type ParseError struct{ error }

func NewNetworkProvider(
	ctx context.Context,
	kokkaiRequestConfig pair.Config,
	handler Handler,
) *NetworkProvider {
	np := &NetworkProvider{
		network: NewNetwork(),
		handler: handler,
	}

	speech, err := pair.NewSpeech(kokkaiRequestConfig)
	if err != nil {
		np.handler.Err(FetchError{err})
		return nil
	}

	// 必要数セット
	np.needKokkaiFetch = uint8(len(speech.GetURLs()))
	// storegeから検索
	if network, found := NManager.Get(speech.GetInitURL()); found {
		np.doneKokkaiCount = np.needKokkaiFetch
		np.network = network
		return np
	}

	// 必要数送信
	go np.handleResp([]*Node{}, []*Edge{})

	// create network
	urlCh := speech.GenerateURL(ctx)
	fetchResultCh := stream.Line[string, *pair.FetchResult](ctx, urlCh, speech.Fetch)
	doneCh := stream.FunIO[*pair.FetchResult, int](ctx, fetchResultCh,
		func(fr *pair.FetchResult) int {
			// フェッチ1回分の処理
			if fr.Error() != nil {
				np.handler.Err(FetchError{fr.Error()})
			}
			speechCh := fr.GenerateSpeech(ctx)
			nounsCh := stream.Line[string, []string](ctx, speechCh, func(s string) []string {
				pr := pair.MAnalytics.Parse(s)
				if pr.Error() != nil {
					np.handler.Err(ParseError{fr.Error()})
				}
				return pr.GetNouns()
			})
			for nouns := range nounsCh {
				np.network.AddNetwork(ctx, nouns...)
			}
			// fetch終了
			return 1
		},
	)

	for range doneCh {
		np.doneKokkaiCount++
		// 完了数送信
		go np.handleResp([]*Node{}, []*Edge{})
	}

	return np
}

type NetworkProvider struct {
	network *Network
	handler Handler

	// 議事録APIフェッチ必要数
	needKokkaiFetch uint8
	// 議事録APIフェッチ済数
	doneKokkaiCount uint8
}

// [nodeID]とそれに関連するnodeとedgeを送信する
func (np *NetworkProvider) StreamNetworksWith(nodeID NodeID) {
	node, found := np.network.Nodes[nodeID]
	if !found {
		np.handler.Err(errors.New("expect node, but not found"))
		return
	}
	nodes, edges := np.network.GetNetworkAround(uint(node.NodeID))
	nodes = append(nodes, node)
	np.handleResp(nodes, edges)
}

// [nodeID]に関連するnodeとedgeを送信する
func (np *NetworkProvider) StreamNetworksAround(nodeID uint) {
	nodes, edges := np.network.GetNetworkAround(nodeID)
	np.handleResp(nodes, edges)
}

func (np *NetworkProvider) handleResp(nodes []*Node, edges []*Edge) {
	resp := &apiv2.ColloNetworkStreamResponse{
		Dones: uint32(np.doneKokkaiCount),
		Needs: uint32(np.needKokkaiFetch),
		Nodes: []*apiv2.Node{},
		Edges: []*apiv2.Edge{},
	}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &apiv2.Node{
			NodeId: uint32(node.NodeID),
			Word:   string(node.Word),
		})
	}
	for _, edge := range edges {
		resp.Edges = append(resp.Edges, &apiv2.Edge{
			EdgeId:  uint32(edge.EdgeID),
			NodeId1: uint32(edge.NodeID1),
			NodeId2: uint32(edge.NodeID2),
			Count:   uint32(edge.Count),
		})
	}
	np.handler.Resp(resp)
}

type Handler struct {
	Err  func(error)
	Done func()
	Resp func(*apiv2.ColloNetworkStreamResponse)
}
