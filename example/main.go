package main

import (
	"log/slog"
	"os"

	logger "github.com/jacute/prettylogger"
)

func main() {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(logger.NewColoredHandler(os.Stdout, &opts))
	logger.Debug("Debug test")
	logger.Info("Info test")
	logger.Warn("Warning test")
	logger.Error("Error test")
}
