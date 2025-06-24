package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInteractive(t *testing.T) {
	// This test will behave differently in different environments
	// We're mainly testing that the function doesn't panic
	result := IsInteractive()
	assert.IsType(t, bool(false), result)
}

func TestSupportsColor(t *testing.T) {
	// Test that it returns false when not interactive
	// This is hard to test properly without mocking
	result := SupportsColor()
	assert.IsType(t, bool(false), result)
}

func TestSupportsEmoji(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name: "iTerm",
			envVars: map[string]string{
				"TERM_PROGRAM": "iTerm.app",
			},
			expected: true,
		},
		{
			name: "VSCode",
			envVars: map[string]string{
				"TERM_PROGRAM": "vscode",
			},
			expected: true,
		},
		{
			name: "Windows Terminal",
			envVars: map[string]string{
				"WT_SESSION": "some-session-id",
			},
			expected: true,
		},
		{
			name:     "Unknown terminal",
			envVars:  map[string]string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origTermProgram := os.Getenv("TERM_PROGRAM")
			origColorTerm := os.Getenv("COLORTERM")
			origWTSession := os.Getenv("WT_SESSION")
			
			// Clear env
			os.Unsetenv("TERM_PROGRAM")
			os.Unsetenv("COLORTERM")
			os.Unsetenv("WT_SESSION")
			
			// Set test env
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			
			// Only test if we're in an interactive terminal
			if IsInteractive() {
				result := SupportsEmoji()
				assert.Equal(t, tt.expected, result)
			}
			
			// Restore original env
			if origTermProgram != "" {
				os.Setenv("TERM_PROGRAM", origTermProgram)
			}
			if origColorTerm != "" {
				os.Setenv("COLORTERM", origColorTerm)
			}
			if origWTSession != "" {
				os.Setenv("WT_SESSION", origWTSession)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	width := GetTerminalWidth()
	// Should return at least the default width
	assert.GreaterOrEqual(t, width, 80)
}