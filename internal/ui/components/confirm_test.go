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