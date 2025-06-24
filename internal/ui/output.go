package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/claude-gate/internal/ui/styles"
	"github.com/yourusername/claude-gate/internal/ui/utils"
)

// Output provides methods for formatted terminal output
type Output struct {
	interactive bool
	colorEnabled bool
	emojiEnabled bool
}

// NewOutput creates a new output handler
func NewOutput() *Output {
	return &Output{
		interactive:  utils.IsInteractive(),
		colorEnabled: utils.SupportsColor(),
		emojiEnabled: utils.SupportsEmoji(),
	}
}

// Success prints a success message
func (o *Output) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if o.colorEnabled {
		fmt.Println(styles.RenderStatus("success", message))
	} else {
		fmt.Printf("✓ %s\n", message)
	}
}

// Error prints an error message
func (o *Output) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if o.colorEnabled {
		fmt.Fprintln(os.Stderr, styles.RenderStatus("error", message))
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", message)
	}
}

// Warning prints a warning message
func (o *Output) Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if o.colorEnabled {
		fmt.Println(styles.RenderStatus("warning", message))
	} else {
		fmt.Printf("⚠ %s\n", message)
	}
}

// Info prints an info message
func (o *Output) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if o.colorEnabled {
		fmt.Println(styles.RenderStatus("info", message))
	} else {
		fmt.Printf("ℹ %s\n", message)
	}
}

// Title prints a title
func (o *Output) Title(text string) {
	if o.colorEnabled {
		fmt.Println(styles.TitleStyle.Render(text))
	} else {
		fmt.Printf("\n%s\n%s\n", text, strings.Repeat("=", len(text)))
	}
}

// Subtitle prints a subtitle
func (o *Output) Subtitle(text string) {
	if o.colorEnabled {
		fmt.Println(styles.SubtitleStyle.Render(text))
	} else {
		fmt.Printf("\n%s\n%s\n", text, strings.Repeat("-", len(text)))
	}
}

// Box prints content in a box
func (o *Output) Box(content string) {
	if o.colorEnabled {
		fmt.Println(styles.BoxStyle.Render(content))
	} else {
		lines := strings.Split(content, "\n")
		maxLen := 0
		for _, line := range lines {
			if len(line) > maxLen {
				maxLen = len(line)
			}
		}
		
		border := "+" + strings.Repeat("-", maxLen+2) + "+"
		fmt.Println(border)
		for _, line := range lines {
			fmt.Printf("| %-*s |\n", maxLen, line)
		}
		fmt.Println(border)
	}
}

// Code prints code or command examples
func (o *Output) Code(code string) {
	if o.colorEnabled {
		fmt.Println(styles.CodeStyle.Render(code))
	} else {
		fmt.Printf("  %s\n", code)
	}
}

// List prints a list of items
func (o *Output) List(items []string) {
	for _, item := range items {
		if o.colorEnabled {
			fmt.Println(styles.ListItemStyle.Render("• " + item))
		} else {
			fmt.Printf("  • %s\n", item)
		}
	}
}

// Table prints a simple table
func (o *Output) Table(headers []string, rows [][]string) {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	
	// Print headers
	headerRow := ""
	separator := ""
	for i, header := range headers {
		if i > 0 {
			headerRow += " | "
			separator += "-+-"
		}
		headerRow += fmt.Sprintf("%-*s", widths[i], header)
		separator += strings.Repeat("-", widths[i])
	}
	
	if o.colorEnabled {
		fmt.Println(styles.TitleStyle.Render(headerRow))
	} else {
		fmt.Println(headerRow)
	}
	fmt.Println(separator)
	
	// Print rows
	for _, row := range rows {
		rowStr := ""
		for i, cell := range row {
			if i > 0 {
				rowStr += " | "
			}
			if i < len(widths) {
				rowStr += fmt.Sprintf("%-*s", widths[i], cell)
			} else {
				rowStr += cell
			}
		}
		fmt.Println(rowStr)
	}
}

// IsInteractive returns true if running in interactive mode
func (o *Output) IsInteractive() bool {
	return o.interactive
}

// SetInteractive overrides the interactive mode detection
func (o *Output) SetInteractive(interactive bool) {
	o.interactive = interactive
}