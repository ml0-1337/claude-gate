package ui

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureOutput captures stdout and stderr during function execution
func captureOutput(t *testing.T, fn func()) (string, string) {
	t.Helper()
	
	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = wOut
	
	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	assert.NoError(t, err)
	os.Stderr = wErr
	
	// Run the function
	fn()
	
	// Close writers
	wOut.Close()
	wErr.Close()
	
	// Read output
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)
	
	// Restore
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	
	return bufOut.String(), bufErr.String()
}

// Test 14: Output methods should format messages correctly
func TestOutput_FormatMessages(t *testing.T) {
	// Prediction: This test will pass - Output methods are straightforward
	
	tests := []struct {
		name       string
		method     string
		format     string
		args       []interface{}
		wantStdout bool
		wantStderr bool
		contains   string
	}{
		{
			name:       "success message",
			method:     "success",
			format:     "Operation %s completed",
			args:       []interface{}{"test"},
			wantStdout: true,
			wantStderr: false,
			contains:   "Operation test completed",
		},
		{
			name:       "error message",
			method:     "error",
			format:     "Error: %s failed",
			args:       []interface{}{"connection"},
			wantStdout: false,
			wantStderr: true,
			contains:   "Error: connection failed",
		},
		{
			name:       "warning message",
			method:     "warning",
			format:     "Warning: %d attempts remaining",
			args:       []interface{}{3},
			wantStdout: true,
			wantStderr: false,
			contains:   "Warning: 3 attempts remaining",
		},
		{
			name:       "info message",
			method:     "info",
			format:     "Processing %s...",
			args:       []interface{}{"data"},
			wantStdout: true,
			wantStderr: false,
			contains:   "Processing data...",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := NewOutput()
			// Disable color for consistent testing
			out.colorEnabled = false
			
			stdout, stderr := captureOutput(t, func() {
				switch tt.method {
				case "success":
					out.Success(tt.format, tt.args...)
				case "error":
					out.Error(tt.format, tt.args...)
				case "warning":
					out.Warning(tt.format, tt.args...)
				case "info":
					out.Info(tt.format, tt.args...)
				}
			})
			
			if tt.wantStdout {
				assert.Contains(t, stdout, tt.contains)
				assert.Empty(t, stderr)
			} else {
				assert.Contains(t, stderr, tt.contains)
				assert.Empty(t, stdout)
			}
		})
	}
}

// Test 15: Table rendering should handle various data
func TestOutput_TableRendering(t *testing.T) {
	// Prediction: This test will pass - Table rendering logic is clear
	
	tests := []struct {
		name     string
		headers  []string
		rows     [][]string
		contains []string
	}{
		{
			name:    "simple table",
			headers: []string{"Name", "Value"},
			rows: [][]string{
				{"Host", "127.0.0.1"},
				{"Port", "5789"},
			},
			contains: []string{
				"Name", "Value",
				"Host", "127.0.0.1",
				"Port", "5789",
				"----", // separator
			},
		},
		{
			name:    "table with varying widths",
			headers: []string{"Short", "Very Long Header"},
			rows: [][]string{
				{"A", "B"},
				{"Long content here", "X"},
			},
			contains: []string{
				"Short", "Very Long Header",
				"Long content here", "X",
			},
		},
		{
			name:    "empty table",
			headers: []string{"Col1", "Col2"},
			rows:    [][]string{},
			contains: []string{
				"Col1", "Col2",
				"----", // separator still shown
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := NewOutput()
			out.colorEnabled = false
			
			stdout, _ := captureOutput(t, func() {
				out.Table(tt.headers, tt.rows)
			})
			
			for _, expected := range tt.contains {
				assert.Contains(t, stdout, expected)
			}
		})
	}
}

// Test 16: Title and subtitle formatting
func TestOutput_TitleSubtitle(t *testing.T) {
	// Prediction: This test will pass
	
	out := NewOutput()
	out.colorEnabled = false
	
	t.Run("title formatting", func(t *testing.T) {
		stdout, _ := captureOutput(t, func() {
			out.Title("Test Title")
		})
		
		assert.Contains(t, stdout, "Test Title")
		assert.Contains(t, stdout, "==========") // underline
	})
	
	t.Run("subtitle formatting", func(t *testing.T) {
		stdout, _ := captureOutput(t, func() {
			out.Subtitle("Test Subtitle")
		})
		
		assert.Contains(t, stdout, "Test Subtitle")
		assert.Contains(t, stdout, "-------------") // underline
	})
}

// Test 17: Box formatting
func TestOutput_Box(t *testing.T) {
	// Prediction: This test will pass
	
	tests := []struct {
		name     string
		content  string
		contains []string
	}{
		{
			name:    "single line box",
			content: "Hello World",
			contains: []string{
				"+-------------+",
				"| Hello World |",
			},
		},
		{
			name:    "multi line box",
			content: "Line 1\nLine 2\nLonger Line 3",
			contains: []string{
				"+---------------+",
				"| Line 1        |",
				"| Line 2        |",
				"| Longer Line 3 |",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := NewOutput()
			out.colorEnabled = false
			
			stdout, _ := captureOutput(t, func() {
				out.Box(tt.content)
			})
			
			for _, expected := range tt.contains {
				assert.Contains(t, stdout, expected)
			}
		})
	}
}

// Test 18: List formatting
func TestOutput_List(t *testing.T) {
	// Prediction: This test will pass
	
	out := NewOutput()
	out.colorEnabled = false
	
	items := []string{"Item 1", "Item 2", "Item 3"}
	
	stdout, _ := captureOutput(t, func() {
		out.List(items)
	})
	
	assert.Contains(t, stdout, "• Item 1")
	assert.Contains(t, stdout, "• Item 2")
	assert.Contains(t, stdout, "• Item 3")
}

// Test 19: Code formatting
func TestOutput_Code(t *testing.T) {
	// Prediction: This test will pass
	
	out := NewOutput()
	out.colorEnabled = false
	
	code := "go test ./..."
	
	stdout, _ := captureOutput(t, func() {
		out.Code(code)
	})
	
	assert.Contains(t, stdout, "  go test ./...")
}

// Test 20: Interactive mode
func TestOutput_InteractiveMode(t *testing.T) {
	// Prediction: This test will pass
	
	out := NewOutput()
	
	// Test default interactive detection
	// This will vary based on environment
	_ = out.IsInteractive()
	
	// Test setting interactive mode
	out.SetInteractive(true)
	assert.True(t, out.IsInteractive())
	
	out.SetInteractive(false)
	assert.False(t, out.IsInteractive())
}

// Test 21: Color mode with formatting
func TestOutput_ColorMode(t *testing.T) {
	// Prediction: This test will pass
	
	t.Run("with color enabled", func(t *testing.T) {
		out := NewOutput()
		out.colorEnabled = true
		
		stdout, _ := captureOutput(t, func() {
			out.Success("Color test")
		})
		
		// Should contain styled output (exact format depends on styles package)
		assert.Contains(t, stdout, "Color test")
	})
	
	t.Run("with color disabled", func(t *testing.T) {
		out := NewOutput()
		out.colorEnabled = false
		
		stdout, _ := captureOutput(t, func() {
			out.Success("No color test")
		})
		
		// Should contain plain output with symbol
		assert.Contains(t, stdout, "✓ No color test")
	})
}