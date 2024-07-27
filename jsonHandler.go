package prettylogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jacute/prettylogger/colors"
)

type JsonHandler struct {
	h    slog.Handler
	b    *bytes.Buffer
	m    *sync.Mutex
	opts *slog.HandlerOptions
}

func (h *JsonHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *JsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JsonHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m, opts: h.opts}
}

func (h *JsonHandler) WithGroup(name string) slog.Handler {
	return &JsonHandler{h: h.h.WithGroup(name), b: h.b, m: h.m, opts: h.opts}
}

func (h *JsonHandler) computeAttr(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("Error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshalling inner handler's Handle attrs: %w", err)
	}
	return attrs, nil
}

func (h *JsonHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = colors.Colorize(colors.LightBlue, level)
	case slog.LevelInfo:
		level = colors.Colorize(colors.LightCyan, level)
	case slog.LevelWarn:
		level = colors.Colorize(colors.Yellow, level)
	case slog.LevelError:
		level = colors.Colorize(colors.Red, level)
	}

	attrs, err := h.computeAttr(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", " ")
	if err != nil {
		return fmt.Errorf("Error when unmarshaling attrs: %v", err)
	}
	fmt.Println(string(bytes))
	return nil
}
