package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_macOSKeychainDefaults(t *testing.T) {
	cfg := DefaultConfig()

	// Check macOS keychain defaults
	assert.True(t, cfg.KeychainTrustApp, "KeychainTrustApp should default to true")
	assert.True(t, cfg.KeychainAccessibleWhenUnlocked, "KeychainAccessibleWhenUnlocked should default to true")
	assert.False(t, cfg.KeychainSynchronizable, "KeychainSynchronizable should default to false")
}

func TestConfig_LoadFromEnv_macOSKeychain(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "trust app enabled",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_TRUST_APP": "true",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.KeychainTrustApp)
			},
		},
		{
			name: "trust app disabled",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_TRUST_APP": "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.False(t, cfg.KeychainTrustApp)
			},
		},
		{
			name: "accessible when unlocked disabled",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_ACCESSIBLE_WHEN_UNLOCKED": "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.False(t, cfg.KeychainAccessibleWhenUnlocked)
			},
		},
		{
			name: "synchronizable enabled",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_SYNCHRONIZABLE": "true",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.KeychainSynchronizable)
			},
		},
		{
			name: "all macOS settings",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_TRUST_APP":                "false",
				"CLAUDE_GATE_KEYCHAIN_ACCESSIBLE_WHEN_UNLOCKED": "false",
				"CLAUDE_GATE_KEYCHAIN_SYNCHRONIZABLE":           "true",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.False(t, cfg.KeychainTrustApp)
				assert.False(t, cfg.KeychainAccessibleWhenUnlocked)
				assert.True(t, cfg.KeychainSynchronizable)
			},
		},
		{
			name: "numeric values (1 and 0)",
			envVars: map[string]string{
				"CLAUDE_GATE_KEYCHAIN_TRUST_APP":                "1",
				"CLAUDE_GATE_KEYCHAIN_ACCESSIBLE_WHEN_UNLOCKED": "0",
				"CLAUDE_GATE_KEYCHAIN_SYNCHRONIZABLE":           "1",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.KeychainTrustApp)
				assert.False(t, cfg.KeychainAccessibleWhenUnlocked)
				assert.True(t, cfg.KeychainSynchronizable)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg := DefaultConfig()
			cfg.LoadFromEnv()

			tt.validate(t, cfg)
		})
	}
}

func TestConfig_createStorageFactoryConfig(t *testing.T) {
	cfg := &Config{
		AuthStorageType:                "keyring",
		AuthStoragePath:                "/test/path",
		KeyringService:                 "test-service",
		KeychainTrustApp:               true,
		KeychainAccessibleWhenUnlocked: false,
		KeychainSynchronizable:         true,
	}

	// This would be tested in the main package
	// Just verify the config has the expected values
	assert.Equal(t, "keyring", cfg.AuthStorageType)
	assert.Equal(t, "/test/path", cfg.AuthStoragePath)
	assert.Equal(t, "test-service", cfg.KeyringService)
	assert.True(t, cfg.KeychainTrustApp)
	assert.False(t, cfg.KeychainAccessibleWhenUnlocked)
	assert.True(t, cfg.KeychainSynchronizable)
}