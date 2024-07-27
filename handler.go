package prettylogger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jacute/coloredlogger/colors"
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

	fmt.Printf(
		"%s %s: %s\n",
		colors.Colorize(colors.LightGray, r.Time.Format(time.RFC3339)),
		level,
		colors.Colorize(colors.LightGreen, r.Message),
	)
	return nil
}
