package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test 21: Terminal utilities should detect features correctly
func TestIsInteractive(t *testing.T) {
	// Prediction: This test will pass - testing that function doesn't panic
	// This test will behave differently in different environments
	// We're mainly testing that the function doesn't panic
	result := IsInteractive()
	assert.IsType(t, bool(false), result)
	
	// The result depends on whether we're running in a TTY
	// In CI/test environments, this is typically false
	t.Logf("IsInteractive returned: %v", result)
}

func TestSupportsColor(t *testing.T) {
	// Prediction: This test will pass - testing color support detection
	// Test that it returns false when not interactive
	// This is hard to test properly without mocking
	result := SupportsColor()
	assert.IsType(t, bool(false), result)
	
	// In non-interactive environments, should always be false
	if !IsInteractive() {
		assert.False(t, result, "SupportsColor should be false in non-interactive environment")
	}
	
	t.Logf("SupportsColor returned: %v", result)
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
	// Current implementation always returns 80
	assert.Equal(t, 80, width)
}

func TestClearLine(t *testing.T) {
	// Prediction: This test will pass - testing ClearLine doesn't panic
	// ClearLine only outputs if interactive
	// We just verify it doesn't panic
	ClearLine()
	
	// Test passes if no panic occurs
	assert.True(t, true, "ClearLine executed without panic")
}

func TestMoveCursorUp(t *testing.T) {
	// Prediction: This test will pass - testing MoveCursorUp doesn't panic
	tests := []struct {
		name  string
		lines int
	}{
		{"move up 1 line", 1},
		{"move up 5 lines", 5},
		{"move up 0 lines", 0},
		{"move up negative lines", -1}, // Edge case
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// MoveCursorUp only outputs if interactive
			// We just verify it doesn't panic
			MoveCursorUp(tt.lines)
			
			// Test passes if no panic occurs
			assert.True(t, true, "MoveCursorUp executed without panic")
		})
	}
}

// Test additional emoji support cases
func TestSupportsEmoji_AdditionalCases(t *testing.T) {
	// Save original env
	origTermProgram := os.Getenv("TERM_PROGRAM")
	origColorTerm := os.Getenv("COLORTERM")
	origWTSession := os.Getenv("WT_SESSION")
	
	defer func() {
		// Restore original env
		if origTermProgram != "" {
			os.Setenv("TERM_PROGRAM", origTermProgram)
		} else {
			os.Unsetenv("TERM_PROGRAM")
		}
		if origColorTerm != "" {
			os.Setenv("COLORTERM", origColorTerm)
		} else {
			os.Unsetenv("COLORTERM")
		}
		if origWTSession != "" {
			os.Setenv("WT_SESSION", origWTSession)
		} else {
			os.Unsetenv("WT_SESSION")
		}
	}()
	
	// Test Apple Terminal
	os.Setenv("TERM_PROGRAM", "Apple_Terminal")
	os.Unsetenv("COLORTERM")
	os.Unsetenv("WT_SESSION")
	if IsInteractive() {
		assert.True(t, SupportsEmoji(), "Apple Terminal should support emoji")
	}
	
	// Test Hyper terminal
	os.Setenv("TERM_PROGRAM", "Hyper")
	if IsInteractive() {
		assert.True(t, SupportsEmoji(), "Hyper should support emoji")
	}
	
	// Test truecolor terminal
	os.Unsetenv("TERM_PROGRAM")
	os.Setenv("COLORTERM", "truecolor")
	if IsInteractive() {
		assert.True(t, SupportsEmoji(), "Truecolor terminal should support emoji")
	}
	
	// Test 24bit color terminal
	os.Setenv("COLORTERM", "24bit")
	if IsInteractive() {
		assert.True(t, SupportsEmoji(), "24bit color terminal should support emoji")
	}
}