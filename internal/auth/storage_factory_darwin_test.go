// +build darwin

package auth

import (
	"runtime"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageFactory_macOSDefaults(t *testing.T) {
	// Only run on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific test")
	}

	tests := []struct {
		name   string
		config StorageFactoryConfig
		verify func(t *testing.T, factory *StorageFactory)
	}{
		{
			name: "applies macOS defaults when fields are zero values",
			config: StorageFactoryConfig{
				Type:        StorageTypeKeyring,
				ServiceName: "test-service",
			},
			verify: func(t *testing.T, factory *StorageFactory) {
				// Should have macOS defaults applied
				assert.True(t, factory.keyringConfig.KeychainTrustApplication)
				assert.True(t, factory.keyringConfig.KeychainAccessibleWhenUnlocked)
				assert.False(t, factory.keyringConfig.KeychainSynchronizable)
			},
		},
		{
			name: "respects custom macOS settings",
			config: StorageFactoryConfig{
				Type:                           StorageTypeKeyring,
				ServiceName:                    "test-service",
				KeychainTrustApp:               false,
				KeychainAccessibleWhenUnlocked: false,
				KeychainSynchronizable:         true,
			},
			verify: func(t *testing.T, factory *StorageFactory) {
				// Should use custom settings
				assert.False(t, factory.keyringConfig.KeychainTrustApplication)
				assert.False(t, factory.keyringConfig.KeychainAccessibleWhenUnlocked)
				assert.True(t, factory.keyringConfig.KeychainSynchronizable)
			},
		},
		{
			name: "mixed custom settings",
			config: StorageFactoryConfig{
				Type:                           StorageTypeKeyring,
				ServiceName:                    "test-service",
				KeychainTrustApp:               true,
				KeychainAccessibleWhenUnlocked: false,
				KeychainSynchronizable:         false,
			},
			verify: func(t *testing.T, factory *StorageFactory) {
				// Should use mixed settings
				assert.True(t, factory.keyringConfig.KeychainTrustApplication)
				assert.False(t, factory.keyringConfig.KeychainAccessibleWhenUnlocked)
				assert.False(t, factory.keyringConfig.KeychainSynchronizable)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewStorageFactory(tt.config)
			require.NotNil(t, factory)
			tt.verify(t, factory)
		})
	}
}

func TestStorageFactory_CreateWithMacOSSettings(t *testing.T) {
	// Only run on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific test")
	}

	// Create factory with custom macOS settings
	config := StorageFactoryConfig{
		Type:                           StorageTypeKeyring,
		ServiceName:                    "test-service",
		KeychainTrustApp:               true,
		KeychainAccessibleWhenUnlocked: true,
		KeychainSynchronizable:         false,
	}

	factory := NewStorageFactory(config)

	// Mock the keyring creation to verify settings are passed through
	originalOpen := openKeyring
	var capturedTrustApp bool
	openKeyring = func(cfg keyring.Config) (keyring.Keyring, error) {
		capturedTrustApp = cfg.KeychainTrustApplication
		return &mockKeyring{}, nil
	}
	defer func() { openKeyring = originalOpen }()

	storage, err := factory.Create()
	require.NoError(t, err)
	require.NotNil(t, storage)

	// Verify the trust app setting was passed through
	assert.True(t, capturedTrustApp)
}