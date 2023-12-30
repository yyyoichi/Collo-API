package provider

import (
	"context"
	"sync"
	"yyyoichi/Collo-API/internal/analyzer"
	apiv2 "yyyoichi/Collo-API/internal/api/v2"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/pkg/stream"
)

type V2MultiRateProvider struct {
	allProvider *V2RateProvider
	providers   []*V2RateProvider
}

func NewV2NultiRateProvider(
	ctx context.Context,
	ndlConfig ndl.Config,
	analyzerConfig analyzer.Config,
	matrixConofig matrix.Config,
	handler Handler[*apiv2.ColloRateWebStreamResponse],
) *V2MultiRateProvider {
	// すべての文書を対象とした共起行列
	allProvder := &V2RateProvider{
		handler: handler,
	}
	var errorHook stream.ErrorHook = func(err error) {
		allProvder.handler.Err(err)
	}
	// fetch -> doc-word matrix
	n, docCh := generateDocument(ctx, errorHook, ndlConfig, analyzerConfig)
	allProvder.handleRespTotalProcess(n * 2) // fetch回分と、allprovider+それぞれのco-matrix計算分
	b := matrix.NewBuilder()
	for doc := range docCh {
		b.AppendDocument(doc)
		allProvder.handleRespProcess()
	}
	numGroups, allMatrix, groupMatrixCh := matrix.NewMultiCoMatrixFromBuilder(ctx, b, matrixConofig)
	// allの進捗を残り1を最後にリセット
	allMatrix.As("all")
	allProvder.m = allMatrix

	var wg sync.WaitGroup
	startMatrixConsume := func(p *V2RateProvider) {
		defer wg.Done()
		for pg := range p.m.ConsumeProgress() {
			switch pg {
			case matrix.ErrDone:
				p.handler.Err(p.m.Error())
			case matrix.ProgressDone:
				p.handleRespProcess()
			default:
				// allの進捗を更新
				allProvder.handleRespProcess()
			}
		}
	}
	// init returns
	multiProvider := &V2MultiRateProvider{
		allProvider: allProvder,
		providers:   []*V2RateProvider{},
	}
	wg.Add(numGroups + 1)
	go startMatrixConsume(multiProvider.allProvider)

	for groupMatrix := range groupMatrixCh {
		p := &V2RateProvider{
			handler:      handler,
			m:            groupMatrix,
			totalProcess: 2,
			doneProcess:  1,
		}
		resp := p.createResp([]*matrix.Node{}, []*matrix.Edge{})
		go p.handler.Resp(resp)
		multiProvider.providers = append(multiProvider.providers, p)
		go startMatrixConsume(p)
	}
	wg.Wait()
	return multiProvider
}

func (p *V2MultiRateProvider) StreamTop3NodesCoOccurrence() {
	top1 := p.allProvider.m.NodeRank(0)
	top2 := p.allProvider.m.NodeRank(1)
	top3 := p.allProvider.m.NodeRank(2)

	var wg sync.WaitGroup
	handleResp := func(provider *V2RateProvider) {
		defer wg.Done()
		nodes, edges := provider.m.CoOccurrences(top1.ID, top2.ID, top3.ID)
		nodes = append(nodes, top1, top2, top3)
		provider.handleResp(nodes, edges)
	}

	wg.Add(1 + len(p.providers))
	go handleResp(p.allProvider)
	for _, provider := range p.providers {
		go handleResp(provider)
	}
}

func (p *V2MultiRateProvider) StreamForcusNodeIDCoOccurrence(nodeID uint) {
	var wg sync.WaitGroup
	handleResp := func(provider *V2RateProvider) {
		defer wg.Done()
		nodes, edges := provider.m.CoOccurrenceRelation(nodeID)
		provider.handleResp(nodes, edges)
	}

	wg.Add(1 + len(p.providers))
	go handleResp(p.allProvider)
	for _, provider := range p.providers {
		go handleResp(provider)
	}
}
