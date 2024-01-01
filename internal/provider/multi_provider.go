package provider

import (
	"context"
	"math"
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
	numDocs, docCh := generateDocument(ctx, errorHook, ndlConfig, analyzerConfig)
	process := newMultiProcess(allProvider, numDocs)
	b := matrix.NewBuilder()
	for doc := range docCh {
		if doc != nil && len(doc.Words) > 0 {
			b.AppendDocument(doc)
		}
		process.doneDoc()
	}
	_, allMatrix, groupMatrixCh := matrix.NewMultiCoMatrixFromBuilder(ctx, b, matrixConofig)
	allMatrix.As("all")
	allProvider.m = allMatrix
	// init returns
	multiProvider := &V2MultiRateProvider{
		providers: []*V2RateProvider{allProvider},
	}
	for m := range groupMatrixCh {
		p := &V2RateProvider{
			handler:      handler,
			m:            m,
			totalProcess: 2,
			doneProcess:  1,
		}
		resp := p.createResp([]*matrix.Node{}, []*matrix.Edge{})
		p.handler.Resp(resp) // send new group
		multiProvider.providers = append(multiProvider.providers, p)
	}
	// provider数セット
	process.setNumProviders(len(multiProvider.providers))

	providerCh := stream.Generator[*V2RateProvider](ctx, multiProvider.providers...)
	doneCh := stream.FunIO[*V2RateProvider, interface{}](ctx, providerCh, func(p *V2RateProvider) interface{} {
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
	})

	for range doneCh {
		// allの進捗を更新
		process.doneProvider()
	}
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
	ctx := context.Background()

	pCh := stream.Generator[*V2RateProvider](ctx, p.providers...)
	for range stream.Line[*V2RateProvider, interface{}](ctx, pCh, handleResp) {
	}
}

type multiProcess struct {
	all      *V2RateProvider
	numDocs  float64
	doneDocs float64

	numProviders  float64
	doneProviders float64
}

func newMultiProcess(allProvider *V2RateProvider, numDocs int) *multiProcess {
	p := &multiProcess{
		all:     allProvider,
		numDocs: float64(numDocs),
	}
	p.sendTotal()
	return p
}

func (p *multiProcess) sendTotal() {
	// docsがすべて送信後、50%の進捗となるようにする。
	p.all.mu.Lock()
	defer p.all.mu.Unlock()
	p.all.totalProcess = 200
	resp := p.all.createResp([]*matrix.Node{}, []*matrix.Edge{})
	p.all.handler.Resp(resp)
}

func (p *multiProcess) doneDoc() {
	p.all.mu.Lock()
	defer p.all.mu.Unlock()
	p.doneDocs += 1.0
	p.all.doneProcess = int(math.Round(p.doneDocs / p.numDocs * 100.0))
	resp := p.all.createResp([]*matrix.Node{}, []*matrix.Edge{})
	p.all.handler.Resp(resp)
}

func (p *multiProcess) setNumProviders(n int) {
	p.all.mu.Lock()
	defer p.all.mu.Unlock()
	p.numProviders = float64(n)
}

func (p *multiProcess) doneProvider() {
	p.all.mu.Lock()
	defer p.all.mu.Unlock()
	p.doneProviders += 1.0
	p.all.doneProcess = 100 + int(math.Round(p.doneProviders/p.numProviders*100.0))
	resp := p.all.createResp([]*matrix.Node{}, []*matrix.Edge{})
	p.all.handler.Resp(resp)
}
