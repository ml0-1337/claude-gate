package auth

import (
	"fmt"
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

func TestClaudeCodeStorage_Get_MissingCredentials(t *testing.T) {
	// Test 2: Should return nil when Claude Code credentials don't exist in keychain
	
	// Create empty mock keyring
	mockKr := NewMockKeyring()
	// No items added - simulates missing credentials
	
	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	// Get token for anthropic provider
	token, err := adapter.Get("anthropic")
	
	// Should succeed with nil token (not found)
	require.NoError(t, err)
	assert.Nil(t, token)
}

func TestClaudeCodeStorage_Get_InvalidJSON(t *testing.T) {
	// Test 3: Should handle invalid JSON format in Claude Code credentials gracefully
	
	// Create mock keyring with invalid JSON
	mockKr := NewMockKeyring()
	mockKr.items["macbook"] = keyring.Item{
		Key: "macbook",
		Data: []byte(`{invalid json content`),
	}
	
	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	// Get token for anthropic provider
	token, err := adapter.Get("anthropic")
	
	// Should return error for invalid JSON
	require.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "failed to unmarshal Claude Code credentials")
}

func TestClaudeCodeStorage_Get_TimestampConversion(t *testing.T) {
	// Test 4: Should convert expiration timestamp from milliseconds to seconds correctly
	
	testCases := []struct {
		name           string
		expiresAtMs    int64
		expectedSec    int64
	}{
		{
			name:        "Standard timestamp",
			expiresAtMs: 1751458199105,
			expectedSec: 1751458199,
		},
		{
			name:        "Zero timestamp",
			expiresAtMs: 0,
			expectedSec: 0,
		},
		{
			name:        "Edge case - rounds down",
			expiresAtMs: 1751458199999,
			expectedSec: 1751458199,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock keyring with test data
			mockKr := NewMockKeyring()
			mockKr.items["macbook"] = keyring.Item{
				Key: "macbook",
				Data: []byte(fmt.Sprintf(`{
					"claudeAiOauth": {
						"accessToken": "test-access",
						"refreshToken": "test-refresh",
						"expiresAt": %d,
						"scopes": ["user:inference"],
						"subscriptionType": "pro"
					}
				}`, tc.expiresAtMs)),
			}
			
			// Create adapter
			adapter := &ClaudeCodeStorage{
				keyring: mockKr,
			}
			
			// Get token
			token, err := adapter.Get("anthropic")
			
			// Should succeed
			require.NoError(t, err)
			require.NotNil(t, token)
			
			// Verify timestamp conversion
			assert.Equal(t, tc.expectedSec, token.ExpiresAt)
		})
	}
}