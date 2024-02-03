package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"yyyoichi/Collo-API/internal/analyzer"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"
	"yyyoichi/Collo-API/pkg/stream"
)

var (
	ErrRequest = errors.New("invalid request")
)

type ProcessHandler struct {
	Err  func(error)
	Resp func(float32)
}

func NewCoMatrixes(ctx context.Context, processHandler ProcessHandler, config Config) matrix.CoMatrixes {
	slog.InfoContext(ctx, "new matrix from original src")
	// エラー発生時Errorを送信する
	var errorHook stream.ErrorHook = func(err error) {
		processHandler.Err(err)
	}
	client := ndl.NewClient(config.ndlConfig)
	// 発言記録
	numRecord, recordCh := client.GenerateNDLResultWithErrorHook(ctx, errorHook)
	// 会議ごとの形態素とその会議情報
	docCh := stream.FunIO[ndl.NDLRecode, matrix.AppendedDocument](
		ctx,
		recordCh,
		func(r ndl.NDLRecode) matrix.AppendedDocument {
			// 形態素解析
			ar := analyzer.Analysis(r.Speeches)
			if ar.Error() != nil {
				errorHook(ar.Error())
			}
			var meta matrix.DocumentMeta
			meta.Key = r.IssueID
			meta.Name = fmt.Sprintf("%s %s %s", r.NameOfHouse, r.NameOfMeeting, r.Issue)
			meta.At = r.Date
			meta.Description = fmt.Sprintf("https://kokkai.ndl.go.jp/#/detail?minId=%s&current=1", r.IssueID)
			return matrix.AppendedDocument{
				Words:        ar.Get(config.analyzerConfig),
				DocumentMeta: meta,
			}
		},
	)

	var prs process
	prs.setNumDoc(numRecord)
	b := matrix.NewBuilder()
	for doc := range docCh {
		if len(doc.Words) > 0 {
			b.Append(doc)
		}
		prs.doneDoc()
		prs.sendProcess(processHandler)
	}
	if prs.doneDocs != prs.numDocs {
		prs.completeDocs()
		prs.sendProcess(processHandler)
	}

	cos := matrix.NewCoMatrixesFromBuilder(ctx, b, config.matrixConfig)
	// 返却総数
	prs.setNumCoMatrixes(len(cos.Data))
	if len(cos.Data) == 0 {
		return cos
	}
	mxCh := stream.Generator[matrix.CoMatrix](ctx, cos.Data...)
	doneCh := stream.FunIO[matrix.CoMatrix, struct{}](ctx, mxCh, func(m matrix.CoMatrix) struct{} {
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
	return cos
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
