package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockKeyring implements a mock keyring for testing
type MockKeyring struct {
	items      map[string]keyring.Item
	keys       []string
	locked     bool
	failNext   bool
	failError  error
}

func NewMockKeyring() *MockKeyring {
	return &MockKeyring{
		items: make(map[string]keyring.Item),
		keys:  []string{},
	}
}

func (m *MockKeyring) Get(key string) (keyring.Item, error) {
	if m.failNext {
		m.failNext = false
		if m.failError != nil {
			return keyring.Item{}, m.failError
		}
		return keyring.Item{}, fmt.Errorf("mock error")
	}
	
	if m.locked {
		return keyring.Item{}, fmt.Errorf("keyring is locked")
	}
	
	item, exists := m.items[key]
	if !exists {
		return keyring.Item{}, keyring.ErrKeyNotFound
	}
	
	return item, nil
}

func (m *MockKeyring) Set(item keyring.Item) error {
	if m.failNext {
		m.failNext = false
		if m.failError != nil {
			return m.failError
		}
		return fmt.Errorf("mock error")
	}
	
	if m.locked {
		return fmt.Errorf("keyring is locked")
	}
	
	m.items[item.Key] = item
	
	// Update keys list
	found := false
	for _, k := range m.keys {
		if k == item.Key {
			found = true
			break
		}
	}
	if !found {
		m.keys = append(m.keys, item.Key)
	}
	
	return nil
}

func (m *MockKeyring) Remove(key string) error {
	if m.failNext {
		m.failNext = false
		if m.failError != nil {
			return m.failError
		}
		return fmt.Errorf("mock error")
	}
	
	if m.locked {
		return fmt.Errorf("keyring is locked")
	}
	
	if _, exists := m.items[key]; !exists {
		return keyring.ErrKeyNotFound
	}
	
	delete(m.items, key)
	
	// Update keys list
	newKeys := []string{}
	for _, k := range m.keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	m.keys = newKeys
	
	return nil
}

func (m *MockKeyring) Keys() ([]string, error) {
	if m.failNext {
		m.failNext = false
		if m.failError != nil {
			return nil, m.failError
		}
		return nil, fmt.Errorf("mock error")
	}
	
	if m.locked {
		return nil, fmt.Errorf("keyring is locked")
	}
	
	return m.keys, nil
}

func (m *MockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	// Mock implementation - metadata not supported
	return keyring.Metadata{}, fmt.Errorf("metadata not supported")
}

// Test helpers
func createTestKeyringStorage(t *testing.T) (*KeyringStorage, *MockKeyring) {
	mock := NewMockKeyring()
	storage := &KeyringStorage{
		keyring: mock,
		config: KeyringConfig{
			ServiceName: "test-service",
		},
		metrics: StorageMetrics{
			Operations: make(map[string]int64),
			Errors:     make(map[string]int64),
			Latencies:  make(map[string]time.Duration),
		},
	}
	return storage, mock
}

func createTestToken() *TokenInfo {
	return &TokenInfo{
		Type:         "oauth",
		RefreshToken: "test-refresh-token",
		AccessToken:  "test-access-token",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
}

// Tests

func TestKeyringStorage_Get(t *testing.T) {
	storage, _ := createTestKeyringStorage(t)
	token := createTestToken()
	
	// Test getting non-existent token
	result, err := storage.Get("anthropic")
	assert.NoError(t, err)
	assert.Nil(t, result)
	
	// Set a token
	err = storage.Set("anthropic", token)
	require.NoError(t, err)
	
	// Get the token
	result, err = storage.Get("anthropic")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token.AccessToken, result.AccessToken)
	assert.Equal(t, token.RefreshToken, result.RefreshToken)
	
	// Check metrics
	assert.Equal(t, int64(2), storage.metrics.Operations["get"]) // Two get operations
	assert.Equal(t, int64(1), storage.metrics.Operations["set"])
}

func TestKeyringStorage_Set(t *testing.T) {
	storage, mock := createTestKeyringStorage(t)
	token := createTestToken()
	
	// Set a token
	err := storage.Set("anthropic", token)
	assert.NoError(t, err)
	
	// Verify it was stored
	assert.Len(t, mock.items, 1)
	item, exists := mock.items["test-service.anthropic"]
	assert.True(t, exists)
	assert.Equal(t, "Claude Gate - anthropic", item.Label)
	
	// Verify metrics
	assert.Equal(t, int64(1), storage.metrics.Operations["set"])
}

func TestKeyringStorage_Remove(t *testing.T) {
	storage, mock := createTestKeyringStorage(t)
	token := createTestToken()
	
	// Set a token
	err := storage.Set("anthropic", token)
	require.NoError(t, err)
	
	// Remove it
	err = storage.Remove("anthropic")
	assert.NoError(t, err)
	
	// Verify it was removed
	assert.Len(t, mock.items, 0)
	
	// Remove non-existent (should not error)
	err = storage.Remove("nonexistent")
	assert.NoError(t, err)
}

func TestKeyringStorage_List(t *testing.T) {
	storage, _ := createTestKeyringStorage(t)
	
	// List empty
	providers, err := storage.List()
	assert.NoError(t, err)
	assert.Empty(t, providers)
	
	// Add some tokens
	token := createTestToken()
	storage.Set("anthropic", token)
	storage.Set("openai", token)
	
	// List providers
	providers, err = storage.List()
	assert.NoError(t, err)
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "anthropic")
	assert.Contains(t, providers, "openai")
}

func TestKeyringStorage_LockedKeychain(t *testing.T) {
	storage, mock := createTestKeyringStorage(t)
	mock.locked = true
	
	// Test all operations with locked keychain
	token := createTestToken()
	
	err := storage.Set("anthropic", token)
	assert.ErrorIs(t, err, ErrKeyringLocked)
	
	_, err = storage.Get("anthropic")
	assert.ErrorIs(t, err, ErrKeyringLocked)
	
	err = storage.Remove("anthropic")
	assert.ErrorIs(t, err, ErrKeyringLocked)
	
	_, err = storage.List()
	assert.ErrorIs(t, err, ErrKeyringLocked)
}

func TestKeyringStorage_ErrorHandling(t *testing.T) {
	storage, mock := createTestKeyringStorage(t)
	
	// Test various error scenarios
	testCases := []struct {
		name      string
		setupMock func()
		operation func() error
		wantError error
	}{
		{
			name: "permission denied",
			setupMock: func() {
				mock.failNext = true
				mock.failError = fmt.Errorf("permission denied")
			},
			operation: func() error {
				return storage.Set("test", createTestToken())
			},
			wantError: ErrKeyringAccessDenied,
		},
		{
			name: "timeout",
			setupMock: func() {
				mock.failNext = true
				mock.failError = fmt.Errorf("operation timeout")
			},
			operation: func() error {
				_, err := storage.Get("test")
				return err
			},
			wantError: ErrKeyringTimeout,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := tc.operation()
			assert.ErrorIs(t, err, tc.wantError)
		})
	}
}

func TestKeyringStorage_Concurrency(t *testing.T) {
	storage, _ := createTestKeyringStorage(t)
	token := createTestToken()
	
	// Run concurrent operations
	done := make(chan bool)
	
	// Writers
	for i := 0; i < 5; i++ {
		go func(id int) {
			provider := fmt.Sprintf("provider%d", id)
			err := storage.Set(provider, token)
			assert.NoError(t, err)
			done <- true
		}(i)
	}
	
	// Readers
	for i := 0; i < 5; i++ {
		go func(id int) {
			provider := fmt.Sprintf("provider%d", id)
			storage.Get(provider)
			done <- true
		}(i)
	}
	
	// Wait for all operations
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all providers were stored
	providers, err := storage.List()
	assert.NoError(t, err)
	assert.Len(t, providers, 5)
}

func TestKeyringStorage_TokenExpiry(t *testing.T) {
	storage, _ := createTestKeyringStorage(t)
	
	// Create expired token
	expiredToken := &TokenInfo{
		Type:         "oauth",
		RefreshToken: "expired-refresh",
		AccessToken:  "expired-access",
		ExpiresAt:    time.Now().Add(-time.Hour).Unix(),
	}
	
	// Store expired token
	err := storage.Set("anthropic", expiredToken)
	assert.NoError(t, err)
	
	// Retrieve and check
	retrieved, err := storage.Get("anthropic")
	assert.NoError(t, err)
	assert.True(t, retrieved.IsExpired())
	assert.True(t, retrieved.NeedsRefresh())
}

func TestKeyringStorage_Metrics(t *testing.T) {
	storage, _ := createTestKeyringStorage(t)
	token := createTestToken()
	
	// Perform operations
	storage.Set("test", token)
	storage.Get("test")
	storage.List()
	storage.Remove("test")
	
	// Check metrics
	assert.Equal(t, int64(1), storage.metrics.Operations["set"])
	assert.Equal(t, int64(1), storage.metrics.Operations["get"])
	assert.Equal(t, int64(1), storage.metrics.Operations["list"])
	assert.Equal(t, int64(1), storage.metrics.Operations["remove"])
	
	// Check latencies were recorded
	assert.Contains(t, storage.metrics.Latencies, "set")
	assert.Contains(t, storage.metrics.Latencies, "get")
	assert.Contains(t, storage.metrics.Latencies, "list")
	assert.Contains(t, storage.metrics.Latencies, "remove")
}