package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/jacute/prettylogger/colors"
)

type JsonHandler struct {
	H    slog.Handler
	B    *bytes.Buffer
	M    *sync.Mutex
	Opts *slog.HandlerOptions
	W    io.Writer
}

func (h *JsonHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Opts.Level.Level()
}

func (h *JsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JsonHandler{H: h.H.WithAttrs(attrs), B: h.B, M: h.M, Opts: h.Opts, W: h.W}
}

func (h *JsonHandler) WithGroup(name string) slog.Handler {
	return &JsonHandler{H: h.H.WithGroup(name), B: h.B, M: h.M, Opts: h.Opts, W: h.W}
}

func (h *JsonHandler) computeAttr(ctx context.Context, r slog.Record) (map[string]any, error) {
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

func (h *JsonHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = colors.Blue(level)
	case slog.LevelInfo:
		level = colors.Cyan(level)
	case slog.LevelWarn:
		level = colors.Yellow(level)
	case slog.LevelError:
		level = colors.Red(level)
	}

	attrs, err := h.computeAttr(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", " ")
	if err != nil {
		return fmt.Errorf("Error when unmarshaling attrs: %v", err)
	}
	fmt.Fprintln(h.W, string(bytes))
	return nil
}
