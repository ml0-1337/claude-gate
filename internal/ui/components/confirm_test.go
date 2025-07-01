package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Test 20: Confirm dialog should handle user input
func TestConfirm_UserInput(t *testing.T) {
	// Prediction: This test will pass - testing confirm dialog behavior
	
	t.Run("initial state", func(t *testing.T) {
		c := NewConfirm("Are you sure?")
		
		assert.Equal(t, "Are you sure?", c.question)
		assert.False(t, c.answer)
		assert.False(t, c.answered)
	})
	
	t.Run("init command", func(t *testing.T) {
		c := NewConfirm("Continue?")
		cmd := c.Init()
		
		// Init should return nil
		assert.Nil(t, cmd)
	})
	
	t.Run("handle yes input", func(t *testing.T) {
		testCases := []string{"y", "Y"}
		
		for _, key := range testCases {
			c := NewConfirm("Proceed?")
			model, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
			
			updatedConfirm := model.(ConfirmModel)
			assert.True(t, updatedConfirm.answer)
			assert.True(t, updatedConfirm.answered)
			assert.NotNil(t, cmd) // Should return quit command
		}
	})
	
	t.Run("handle no input", func(t *testing.T) {
		testCases := []string{"n", "N"}
		
		for _, key := range testCases {
			c := NewConfirm("Delete?")
			model, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
			
			updatedConfirm := model.(ConfirmModel)
			assert.False(t, updatedConfirm.answer)
			assert.True(t, updatedConfirm.answered)
			assert.NotNil(t, cmd) // Should return quit command
		}
	})
	
	t.Run("handle escape/cancel", func(t *testing.T) {
		testCases := []string{"ctrl+c", "esc"}
		
		for _, key := range testCases {
			c := NewConfirm("Cancel?")
			model, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
			
			updatedConfirm := model.(ConfirmModel)
			assert.False(t, updatedConfirm.answer)
			assert.True(t, updatedConfirm.answered)
			assert.NotNil(t, cmd) // Should return quit command
		}
	})
	
	t.Run("handle other keys", func(t *testing.T) {
		c := NewConfirm("Confirm?")
		model, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
		
		updatedConfirm := model.(ConfirmModel)
		assert.False(t, updatedConfirm.answer)
		assert.False(t, updatedConfirm.answered)
		assert.Nil(t, cmd) // Should not quit
	})
	
	t.Run("view rendering", func(t *testing.T) {
		c := NewConfirm("Delete all files?")
		view := c.View()
		
		assert.Contains(t, view, "Delete all files?")
		assert.Contains(t, view, "(y/N)")
	})
	
	t.Run("view changes after answer", func(t *testing.T) {
		c := NewConfirm("Delete file?")
		
		// Before answering
		view := c.View()
		assert.Contains(t, view, "(y/N)")
		
		// After answering yes
		c.answer = true
		c.answered = true
		view = c.View()
		assert.Contains(t, view, "Yes")
		assert.NotContains(t, view, "(y/N)")
		
		// After answering no
		c2 := NewConfirm("Delete file?")
		c2.answer = false
		c2.answered = true
		view2 := c2.View()
		assert.Contains(t, view2, "No")
	})
}

// Test Confirm helper function (limited testing due to TTY requirements)
func TestConfirm_HelperFunction(t *testing.T) {
	// Prediction: This test will pass - just verifying the function exists
	// We can't fully test Confirm() without TTY
	
	// Verify the function signature
	var fn func(string) bool = Confirm
	assert.NotNil(t, fn)
}

// Test 21: ConfirmDefaultModel should handle default values
func TestConfirmDefaultModel_DefaultValues(t *testing.T) {
	// Prediction: This test will pass - testing default value handling
	
	t.Run("default yes", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (Y/n)",
				answer:   true,
				answered: false,
			},
			defaultYes: true,
		}
		
		// Test pressing Enter uses default
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.True(t, confirm.answer)
		assert.True(t, confirm.answered)
		assert.NotNil(t, cmd) // Should quit
	})
	
	t.Run("default no", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (y/N)",
				answer:   false,
				answered: false,
			},
			defaultYes: false,
		}
		
		// Test pressing Enter uses default
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.False(t, confirm.answer)
		assert.True(t, confirm.answered)
		assert.NotNil(t, cmd) // Should quit
	})
	
	t.Run("explicit yes overrides default", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (y/N)",
				answer:   false,
				answered: false,
			},
			defaultYes: false,
		}
		
		// Test pressing Y overrides default
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Y")})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.True(t, confirm.answer)
		assert.True(t, confirm.answered)
		assert.NotNil(t, cmd)
	})
	
	t.Run("explicit no overrides default", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (Y/n)",
				answer:   true,
				answered: false,
			},
			defaultYes: true,
		}
		
		// Test pressing N overrides default
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("N")})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.False(t, confirm.answer)
		assert.True(t, confirm.answered)
		assert.NotNil(t, cmd)
	})
	
	t.Run("escape cancels", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (Y/n)",
				answer:   true,
				answered: false,
			},
			defaultYes: true,
		}
		
		// Test escape key cancels
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.False(t, confirm.answer) // Always false on cancel
		assert.True(t, confirm.answered)
		assert.NotNil(t, cmd)
	})
	
	t.Run("other keys ignored", func(t *testing.T) {
		model := &ConfirmDefaultModel{
			ConfirmModel: ConfirmModel{
				question: "Continue? (Y/n)",
				answer:   true,
				answered: false,
			},
			defaultYes: true,
		}
		
		// Test other keys are ignored
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		confirm := updatedModel.(*ConfirmDefaultModel)
		
		assert.True(t, confirm.answer) // Unchanged
		assert.False(t, confirm.answered) // Not answered
		assert.Nil(t, cmd)
	})
}

// Test ConfirmWithDefault helper function
func TestConfirmWithDefault_HelperFunction(t *testing.T) {
	// Prediction: This test will pass - just verifying the function exists
	// We can't fully test ConfirmWithDefault() without TTY
	
	// Verify the function signature
	var fn func(string, bool) bool = ConfirmWithDefault
	assert.NotNil(t, fn)
}