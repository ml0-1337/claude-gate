package auth

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Get returns transformed token when Claude Code credentials exist
func TestClaudeCodeStorage_Get_ReturnsTransformedToken(t *testing.T) {
	// Create mock keyring with Claude Code credentials (stored under username)
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key: "testuser",
				Data: []byte(`{
					"claudeAiOauth": {
						"accessToken": "sk-ant-oat01-test-access",
						"refreshToken": "sk-ant-ort01-test-refresh",
						"expiresAt": 1751458199105,
						"scopes": ["user:inference", "user:profile"],
						"subscriptionType": "max"
					}
				}`),
			},
		},
	}

	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	// Get token
	token, err := adapter.Get("anthropic")
	
	// Verify
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, "oauth", token.Type)
	assert.Equal(t, "sk-ant-oat01-test-access", token.AccessToken)
	assert.Equal(t, "sk-ant-ort01-test-refresh", token.RefreshToken)
	assert.Equal(t, int64(1751458199), token.ExpiresAt) // milliseconds converted to seconds
}

// mockKeyringForClaudeCode implements keyring.Keyring for testing
type mockKeyringForClaudeCode struct {
	items map[string]keyring.Item
	err   error
}

func (m *mockKeyringForClaudeCode) Get(key string) (keyring.Item, error) {
	if m.err != nil {
		return keyring.Item{}, m.err
	}
	item, ok := m.items[key]
	if !ok {
		return keyring.Item{}, keyring.ErrKeyNotFound
	}
	return item, nil
}

func (m *mockKeyringForClaudeCode) Set(item keyring.Item) error {
	return nil
}

func (m *mockKeyringForClaudeCode) Remove(key string) error {
	return nil
}

func (m *mockKeyringForClaudeCode) Keys() ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	keys := make([]string, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *mockKeyringForClaudeCode) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

// Test 2: Get returns nil when Claude Code credentials don't exist
func TestClaudeCodeStorage_Get_ReturnsNilWhenNotFound(t *testing.T) {
	// Create mock keyring with no items
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{},
	}

	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	// Get token
	token, err := adapter.Get("anthropic")
	
	// Verify
	require.NoError(t, err)
	assert.Nil(t, token)
}

// Test 3: Get handles invalid JSON in keychain gracefully
func TestClaudeCodeStorage_Get_HandlesInvalidJSON(t *testing.T) {
	// Create mock keyring with invalid JSON
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key:  "testuser",
				Data: []byte(`{invalid json`),
			},
		},
	}

	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	// Get token
	token, err := adapter.Get("anthropic")
	
	// Verify
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse Claude Code credentials")
	assert.Nil(t, token)
}

// Test 4: Get converts milliseconds to seconds for expiry time
func TestClaudeCodeStorage_Get_ConvertsMillisecondsToSeconds(t *testing.T) {
	// Create mock keyring with specific expiry time
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key: "testuser",
				Data: []byte(`{
					"claudeAiOauth": {
						"accessToken": "token",
						"refreshToken": "refresh",
						"expiresAt": 1234567890123
					}
				}`),
			},
		},
	}

	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	// Get token
	token, err := adapter.Get("anthropic")
	
	// Verify
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, int64(1234567890), token.ExpiresAt) // 1234567890123 / 1000
}

// Test 5: Set returns error (read-only adapter)
func TestClaudeCodeStorage_Set_ReturnsError(t *testing.T) {
	adapter := &ClaudeCodeStorage{
		keyring: &mockKeyringForClaudeCode{},
	}

	token := &TokenInfo{
		Type:        "oauth",
		AccessToken: "test",
	}

	err := adapter.Set("anthropic", token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

// Test 6: Remove returns error (read-only adapter)
func TestClaudeCodeStorage_Remove_ReturnsError(t *testing.T) {
	adapter := &ClaudeCodeStorage{
		keyring: &mockKeyringForClaudeCode{},
	}

	err := adapter.Remove("anthropic")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

// Test 7: List returns ["anthropic"] when credentials exist
func TestClaudeCodeStorage_List_ReturnsAnthropicWhenExists(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key:  "testuser",
				Data: []byte(`{"claudeAiOauth": {}}`),
			},
		},
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	providers, err := adapter.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"anthropic"}, providers)
}

// Test 8: List returns empty array when no credentials exist
func TestClaudeCodeStorage_List_ReturnsEmptyWhenNotExists(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{},
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	providers, err := adapter.List()
	require.NoError(t, err)
	assert.Empty(t, providers)
}

// Test 9: IsAvailable returns true when keychain is accessible
func TestClaudeCodeStorage_IsAvailable_ReturnsTrue(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	assert.True(t, adapter.IsAvailable())
}

// Test 10: IsAvailable returns false when keychain access fails
func TestClaudeCodeStorage_IsAvailable_ReturnsFalse(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		err: keyring.ErrKeyNotFound,
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	assert.False(t, adapter.IsAvailable())
}

// Test 11: Name returns descriptive identifier "claude-code-adapter"
func TestClaudeCodeStorage_Name(t *testing.T) {
	adapter := &ClaudeCodeStorage{}
	assert.Equal(t, "claude-code-adapter", adapter.Name())
}

// Test 12: Get handles missing nested claudeAiOauth object
func TestClaudeCodeStorage_Get_HandlesMissingNestedObject(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key:  "testuser",
				Data: []byte(`{}`), // Missing claudeAiOauth
			},
		},
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	token, err := adapter.Get("anthropic")
	
	// Should not error, but return empty token
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Empty(t, token.AccessToken)
	assert.Empty(t, token.RefreshToken)
}

// Test 13: Get preserves all token fields during transformation
func TestClaudeCodeStorage_Get_PreservesAllFields(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"testuser": {
				Key: "testuser",
				Data: []byte(`{
					"claudeAiOauth": {
						"accessToken": "sk-ant-oat01-full-test",
						"refreshToken": "sk-ant-ort01-full-test",
						"expiresAt": 1700000000000,
						"scopes": ["user:inference", "user:profile"],
						"subscriptionType": "pro"
					}
				}`),
			},
		},
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	token, err := adapter.Get("anthropic")
	
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, "oauth", token.Type)
	assert.Equal(t, "sk-ant-oat01-full-test", token.AccessToken)
	assert.Equal(t, "sk-ant-ort01-full-test", token.RefreshToken)
	assert.Equal(t, int64(1700000000), token.ExpiresAt)
}

// Test 14: Get handles keychain read errors appropriately
func TestClaudeCodeStorage_Get_HandlesKeychainError(t *testing.T) {
	mockKeyring := &mockKeyringForClaudeCode{
		err: keyring.ErrKeyNotFound,
	}

	adapter := &ClaudeCodeStorage{
		keyring: mockKeyring,
	}

	// Test with provider that doesn't exist - should return nil, no error
	token, err := adapter.Get("anthropic")
	require.NoError(t, err)
	assert.Nil(t, token)
}