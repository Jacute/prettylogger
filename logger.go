package prettylogger

import (
	"bytes"
	"io"
	"log/slog"
	"sync"
)

func NewColoredHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		h:    slog.Default().Handler(),
		b:    &bytes.Buffer{},
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
