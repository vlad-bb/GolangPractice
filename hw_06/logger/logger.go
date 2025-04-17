package logger

import (
	"log/slog"
	"os"
)

var appLogger *slog.Logger

func init() {
	consoleHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug})

	appLogger = slog.New(consoleHandler)
}

func GetLogger() *slog.Logger {
	return appLogger
}
