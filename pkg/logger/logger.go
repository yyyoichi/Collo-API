package logger

import (
	"context"
	"errors"
	"io"
	"log/slog"
)

type (
	LogHandler struct {
		slog.Handler
	}
	LogContext struct {
		RequestID string
	}
)

func NewLogHandler(w io.Writer, opts *slog.HandlerOptions) *LogHandler {
	var h LogHandler
	h.Handler = slog.NewJSONHandler(w, opts)
	return &h
}
func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	var l LogContext
	l.get(ctx)
	r.AddAttrs(
		slog.String("requestID", l.RequestID),
	)
	return h.Handler.Handle(ctx, r)
}

type (
	lk string
)

const logkey lk = "logcontextkey"

func (l *LogContext) Set(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, logkey, l)
}
func (l *LogContext) get(ctx context.Context) error {
	v := ctx.Value(logkey)
	eplog, ok := v.(*LogContext)
	if !ok {
		return errors.New("not found")
	}
	*l = *eplog
	return nil
}
