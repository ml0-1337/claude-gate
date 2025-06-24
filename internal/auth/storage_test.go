package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenStorage(t *testing.T) {
	// Use temp directory for tests
	tempDir := t.TempDir()
	
	t.Run("creates storage with correct path", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, ".claude-gate", "auth.json"))
		assert.NotNil(t, storage)
		assert.Contains(t, storage.path, ".claude-gate")
		assert.Contains(t, storage.path, "auth.json")
	})
	
	t.Run("stores and retrieves OAuth token", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, "test1", "auth.json"))
		
		token := &TokenInfo{
			Type:         "oauth",
			RefreshToken: "test-refresh-token",
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		}
		
		err := storage.Set("anthropic", token)
		require.NoError(t, err)
		
		retrieved, err := storage.Get("anthropic")
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "oauth", retrieved.Type)
		assert.Equal(t, "test-refresh-token", retrieved.RefreshToken)
		assert.Equal(t, "test-access-token", retrieved.AccessToken)
		assert.Equal(t, token.ExpiresAt, retrieved.ExpiresAt)
	})
	
	t.Run("stores and retrieves API key", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, "test2", "auth.json"))
		
		token := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-key",
		}
		
		err := storage.Set("anthropic", token)
		require.NoError(t, err)
		
		retrieved, err := storage.Get("anthropic")
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "api", retrieved.Type)
		assert.Equal(t, "sk-test-key", retrieved.APIKey)
	})
	
	t.Run("returns nil for non-existent provider", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, "test3", "auth.json"))
		
		retrieved, err := storage.Get("non-existent")
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
	
	t.Run("removes token", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, "test4", "auth.json"))
		
		token := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-key",
		}
		
		err := storage.Set("anthropic", token)
		require.NoError(t, err)
		
		err = storage.Remove("anthropic")
		require.NoError(t, err)
		
		retrieved, err := storage.Get("anthropic")
		require.NoError(t, err)
		assert.Nil(t, retrieved)
	})
	
	t.Run("creates directory if not exists", func(t *testing.T) {
		authPath := filepath.Join(tempDir, "new-dir", "auth.json")
		storage := NewTokenStorage(authPath)
		
		token := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-key",
		}
		
		err := storage.Set("anthropic", token)
		require.NoError(t, err)
		
		// Check directory was created
		_, err = os.Stat(filepath.Dir(authPath))
		assert.NoError(t, err)
	})
	
	t.Run("sets secure file permissions", func(t *testing.T) {
		authPath := filepath.Join(tempDir, "perms", "auth.json")
		storage := NewTokenStorage(authPath)
		
		token := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-key",
		}
		
		err := storage.Set("anthropic", token)
		require.NoError(t, err)
		
		// Check file permissions (owner read/write only)
		info, err := os.Stat(authPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})
	
	t.Run("handles multiple providers", func(t *testing.T) {
		storage := NewTokenStorage(filepath.Join(tempDir, "multi", "auth.json"))
		
		token1 := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-1",
		}
		token2 := &TokenInfo{
			Type:   "api",
			APIKey: "sk-test-2",
		}
		
		err := storage.Set("anthropic", token1)
		require.NoError(t, err)
		err = storage.Set("openai", token2)
		require.NoError(t, err)
		
		retrieved1, err := storage.Get("anthropic")
		require.NoError(t, err)
		assert.Equal(t, "sk-test-1", retrieved1.APIKey)
		
		retrieved2, err := storage.Get("openai")
		require.NoError(t, err)
		assert.Equal(t, "sk-test-2", retrieved2.APIKey)
	})
}

func TestTokenInfo(t *testing.T) {
	t.Run("checks if token is expired", func(t *testing.T) {
		token := &TokenInfo{
			Type:        "oauth",
			AccessToken: "test",
			ExpiresAt:   time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		}
		assert.True(t, token.IsExpired())
		
		token.ExpiresAt = time.Now().Add(time.Hour).Unix() // Expires in 1 hour
		assert.False(t, token.IsExpired())
	})
	
	t.Run("needs refresh when close to expiry", func(t *testing.T) {
		token := &TokenInfo{
			Type:        "oauth",
			AccessToken: "test",
			ExpiresAt:   time.Now().Add(30 * time.Second).Unix(), // Expires in 30 seconds
		}
		assert.True(t, token.NeedsRefresh())
		
		token.ExpiresAt = time.Now().Add(10 * time.Minute).Unix() // Expires in 10 minutes
		assert.False(t, token.NeedsRefresh())
	})
}