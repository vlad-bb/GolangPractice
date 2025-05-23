package llogger

import (
	"context"
	"log/slog"
	"os"
)

type LevelBasedHandler struct {
	debugHandler slog.Handler
	teeHandler   slog.Handler
}

func NewLevelBasedHandler(debugHandler slog.Handler, teeHandler slog.Handler) *LevelBasedHandler {
	return &LevelBasedHandler{
		debugHandler: debugHandler,
		teeHandler:   teeHandler,
	}
}

func (h *LevelBasedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.debugHandler.Enabled(ctx, level) || h.teeHandler.Enabled(ctx, level)
}

func (h *LevelBasedHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < slog.LevelInfo {
		return h.debugHandler.Handle(ctx, r)
	}
	if err := h.debugHandler.Handle(ctx, r); err != nil {
		return err
	}
	return h.teeHandler.Handle(ctx, r)
}

func (h *LevelBasedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LevelBasedHandler{
		debugHandler: h.debugHandler.WithAttrs(attrs),
		teeHandler:   h.teeHandler.WithAttrs(attrs),
	}
}

func (h *LevelBasedHandler) WithGroup(name string) slog.Handler {
	return &LevelBasedHandler{
		debugHandler: h.debugHandler.WithGroup(name),
		teeHandler:   h.teeHandler.WithGroup(name),
	}
}

func SetupLogger() *slog.Logger {
	// Console handler (e.g., os.Stdout)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// File handler
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Wrap them with LevelBasedHandler
	customHandler := NewLevelBasedHandler(consoleHandler, fileHandler)

	logger := slog.New(customHandler)
	slog.SetDefault(logger)
	return logger
}
