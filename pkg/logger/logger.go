package logger

import (
	"log/slog"
	"os"
	"strings"
)

var l *slog.Logger

func init() {
	l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func SetLevel(mode string) {
	level := slog.LevelInfo
	switch strings.ToLower(mode) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func Info(msg string, args ...any)  { l.Info(msg, args...) }
func Warn(msg string, args ...any)  { l.Warn(msg, args...) }
func Error(msg string, args ...any) { l.Error(msg, args...) }
func Debug(msg string, args ...any) { l.Debug(msg, args...) }
