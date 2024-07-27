package main

import (
	"fmt"
	"log/slog"

	prettylogger "github.com/jacute/prettylogger"
)

func main() {
	logger := slog.New(prettylogger.NewColoredHandler(nil))
	logger.Debug("Debug test", prettylogger.Err(fmt.Errorf("Aboba")))
	logger.Info("Info test")
	logger.Warn("Warning test")
	logger.Error("Error test")
}
