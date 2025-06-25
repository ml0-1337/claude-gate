package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const loggerKey contextKey = "logger"

// LogLevel represents the logging level
type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

// New creates a new structured logger with the specified level
func New(level LogLevel) *slog.Logger {
	var slogLevel slog.Level
	
	switch level {
	case DEBUG:
		slogLevel = slog.LevelDebug
	case INFO:
		slogLevel = slog.LevelInfo
	case WARNING:
		slogLevel = slog.LevelWarn
	case ERROR:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	
	opts := &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize time format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z07:00")),
				}
			}
			return a
		},
	}
	
	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler)
}

// ParseLevel converts a string to LogLevel
func ParseLevel(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING", "WARN":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// WithContext adds the logger to the context
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from the context
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	// Return a default logger if none in context
	return slog.Default()
}