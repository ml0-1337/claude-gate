package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ml0-1337/claude-gate/internal/ui/styles"
)

// TimerModel represents a countdown timer
type TimerModel struct {
	duration  time.Duration
	remaining time.Duration
	message   string
	expired   bool
	quitting  bool
}

// NewTimer creates a new countdown timer
func NewTimer(duration time.Duration, message string) TimerModel {
	return TimerModel{
		duration:  duration,
		remaining: duration,
		message:   message,
		expired:   false,
		quitting:  false,
	}
}

// Init initializes the timer
func (m TimerModel) Init() tea.Cmd {
	return tickCmd()
}

// tickMsg is sent every second
type tickMsg time.Time

// tickCmd returns a command that sends a tick every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ExpiredMsg is sent when the timer expires
type ExpiredMsg struct{}

// Update handles timer updates
func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		m.remaining -= time.Second
		if m.remaining <= 0 {
			m.expired = true
			m.quitting = true
			return m, tea.Batch(
				tea.Quit,
				func() tea.Msg { return ExpiredMsg{} },
			)
		}
		return m, tickCmd()

	case StatusMsg:
		if msg.Status == "success" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the timer
func (m TimerModel) View() string {
	if m.quitting && !m.expired {
		return ""
	}

	if m.expired {
		return styles.ErrorStyle.Render("⏰ Timer expired!")
	}

	minutes := int(m.remaining.Minutes())
	seconds := int(m.remaining.Seconds()) % 60

	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
	
	// Create a progress bar based on remaining time
	progress := float64(m.remaining) / float64(m.duration)
	width := 30
	filled := int(float64(width) * (1 - progress))
	
	progressBar := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().
		Foreground(styles.Muted).
		Render(strings.Repeat("░", width-filled))

	return fmt.Sprintf(
		"%s\n%s %s\n",
		m.message,
		progressBar,
		styles.InfoStyle.Render(timeStr),
	)
}

// CountdownTimer runs a countdown timer with a message
func CountdownTimer(duration time.Duration, message string) tea.Model {
	return NewTimer(duration, message)
}

// TimerWithCallback runs a timer and executes a callback when done
func TimerWithCallback(duration time.Duration, message string, onExpired func()) {
	model := NewTimer(duration, message)
	p := tea.NewProgram(model)
	
	go func() {
		finalModel, err := p.Run()
		if err == nil {
			// Check if timer expired
			if m, ok := finalModel.(TimerModel); ok && m.expired {
				onExpired()
			}
		}
	}()
}