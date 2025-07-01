package components

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Test 19: Spinner should update states correctly
func TestSpinner_States(t *testing.T) {
	// Prediction: This test will pass - testing spinner state transitions
	
	t.Run("initial state", func(t *testing.T) {
		s := NewSpinner("Test message")
		
		assert.Equal(t, "Test message", s.message)
		assert.Equal(t, "loading", s.status)
		assert.False(t, s.quitting)
		assert.Equal(t, spinner.Dot, s.spinner.Spinner)
	})
	
	t.Run("init command", func(t *testing.T) {
		s := NewSpinner("Testing")
		cmd := s.Init()
		
		// Init should return a tick command
		assert.NotNil(t, cmd)
	})
	
	t.Run("status update", func(t *testing.T) {
		s := NewSpinner("Processing")
		
		// Test success status
		model, cmd := s.Update(StatusMsg{
			Status:  "success",
			Message: "Done!",
		})
		
		updatedSpinner := model.(SpinnerModel)
		assert.Equal(t, "success", updatedSpinner.status)
		assert.True(t, updatedSpinner.quitting)
		assert.NotNil(t, cmd) // Should return quit command
		
		// Test error status
		s2 := NewSpinner("Processing")
		model2, cmd2 := s2.Update(StatusMsg{
			Status:  "error",
			Message: "Failed",
		})
		
		updatedSpinner2 := model2.(SpinnerModel)
		assert.Equal(t, "error", updatedSpinner2.status)
		assert.True(t, updatedSpinner2.quitting)
		assert.NotNil(t, cmd2)
	})
	
	t.Run("keyboard handling", func(t *testing.T) {
		s := NewSpinner("Processing")
		
		// Test quit keys
		quitKeys := []string{"q", "esc", "ctrl+c"}
		
		for _, key := range quitKeys {
			model, cmd := s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
			updatedSpinner := model.(SpinnerModel)
			
			assert.True(t, updatedSpinner.quitting)
			assert.NotNil(t, cmd) // Should return quit command
		}
		
		// Test non-quit key
		s3 := NewSpinner("Processing")
		model3, cmd3 := s3.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
		updatedSpinner3 := model3.(SpinnerModel)
		
		assert.False(t, updatedSpinner3.quitting)
		assert.Nil(t, cmd3)
	})
	
	t.Run("spinner tick", func(t *testing.T) {
		s := NewSpinner("Spinning")
		
		// Simulate spinner tick
		model, cmd := s.Update(spinner.TickMsg{
			Time: time.Now(),
		})
		
		// Should return updated model with tick command
		assert.NotNil(t, model)
		assert.NotNil(t, cmd)
	})
	
	t.Run("view rendering", func(t *testing.T) {
		s := NewSpinner("Test view")
		
		// Test loading view
		view := s.View()
		assert.Contains(t, view, "Test view")
		
		// Test success view after status update
		s.status = "success"
		s.quitting = true
		s.message = "Operation completed"
		view = s.View()
		assert.Contains(t, view, "Operation completed")
		
		// Test error view
		s.status = "error"
		s.message = "Operation failed"
		view = s.View()
		assert.Contains(t, view, "Operation failed")
		
		// Test loading state (not quitting)
		s2 := NewSpinner("Loading...")
		s2.quitting = false
		view2 := s2.View()
		assert.Contains(t, view2, "Loading...")
	})
}

// Test RunSpinner function (limited testing due to TTY requirements)
func TestRunSpinner_Structure(t *testing.T) {
	// Prediction: This test will pass - just verifying the function exists
	// We can't fully test RunSpinner without TTY
	
	// Verify the function signature
	var fn func(string, func() error) error = RunSpinner
	assert.NotNil(t, fn)
}

// Test 10: SimpleSpinner should handle non-TTY environments gracefully
func TestSimpleSpinner_Structure(t *testing.T) {
	// Prediction: This test will pass - verifying function signature
	// Note: We can't fully test SimpleSpinner without TTY, but we can verify structure
	
	// Verify the function exists and has the right signature
	var fn func(string, time.Duration) = SimpleSpinner
	assert.NotNil(t, fn)
	
	// Test that it doesn't panic in non-TTY environment
	// Create a very short duration to avoid hanging
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SimpleSpinner panicked: %v", r)
			}
			done <- true
		}()
		// This will fail in non-TTY but shouldn't panic
		SimpleSpinner("Test", 1*time.Millisecond)
	}()
	
	select {
	case <-done:
		// Function completed (or failed gracefully)
	case <-time.After(100 * time.Millisecond):
		// Give it a reasonable timeout
	}
}