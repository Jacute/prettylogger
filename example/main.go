package main

import (
	"fmt"
	"log/slog"
	"os"

	prettylogger "github.com/jacute/prettylogger"
)

func main() {
	consoleLogger := slog.New(prettylogger.NewColoredHandler(os.Stdout, nil))

	consoleLogger.Debug("Debug test", prettylogger.Err(fmt.Errorf("test error")))
	consoleLogger.Info("Info test")
	consoleLogger.Warn("Warning test")
	consoleLogger.Error("Error test")

	consoleLogger.With(slog.String("test", "test")).Info("Test")

	file, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileLogger := slog.New(prettylogger.NewJsonHandler(file, nil))
	fileLogger.Debug("Debug test", prettylogger.Err(fmt.Errorf("test error")))
	fileLogger.Error("Error test")

	discardLogger := slog.New(prettylogger.NewDiscardHandler())
	discardLogger.Info("Nothing")
}
