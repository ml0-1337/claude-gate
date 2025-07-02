package auth

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeCodeStorage_Get_ValidCredentials(t *testing.T) {
	// Test 1: Should retrieve valid Claude Code credentials and return in claude-gate format
	
	// Create mock keyring with Claude Code credentials
	mockKr := NewMockKeyring()
	mockKr.items["macbook"] = keyring.Item{
		Key: "macbook",
		Data: []byte(`{
			"claudeAiOauth": {
				"accessToken": "sk-ant-oat01-test-access",
				"refreshToken": "sk-ant-ort01-test-refresh",
				"expiresAt": 1751458199105,
				"scopes": ["user:inference", "user:profile"],
				"subscriptionType": "max"
			}
		}`),
	}
	
	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	// Get token for anthropic provider
	token, err := adapter.Get("anthropic")
	
	// Should succeed
	require.NoError(t, err)
	require.NotNil(t, token)
	
	// Verify transformed fields
	assert.Equal(t, "oauth", token.Type)
	assert.Equal(t, "sk-ant-oat01-test-access", token.AccessToken)
	assert.Equal(t, "sk-ant-ort01-test-refresh", token.RefreshToken)
	assert.Equal(t, int64(1751458199), token.ExpiresAt) // Converted from milliseconds to seconds
}