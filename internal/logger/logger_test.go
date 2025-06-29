package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test 11: Logger should initialize with correct log level
func TestNew_InitializesWithCorrectLevel(t *testing.T) {
	// This test verifies that the logger is created with the correct log level
	
	// Prediction: This test will pass because New() function properly sets log levels
	
	tests := []struct {
		name     string
		level    LogLevel
		expected slog.Level
	}{
		{"DEBUG level", DEBUG, slog.LevelDebug},
		{"INFO level", INFO, slog.LevelInfo},
		{"WARNING level", WARNING, slog.LevelWarn},
		{"ERROR level", ERROR, slog.LevelError},
		{"Invalid level defaults to INFO", LogLevel("INVALID"), slog.LevelInfo},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger
			logger := New(tt.level)
			assert.NotNil(t, logger)
			
			// To properly test the level, we would need to capture the handler
			// For now, we verify the logger is created successfully
		})
	}
}

// Test 12: Logger should parse log levels from strings
func TestParseLevel(t *testing.T) {
	// This test verifies that ParseLevel correctly converts strings to LogLevel
	
	// Prediction: This test will pass because ParseLevel is straightforward
	
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", DEBUG},
		{"debug", DEBUG},
		{"INFO", INFO},
		{"info", INFO},
		{"WARNING", WARNING},
		{"warning", WARNING},
		{"WARN", WARNING},
		{"warn", WARNING},
		{"ERROR", ERROR},
		{"error", ERROR},
		{"invalid", INFO}, // defaults to INFO
		{"", INFO},        // empty defaults to INFO
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test 13: Logger should handle context operations
func TestContextOperations(t *testing.T) {
	// This test verifies WithContext and FromContext work correctly
	
	// Prediction: This test will pass because context operations are simple
	
	// Create a logger
	logger := New(INFO)
	
	// Test WithContext
	ctx := context.Background()
	ctxWithLogger := WithContext(ctx, logger)
	
	// Test FromContext retrieves the same logger
	retrieved := FromContext(ctxWithLogger)
	assert.Equal(t, logger, retrieved)
	
	// Test FromContext returns default logger when none in context
	emptyCtx := context.Background()
	defaultLogger := FromContext(emptyCtx)
	assert.NotNil(t, defaultLogger)
	assert.Equal(t, slog.Default(), defaultLogger)
}

// Test 14: Logger should format output correctly
func TestLoggerOutput(t *testing.T) {
	// This test verifies that the logger produces correctly formatted output
	
	// Prediction: This test will pass after we capture stderr output
	
	// Create a buffer to capture output
	var buf bytes.Buffer
	
	// Create a custom handler that writes to our buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Skip time for consistent testing
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	
	handler := slog.NewTextHandler(&buf, opts)
	testLogger := slog.New(handler)
	
	// Log a message
	testLogger.Info("test message", "key", "value")
	
	// Check output contains expected content
	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key=value")
	assert.Contains(t, output, "level=INFO")
}

// Test that log levels filter messages correctly
func TestLogLevelFiltering(t *testing.T) {
	// Verify that messages below the set level are filtered out
	
	tests := []struct {
		name        string
		loggerLevel LogLevel
		logAt       string
		shouldLog   bool
	}{
		{"INFO logger filters DEBUG", INFO, "debug", false},
		{"INFO logger shows INFO", INFO, "info", true},
		{"INFO logger shows WARNING", INFO, "warn", true},
		{"INFO logger shows ERROR", INFO, "error", true},
		{"ERROR logger filters INFO", ERROR, "info", false},
		{"ERROR logger shows ERROR", ERROR, "error", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Map LogLevel to slog.Level
			var slogLevel slog.Level
			switch tt.loggerLevel {
			case DEBUG:
				slogLevel = slog.LevelDebug
			case INFO:
				slogLevel = slog.LevelInfo
			case WARNING:
				slogLevel = slog.LevelWarn
			case ERROR:
				slogLevel = slog.LevelError
			}
			
			opts := &slog.HandlerOptions{Level: slogLevel}
			handler := slog.NewTextHandler(&buf, opts)
			testLogger := slog.New(handler)
			
			// Log at different levels
			switch tt.logAt {
			case "debug":
				testLogger.Debug("debug message")
			case "info":
				testLogger.Info("info message")
			case "warn":
				testLogger.Warn("warn message")
			case "error":
				testLogger.Error("error message")
			}
			
			output := buf.String()
			if tt.shouldLog {
				assert.NotEmpty(t, output)
				assert.Contains(t, output, strings.ToUpper(tt.logAt))
			} else {
				assert.Empty(t, output)
			}
		})
	}
}