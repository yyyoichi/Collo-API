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
	providers []*V2RateProvider
}

func NewV2NultiRateProvider(
	ctx context.Context,
	ndlConfig ndl.Config,
	analyzerConfig analyzer.Config,
	matrixConofig matrix.Config,
	handler Handler[*apiv2.ColloRateWebStreamResponse],
) *V2MultiRateProvider {
	// すべての文書を対象とした共起行列
	allProvider := &V2RateProvider{
		handler: handler,
	}
	var errorHook stream.ErrorHook = func(err error) {
		allProvider.handler.Err(err)
	}
	// fetch -> doc-word matrix
	n, docCh := generateDocument(ctx, errorHook, ndlConfig, analyzerConfig)
	allProvider.handleRespTotalProcess(n * 2) // fetch回分と、allprovider+それぞれのco-matrix計算分
	b := matrix.NewBuilder()
	for doc := range docCh {
		b.AppendDocument(doc)
		allProvider.handleRespProcess()
	}
	_, allMatrix, groupMatrixCh := matrix.NewMultiCoMatrixFromBuilder(ctx, b, matrixConofig)
	allMatrix.As("all")
	allProvider.m = allMatrix

	startMatrixConsume := func(p *V2RateProvider) interface{} {
		for pg := range p.m.ConsumeProgress() {
			switch pg {
			case matrix.ErrDone:
				p.handler.Err(p.m.Error())
			case matrix.ProgressDone:
				p.handleRespProcess()
			default:
			}
		}
		return struct{}{}
	}
	// init returns
	multiProvider := &V2MultiRateProvider{
		providers: []*V2RateProvider{allProvider},
	}

	// all 用
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		startMatrixConsume(allProvider)
	}()

	for range stream.FunIO[*matrix.CoMatrix, interface{}](ctx, groupMatrixCh, func(cm *matrix.CoMatrix) interface{} {
		p := &V2RateProvider{
			handler:      handler,
			m:            cm,
			totalProcess: 2,
			doneProcess:  1,
		}
		resp := p.createResp([]*matrix.Node{}, []*matrix.Edge{})
		p.handler.Resp(resp)
		multiProvider.providers = append(multiProvider.providers, p)
		return startMatrixConsume(p)
	}) {
		// allの進捗を更新
		allProvider.handleRespProcess()
	}
	// all 用
	wg.Wait()
	return multiProvider
}

func (p *V2MultiRateProvider) StreamTop3NodesCoOccurrence() {
	if len(p.providers) == 0 {
		return
	}
	top1 := p.providers[0].m.NodeRank(0)
	top2 := p.providers[0].m.NodeRank(1)
	top3 := p.providers[0].m.NodeRank(2)
	if len(p.providers) == 1 {
		p.providers[0].StreamTop3NodesCoOccurrence()
		return
	}

	handleResp := func(provider *V2RateProvider) interface{} {
		nodes, edges := provider.m.CoOccurrences(top1.ID, top2.ID, top3.ID)
		nodes = append(nodes, top1, top2, top3)
		provider.handleResp(nodes, edges)
		return struct{}{}
	}
	ctx := context.Background()

	pCh := stream.Generator[*V2RateProvider](ctx, p.providers...)
	for range stream.Line[*V2RateProvider, interface{}](ctx, pCh, handleResp) {
	}
}

func (p *V2MultiRateProvider) StreamForcusNodeIDCoOccurrence(nodeID uint) {
	if len(p.providers) == 0 {
		return
	}
	if len(p.providers) == 1 {
		p.providers[0].StreamForcusNodeIDCoOccurrence(nodeID)
		return
	}
	handleResp := func(provider *V2RateProvider) interface{} {
		nodes, edges := provider.m.CoOccurrenceRelation(nodeID)
		provider.handleResp(nodes, edges)
		return struct{}{}
	}

	for _, provider := range p.providers {
		go handleResp(provider)
	}
	ctx := context.Background()

	pCh := stream.Generator[*V2RateProvider](ctx, p.providers...)
	for range stream.Line[*V2RateProvider, interface{}](ctx, pCh, handleResp) {
	}
}
