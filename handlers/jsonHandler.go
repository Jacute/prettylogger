package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jacute/prettylogger/colors"
)

type JsonHandler struct {
	*BaseHandler
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
