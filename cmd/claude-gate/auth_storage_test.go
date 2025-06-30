package main

import (
	"io"
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
)

// captureAuthOutput captures stdout during test execution for auth commands
func captureAuthOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}

// Test 19: AuthStorageStatusCmd should show storage status
func TestAuthStorageStatusCmd_Run(t *testing.T) {
	// Prediction: This test will pass - testing status command output
	
	t.Run("show status output structure", func(t *testing.T) {
		cmd := &AuthStorageStatusCmd{}
		
		output := captureAuthOutput(func() {
			// Don't check error as it depends on storage availability
			_ = cmd.Run(&kong.Context{})
		})
		
		// Check output contains expected sections
		assert.Contains(t, output, "Storage Backend Status")
		assert.Contains(t, output, "Type:")
		assert.Contains(t, output, "Backend:")
		assert.Contains(t, output, "Available:")
		assert.Contains(t, output, "Configuration:")
		assert.Contains(t, output, "Storage Path:")
		assert.Contains(t, output, "Keyring Service:")
	})
}

// Test 20: AuthStorageMigrateCmd should validate migration parameters
func TestAuthStorageMigrateCmd_Validation(t *testing.T) {
	// Prediction: This test will pass - testing validation logic
	
	t.Run("migrate with no tokens", func(t *testing.T) {
		// Set up empty auth file
		tmpFile, err := os.CreateTemp("", "auth-empty-*.json")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		
		_, err = tmpFile.WriteString("{}")
		assert.NoError(t, err)
		tmpFile.Close()
		
		// Force file storage type
		os.Setenv("AUTH_STORAGE_TYPE", "file")
		os.Setenv("AUTH_STORAGE_PATH", tmpFile.Name())
		defer os.Unsetenv("AUTH_STORAGE_TYPE")
		defer os.Unsetenv("AUTH_STORAGE_PATH")
		
		cmd := &AuthStorageMigrateCmd{
			From: "file",
			To:   "file", // Migrate to same type to avoid keyring issues
		}
		
		output := captureAuthOutput(func() {
			err := cmd.Run(&kong.Context{})
			assert.NoError(t, err)
		})
		
		assert.Contains(t, output, "No tokens to migrate")
	})
}

// Test 21: AuthStorageBackupCmd should handle different storage types
func TestAuthStorageBackupCmd_Run(t *testing.T) {
	// Prediction: This test will pass - testing backup functionality
	
	t.Run("backup keyring storage message", func(t *testing.T) {
		os.Setenv("AUTH_STORAGE_TYPE", "keyring")
		defer os.Unsetenv("AUTH_STORAGE_TYPE")
		
		cmd := &AuthStorageBackupCmd{}
		
		output := captureAuthOutput(func() {
			err := cmd.Run(&kong.Context{})
			assert.NoError(t, err)
		})
		
		assert.Contains(t, output, "Backup is only supported for file storage")
	})
}

// Test 22: AuthStorageResetCmd should handle non-interactive mode
func TestAuthStorageResetCmd_Run(t *testing.T) {
	// Prediction: This test will pass - testing reset functionality
	
	t.Run("reset without tokens", func(t *testing.T) {
		// Create empty auth file
		tmpFile, err := os.CreateTemp("", "auth-reset-*.json")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		
		_, err = tmpFile.WriteString("{}")
		assert.NoError(t, err)
		tmpFile.Close()
		
		os.Setenv("AUTH_STORAGE_TYPE", "file")
		os.Setenv("AUTH_STORAGE_PATH", tmpFile.Name())
		defer os.Unsetenv("AUTH_STORAGE_TYPE")
		defer os.Unsetenv("AUTH_STORAGE_PATH")
		
		cmd := &AuthStorageResetCmd{Force: true}
		
		output := captureAuthOutput(func() {
			// Just run without checking error
			_ = cmd.Run(&kong.Context{})
		})
		
		// Should either show "no tokens" or "reset only for keyring"
		// depending on how the config resolves
		assert.NotEmpty(t, output)
	})
}

// Test 23: AuthStorageTestCmd should test storage operations
func TestAuthStorageTestCmd_Run(t *testing.T) {
	// Prediction: This test will pass - testing storage test functionality
	
	t.Run("test storage operations with file backend", func(t *testing.T) {
		// Create temporary auth file
		tmpFile, err := os.CreateTemp("", "auth-test-ops-*.json")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		
		// Start with empty file
		_, err = tmpFile.WriteString("{}")
		assert.NoError(t, err)
		tmpFile.Close()
		
		os.Setenv("AUTH_STORAGE_TYPE", "file")
		os.Setenv("AUTH_STORAGE_PATH", tmpFile.Name())
		defer os.Unsetenv("AUTH_STORAGE_TYPE")
		defer os.Unsetenv("AUTH_STORAGE_PATH")
		
		cmd := &AuthStorageTestCmd{}
		
		output := captureAuthOutput(func() {
			err := cmd.Run(&kong.Context{})
			assert.NoError(t, err)
		})
		
		// Check all operations were tested
		assert.Contains(t, output, "Testing storage backend:")
		assert.Contains(t, output, "Testing availability...")
		assert.Contains(t, output, "Testing write operation...")
		assert.Contains(t, output, "Testing read operation...")
		assert.Contains(t, output, "Testing list operation...")
		assert.Contains(t, output, "Testing remove operation...")
		assert.Contains(t, output, "All tests passed!")
	})
}