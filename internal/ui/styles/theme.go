package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Base colors
	Primary   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7C7AFF"}
	Secondary = lipgloss.AdaptiveColor{Light: "#6C6CA0", Dark: "#9E9ED1"}
	Success   = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02D69F"}
	Warning   = lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#FBBF24"}
	Error     = lipgloss.AdaptiveColor{Light: "#E11D48", Dark: "#F43F5E"}
	Info      = lipgloss.AdaptiveColor{Light: "#0EA5E9", Dark: "#38BDF8"}
	Muted     = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			MarginBottom(1)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(Muted)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Progress bar styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Background(lipgloss.Color("237"))

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("237"))

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(Primary).
			Padding(0, 2).
			MarginRight(1)

	ButtonInactiveStyle = lipgloss.NewStyle().
				Foreground(Muted).
				Background(lipgloss.Color("237")).
				Padding(0, 2).
				MarginRight(1)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	InputInactiveStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(Muted).
				Padding(0, 1)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	// Code styles
	CodeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Padding(0, 1)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedListItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(Primary).
				Bold(true)
)

// StatusIcon returns an icon for a given status
func StatusIcon(status string) string {
	switch status {
	case "success":
		return SuccessStyle.Render("✓")
	case "warning":
		return WarningStyle.Render("⚠")
	case "error":
		return ErrorStyle.Render("✗")
	case "info":
		return InfoStyle.Render("ℹ")
	case "loading":
		return InfoStyle.Render("◐")
	default:
		return ""
	}
}

// RenderStatus renders a status message with an icon
func RenderStatus(status, message string) string {
	icon := StatusIcon(status)
	switch status {
	case "success":
		return icon + " " + SuccessStyle.Render(message)
	case "warning":
		return icon + " " + WarningStyle.Render(message)
	case "error":
		return icon + " " + ErrorStyle.Render(message)
	case "info":
		return icon + " " + InfoStyle.Render(message)
	default:
		return message
	}
}