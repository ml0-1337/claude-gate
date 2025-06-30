package styles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test 5: theme.RenderStatus should format status messages correctly
func TestRenderStatus(t *testing.T) {
	// Prediction: This test will pass - RenderStatus is straightforward
	
	tests := []struct {
		name        string
		status      string
		message     string
		wantIcon    string
		wantMessage string
	}{
		{
			name:        "success status",
			status:      "success",
			message:     "Operation completed",
			wantIcon:    "✓",
			wantMessage: "Operation completed",
		},
		{
			name:        "warning status",
			status:      "warning",
			message:     "This might cause issues",
			wantIcon:    "⚠",
			wantMessage: "This might cause issues",
		},
		{
			name:        "error status",
			status:      "error",
			message:     "Something went wrong",
			wantIcon:    "✗",
			wantMessage: "Something went wrong",
		},
		{
			name:        "info status",
			status:      "info",
			message:     "For your information",
			wantIcon:    "ℹ",
			wantMessage: "For your information",
		},
		{
			name:        "unknown status",
			status:      "unknown",
			message:     "Unknown status message",
			wantIcon:    "",
			wantMessage: "Unknown status message",
		},
		{
			name:        "empty status",
			status:      "",
			message:     "No status provided",
			wantIcon:    "",
			wantMessage: "No status provided",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderStatus(tt.status, tt.message)
			
			// Check that the message is included
			assert.Contains(t, result, tt.wantMessage)
			
			// Check that the icon is included if expected
			if tt.wantIcon != "" {
				assert.Contains(t, result, tt.wantIcon)
			}
			
			// For unknown status, result should be just the message
			if tt.status == "unknown" || tt.status == "" {
				assert.Equal(t, tt.wantMessage, result)
			}
		})
	}
}

// Test 6: StatusIcon should return correct icons for each status
func TestStatusIcon(t *testing.T) {
	// Prediction: This test will pass - StatusIcon has clear mapping
	
	tests := []struct {
		name     string
		status   string
		wantIcon string
	}{
		{
			name:     "success icon",
			status:   "success",
			wantIcon: "✓",
		},
		{
			name:     "warning icon",
			status:   "warning",
			wantIcon: "⚠",
		},
		{
			name:     "error icon",
			status:   "error",
			wantIcon: "✗",
		},
		{
			name:     "info icon",
			status:   "info",
			wantIcon: "ℹ",
		},
		{
			name:     "loading icon",
			status:   "loading",
			wantIcon: "◐",
		},
		{
			name:     "unknown status returns empty",
			status:   "unknown",
			wantIcon: "",
		},
		{
			name:     "empty status returns empty",
			status:   "",
			wantIcon: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusIcon(tt.status)
			
			if tt.wantIcon == "" {
				assert.Empty(t, result)
			} else {
				// The icon will be styled, so we just check it contains the icon
				assert.Contains(t, result, tt.wantIcon)
			}
		})
	}
}

// Test 7: Style definitions should be properly configured
func TestStyleDefinitions(t *testing.T) {
	// Prediction: This test will pass - verifying style objects exist
	
	// Test that all style variables are not nil
	assert.NotNil(t, TitleStyle)
	assert.NotNil(t, SubtitleStyle)
	assert.NotNil(t, DescriptionStyle)
	assert.NotNil(t, SuccessStyle)
	assert.NotNil(t, WarningStyle)
	assert.NotNil(t, ErrorStyle)
	assert.NotNil(t, InfoStyle)
	assert.NotNil(t, BoxStyle)
	assert.NotNil(t, ProgressBarStyle)
	assert.NotNil(t, ProgressEmptyStyle)
	assert.NotNil(t, ButtonStyle)
	assert.NotNil(t, ButtonInactiveStyle)
	assert.NotNil(t, InputStyle)
	assert.NotNil(t, InputInactiveStyle)
	assert.NotNil(t, HelpStyle)
	assert.NotNil(t, CodeStyle)
	assert.NotNil(t, ListItemStyle)
	assert.NotNil(t, SelectedListItemStyle)
	
	// Test that styles can render text without panic
	assert.NotPanics(t, func() {
		_ = TitleStyle.Render("Test Title")
		_ = SubtitleStyle.Render("Test Subtitle")
		_ = DescriptionStyle.Render("Test Description")
		_ = SuccessStyle.Render("Success")
		_ = WarningStyle.Render("Warning")
		_ = ErrorStyle.Render("Error")
		_ = InfoStyle.Render("Info")
		_ = BoxStyle.Render("Box Content")
		_ = CodeStyle.Render("code")
		_ = ListItemStyle.Render("• Item")
		_ = SelectedListItemStyle.Render("▸ Selected")
	})
}

// Test 8: Colors should be properly defined
func TestColorDefinitions(t *testing.T) {
	// Prediction: This test will pass - verifying color variables
	
	// Test that adaptive colors are defined
	colors := map[string]interface{}{
		"Primary":   Primary,
		"Secondary": Secondary,
		"Success":   Success,
		"Warning":   Warning,
		"Error":     Error,
		"Info":      Info,
		"Muted":     Muted,
	}
	
	for name, color := range colors {
		t.Run(name, func(t *testing.T) {
			assert.NotNil(t, color, "%s color should be defined", name)
		})
	}
}

// Test 9: RenderStatus edge cases
func TestRenderStatus_EdgeCases(t *testing.T) {
	// Prediction: This test will pass - testing edge cases
	
	tests := []struct {
		name    string
		status  string
		message string
	}{
		{
			name:    "empty message",
			status:  "success",
			message: "",
		},
		{
			name:    "very long message",
			status:  "error",
			message: strings.Repeat("a", 1000),
		},
		{
			name:    "message with special characters",
			status:  "info",
			message: "Message with 特殊字符 and émojis 🎉",
		},
		{
			name:    "message with newlines",
			status:  "warning",
			message: "Line 1\nLine 2\nLine 3",
		},
		{
			name:    "message with ANSI codes",
			status:  "success",
			message: "\033[31mRed text\033[0m",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with any input
			assert.NotPanics(t, func() {
				result := RenderStatus(tt.status, tt.message)
				// Result should contain the message (even if empty)
				if tt.message != "" {
					assert.Contains(t, result, tt.message)
				}
			})
		})
	}
}

// Test 10: Style composition and chaining
func TestStyleComposition(t *testing.T) {
	// Prediction: This test will pass - testing style methods work correctly
	
	// Test that styles can be composed without issues
	t.Run("TitleStyle composition", func(t *testing.T) {
		// TitleStyle should be bold with margin
		rendered := TitleStyle.Render("Title")
		assert.NotEmpty(t, rendered)
	})
	
	t.Run("BoxStyle with padding and border", func(t *testing.T) {
		rendered := BoxStyle.Render("Content")
		assert.NotEmpty(t, rendered)
		// Box style includes padding, so content should be longer than original
		assert.Greater(t, len(rendered), len("Content"))
	})
	
	t.Run("ButtonStyle with background", func(t *testing.T) {
		rendered := ButtonStyle.Render("Click Me")
		assert.NotEmpty(t, rendered)
	})
	
	t.Run("InputStyle with border", func(t *testing.T) {
		rendered := InputStyle.Render("Input Text")
		assert.NotEmpty(t, rendered)
	})
}