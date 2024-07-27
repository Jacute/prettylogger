package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type BaseHandler struct {
	H    slog.Handler
	B    *bytes.Buffer
	M    *sync.Mutex
	Opts *slog.HandlerOptions
	W    io.Writer
}

func (h *BaseHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Opts.Level.Level()
}

func (h *BaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BaseHandler{H: h.H.WithAttrs(attrs), B: h.B, M: h.M, Opts: h.Opts, W: h.W}
}

func (h *BaseHandler) WithGroup(name string) slog.Handler {
	return &BaseHandler{H: h.H.WithGroup(name), B: h.B, M: h.M, Opts: h.Opts, W: h.W}
}

func (h *BaseHandler) computeAttr(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.M.Lock()
	defer func() {
		h.B.Reset()
		h.M.Unlock()
	}()
	if err := h.H.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("Error when calling inner handler's Handle: %W", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.B.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshalling inner handler's Handle attrs: %W", err)
	}
	return attrs, nil
}

// Define the Handle method but leave it unimplemented
func (H *BaseHandler) Handle(ctx context.Context, r slog.Record) error {
	panic("Handle method not implemented")
}

// Skip default attrs
func SupressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
