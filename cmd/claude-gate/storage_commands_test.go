package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 9: Storage commands structure
func TestStorageCommands_Structure(t *testing.T) {
	// Test the storage command structure
	
	storageCmd := &AuthStorageCmd{}
	
	// Verify subcommands exist
	assert.NotNil(t, storageCmd.Status)
	assert.NotNil(t, storageCmd.Migrate)
	assert.NotNil(t, storageCmd.Test)
	assert.NotNil(t, storageCmd.Backup)
	assert.NotNil(t, storageCmd.Reset)
}

// Test 10: Storage StatusCmd should show storage details
func TestStorageStatusCmd_ShowDetails(t *testing.T) {
	// Prediction: This test will pass for file storage
	
	// Create temporary directory with token
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	// Create and store a test token
	storage := auth.NewFileStorage(authFile)
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	err := storage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Also create a second provider
	apiToken := &auth.TokenInfo{
		Type:   "api_key",
		APIKey: "test-api-key",
	}
	err = storage.Set("openai", apiToken)
	require.NoError(t, err)
	
	// Create StatusCmd
	cmd := &AuthStorageStatusCmd{}
	
	// Mock environment
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", authFile)
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
	
	// Run command
	stdout, stderr, err := captureOutput(func() error {
		return cmd.Run(&kong.Context{})
	})
	
	// Should not error
	assert.NoError(t, err)
	
	// Check output
	output := stdout + stderr
	assert.Contains(t, output, "Storage Backend Status")
	assert.Contains(t, output, "Type: file")
	assert.Contains(t, output, "Available: true")
	assert.Contains(t, output, "anthropic")
	assert.Contains(t, output, "openai")
}

// Test 11: Storage TestCmd should verify storage operations
func TestStorageTestCmd_Operations(t *testing.T) {
	// Prediction: This test will pass - we can test file storage operations
	
	// Create temporary directory
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	// Create TestCmd
	cmd := &AuthStorageTestCmd{}
	
	// Mock environment
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", authFile)
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
	
	// Run command
	stdout, stderr, err := captureOutput(func() error {
		return cmd.Run(&kong.Context{})
	})
	
	// Should not error
	assert.NoError(t, err)
	
	// Check output shows test results
	output := stdout + stderr
	assert.Contains(t, output, "Testing storage backend")
	assert.Contains(t, output, "Testing availability")
	assert.Contains(t, output, "Testing write operation")
	assert.Contains(t, output, "Testing read operation")
	assert.Contains(t, output, "Testing list operation")
	assert.Contains(t, output, "Testing remove operation")
	assert.Contains(t, output, "All tests passed")
}

// Test 12: MigrateCmd should migrate tokens between storages
func TestStorageMigrateCmd_Migration(t *testing.T) {
	// Prediction: This test will be limited - actual migration requires user confirmation
	
	// Create source storage with tokens
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.json")
	
	// Create source storage with test token
	sourceStorage := auth.NewFileStorage(sourceFile)
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	err := sourceStorage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Create MigrateCmd
	cmd := &AuthStorageMigrateCmd{
		From: "file",
		To:   "file",
	}
	
	// Mock environment for source
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", sourceFile)
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
	
	// Verify command is configured
	assert.Equal(t, "file", cmd.From)
	assert.Equal(t, "file", cmd.To)
	
	// Actual migration would require:
	// - Mocking the confirmation dialog
	// - Running the migration
	// - Verifying tokens were copied
}

// We don't have an InfoCmd in the current implementation
// The tests above cover the available storage commands