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

func TestClaudeCodeStorage_Get_FieldMapping(t *testing.T) {
	// Test 5: Should map all token fields correctly (accessToken→access, refreshToken→refresh)
	
	// Create mock keyring with complete Claude Code credentials
	mockKr := NewMockKeyring()
	mockKr.items["macbook"] = keyring.Item{
		Key: "macbook",
		Data: []byte(`{
			"claudeAiOauth": {
				"accessToken": "sk-ant-oat01-mapped-access",
				"refreshToken": "sk-ant-ort01-mapped-refresh",
				"expiresAt": 1751458199000,
				"scopes": ["user:inference", "user:profile"],
				"subscriptionType": "max"
			}
		}`),
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
	
	// Verify all field mappings
	assert.Equal(t, "oauth", token.Type, "Should add oauth type")
	assert.Equal(t, "sk-ant-oat01-mapped-access", token.AccessToken, "Should map accessToken to AccessToken")
	assert.Equal(t, "sk-ant-ort01-mapped-refresh", token.RefreshToken, "Should map refreshToken to RefreshToken")
	assert.Equal(t, int64(1751458199), token.ExpiresAt, "Should map expiresAt and convert to seconds")
	
	// Verify unused fields are empty
	assert.Empty(t, token.APIKey, "APIKey should be empty for OAuth tokens")
}

func TestClaudeCodeStorage_List(t *testing.T) {
	// Test 7: Should return appropriate provider name ("anthropic") in List operation
	
	t.Run("With credentials", func(t *testing.T) {
		// Create mock keyring with Claude Code credentials
		mockKr := NewMockKeyring()
		mockKr.items["macbook"] = keyring.Item{
			Key: "macbook",
			Data: []byte(`{
				"claudeAiOauth": {
					"accessToken": "test",
					"refreshToken": "test",
					"expiresAt": 1751458199000
				}
			}`),
		}
		
		// Create adapter
		adapter := &ClaudeCodeStorage{
			keyring: mockKr,
		}
		
		// List providers
		providers, err := adapter.List()
		
		// Should succeed
		require.NoError(t, err)
		assert.Equal(t, []string{"anthropic"}, providers)
	})
	
	t.Run("Without credentials", func(t *testing.T) {
		// Create empty mock keyring
		mockKr := NewMockKeyring()
		
		// Create adapter
		adapter := &ClaudeCodeStorage{
			keyring: mockKr,
		}
		
		// List providers
		providers, err := adapter.List()
		
		// Should succeed with empty list
		require.NoError(t, err)
		assert.Empty(t, providers)
	})
}

func TestClaudeCodeStorage_IsAvailable(t *testing.T) {
	// Test 8: Should report as available when keychain is accessible
	
	// Create mock keyring that reports as available
	mockKr := NewMockKeyring()
	
	// Create adapter
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	// Check availability
	available := adapter.IsAvailable()
	
	// Should be available
	assert.True(t, available)
	
	// Test with failing keyring
	mockKr.failNext = true
	mockKr.failError = fmt.Errorf("keyring locked")
	
	// Check availability again
	available = adapter.IsAvailable()
	
	// Should be unavailable
	assert.False(t, available)
}

func TestClaudeCodeStorage_Set(t *testing.T) {
	// Test 9: Should handle Set operation as no-op (read-only adapter)
	
	mockKr := NewMockKeyring()
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	token := &TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    1751458199,
	}
	
	// Set should succeed (no-op)
	err := adapter.Set("anthropic", token)
	assert.NoError(t, err)
	
	// Verify nothing was written to keyring
	assert.Empty(t, mockKr.items)
}

func TestClaudeCodeStorage_Remove(t *testing.T) {
	// Test 10: Should handle Remove operation as no-op (read-only adapter)
	
	mockKr := NewMockKeyring()
	// Add an item to verify it's not removed
	mockKr.items["test"] = keyring.Item{Key: "test", Data: []byte("data")}
	
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	// Remove should succeed (no-op)
	err := adapter.Remove("anthropic")
	assert.NoError(t, err)
	
	// Verify item still exists
	assert.Len(t, mockKr.items, 1)
}

func TestClaudeCodeStorage_Name(t *testing.T) {
	// Test 11: Should return correct adapter name for identification
	
	mockKr := NewMockKeyring()
	adapter := &ClaudeCodeStorage{
		keyring: mockKr,
	}
	
	name := adapter.Name()
	assert.Equal(t, "claude-code-adapter", name)
}

func TestClaudeCodeStorage_Get_MissingFields(t *testing.T) {
	// Test 12: Should handle missing required fields in Claude Code JSON
	
	testCases := []struct {
		name        string
		jsonData    string
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Missing accessToken",
			jsonData: `{
				"claudeAiOauth": {
					"refreshToken": "test-refresh",
					"expiresAt": 1751458199000
				}
			}`,
			shouldError: false, // Still creates token with empty access token
		},
		{
			name: "Missing refreshToken",
			jsonData: `{
				"claudeAiOauth": {
					"accessToken": "test-access",
					"expiresAt": 1751458199000
				}
			}`,
			shouldError: false, // Still creates token with empty refresh token
		},
		{
			name: "Missing claudeAiOauth object",
			jsonData: `{
				"someOtherField": {}
			}`,
			shouldError: false, // Returns token with empty fields
		},
		{
			name:        "Empty JSON object",
			jsonData:    `{}`,
			shouldError: false, // Returns token with empty fields
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockKr := NewMockKeyring()
			mockKr.items["macbook"] = keyring.Item{
				Key:  "macbook",
				Data: []byte(tc.jsonData),
			}
			
			adapter := &ClaudeCodeStorage{
				keyring: mockKr,
			}
			
			token, err := adapter.Get("anthropic")
			
			if tc.shouldError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, token)
			} else {
				require.NoError(t, err)
				require.NotNil(t, token)
				assert.Equal(t, "oauth", token.Type)
			}
		})
	}
}