package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger/colors"
)

type TextHandler struct {
	*BaseHandler
}

func (h *TextHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = colors.Cyan(level)
	case slog.LevelInfo:
		level = colors.Magenta(level)
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
	fmt.Fprintf(
		h.W,
		"%s %s: %s %s\n",
		colors.White(r.Time.Format(time.RFC3339)),
		level,
		colors.Blue(r.Message),
		colors.HiWhite(string(bytes)),
	)
	return nil
}
