// +build integration

package integration

import (
	"testing"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeCodeStorageIntegration(t *testing.T) {
	// Skip if not on a system with Claude Code installed
	t.Skip("Manual test - requires Claude Code to be installed with valid credentials")
	
	// Create Claude Code storage adapter
	storage, err := auth.NewClaudeCodeStorage()
	require.NoError(t, err)
	
	// Test availability
	assert.True(t, storage.IsAvailable())
	
	// Test name
	assert.Equal(t, "claude-code-adapter", storage.Name())
	
	// Test listing providers
	providers, err := storage.List()
	require.NoError(t, err)
	
	// If Claude Code has credentials, should return anthropic
	if len(providers) > 0 {
		assert.Contains(t, providers, "anthropic")
		
		// Test getting token
		token, err := storage.Get("anthropic")
		require.NoError(t, err)
		require.NotNil(t, token)
		
		// Verify token fields
		assert.Equal(t, "oauth", token.Type)
		assert.NotEmpty(t, token.AccessToken)
		assert.NotEmpty(t, token.RefreshToken)
		assert.Greater(t, token.ExpiresAt, int64(0))
	}
	
	// Test that Set is a no-op
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test",
		RefreshToken: "test",
		ExpiresAt:    1234567890,
	}
	err = storage.Set("anthropic", testToken)
	assert.NoError(t, err)
	
	// Test that Remove is a no-op
	err = storage.Remove("anthropic")
	assert.NoError(t, err)
}

func TestStorageFactoryWithClaudeCode(t *testing.T) {
	// Test that storage factory can create Claude Code storage
	factory := auth.NewStorageFactory(auth.StorageFactoryConfig{
		Type: auth.StorageTypeClaudeCode,
	})
	
	storage, err := factory.Create()
	require.NoError(t, err)
	require.NotNil(t, storage)
	
	// Verify it's the right type
	assert.Equal(t, "claude-code-adapter", storage.Name())
}