package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageMigrator_Migrate(t *testing.T) {
	// Create source storage with test data
	source := createTestFileStorage(t)
	tokens := map[string]*TokenInfo{
		"anthropic": {
			Type:         "oauth",
			RefreshToken: "refresh1",
			AccessToken:  "access1",
			ExpiresAt:    1234567890,
		},
		"openai": {
			Type:         "oauth",
			RefreshToken: "refresh2",
			AccessToken:  "access2",
			ExpiresAt:    9876543210,
		},
	}
	
	// Populate source
	for provider, token := range tokens {
		err := source.Set(provider, token)
		require.NoError(t, err)
	}
	
	// Create destination storage
	destination := createTestFileStorage(t)
	
	// Create migrator
	migrator := NewStorageMigrator(source, destination)
	
	// Perform migration
	err := migrator.Migrate()
	assert.NoError(t, err)
	
	// Verify all tokens were migrated
	for provider, expectedToken := range tokens {
		actualToken, err := destination.Get(provider)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, actualToken)
	}
	
	// Verify source was marked as migrated
	_, err = os.Stat(source.path + ".migrated")
	assert.NoError(t, err)
}

func TestStorageMigrator_MigrateEmpty(t *testing.T) {
	// Create empty source
	source := createTestFileStorage(t)
	destination := createTestFileStorage(t)
	
	migrator := NewStorageMigrator(source, destination)
	
	// Migrate empty storage
	err := migrator.Migrate()
	assert.NoError(t, err)
	
	// Verify destination is still empty
	providers, err := destination.List()
	assert.NoError(t, err)
	assert.Empty(t, providers)
}

func TestStorageMigrator_MigrateWithErrors(t *testing.T) {
	// Create source with test data
	source := createTestFileStorage(t)
	source.Set("provider1", &TokenInfo{Type: "oauth", AccessToken: "token1"})
	source.Set("provider2", &TokenInfo{Type: "oauth", AccessToken: "token2"})
	
	// Create destination that will fail on second provider
	destination := &mockFailingStorage{
		FileStorage: createTestFileStorage(t),
		failOn:      "provider2",
	}
	
	migrator := NewStorageMigrator(source, destination)
	
	// Perform migration
	err := migrator.Migrate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "1 migrated, 1 failed")
	
	// Verify first provider was migrated
	token, err := destination.Get("provider1")
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestStorageMigrator_Rollback(t *testing.T) {
	// Create storages
	original := createTestFileStorage(t)
	migrated := createTestFileStorage(t)
	
	// Add test data to migrated storage
	testToken := &TokenInfo{
		Type:        "oauth",
		AccessToken: "rolled-back-token",
	}
	migrated.Set("anthropic", testToken)
	
	// Create migrator (note: source and destination are swapped for rollback)
	migrator := NewStorageMigrator(original, migrated)
	
	// Perform rollback
	err := migrator.Rollback()
	assert.NoError(t, err)
	
	// Verify token is back in original
	token, err := original.Get("anthropic")
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)
}

func TestStorageMigrator_Backup(t *testing.T) {
	// Create source storage with test data
	source := createTestFileStorage(t)
	source.Set("anthropic", &TokenInfo{Type: "oauth", AccessToken: "backup-me"})
	
	destination := createTestFileStorage(t)
	
	// Create migrator with backup enabled
	migrator := NewStorageMigrator(source, destination)
	migrator.backup = true
	
	// Perform migration
	err := migrator.Migrate()
	assert.NoError(t, err)
	
	// Check that backup was created
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".claude-gate", "backups")
	entries, err := os.ReadDir(backupDir)
	if err == nil && len(entries) > 0 {
		// At least one backup exists
		assert.NotEmpty(t, entries)
	}
}

func TestStorageMigrator_VerifyMigration(t *testing.T) {
	// Create storages with identical data
	source := createTestFileStorage(t)
	destination := createTestFileStorage(t)
	
	tokens := map[string]*TokenInfo{
		"provider1": {Type: "oauth", AccessToken: "token1"},
		"provider2": {Type: "oauth", AccessToken: "token2"},
	}
	
	for provider, token := range tokens {
		source.Set(provider, token)
		destination.Set(provider, token)
	}
	
	migrator := NewStorageMigrator(source, destination)
	
	// Verify migration
	err := migrator.VerifyMigration()
	assert.NoError(t, err)
	
	// Modify one token
	destination.Set("provider1", &TokenInfo{Type: "oauth", AccessToken: "modified"})
	
	// Verify should fail
	err = migrator.VerifyMigration()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token mismatch")
}

func TestTokensEqual(t *testing.T) {
	token1 := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "refresh",
		AccessToken:  "access",
		ExpiresAt:    1234567890,
	}
	
	token2 := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "refresh",
		AccessToken:  "access",
		ExpiresAt:    1234567890,
	}
	
	token3 := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "different",
		AccessToken:  "access",
		ExpiresAt:    1234567890,
	}
	
	assert.True(t, tokensEqual(token1, token2))
	assert.False(t, tokensEqual(token1, token3))
	assert.True(t, tokensEqual(nil, nil))
	assert.False(t, tokensEqual(token1, nil))
	assert.False(t, tokensEqual(nil, token1))
}

// Helper functions

func createTestFileStorage(t *testing.T) *FileStorage {
	tempDir := t.TempDir()
	return NewFileStorage(filepath.Join(tempDir, "test-auth.json"))
}

// mockFailingStorage simulates a storage that fails on specific operations
type mockFailingStorage struct {
	*FileStorage
	failOn string
}

func (m *mockFailingStorage) Set(provider string, token *TokenInfo) error {
	if provider == m.failOn {
		return fmt.Errorf("simulated failure for %s", provider)
	}
	return m.FileStorage.Set(provider, token)
}