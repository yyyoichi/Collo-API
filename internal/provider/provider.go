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

type V2RateProvider struct {
	handler Handler[*apiv2.ColloRateWebStreamResponse]
	m       *matrix.CoMatrix

	// 完了プロセス用
	mu sync.Mutex
	// 必要process
	totalProcess int
	// 完了プロセス
	doneProcess int
}

func NewV2RateProvider(
	ctx context.Context,
	ndlConfig ndl.Config,
	analyzerConfig analyzer.Config,
	handler Handler[*apiv2.ColloRateWebStreamResponse],
) *V2RateProvider {
	p := &V2RateProvider{
		handler: handler,
	}
	b := p.getWightDocMatrix(ctx, ndlConfig, analyzerConfig)
	m := matrix.NewCoMatrixFromBuilder(b, matrix.Config{})
	for pg := range m.ConsumeProgress() {
		switch pg {
		case matrix.ErrDone:
			p.handler.Err(m.Error())
		case matrix.ProgressDone:
			p.handleRespProcess()
		default:
		}
	}
	p.m = m
	return p
}

func (p *V2RateProvider) getWightDocMatrix(
	ctx context.Context,
	ndlConfig ndl.Config,
	analyzerConfig analyzer.Config,
) *matrix.Builder {
	m := ndl.NewMeeting(ndlConfig)

	// エラー発生時Errorを送信する
	var errorHook stream.ErrorHook = func(err error) {
		p.handler.Err(err)
	}

	// 総処理数送信
	go p.handleRespTotalProcess(m.GetNumberOfRecords())
	// 会議APIから結果取得
	meetingResultCh := m.GenerateMeeting(ctx)
	// 会議ごとの発言
	meetingCh := stream.DemultiWithErrorHook[*ndl.MeetingResult, string](
		ctx,
		errorHook,
		meetingResultCh,
		func(mr *ndl.MeetingResult) []string {
			return mr.GetSpeechsPerMeeting()
		})
	// 形態素解析
	analysisResultCh := stream.FunIO[string, *analyzer.AnalysisResult](
		ctx,
		meetingCh,
		analyzer.Analysis,
	)
	// 会議ごとの単語
	wordsCh := stream.LineWithErrorHook[*analyzer.AnalysisResult, []string](
		ctx,
		errorHook,
		analysisResultCh,
		func(ar *analyzer.AnalysisResult) []string {
			return ar.Get(analyzerConfig)
		})

	builder := matrix.NewBuilder()
	for words := range wordsCh {
		builder.AppendDoc(words)
		// 処理完了済み通知
		go p.handleRespProcess()
	}
	return builder
}

func (p *V2RateProvider) StreamTop3NodesCoOccurrence() {
	top1 := p.m.NodeRank(0)
	top2 := p.m.NodeRank(1)
	top3 := p.m.NodeRank(2)
	nodes, edges := p.m.CoOccurrences(top1.ID, top2.ID, top3.ID)
	nodes = append(nodes, top1, top2, top3)
	p.handleResp(nodes, edges)
}

func (p *V2RateProvider) StreamForcusNodeIDCoOccurrence(nodeID uint) {
	nodes, edges := p.m.CoOccurrenceRelation(nodeID)
	p.handleResp(nodes, edges)
}

func (p *V2RateProvider) handleResp(nodes []*matrix.Node, edges []*matrix.Edge) {
	resp := &apiv2.ColloRateWebStreamResponse{
		Dones: uint32(p.doneProcess),
		Needs: uint32(p.totalProcess),
		Nodes: []*apiv2.RateNode{},
		Edges: []*apiv2.RateEdge{},
	}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &apiv2.RateNode{
			NodeId: uint32(node.ID),
			Word:   string(node.Word),
			Rate:   float32(node.Rate),
		})
	}
	for _, edge := range edges {
		resp.Edges = append(resp.Edges, &apiv2.RateEdge{
			EdgeId:  uint32(edge.ID),
			NodeId1: uint32(edge.Node1ID),
			NodeId2: uint32(edge.Node2ID),
			Rate:    float32(edge.Rate),
		})
	}
	p.handler.Resp(resp)
}

// [numFetch]のリクエスト必要数をセットして必要処理数を送信する
func (p *V2RateProvider) handleRespTotalProcess(numFetch int) {
	// +2 行列計算用
	p.totalProcess = numFetch + 1
	p.handleResp([]*matrix.Node{}, []*matrix.Edge{})
}

// 処理済みをカウントアップして送信する
func (p *V2RateProvider) handleRespProcess() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.doneProcess < p.totalProcess {
		p.doneProcess += 1
	}
	p.handleResp([]*matrix.Node{}, []*matrix.Edge{})
}
