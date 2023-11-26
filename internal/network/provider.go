package network

import (
	"context"
	"yyyoichi/Collo-API/internal/pair"
	"yyyoichi/Collo-API/pkg/stream"
)

type FetchError struct{ error }
type ParseError struct{ error }

func NewNetworkProvider(
	ctx context.Context,
	kokkaiRequestConfig pair.Config,
	handler Handler,
) *NetwrokProvider {
	np := &NetwrokProvider{
		network: NewNetwork(),
	}

	speech, err := pair.NewSpeech(kokkaiRequestConfig)
	if err != nil {
		np.handler.Err(FetchError{err})
		return nil
	}
	urls := speech.GetURLs()
	// 必要数セット
	np.needKokkaiFetch = uint8(len(urls))
	// 必要数送信
	go np.handler.Resp(0)

	urlCh := stream.Generator[string](ctx, urls...)
	fetchResultCh := stream.Line[string, *pair.FetchResult](ctx, urlCh, speech.Fetch)
	doneCh := stream.FunIO[*pair.FetchResult, int](ctx, fetchResultCh,
		func(fr *pair.FetchResult) int {
			// フェッチ1回分の処理
			if fr.Error() != nil {
				np.handler.Err(ParseError{fr.Error()})
			}
			speechCh := fr.GenerateSpeech(ctx)
			parseResultCh := stream.Line[string, *pair.ParseResult](ctx, speechCh, func(s string) *pair.ParseResult {
				pr := pair.MAnalytics.Parse(s)
				if pr.Error() != nil {
					np.handler.Err(ParseError{fr.Error()})
				}
				return pr
			})
			stream.Line[*pair.ParseResult, struct{}](ctx, parseResultCh, func(pr *pair.ParseResult) struct{} {
				nouns := pr.GetNouns()
				np.network.AddNetwork(ctx, nouns...)
				return struct{}{}
			})
			// fetch終了
			return 1
		},
	)

	stream.Line[int, int](ctx, doneCh, func(int) int {
		np.doneKokkaiCount++
		// 完了数送信
		go np.handler.Resp(0)
		return 0
	})

	// リクエストされた単語に関連するnodeとedgeを送信する
	go np.streamNetworksWith(kokkaiRequestConfig.Search.Any)
	return np
}

type NetwrokProvider struct {
	network *Network
	handler Handler

	// 議事録APIフェッチ必要数
	needKokkaiFetch uint8
	// 議事録APIフェッチ済数
	doneKokkaiCount uint8
}

// [word]とそれに関連するnodeとedgeを送信する
func (np *NetwrokProvider) streamNetworksWith(word string) {
	node := np.network.nodesByWord[NodeWord(word)]
	_, _ = np.network.GetNetworkAround(uint(node.nodeID))
	// add 'node'
	np.handler.Resp(0)
}

// [nodeID]に関連するnodeとedgeを送信する
func (np *NetwrokProvider) StreamNetworksAround(nodeID uint) {
	np.network.GetNetworkAround(nodeID)
	np.handler.Resp(0)
}

type Handler struct {
	Err  func(error)
	Done func()
	Resp func(any)
}
