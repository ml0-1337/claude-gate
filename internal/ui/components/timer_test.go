package components

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Test 15: TimerModel should initialize with correct state
func TestTimerModel_Initialization(t *testing.T) {
	// Prediction: This test will pass - testing initial state
	
	duration := 5 * time.Minute
	message := "Test Timer"
	model := NewTimer(duration, message)
	
	// Check initial state
	assert.Equal(t, duration, model.duration)
	assert.Equal(t, duration, model.remaining)
	assert.Equal(t, message, model.message)
	assert.False(t, model.expired)
	assert.False(t, model.quitting)
	
	// Test Init command
	cmd := model.Init()
	assert.NotNil(t, cmd) // Should return tick command
}

// Test 16: TimerModel.Update should handle tick events
func TestTimerModel_UpdateTickEvents(t *testing.T) {
	// Prediction: This test will pass - testing countdown behavior
	
	model := NewTimer(3*time.Second, "Test")
	
	t.Run("first tick", func(t *testing.T) {
		msg := tickMsg(time.Now())
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		// Should decrease remaining time by 1 second
		assert.Equal(t, 2*time.Second, timerModel.remaining)
		assert.False(t, timerModel.expired)
		assert.NotNil(t, cmd) // Should return next tick command
	})
	
	t.Run("second tick", func(t *testing.T) {
		model.remaining = 2 * time.Second
		msg := tickMsg(time.Now())
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.Equal(t, 1*time.Second, timerModel.remaining)
		assert.False(t, timerModel.expired)
		assert.NotNil(t, cmd) // Should return next tick command
	})
	
	t.Run("final tick expires timer", func(t *testing.T) {
		model.remaining = 1 * time.Second
		msg := tickMsg(time.Now())
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.Equal(t, time.Duration(0), timerModel.remaining)
		assert.True(t, timerModel.expired)
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd) // Should return batch command with Quit
	})
}

// Test 17: TimerModel.Update should handle other messages
func TestTimerModel_UpdateOtherMessages(t *testing.T) {
	// Prediction: This test will pass - testing message handling
	
	model := NewTimer(5*time.Minute, "Test")
	
	t.Run("handle quit key", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd) // Should return Quit command
	})
	
	t.Run("handle escape key", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyEsc}
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd) // Should return Quit command
	})
	
	t.Run("handle ctrl+c", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd) // Should return Quit command
	})
	
	t.Run("handle success status", func(t *testing.T) {
		msg := StatusMsg{Status: "success"}
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd) // Should return Quit command
	})
	
	t.Run("handle other status", func(t *testing.T) {
		msg := StatusMsg{Status: "error"}
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.False(t, timerModel.quitting)
		assert.Nil(t, cmd) // Should not quit
	})
	
	t.Run("handle unknown message", func(t *testing.T) {
		msg := struct{ Unknown string }{Unknown: "test"}
		updatedModel, cmd := model.Update(msg)
		
		assert.Equal(t, model, updatedModel)
		assert.Nil(t, cmd)
	})
}

// Test 18: TimerModel.View should render timer correctly
func TestTimerModel_View(t *testing.T) {
	// Prediction: This test will pass - testing view rendering
	
	t.Run("active timer", func(t *testing.T) {
		model := NewTimer(2*time.Minute+30*time.Second, "Test Timer")
		view := model.View()
		
		// Should show message
		assert.Contains(t, view, "Test Timer")
		// Should show time remaining
		assert.Contains(t, view, "02:30")
		// Should have progress bar
		assert.Contains(t, view, "░") // Empty progress
	})
	
	t.Run("timer with progress", func(t *testing.T) {
		model := NewTimer(2*time.Minute, "Test")
		model.remaining = 1 * time.Minute // 50% remaining
		view := model.View()
		
		assert.Contains(t, view, "01:00")
		// Should have both filled and empty progress
		assert.Contains(t, view, "█") // Filled progress
		assert.Contains(t, view, "░") // Empty progress
	})
	
	t.Run("expired timer", func(t *testing.T) {
		model := NewTimer(1*time.Minute, "Test")
		model.expired = true
		view := model.View()
		
		assert.Contains(t, view, "Timer expired!")
		assert.Contains(t, view, "⏰")
	})
	
	t.Run("quitting without expiration", func(t *testing.T) {
		model := NewTimer(1*time.Minute, "Test")
		model.quitting = true
		model.expired = false
		view := model.View()
		
		assert.Empty(t, view) // Should return empty string
	})
	
	t.Run("timer with seconds only", func(t *testing.T) {
		model := NewTimer(45*time.Second, "Test")
		view := model.View()
		
		assert.Contains(t, view, "00:45")
	})
}

// Test edge cases and special scenarios
func TestTimerModel_EdgeCases(t *testing.T) {
	// Prediction: This test will pass - testing edge cases
	
	t.Run("zero duration timer", func(t *testing.T) {
		model := NewTimer(0, "Instant Timer")
		
		// First tick should expire immediately
		msg := tickMsg(time.Now())
		updatedModel, cmd := model.Update(msg)
		timerModel := updatedModel.(TimerModel)
		
		assert.True(t, timerModel.expired)
		assert.True(t, timerModel.quitting)
		assert.NotNil(t, cmd)
	})
	
	t.Run("negative remaining time", func(t *testing.T) {
		model := NewTimer(1*time.Minute, "Test")
		model.remaining = -5 * time.Second
		
		// View should handle negative time gracefully
		assert.NotPanics(t, func() {
			view := model.View()
			assert.NotEmpty(t, view)
		})
	})
	
	t.Run("very long duration", func(t *testing.T) {
		model := NewTimer(99*time.Hour+59*time.Minute+59*time.Second, "Long Timer")
		view := model.View()
		
		// Should show large minutes value
		assert.Contains(t, view, "5999:59") // 99:59:59 in minutes:seconds
	})
}

// Test helper functions
func TestTimerHelperFunctions(t *testing.T) {
	// Prediction: This test will pass - testing helper functions
	
	t.Run("CountdownTimer creates timer", func(t *testing.T) {
		duration := 3 * time.Minute
		message := "Countdown"
		
		model := CountdownTimer(duration, message)
		timerModel, ok := model.(TimerModel)
		
		assert.True(t, ok)
		assert.Equal(t, duration, timerModel.duration)
		assert.Equal(t, message, timerModel.message)
	})
	
	t.Run("tickCmd function", func(t *testing.T) {
		cmd := tickCmd()
		assert.NotNil(t, cmd)
		// Can't test the actual tick behavior without running the command
	})
}

// Test callback functionality (limited testing without running the program)
func TestTimerWithCallback(t *testing.T) {
	// Prediction: This test will pass - but limited without running the program
	
	t.Run("function exists", func(t *testing.T) {
		// Verify the function can be called without panic
		assert.NotPanics(t, func() {
			// We can't actually test the callback execution without TTY
			var fn func(time.Duration, string, func()) = TimerWithCallback
			assert.NotNil(t, fn)
		})
	})
}