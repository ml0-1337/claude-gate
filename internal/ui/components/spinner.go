package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/claude-gate/internal/ui/styles"
)

// SpinnerModel represents a spinner with a message
type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	status   string
	quitting bool
}

// NewSpinner creates a new spinner model
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.InfoStyle
	return SpinnerModel{
		spinner: s,
		message: message,
		status:  "loading",
	}
}

// Init initializes the spinner
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles spinner updates
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case StatusMsg:
		m.status = msg.Status
		m.message = msg.Message
		if msg.Status != "loading" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	default:
		return m, nil
	}
}

// View renders the spinner
func (m SpinnerModel) View() string {
	if m.quitting && m.status != "loading" {
		return styles.RenderStatus(m.status, m.message) + "\n"
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// StatusMsg is used to update the spinner status
type StatusMsg struct {
	Status  string
	Message string
}

// RunSpinner runs a spinner while executing a function
func RunSpinner(message string, fn func() error) error {
	// Create the spinner model
	model := NewSpinner(message)
	
	// Create a channel to receive the result
	done := make(chan error, 1)
	
	// Create the Bubble Tea program
	p := tea.NewProgram(model)
	
	// Run the function in a goroutine
	go func() {
		err := fn()
		done <- err
		
		// Send status update to the spinner
		if err != nil {
			p.Send(StatusMsg{
				Status:  "error",
				Message: err.Error(),
			})
		} else {
			p.Send(StatusMsg{
				Status:  "success",
				Message: "Done!",
			})
		}
	}()
	
	// Run the spinner
	if _, err := p.Run(); err != nil {
		return err
	}
	
	// Wait for the function to complete
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
		return fmt.Errorf("operation timed out")
	}
}

// SimpleSpinner shows a spinner for a duration or until interrupted
func SimpleSpinner(message string, duration time.Duration) {
	model := NewSpinner(message)
	p := tea.NewProgram(model)
	
	// Auto-quit after duration
	go func() {
		time.Sleep(duration)
		p.Send(StatusMsg{
			Status:  "success",
			Message: "Done!",
		})
	}()
	
	p.Run()
}