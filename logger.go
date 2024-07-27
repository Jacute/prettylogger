package prettylogger

import (
	"bytes"
	"log/slog"
	"sync"
)

func NewColoredHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	b := &bytes.Buffer{}
	return &Handler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: supressDefaults(opts.ReplaceAttr),
		}),
		b:    b,
		m:    &sync.Mutex{},
		opts: opts,
	}
}

func NewJsonHandler(opts *slog.HandlerOptions) *JsonHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	b := &bytes.Buffer{}
	return &JsonHandler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: opts.ReplaceAttr,
		}),
		b:    b,
		m:    &sync.Mutex{},
		opts: opts,
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
