package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ml0-1337/claude-gate/internal/ui/styles"
)

// AuthFlowStep represents a step in the OAuth flow
type AuthFlowStep struct {
	Name      string
	Completed bool
	Current   bool
	Failed    bool
}

// AuthFlowModel represents the OAuth authentication flow UI
type AuthFlowModel struct {
	steps      []AuthFlowStep
	authURL    string
	textInput  textinput.Model
	err        error
	width      int
	showingURL bool
	code       string
	done       bool
}

// NewAuthFlow creates a new OAuth flow UI
func NewAuthFlow() AuthFlowModel {
	ti := textinput.New()
	ti.Placeholder = "Enter authorization code"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return AuthFlowModel{
		steps: []AuthFlowStep{
			{Name: "Generate authorization URL", Completed: false, Current: true},
			{Name: "Open browser for authentication", Completed: false, Current: false},
			{Name: "Enter authorization code", Completed: false, Current: false},
			{Name: "Exchange code for tokens", Completed: false, Current: false},
			{Name: "Save tokens securely", Completed: false, Current: false},
		},
		textInput: ti,
		width:     80,
	}
}

// Init initializes the auth flow
func (m AuthFlowModel) Init() tea.Cmd {
	return textinput.Blink
}

// AuthURLMsg contains the authorization URL
type AuthURLMsg struct {
	URL string
}

// AuthCodeMsg indicates the code was submitted
type AuthCodeMsg struct {
	Code string
}

// AuthStepMsg updates the current step
type AuthStepMsg struct {
	Step      int
	Completed bool
	Failed    bool
	Error     error
}

// AuthCompleteMsg indicates authentication is complete
type AuthCompleteMsg struct{}

// Update handles auth flow updates
func (m AuthFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.showingURL && m.textInput.Value() != "" {
				m.code = m.textInput.Value()
				return m, func() tea.Msg {
					return AuthCodeMsg{Code: m.code}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case AuthURLMsg:
		m.authURL = msg.URL
		m.showingURL = true
		m.steps[0].Completed = true
		m.steps[0].Current = false
		m.steps[1].Current = true
		// Auto-advance to step 2
		return m, func() tea.Msg {
			return AuthStepMsg{Step: 2, Completed: true}
		}

	case AuthStepMsg:
		if msg.Step < len(m.steps) {
			if msg.Failed {
				m.steps[msg.Step].Failed = true
				m.steps[msg.Step].Current = false
				m.err = msg.Error
			} else if msg.Completed {
				m.steps[msg.Step].Completed = true
				m.steps[msg.Step].Current = false
				if msg.Step+1 < len(m.steps) {
					m.steps[msg.Step+1].Current = true
				}
			} else {
				// Just update current
				for i := range m.steps {
					m.steps[i].Current = i == msg.Step
				}
			}
		}

	case AuthCompleteMsg:
		m.done = true
		return m, tea.Quit

	case StatusMsg:
		if msg.Status == "error" {
			// Find current step and mark as failed
			for i, step := range m.steps {
				if step.Current {
					m.steps[i].Failed = true
					m.err = fmt.Errorf(msg.Message)
					break
				}
			}
		}
	}

	// Update text input
	if m.showingURL && m.steps[2].Current {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

// View renders the auth flow UI
func (m AuthFlowModel) View() string {
	if m.done {
		return styles.SuccessStyle.Render("✓ Authentication complete!")
	}

	var s strings.Builder

	// Title
	title := styles.TitleStyle.Render("🔐 OAuth Authentication Flow")
	s.WriteString(title + "\n\n")

	// Progress steps
	for _, step := range m.steps {
		var icon string
		var style lipgloss.Style

		if step.Failed {
			icon = "✗"
			style = styles.ErrorStyle
		} else if step.Completed {
			icon = "✓"
			style = styles.SuccessStyle
		} else if step.Current {
			icon = "◐"
			style = styles.InfoStyle
		} else {
			icon = "○"
			style = styles.DescriptionStyle
		}

		stepText := fmt.Sprintf("%s %s", icon, step.Name)
		s.WriteString(style.Render(stepText) + "\n")
	}

	s.WriteString("\n")

	// Show content based on current step
	if m.showingURL && m.authURL != "" {
		s.WriteString(styles.SubtitleStyle.Render("Please visit this URL to authorize:") + "\n\n")
		
		// Create a nice box for the URL
		urlBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Primary).
			Padding(1, 2).
			Width(min(len(m.authURL)+4, m.width-4))
		
		s.WriteString(urlBox.Render(m.authURL) + "\n\n")

		// Show input field if we're on step 2
		if m.steps[2].Current {
			s.WriteString(styles.InfoStyle.Render("After authorizing, enter the code below:") + "\n")
			s.WriteString(m.textInput.View() + "\n\n")
			s.WriteString(styles.HelpStyle.Render("Press Enter to submit, Esc to cancel"))
		}
	}

	// Show error if any
	if m.err != nil {
		s.WriteString("\n" + styles.ErrorStyle.Render("Error: "+m.err.Error()))
	}

	return s.String()
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AuthFlowUI runs the interactive OAuth flow UI
type AuthFlowUI struct {
	model   AuthFlowModel
	program *tea.Program
}

// NewAuthFlowUI creates a new auth flow UI runner
func NewAuthFlowUI() *AuthFlowUI {
	model := NewAuthFlow()
	return &AuthFlowUI{
		model:   model,
		program: tea.NewProgram(model),
	}
}

// Start starts the auth flow UI
func (a *AuthFlowUI) Start() error {
	go func() {
		if _, err := a.program.Run(); err != nil {
			fmt.Printf("Error running auth flow: %v\n", err)
		}
	}()
	return nil
}

// SetAuthURL sets the authorization URL
func (a *AuthFlowUI) SetAuthURL(url string) {
	a.program.Send(AuthURLMsg{URL: url})
}

// UpdateStep updates the current step
func (a *AuthFlowUI) UpdateStep(step int, completed bool, failed bool, err error) {
	a.program.Send(AuthStepMsg{
		Step:      step,
		Completed: completed,
		Failed:    failed,
		Error:     err,
	})
}

// GetCode waits for and returns the authorization code
func (a *AuthFlowUI) GetCode() (string, error) {
	// This would need to be implemented with channels
	// For now, return empty
	return "", nil
}

// Complete marks the flow as complete
func (a *AuthFlowUI) Complete() {
	a.program.Send(AuthCompleteMsg{})
}

// Quit quits the UI
func (a *AuthFlowUI) Quit() {
	a.program.Quit()
}