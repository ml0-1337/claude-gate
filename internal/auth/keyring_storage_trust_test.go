package auth

import (
	"testing"
	"time"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyringStorage_SetWithTrustSettings(t *testing.T) {
	tests := []struct {
		name                               string
		token                              *TokenInfo
		expectedKeyPrefix                  string
		expectedNotTrustApplication        bool
		expectedNotSynchronizable          bool
	}{
		{
			name: "OAuth token with trust settings",
			token: &TokenInfo{
				Type:         "oauth",
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
				ExpiresAt:    1234567890,
			},
			expectedKeyPrefix:           "test-service.anthropic",
			expectedNotTrustApplication: false, // false = trust the app
			expectedNotSynchronizable:   true,  // true = don't sync to iCloud
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock keyring that captures the item
			var capturedItem keyring.Item
			mockKeyring := &mockKeyringWithCapture{
				captureSet: func(item keyring.Item) error {
					capturedItem = item
					return nil
				},
			}

			// Create storage with mock
			storage := &KeyringStorage{
				keyring: mockKeyring,
				config: KeyringConfig{
					ServiceName: "test-service",
				},
				metrics: StorageMetrics{
					Operations: make(map[string]int64),
					Errors:     make(map[string]int64),
					Latencies:  make(map[string]time.Duration),
				},
			}

			// Set the token
			err := storage.Set("anthropic", tt.token)
			require.NoError(t, err)

			// Verify the item was created with correct trust settings
			assert.Equal(t, tt.expectedKeyPrefix, capturedItem.Key)
			assert.Equal(t, tt.expectedNotTrustApplication, capturedItem.KeychainNotTrustApplication,
				"KeychainNotTrustApplication should be false to trust the app")
			assert.Equal(t, tt.expectedNotSynchronizable, capturedItem.KeychainNotSynchronizable,
				"KeychainNotSynchronizable should be true to prevent iCloud sync")
			assert.Equal(t, "Claude Gate - anthropic", capturedItem.Label)
			assert.Equal(t, "OAuth token for anthropic", capturedItem.Description)
		})
	}
}

// mockKeyringWithCapture allows capturing the item passed to Set
type mockKeyringWithCapture struct {
	captureSet func(item keyring.Item) error
}

func (m *mockKeyringWithCapture) Get(key string) (keyring.Item, error) {
	return keyring.Item{}, keyring.ErrKeyNotFound
}

func (m *mockKeyringWithCapture) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (m *mockKeyringWithCapture) Set(item keyring.Item) error {
	if m.captureSet != nil {
		return m.captureSet(item)
	}
	return nil
}

func (m *mockKeyringWithCapture) Remove(key string) error {
	return nil
}

func (m *mockKeyringWithCapture) Keys() ([]string, error) {
	return []string{}, nil
}