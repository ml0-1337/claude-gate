package utils

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

// IsInteractive returns true if we're running in an interactive terminal
func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// SupportsColor returns true if the terminal supports color output
func SupportsColor() bool {
	if !IsInteractive() {
		return false
	}
	return termenv.ColorProfile() != termenv.Ascii
}

// SupportsEmoji returns true if the terminal likely supports emoji
func SupportsEmoji() bool {
	if !IsInteractive() {
		return false
	}
	// Check if we're in a known good terminal
	term := os.Getenv("TERM_PROGRAM")
	switch term {
	case "iTerm.app", "Apple_Terminal", "vscode", "Hyper":
		return true
	}
	// Also check for modern terminal emulators on Linux
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return true
	}
	// Windows Terminal supports emoji
	if os.Getenv("WT_SESSION") != "" {
		return true
	}
	return false
}

// GetTerminalWidth returns the terminal width or a default value
func GetTerminalWidth() int {
	// For now, return a sensible default
	// We can enhance this later with proper terminal size detection
	return 80
}

// ClearLine clears the current line in the terminal
func ClearLine() {
	if IsInteractive() {
		fmt.Print("\r\033[K")
	}
}

// MoveCursorUp moves the cursor up n lines
func MoveCursorUp(n int) {
	if IsInteractive() {
		fmt.Printf("\033[%dA", n)
	}
}