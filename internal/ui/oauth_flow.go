package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ml0-1337/claude-gate/internal/ui/styles"
)

// OAuthStep represents a step in the OAuth flow
type OAuthStep int

const (
	StepGenerateURL OAuthStep = iota
	StepOpenBrowser
	StepEnterCode
	StepExchangeToken
	StepSaveToken
)

// OAuthFlowModel represents the OAuth flow UI
type OAuthFlowModel struct {
	currentStep OAuthStep
	authURL     string
	textInput   textinput.Model
	code        string
	err         error
	done        bool
	canceled    bool
	
	// Channels for communication
	codeChan chan string
	errChan  chan error
}

// NewOAuthFlow creates a new OAuth flow model
func NewOAuthFlow() *OAuthFlowModel {
	ti := textinput.New()
	ti.Placeholder = "Enter the authorization code"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 60

	return &OAuthFlowModel{
		currentStep: StepGenerateURL,
		textInput:   ti,
		codeChan:    make(chan string, 1),
		errChan:     make(chan error, 1),
	}
}

// Init initializes the OAuth flow
func (m *OAuthFlowModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles OAuth flow updates
func (m *OAuthFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.canceled = true
			close(m.codeChan)
			return m, tea.Quit
		case tea.KeyEnter:
			if m.currentStep == StepEnterCode && m.textInput.Value() != "" {
				m.code = m.textInput.Value()
				m.codeChan <- m.code
				m.currentStep = StepExchangeToken
				return m, nil
			}
		}

	case AuthURLMsg:
		m.authURL = msg.URL
		m.currentStep = StepOpenBrowser
		// Try to open browser
		go TryOpenBrowser(msg.URL)
		// Auto-advance after showing URL
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return AdvanceStepMsg{}
		})

	case AdvanceStepMsg:
		if m.currentStep == StepOpenBrowser {
			m.currentStep = StepEnterCode
		}

	case AuthProgressMsg:
		m.currentStep = msg.Step

	case AuthErrorMsg:
		m.err = msg.Error
		m.errChan <- msg.Error
		return m, tea.Quit

	case AuthCompleteMsg:
		m.done = true
		close(m.codeChan)
		return m, tea.Quit
	}

	// Update text input when entering code
	if m.currentStep == StepEnterCode {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

// View renders the OAuth flow UI
func (m *OAuthFlowModel) View() string {
	if m.done {
		return styles.SuccessStyle.Render("\n✓ Authentication complete!")
	}

	if m.canceled {
		return styles.WarningStyle.Render("\nAuthentication canceled.")
	}

	var s strings.Builder

	// Title
	title := styles.TitleStyle.Render("🔐 Claude Pro/Max OAuth Authentication")
	s.WriteString("\n" + title + "\n\n")

	// Progress indicator
	steps := []struct {
		step OAuthStep
		name string
	}{
		{StepGenerateURL, "Generate authorization URL"},
		{StepOpenBrowser, "Open browser for authentication"},
		{StepEnterCode, "Enter authorization code"},
		{StepExchangeToken, "Exchange code for tokens"},
		{StepSaveToken, "Save tokens securely"},
	}

	for _, step := range steps {
		var icon string
		var style lipgloss.Style

		if step.step < m.currentStep {
			icon = "✓"
			style = styles.SuccessStyle
		} else if step.step == m.currentStep {
			if m.err != nil && step.step == m.currentStep {
				icon = "✗"
				style = styles.ErrorStyle
			} else {
				icon = "◐"
				style = styles.InfoStyle
			}
		} else {
			icon = "○"
			style = styles.DescriptionStyle
		}

		s.WriteString(fmt.Sprintf("%s %s\n", style.Render(icon), style.Render(step.name)))
	}

	s.WriteString("\n")

	// Show content based on current step
	switch m.currentStep {
	case StepOpenBrowser:
		if m.authURL != "" {
			s.WriteString(styles.InfoStyle.Render("Opening browser to authorization page...") + "\n")
			s.WriteString(styles.DescriptionStyle.Render("If the browser doesn't open, please visit this URL manually:") + "\n\n")
			
			// URL box
			urlBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.Primary).
				Padding(1, 2).
				Width(min(len(m.authURL)+4, 100))
			
			s.WriteString(urlBox.Render(m.authURL) + "\n\n")
		}

	case StepEnterCode:
		if m.authURL != "" {
			s.WriteString(styles.SubtitleStyle.Render("Waiting for authorization...") + "\n\n")
			s.WriteString(styles.InfoStyle.Render("After authorizing, you'll receive a code.") + "\n")
			s.WriteString(styles.InfoStyle.Render("Enter it below:") + "\n\n")
			s.WriteString(m.textInput.View() + "\n\n")
			s.WriteString(styles.HelpStyle.Render("Press Enter to submit, Esc to cancel"))
		}

	case StepExchangeToken:
		s.WriteString(styles.InfoStyle.Render("Exchanging authorization code for tokens..."))

	case StepSaveToken:
		s.WriteString(styles.InfoStyle.Render("Saving tokens securely..."))
	}

	// Show error if any
	if m.err != nil {
		s.WriteString("\n\n" + styles.ErrorStyle.Render("Error: "+m.err.Error()))
	}

	return s.String()
}

// Message types

type AuthURLMsg struct {
	URL string
}

type AdvanceStepMsg struct{}

type AuthProgressMsg struct {
	Step OAuthStep
}

type AuthErrorMsg struct {
	Error error
}

type AuthCompleteMsg struct{}

// RunOAuthFlow runs the interactive OAuth flow and returns the authorization code
func RunOAuthFlow(authURL string) (string, error) {
	model := NewOAuthFlow()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Send the auth URL to the model
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure program is ready
		p.Send(AuthURLMsg{URL: authURL})
	}()

	// Run the program in a goroutine
	go func() {
		if _, err := p.Run(); err != nil {
			model.errChan <- err
		}
	}()

	// Wait for either code or error
	select {
	case code := <-model.codeChan:
		return code, nil
	case err := <-model.errChan:
		return "", err
	}
}

// Helper to update progress from outside
func UpdateOAuthProgress(p *tea.Program, step OAuthStep) {
	p.Send(AuthProgressMsg{Step: step})
}

// Helper to signal completion
func CompleteOAuthFlow(p *tea.Program) {
	p.Send(AuthCompleteMsg{})
}

// Helper to signal error
func ErrorOAuthFlow(p *tea.Program, err error) {
	p.Send(AuthErrorMsg{Error: err})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}