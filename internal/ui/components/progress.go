package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ml0-1337/claude-gate/internal/ui/styles"
	"github.com/ml0-1337/claude-gate/internal/ui/utils"
)

// ProgressModel represents a progress bar
type ProgressModel struct {
	progress progress.Model
	title    string
	percent  float64
	width    int
}

// NewProgress creates a new progress bar model
func NewProgress(title string) ProgressModel {
	width := utils.GetTerminalWidth() - 20 // Leave some margin
	if width < 40 {
		width = 40
	}
	
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
		progress.WithoutPercentage(),
	)
	
	return ProgressModel{
		progress: p,
		title:    title,
		percent:  0,
		width:    width,
	}
}

// Init initializes the progress bar
func (m ProgressModel) Init() tea.Cmd {
	return nil
}

// Update handles progress updates
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 20
		if m.progress.Width < 40 {
			m.progress.Width = 40
		}
		return m, nil

	case ProgressMsg:
		m.percent = msg.Percent
		if msg.Title != "" {
			m.title = msg.Title
		}
		cmd := m.progress.SetPercent(float64(m.percent))
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

// View renders the progress bar
func (m ProgressModel) View() string {
	// Calculate padding, ensuring it's never negative
	padWidth := m.width - len(m.title)
	if padWidth < 0 {
		padWidth = 0
	}
	pad := strings.Repeat(" ", padWidth)
	
	title := styles.SubtitleStyle.Render(m.title)
	percent := styles.InfoStyle.Render(fmt.Sprintf("%.0f%%", m.percent*100))
	
	return fmt.Sprintf("%s%s%s\n%s", title, pad, percent, m.progress.View())
}

// ProgressMsg is used to update progress
type ProgressMsg struct {
	Percent float64
	Title   string
}

// ProgressTracker provides a simple interface for tracking progress
type ProgressTracker struct {
	program *tea.Program
	total   int
	current int
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(title string, total int) *ProgressTracker {
	model := NewProgress(title)
	p := tea.NewProgram(model)
	
	tracker := &ProgressTracker{
		program: p,
		total:   total,
		current: 0,
	}
	
	// Start the program in a goroutine
	go p.Run()
	
	return tracker
}

// Increment increments the progress
func (t *ProgressTracker) Increment(title string) {
	t.current++
	percent := float64(t.current) / float64(t.total)
	t.program.Send(ProgressMsg{
		Percent: percent,
		Title:   title,
	})
}

// Finish completes the progress
func (t *ProgressTracker) Finish() {
	t.program.Send(ProgressMsg{Percent: 1.0})
	time.Sleep(500 * time.Millisecond) // Brief pause to show completion
	t.program.Quit()
}

// SimpleProgress shows a progress bar for a simple operation
func SimpleProgress(title string, steps []string, stepDuration time.Duration) {
	tracker := NewProgressTracker(title, len(steps))
	
	for _, step := range steps {
		tracker.Increment(step)
		time.Sleep(stepDuration)
	}
	
	tracker.Finish()
}