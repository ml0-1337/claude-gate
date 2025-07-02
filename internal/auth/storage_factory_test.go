package auth

import (
	"runtime"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageFactory_Create(t *testing.T) {
	// Test file storage creation
	t.Run("file storage", func(t *testing.T) {
		testPath := filepath.Join(os.TempDir(), "test-auth.json")
		factory := NewStorageFactory(StorageFactoryConfig{
			Type:     StorageTypeFile,
			FilePath: testPath,
		})
		
		storage, err := factory.Create()
		assert.NoError(t, err)
		assert.NotNil(t, storage)
		assert.IsType(t, &FileStorage{}, storage)
		assert.Equal(t, "file:"+testPath, storage.Name())
	})
	
	// Test auto storage creation
	t.Run("auto storage", func(t *testing.T) {
		testPath := filepath.Join(os.TempDir(), "test-auth.json")
		factory := NewStorageFactory(StorageFactoryConfig{
			Type:     StorageTypeAuto,
			FilePath: testPath,
		})
		
		storage, err := factory.Create()
		assert.NoError(t, err)
		assert.NotNil(t, storage)
		// Auto mode will select best available backend
		// Could be either keyring or file depending on system
		assert.True(t, storage.Name() == "keyring:claude-gate" || strings.Contains(storage.Name(), "file:"))
	})
	
	// Test Claude Code storage creation
	t.Run("claude code storage", func(t *testing.T) {
		factory := NewStorageFactory(StorageFactoryConfig{
			Type: StorageTypeClaudeCode,
		})
		
		storage, err := factory.Create()
		assert.NoError(t, err)
		assert.NotNil(t, storage)
		
		// Verify it's the right type
		ccs, ok := storage.(*ClaudeCodeStorage)
		assert.True(t, ok, "Expected ClaudeCodeStorage type")
		assert.NotNil(t, ccs)
		assert.Equal(t, "claude-code-adapter", storage.Name())
	})
	
	// Test unknown storage type
	t.Run("unknown storage type", func(t *testing.T) {
		factory := &StorageFactory{
			storageType: StorageType("unknown"),
		}
		
		storage, err := factory.Create()
		assert.Error(t, err)
		assert.Nil(t, storage)
		assert.Contains(t, err.Error(), "unknown storage type")
	})
}

func TestStorageFactory_CreateWithMigration(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "auth.json")
	
	// Create file storage with test data
	fileStorage := NewFileStorage(jsonPath)
	testToken := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "test-refresh",
		AccessToken:  "test-access",
		ExpiresAt:    1234567890,
	}
	
	err := fileStorage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Create factory that will migrate to file storage (simulating keyring)
	factory := NewStorageFactory(StorageFactoryConfig{
		Type:     StorageTypeFile,
		FilePath: filepath.Join(tempDir, "auth-migrated.json"),
	})
	
	// Create with migration
	storage, err := factory.CreateWithMigration()
	assert.NoError(t, err)
	assert.NotNil(t, storage)
	
	// Original file should still exist
	_, err = os.Stat(jsonPath)
	assert.NoError(t, err)
}

func TestStorageFactory_Defaults(t *testing.T) {
	// Test with empty config
	factory := NewStorageFactory(StorageFactoryConfig{})
	
	assert.Equal(t, StorageTypeAuto, factory.storageType)
	assert.Equal(t, "claude-gate", factory.keyringConfig.ServiceName)
	assert.Contains(t, factory.filePath, ".claude-gate/auth.json")
	assert.NotNil(t, factory.passwordPrompt)
}

func TestIsKeyringAvailable(t *testing.T) {
	// This test is platform-specific
	available := isKeyringAvailable()
	
	switch runtime.GOOS {
	case "darwin":
		// Should always be available on macOS
		assert.True(t, available)
	case "linux":
		// Depends on display environment
		hasDisplay := os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
		assert.Equal(t, hasDisplay, available)
	default:
		// Other platforms should return false
		assert.False(t, available)
	}
}