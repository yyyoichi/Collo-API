package provider

import (
	"context"
	"fmt"
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
	_, m, _ := matrix.NewMultiCoMatrixFromBuilder(ctx, b, matrix.Config{})
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
	meetingCh := stream.DemultiWithErrorHook[*ndl.MeetingResult, *ndl.MeetingRecode](
		ctx,
		errorHook,
		meetingResultCh,
		ndl.NewMeetingRecodes)
	// 会議ごとの形態素とその会議情報
	docCh := stream.FunIO[*ndl.MeetingRecode, *matrix.Document](
		ctx,
		meetingCh,
		func(mr *ndl.MeetingRecode) *matrix.Document {
			// 形態素解析
			ar := analyzer.Analysis(mr.Speeches)
			if ar.Error() != nil {
				errorHook(ar.Error())
			}
			doc := &matrix.Document{}
			doc.Key = mr.IssueID
			doc.Name = fmt.Sprintf("%s %s %s", mr.NameOfHouse, mr.NameOfMeeting, mr.Issue)
			doc.At = mr.Date
			doc.Description = fmt.Sprintf("https://kokkai.ndl.go.jp/#/detail?minId=%s&current=1", mr.IssueID)
			doc.Words = ar.Get(analyzerConfig)
			return doc
		},
	)

	builder := matrix.NewBuilder()
	for doc := range docCh {
		builder.AppendDocument(doc)
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
