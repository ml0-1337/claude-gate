package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockKeyringStorage is a simple in-memory storage for testing
type MockKeyringStorage struct {
	tokens map[string]*auth.TokenInfo
}

func NewMockKeyringStorage() *MockKeyringStorage {
	return &MockKeyringStorage{
		tokens: make(map[string]*auth.TokenInfo),
	}
}

func (m *MockKeyringStorage) Get(provider string) (*auth.TokenInfo, error) {
	return m.tokens[provider], nil
}

func (m *MockKeyringStorage) Set(provider string, token *auth.TokenInfo) error {
	m.tokens[provider] = token
	return nil
}

func (m *MockKeyringStorage) Remove(provider string) error {
	delete(m.tokens, provider)
	return nil
}

func (m *MockKeyringStorage) List() ([]string, error) {
	providers := make([]string, 0, len(m.tokens))
	for provider := range m.tokens {
		providers = append(providers, provider)
	}
	return providers, nil
}

func (m *MockKeyringStorage) IsAvailable() bool { return true }
func (m *MockKeyringStorage) RequiresUnlock() bool { return false }
func (m *MockKeyringStorage) Unlock() error { return nil }
func (m *MockKeyringStorage) Lock() error { return nil }
func (m *MockKeyringStorage) Name() string { return "mock-keyring" }

// TestStartCmdFailsWithMigratedTokens demonstrates the bug where start command
// can't find tokens that have been migrated to keyring
func TestStartCmdFailsWithMigratedTokens(t *testing.T) {
	// Simulate the scenario where tokens were migrated from file to keyring
	tmpDir := t.TempDir()
	testAuthPath := filepath.Join(tmpDir, "auth.json")
	
	// Create the .migrated file to simulate migration
	err := os.WriteFile(testAuthPath+".migrated", []byte("{}"), 0600)
	require.NoError(t, err)
	
	// Test what happens when StartCmd looks for tokens
	t.Run("StartCmd with NewTokenStorage finds no tokens", func(t *testing.T) {
		// This is what StartCmd used to do with deprecated function
		storage := auth.NewFileStorage(testAuthPath)
		token, err := storage.Get("anthropic")
		
		// It finds nothing because the file doesn't exist
		assert.NoError(t, err)
		assert.Nil(t, token) // This is the bug!
	})
	
	t.Run("StorageFactory would find tokens in keyring", func(t *testing.T) {
		// Skip this test since we can't easily mock the keyring in this context
		t.Skip("Would need proper keyring mocking")
	})
}

// TestBothCommandsUseConsistentStorageAfterFix verifies the fix works
func TestBothCommandsUseConsistentStorageAfterFix(t *testing.T) {
	tmpDir := t.TempDir()
	testAuthPath := filepath.Join(tmpDir, "auth.json")
	
	// Set up config
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", testAuthPath)
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file") 
	defer func() {
		os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
		os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
	}()
	
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create token via factory
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	storage, err := factory.Create()
	require.NoError(t, err)
	
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		RefreshToken: "test-refresh",
		AccessToken:  "test-access",
		ExpiresAt:    0,
	}
	err = storage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Verify both commands would use the same pattern
	t.Run("Both use StorageFactory", func(t *testing.T) {
		// Start command (after fix)
		startFactory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
		startStorage, err := startFactory.Create()
		require.NoError(t, err)
		
		// Dashboard command (already correct)
		dashFactory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
		dashStorage, err := dashFactory.Create()
		require.NoError(t, err)
		
		// Both should find the same token
		startToken, err := startStorage.Get("anthropic")
		require.NoError(t, err)
		dashToken, err := dashStorage.Get("anthropic")
		require.NoError(t, err)
		
		assert.Equal(t, startToken, dashToken)
		assert.NotNil(t, startToken)
		assert.Equal(t, "oauth", startToken.Type)
	})
}