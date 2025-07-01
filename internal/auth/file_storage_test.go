package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: NewFileStorage creates storage with correct path
func TestNewFileStorage_Initialization(t *testing.T) {
	// Prediction: This test will pass - testing initialization
	
	path := "/tmp/test-tokens.json"
	storage := NewFileStorage(path)
	
	assert.NotNil(t, storage)
	assert.Equal(t, path, storage.path)
	assert.NotNil(t, storage.metrics.Operations)
	assert.NotNil(t, storage.metrics.Errors)
	assert.NotNil(t, storage.metrics.Latencies)
}

// Test 2: Save and Get token roundtrip works correctly
func TestFileStorage_SaveAndGet(t *testing.T) {
	// Prediction: This test will pass - testing basic save/get functionality
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "tokens.json"))
	
	token := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "test-refresh",
		AccessToken:  "test-access",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	
	// Save token
	err := storage.Set("claude", token)
	require.NoError(t, err)
	
	// Get token back
	retrieved, err := storage.Get("claude")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	
	assert.Equal(t, token.Type, retrieved.Type)
	assert.Equal(t, token.RefreshToken, retrieved.RefreshToken)
	assert.Equal(t, token.AccessToken, retrieved.AccessToken)
	assert.Equal(t, token.ExpiresAt, retrieved.ExpiresAt)
	
	// Verify metrics were recorded
	assert.Greater(t, storage.metrics.Operations["set"], int64(0))
	assert.Greater(t, storage.metrics.Operations["get"], int64(0))
}

// Test 3: Save handles directory creation
func TestFileStorage_SaveCreatesDirectory(t *testing.T) {
	// Prediction: This test will pass - testing directory creation
	
	tmpDir := t.TempDir()
	deepPath := filepath.Join(tmpDir, "deep", "nested", "dir", "tokens.json")
	storage := NewFileStorage(deepPath)
	
	token := &TokenInfo{
		Type:   "api",
		APIKey: "test-key",
	}
	
	err := storage.Set("test", token)
	require.NoError(t, err)
	
	// Verify directory was created
	assert.DirExists(t, filepath.Dir(deepPath))
	
	// Verify file was created with correct permissions
	info, err := os.Stat(deepPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

// Test 4: Get handles file not found error
func TestFileStorage_GetFileNotFound(t *testing.T) {
	// Prediction: This test will pass - testing missing file handling
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "nonexistent.json"))
	
	// Get from non-existent file should return nil, nil
	token, err := storage.Get("claude")
	assert.NoError(t, err)
	assert.Nil(t, token)
}

// Test 5: Get handles corrupted file
func TestFileStorage_GetCorruptedFile(t *testing.T) {
	// Prediction: This test will pass - testing error handling for corrupted data
	
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "corrupted.json")
	
	// Write corrupted JSON
	err := os.WriteFile(path, []byte("{invalid json"), 0600)
	require.NoError(t, err)
	
	storage := NewFileStorage(path)
	token, err := storage.Get("claude")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
	assert.Nil(t, token)
	
	// Verify error was recorded
	assert.Greater(t, storage.metrics.Errors["get"], int64(0))
}

// Test 6: Remove deletes token successfully
func TestFileStorage_Remove(t *testing.T) {
	// Prediction: This test will pass - testing token removal
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "tokens.json"))
	
	// Add two tokens
	token1 := &TokenInfo{Type: "oauth", AccessToken: "token1"}
	token2 := &TokenInfo{Type: "oauth", AccessToken: "token2"}
	
	err := storage.Set("provider1", token1)
	require.NoError(t, err)
	err = storage.Set("provider2", token2)
	require.NoError(t, err)
	
	// Remove one token
	err = storage.Remove("provider1")
	require.NoError(t, err)
	
	// Verify it's gone
	retrieved, err := storage.Get("provider1")
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
	
	// Verify other token still exists
	retrieved, err = storage.Get("provider2")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "token2", retrieved.AccessToken)
}

// Test 7: Remove handles missing file gracefully
func TestFileStorage_RemoveMissingFile(t *testing.T) {
	// Prediction: This test will pass - testing remove on non-existent file
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "missing.json"))
	
	// Remove from non-existent file should succeed
	err := storage.Remove("provider")
	assert.NoError(t, err)
}

// Test 8: List returns all providers
func TestFileStorage_List(t *testing.T) {
	// Prediction: This test will pass - testing list functionality
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "tokens.json"))
	
	// Add multiple tokens
	providers := []string{"claude", "openai", "gemini"}
	for _, provider := range providers {
		token := &TokenInfo{Type: "oauth", AccessToken: provider}
		err := storage.Set(provider, token)
		require.NoError(t, err)
	}
	
	// List all providers
	listed, err := storage.List()
	require.NoError(t, err)
	
	assert.Len(t, listed, 3)
	for _, provider := range providers {
		assert.Contains(t, listed, provider)
	}
}

// Test 9: Clear removes all token files
func TestFileStorage_RemoveLastToken(t *testing.T) {
	// Prediction: This test will pass - testing file removal when last token deleted
	
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "tokens.json")
	storage := NewFileStorage(path)
	
	// Add one token
	token := &TokenInfo{Type: "oauth", AccessToken: "only-token"}
	err := storage.Set("provider", token)
	require.NoError(t, err)
	
	// Verify file exists
	assert.FileExists(t, path)
	
	// Remove the only token
	err = storage.Remove("provider")
	require.NoError(t, err)
	
	// Verify file was removed
	assert.NoFileExists(t, path)
}

// Test 10: IsAvailable always returns true
func TestFileStorage_IsAvailable(t *testing.T) {
	// Prediction: This test will pass - file storage is always available
	
	storage := NewFileStorage("/tmp/test.json")
	assert.True(t, storage.IsAvailable())
}

// Test additional methods
func TestFileStorage_AdditionalMethods(t *testing.T) {
	// Prediction: This test will pass - testing utility methods
	
	storage := NewFileStorage("/tmp/test.json")
	
	t.Run("RequiresUnlock returns false", func(t *testing.T) {
		assert.False(t, storage.RequiresUnlock())
	})
	
	t.Run("Unlock returns nil", func(t *testing.T) {
		assert.NoError(t, storage.Unlock())
	})
	
	t.Run("Lock returns nil", func(t *testing.T) {
		assert.NoError(t, storage.Lock())
	})
	
	t.Run("Name returns correct format", func(t *testing.T) {
		assert.Equal(t, "file:/tmp/test.json", storage.Name())
	})
}

// Test TokenInfo methods
func TestTokenInfo_Methods(t *testing.T) {
	// Prediction: This test will pass - testing TokenInfo helper methods
	
	t.Run("IsExpired for oauth token", func(t *testing.T) {
		// Expired token
		token := &TokenInfo{
			Type:      "oauth",
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		}
		assert.True(t, token.IsExpired())
		
		// Valid token
		token.ExpiresAt = time.Now().Add(time.Hour).Unix()
		assert.False(t, token.IsExpired())
		
		// No expiry
		token.ExpiresAt = 0
		assert.False(t, token.IsExpired())
	})
	
	t.Run("IsExpired for API key", func(t *testing.T) {
		token := &TokenInfo{
			Type:   "api",
			APIKey: "test-key",
		}
		assert.False(t, token.IsExpired())
	})
	
	t.Run("NeedsRefresh for oauth token", func(t *testing.T) {
		// Needs refresh (less than 5 minutes)
		token := &TokenInfo{
			Type:      "oauth",
			ExpiresAt: time.Now().Add(3 * time.Minute).Unix(),
		}
		assert.True(t, token.NeedsRefresh())
		
		// Doesn't need refresh
		token.ExpiresAt = time.Now().Add(10 * time.Minute).Unix()
		assert.False(t, token.NeedsRefresh())
		
		// Already expired
		token.ExpiresAt = time.Now().Add(-time.Hour).Unix()
		assert.True(t, token.NeedsRefresh())
	})
	
	t.Run("NeedsRefresh for API key", func(t *testing.T) {
		token := &TokenInfo{
			Type:   "api",
			APIKey: "test-key",
		}
		assert.False(t, token.NeedsRefresh())
	})
}

// Test concurrent access
func TestFileStorage_ConcurrentAccess(t *testing.T) {
	// Prediction: This test will pass - testing thread safety
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "concurrent.json"))
	
	// Run concurrent operations
	done := make(chan bool, 3)
	
	// Writer 1
	go func() {
		for i := 0; i < 10; i++ {
			token := &TokenInfo{Type: "oauth", AccessToken: fmt.Sprintf("token1-%d", i)}
			storage.Set("provider1", token)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()
	
	// Writer 2
	go func() {
		for i := 0; i < 10; i++ {
			token := &TokenInfo{Type: "oauth", AccessToken: fmt.Sprintf("token2-%d", i)}
			storage.Set("provider2", token)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()
	
	// Reader
	go func() {
		for i := 0; i < 20; i++ {
			storage.Get("provider1")
			storage.List()
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Verify both tokens exist
	token1, err := storage.Get("provider1")
	assert.NoError(t, err)
	assert.NotNil(t, token1)
	
	token2, err := storage.Get("provider2")
	assert.NoError(t, err)
	assert.NotNil(t, token2)
}

// Test edge cases
func TestFileStorage_EdgeCases(t *testing.T) {
	// Prediction: This test will pass - testing edge cases
	
	tmpDir := t.TempDir()
	
	t.Run("Empty provider name", func(t *testing.T) {
		storage := NewFileStorage(filepath.Join(tmpDir, "edge1.json"))
		token := &TokenInfo{Type: "oauth", AccessToken: "test"}
		
		// Can save with empty provider
		err := storage.Set("", token)
		assert.NoError(t, err)
		
		// Can retrieve with empty provider
		retrieved, err := storage.Get("")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
	})
	
	t.Run("Special characters in provider name", func(t *testing.T) {
		storage := NewFileStorage(filepath.Join(tmpDir, "edge2.json"))
		token := &TokenInfo{Type: "oauth", AccessToken: "test"}
		
		specialProvider := "test@provider/with:special*chars"
		err := storage.Set(specialProvider, token)
		assert.NoError(t, err)
		
		retrieved, err := storage.Get(specialProvider)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
	})
	
	t.Run("Very large token data", func(t *testing.T) {
		storage := NewFileStorage(filepath.Join(tmpDir, "edge3.json"))
		
		// Create large token
		largeToken := &TokenInfo{
			Type:         "oauth",
			RefreshToken: string(make([]byte, 10000)), // 10KB token
			AccessToken:  string(make([]byte, 10000)),
		}
		
		err := storage.Set("large", largeToken)
		assert.NoError(t, err)
		
		retrieved, err := storage.Get("large")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Len(t, retrieved.RefreshToken, 10000)
	})
}

// Test metrics recording
func TestFileStorage_Metrics(t *testing.T) {
	// Prediction: This test will pass - testing metrics functionality
	
	tmpDir := t.TempDir()
	storage := NewFileStorage(filepath.Join(tmpDir, "metrics.json"))
	
	// Perform operations
	token := &TokenInfo{Type: "oauth", AccessToken: "test"}
	storage.Set("test", token)
	storage.Get("test")
	storage.List()
	storage.Remove("test")
	
	// Check operation counts
	assert.Equal(t, int64(1), storage.metrics.Operations["set"])
	assert.Equal(t, int64(1), storage.metrics.Operations["get"])
	assert.Equal(t, int64(1), storage.metrics.Operations["list"])
	assert.Equal(t, int64(1), storage.metrics.Operations["remove"])
	
	// Check latencies were recorded
	assert.Greater(t, storage.metrics.Latencies["set"], time.Duration(0))
	assert.Greater(t, storage.metrics.Latencies["get"], time.Duration(0))
	
	// Force an error
	storage.path = "/invalid/path/that/cannot/exist/tokens.json"
	err := storage.Set("fail", token)
	assert.Error(t, err)
	assert.Greater(t, storage.metrics.Errors["set_mkdir"], int64(0))
	assert.NotNil(t, storage.metrics.LastError)
}

// Test malformed token data handling
func TestFileStorage_MalformedTokenData(t *testing.T) {
	// Prediction: This test will fail initially - testing handling of malformed data
	// The unmarshal will fail because of type mismatch
	
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "malformed.json")
	
	// Write valid JSON but with wrong structure
	data := map[string]interface{}{
		"provider": map[string]interface{}{
			"type": 123, // Should be string
			"access": true, // Should be string
		},
	}
	
	jsonData, _ := json.Marshal(data)
	err := os.WriteFile(path, jsonData, 0600)
	require.NoError(t, err)
	
	storage := NewFileStorage(path)
	token, err := storage.Get("provider")
	
	// The unmarshal will fail due to type mismatch
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
	assert.Nil(t, token)
	
	// Test with valid but minimal data
	data2 := map[string]interface{}{
		"provider2": map[string]interface{}{
			"type": "oauth",
		},
	}
	
	jsonData2, _ := json.Marshal(data2)
	err = os.WriteFile(path, jsonData2, 0600)
	require.NoError(t, err)
	
	token2, err := storage.Get("provider2")
	assert.NoError(t, err)
	assert.NotNil(t, token2)
	assert.Equal(t, "oauth", token2.Type)
}