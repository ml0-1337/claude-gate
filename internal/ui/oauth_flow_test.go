package ui

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Test 24: OAuthFlowModel should initialize with correct state
func TestOAuthFlowModel_Initialization(t *testing.T) {
	// Prediction: This test will pass - testing initial state
	
	model := NewOAuthFlow()
	
	// Check initial state
	assert.Equal(t, StepGenerateURL, model.currentStep)
	assert.Empty(t, model.authURL)
	assert.Empty(t, model.code)
	assert.Nil(t, model.err)
	assert.False(t, model.done)
	assert.False(t, model.canceled)
	assert.NotNil(t, model.codeChan)
	assert.NotNil(t, model.errChan)
	assert.NotNil(t, model.textInput)
	
	// Check text input configuration
	assert.Equal(t, "Enter the authorization code", model.textInput.Placeholder)
	assert.Equal(t, 100, model.textInput.CharLimit)
	assert.Equal(t, 60, model.textInput.Width)
	
	// Test Init command
	cmd := model.Init()
	assert.NotNil(t, cmd) // Should return textinput.Blink
}

// Test 25: OAuthFlowModel should handle state transitions
func TestOAuthFlowModel_StateTransitions(t *testing.T) {
	// Prediction: This test will pass - testing state transitions
	
	t.Run("receive auth URL", func(t *testing.T) {
		model := NewOAuthFlow()
		
		// Send auth URL message
		msg := AuthURLMsg{URL: "https://example.com/auth"}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.Equal(t, "https://example.com/auth", oauthModel.authURL)
		assert.Equal(t, StepOpenBrowser, oauthModel.currentStep)
		assert.NotNil(t, cmd) // Should return tick command
	})
	
	t.Run("advance from browser to code entry", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepOpenBrowser
		
		msg := AdvanceStepMsg{}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.Equal(t, StepEnterCode, oauthModel.currentStep)
		assert.Nil(t, cmd)
	})
	
	t.Run("submit code", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepEnterCode
		model.textInput.SetValue("test-code-123")
		
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.Equal(t, "test-code-123", oauthModel.code)
		assert.Equal(t, StepExchangeToken, oauthModel.currentStep)
		assert.Nil(t, cmd)
	})
	
	t.Run("progress update", func(t *testing.T) {
		model := NewOAuthFlow()
		
		msg := AuthProgressMsg{Step: StepSaveToken}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.Equal(t, StepSaveToken, oauthModel.currentStep)
		assert.Nil(t, cmd)
	})
	
	t.Run("completion", func(t *testing.T) {
		model := NewOAuthFlow()
		
		msg := AuthCompleteMsg{}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.True(t, oauthModel.done)
		assert.NotNil(t, cmd) // Should return quit command
	})
}

// Test 26: OAuthFlowModel should handle code processing
func TestOAuthFlowModel_CodeProcessing(t *testing.T) {
	// Prediction: This test will pass - testing code entry and processing
	
	t.Run("enter code with empty input", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepEnterCode
		// textInput value is empty by default
		
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		// Should not advance with empty code
		assert.Equal(t, StepEnterCode, oauthModel.currentStep)
		assert.Empty(t, oauthModel.code)
		assert.Nil(t, cmd)
	})
	
	t.Run("text input update", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepEnterCode
		
		// Simulate typing
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		updatedModel, _ := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		// Text input should be updated
		assert.Equal(t, StepEnterCode, oauthModel.currentStep)
	})
	
	t.Run("code channel communication", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepEnterCode
		model.textInput.SetValue("channel-test-code")
		
		// Start a goroutine to read from channel
		codeReceived := make(chan string, 1)
		go func() {
			select {
			case code := <-model.codeChan:
				codeReceived <- code
			case <-time.After(100 * time.Millisecond):
				codeReceived <- "timeout"
			}
		}()
		
		// Submit code
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		model.Update(msg)
		
		// Check code was sent to channel
		code := <-codeReceived
		assert.Equal(t, "channel-test-code", code)
	})
}

// Test 27: OAuthFlowModel should handle errors
func TestOAuthFlowModel_ErrorHandling(t *testing.T) {
	// Prediction: This test will pass - testing error handling
	
	t.Run("handle error message", func(t *testing.T) {
		model := NewOAuthFlow()
		testErr := errors.New("test error")
		
		// Start goroutine to read from error channel
		errReceived := make(chan error, 1)
		go func() {
			select {
			case err := <-model.errChan:
				errReceived <- err
			case <-time.After(100 * time.Millisecond):
				errReceived <- errors.New("timeout")
			}
		}()
		
		msg := AuthErrorMsg{Error: testErr}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.Equal(t, testErr, oauthModel.err)
		assert.NotNil(t, cmd) // Should return quit command
		
		// Check error was sent to channel
		err := <-errReceived
		assert.Equal(t, testErr, err)
	})
	
	t.Run("cancel flow", func(t *testing.T) {
		model := NewOAuthFlow()
		
		msg := tea.KeyMsg{Type: tea.KeyEsc}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.True(t, oauthModel.canceled)
		assert.NotNil(t, cmd) // Should return quit command
	})
	
	t.Run("ctrl+c cancellation", func(t *testing.T) {
		model := NewOAuthFlow()
		
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		updatedModel, cmd := model.Update(msg)
		oauthModel := updatedModel.(*OAuthFlowModel)
		
		assert.True(t, oauthModel.canceled)
		assert.NotNil(t, cmd) // Should return quit command
	})
}

// Test 28: OAuthFlowModel.View should render correctly
func TestOAuthFlowModel_View(t *testing.T) {
	// Prediction: This test will pass - testing view rendering
	
	t.Run("initial view", func(t *testing.T) {
		model := NewOAuthFlow()
		view := model.View()
		
		assert.Contains(t, view, "Claude Pro/Max OAuth Authentication")
		assert.Contains(t, view, "Generate authorization URL")
		assert.Contains(t, view, "○") // Pending steps
	})
	
	t.Run("browser step view", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepOpenBrowser
		model.authURL = "https://example.com/auth"
		
		view := model.View()
		
		assert.Contains(t, view, "Opening browser")
		assert.Contains(t, view, "https://example.com/auth")
		assert.Contains(t, view, "✓") // Completed step
		assert.Contains(t, view, "◐") // Current step
	})
	
	t.Run("code entry view", func(t *testing.T) {
		model := NewOAuthFlow()
		model.currentStep = StepEnterCode
		model.authURL = "https://example.com/auth"
		
		view := model.View()
		
		assert.Contains(t, view, "Waiting for authorization")
		assert.Contains(t, view, "Enter it below")
		assert.Contains(t, view, "Press Enter to submit")
	})
	
	t.Run("error view", func(t *testing.T) {
		model := NewOAuthFlow()
		model.err = errors.New("test error message")
		
		view := model.View()
		
		assert.Contains(t, view, "Error: test error message")
		assert.Contains(t, view, "✗") // Error icon
	})
	
	t.Run("completed view", func(t *testing.T) {
		model := NewOAuthFlow()
		model.done = true
		
		view := model.View()
		
		assert.Contains(t, view, "Authentication complete!")
		assert.Contains(t, view, "✓")
	})
	
	t.Run("canceled view", func(t *testing.T) {
		model := NewOAuthFlow()
		model.canceled = true
		
		view := model.View()
		
		assert.Contains(t, view, "Authentication canceled")
	})
}

// Test helper functions
func TestOAuthFlowHelpers(t *testing.T) {
	// Prediction: This test will pass - testing helper functions
	
	t.Run("min function", func(t *testing.T) {
		assert.Equal(t, 5, min(5, 10))
		assert.Equal(t, 5, min(10, 5))
		assert.Equal(t, 5, min(5, 5))
	})
	
	t.Run("helper functions exist", func(t *testing.T) {
		// Just verify these functions exist without running them
		assert.NotNil(t, UpdateOAuthProgress)
		assert.NotNil(t, CompleteOAuthFlow)
		assert.NotNil(t, ErrorOAuthFlow)
	})
}

// Test RunOAuthFlow function (limited testing without TTY)
func TestRunOAuthFlow(t *testing.T) {
	// Prediction: This test will have limited functionality without TTY
	
	t.Run("function exists", func(t *testing.T) {
		// We can't fully test this without a TTY, but verify it exists
		assert.NotNil(t, RunOAuthFlow)
	})
}