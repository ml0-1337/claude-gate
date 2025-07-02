package auth

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Get returns transformed token when Claude Code credentials exist
func TestClaudeCodeStorage_Get_ReturnsTransformedToken(t *testing.T) {
	// Create mock keyring with Claude Code credentials
	mockKeyring := &mockKeyringForClaudeCode{
		items: map[string]keyring.Item{
			"Claude Code-credentials": {
				Key: "Claude Code-credentials",
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
	keys := make([]string, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *mockKeyringForClaudeCode) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}