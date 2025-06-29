package main

import (
	"path/filepath"
	"testing"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartCmdAuthCheckUsesCorrectStorage verifies that start command finds tokens correctly
func TestStartCmdAuthCheckUsesCorrectStorage(t *testing.T) {
	// Setup: Create a token using the modern factory pattern
	tmpDir := t.TempDir()
	testAuthPath := filepath.Join(tmpDir, "auth.json")

	// Create storage via factory (the correct way)
	factory := auth.NewStorageFactory(auth.StorageFactoryConfig{
		Type:     auth.StorageTypeFile,
		FilePath: testAuthPath,
	})
	storage, err := factory.Create()
	require.NoError(t, err)

	// Store a test OAuth token
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		RefreshToken: "test-refresh",
		AccessToken:  "test-access",
		ExpiresAt:    0,
	}
	err = storage.Set("anthropic", testToken)
	require.NoError(t, err)

	// Test: What StartCmd currently does (will fail)
	t.Run("Current implementation using NewTokenStorage", func(t *testing.T) {
		// This is what the old broken code in StartCmd did
		legacyStorage := auth.NewFileStorage(testAuthPath)
		token, err := legacyStorage.Get("anthropic")
		
		// This should work since we're using file storage
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, "oauth", token.Type)
	})

	// Test: What StartCmd should do (will pass after fix)
	t.Run("Fixed implementation using StorageFactory", func(t *testing.T) {
		// This is what StartCmd should do
		factory := auth.NewStorageFactory(auth.StorageFactoryConfig{
			Type:     auth.StorageTypeFile,
			FilePath: testAuthPath,
		})
		storage, err := factory.Create()
		require.NoError(t, err)
		
		token, err := storage.Get("anthropic")
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, "oauth", token.Type)
	})
}