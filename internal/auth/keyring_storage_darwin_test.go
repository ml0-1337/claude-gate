// +build darwin

package auth

import (
	"runtime"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyringStorage_macOSConfig(t *testing.T) {
	// Only run on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific test")
	}

	tests := []struct {
		name                           string
		config                         KeyringConfig
		expectedTrustApp               bool
		expectedAccessibleWhenUnlocked bool
		expectedSynchronizable         bool
	}{
		{
			name: "custom macOS settings",
			config: KeyringConfig{
				ServiceName:                    "test-service",
				KeychainTrustApplication:       true,
				KeychainAccessibleWhenUnlocked: true,
				KeychainSynchronizable:         false,
			},
			expectedTrustApp:               true,
			expectedAccessibleWhenUnlocked: true,
			expectedSynchronizable:         false,
		},
		{
			name: "all false settings",
			config: KeyringConfig{
				ServiceName:                    "test-service",
				KeychainTrustApplication:       false,
				KeychainAccessibleWhenUnlocked: false,
				KeychainSynchronizable:         false,
			},
			expectedTrustApp:               false,
			expectedAccessibleWhenUnlocked: false,
			expectedSynchronizable:         false,
		},
		{
			name: "sync enabled",
			config: KeyringConfig{
				ServiceName:                    "test-service",
				KeychainTrustApplication:       true,
				KeychainAccessibleWhenUnlocked: true,
				KeychainSynchronizable:         true,
			},
			expectedTrustApp:               true,
			expectedAccessibleWhenUnlocked: true,
			expectedSynchronizable:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the keyring creation to verify the config
			var capturedConfig keyring.Config
			originalOpen := openKeyring
			openKeyring = func(cfg keyring.Config) (keyring.Keyring, error) {
				capturedConfig = cfg
				return &mockKeyring{}, nil
			}
			defer func() { openKeyring = originalOpen }()

			_, err := NewKeyringStorage(tt.config)
			require.NoError(t, err)

			// Verify macOS-specific settings were applied
			assert.Equal(t, tt.expectedTrustApp, capturedConfig.KeychainTrustApplication)
			assert.Equal(t, tt.expectedAccessibleWhenUnlocked, capturedConfig.KeychainAccessibleWhenUnlocked)
			assert.Equal(t, tt.expectedSynchronizable, capturedConfig.KeychainSynchronizable)
		})
	}
}

// Mock keyring for testing
type mockKeyring struct{}

func (m *mockKeyring) Get(key string) (keyring.Item, error) {
	return keyring.Item{}, keyring.ErrKeyNotFound
}

func (m *mockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (m *mockKeyring) Set(item keyring.Item) error {
	return nil
}

func (m *mockKeyring) Remove(key string) error {
	return nil
}

func (m *mockKeyring) Keys() ([]string, error) {
	return []string{}, nil
}

