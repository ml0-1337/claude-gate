package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageFactory_CreateClaudeCodeStorage(t *testing.T) {
	// Test that storage factory can create Claude Code storage adapter
	
	factory := NewStorageFactory(StorageFactoryConfig{
		Type: StorageTypeClaudeCode,
	})
	
	// Create storage
	storage, err := factory.Create()
	
	// Should succeed
	require.NoError(t, err)
	require.NotNil(t, storage)
	
	// Verify it's the correct type
	assert.Equal(t, "claude-code-adapter", storage.Name())
	
	// Verify it implements StorageBackend interface
	var _ StorageBackend = storage
}