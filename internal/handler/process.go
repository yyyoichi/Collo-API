package handler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"yyyoichi/Collo-API/internal/analyzer"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/pkg/stream"
)

var ErrRequest = errors.New("invalid request")

type ProcessHandler struct {
	Err  func(error)
	Resp func(float32)
}

func NewCoMatrixes(ctx context.Context, processHandler ProcessHandler, config Config) CoMatrixes {
	// エラー発生時Errorを送信する
	var errorHook stream.ErrorHook = func(err error) {
		processHandler.Err(err)
	}
	client := ndl.NewClient(config.ndlConfig)
	// 発言記録
	numRecord, recordCh := client.GenerateNDLResultWithErrorHook(ctx, errorHook)
	// 会議ごとの形態素とその会議情報
	docCh := stream.FunIO[*ndl.NDLRecode, *matrix.Document](
		ctx,
		recordCh,
		func(r *ndl.NDLRecode) *matrix.Document {
			// 形態素解析
			ar := analyzer.Analysis(r.Speeches)
			if ar.Error() != nil {
				errorHook(ar.Error())
			}
			doc := &matrix.Document{}
			doc.Key = r.IssueID
			doc.Name = fmt.Sprintf("%s %s %s", r.NameOfHouse, r.NameOfMeeting, r.Issue)
			doc.At = r.Date
			doc.Description = fmt.Sprintf("https://kokkai.ndl.go.jp/#/detail?minId=%s&current=1", r.IssueID)
			doc.Words = ar.Get(config.analyzerConfig)
			return doc
		},
	)

	var prs process
	prs.setNumDoc(numRecord)
	b := matrix.NewBuilder()
	for doc := range docCh {
		if doc != nil && len(doc.Words) > 0 {
			b.AppendDocument(doc)
		}
		prs.doneDoc()
		prs.sendProcess(processHandler)
	}
	if prs.doneDocs != prs.numDocs {
		prs.completeDocs()
		prs.sendProcess(processHandler)
	}

	numCoMatrix, totalMatrix, groupMatrixCh := matrix.NewMultiCoMatrixFromBuilder(ctx, b, config.matrixConfig)
	totalMatrix.As("total")

	if config.matrixConfig.AtGroupID == "" {
		resp := []*matrix.CoMatrix{totalMatrix}
		for m := range groupMatrixCh {
			resp = append(resp, m)
		}
		// 返却総数
		prs.setNumCoMatrixes(len(resp))

		coMatrixCh := stream.Generator[*matrix.CoMatrix](ctx, resp...)
		doneCh := stream.FunIO[*matrix.CoMatrix, interface{}](ctx, coMatrixCh, func(m *matrix.CoMatrix) interface{} {
			for pg := range m.ConsumeProgress() {
				switch pg {
				case matrix.ErrDone:
					errorHook(m.Error())
				}
			}
			return struct{}{}
		})
		for range doneCh {
			prs.doneCoMatrix()
			prs.sendProcess(processHandler)
		}
		return resp
	}

	prs.setNumCoMatrixes(1)
	var cm *matrix.CoMatrix
	if config.matrixConfig.AtGroupID == "total" {
		cm = totalMatrix
	} else {
		if numCoMatrix < 1 {
			errorHook(ErrRequest)
			return CoMatrixes{}
		}
		cm = <-groupMatrixCh
	}
	for pg := range cm.ConsumeProgress() {
		switch pg {
		case matrix.ErrDone:
			errorHook(cm.Error())
		}
	}
	prs.doneCoMatrix()
	prs.sendProcess(processHandler)
	return CoMatrixes{cm}
}

type process struct {
	numDocs  float64
	doneDocs float64

	numCoMatrixes  float64
	doneCoMatrixes float64

	mu sync.Mutex
}

// 進捗を百分率で返す。
func (p *process) sendProcess(h ProcessHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var docp float64
	if p.numDocs > 0 {
		docp = p.doneDocs / p.numDocs / 2.0
	}
	var comp float64
	if p.numCoMatrixes > 0 {
		comp = p.doneCoMatrixes / p.numCoMatrixes / 2.0
	}

	h.Resp(float32(docp + comp))
}

func (p *process) setNumDoc(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.numDocs = float64(n)
}
func (p *process) doneDoc() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.doneDocs += 1.0
}
func (p *process) completeDocs() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.doneDocs = p.numDocs
}
func (p *process) setNumCoMatrixes(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.numCoMatrixes = float64(n)
}
func (p *process) doneCoMatrix() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.doneCoMatrixes += 1.0
}