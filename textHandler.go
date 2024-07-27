package prettylogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jacute/prettylogger/colors"
)

type Handler struct {
	h    slog.Handler
	b    *bytes.Buffer
	m    *sync.Mutex
	opts *slog.HandlerOptions
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m, opts: h.opts}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m, opts: h.opts}
}

func (h *Handler) computeAttr(ctx context.Context, r slog.Record) (map[string]any, error) {
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

func supressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
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

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
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
	fmt.Printf(
		"%s %s: %s %s\n",
		colors.Colorize(colors.LightGray, r.Time.Format(time.RFC3339)),
		level,
		colors.Colorize(colors.White, r.Message),
		colors.Colorize(colors.DarkGray, string(bytes)),
	)
	return nil
}
