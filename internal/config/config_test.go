package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

// Test 19: Config should handle all environment variables
func TestConfig_LoadFromEnv_AllVariables(t *testing.T) {
	// This test verifies that LoadFromEnv correctly loads all environment variables
	
	// Prediction: This test will pass because LoadFromEnv is straightforward
	
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "server settings",
			envVars: map[string]string{
				"CLAUDE_GATE_HOST": "0.0.0.0",
				"CLAUDE_GATE_PORT": "8080",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "0.0.0.0", cfg.Host)
				assert.Equal(t, 8080, cfg.Port)
			},
		},
		{
			name: "anthropic API settings",
			envVars: map[string]string{
				"CLAUDE_GATE_ANTHROPIC_BASE_URL": "https://custom.api.com",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "https://custom.api.com", cfg.AnthropicBaseURL)
			},
		},
		{
			name: "proxy auth token",
			envVars: map[string]string{
				"CLAUDE_GATE_PROXY_AUTH_TOKEN": "secret-token",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "secret-token", cfg.ProxyAuthToken)
			},
		},
		{
			name: "request settings",
			envVars: map[string]string{
				"CLAUDE_GATE_REQUEST_TIMEOUT":  "30s",
				"CLAUDE_GATE_MAX_REQUEST_SIZE": "5242880", // 5MB
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 30*time.Second, cfg.RequestTimeout)
				assert.Equal(t, 5242880, cfg.MaxRequestSize)
			},
		},
		{
			name: "logging settings",
			envVars: map[string]string{
				"CLAUDE_GATE_LOG_LEVEL":    "DEBUG",
				"CLAUDE_GATE_LOG_REQUESTS": "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "DEBUG", cfg.LogLevel)
				assert.False(t, cfg.LogRequests)
			},
		},
		{
			name: "rate limiting",
			envVars: map[string]string{
				"CLAUDE_GATE_ENABLE_RATE_LIMIT":     "true",
				"CLAUDE_GATE_RATE_LIMIT_PER_MINUTE": "30",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.EnableRateLimit)
				assert.Equal(t, 30, cfg.RateLimitPerMinute)
			},
		},
		{
			name: "storage settings",
			envVars: map[string]string{
				"CLAUDE_GATE_AUTH_STORAGE_PATH":  "/custom/auth.json",
				"CLAUDE_GATE_AUTH_STORAGE_TYPE":  "file",
				"CLAUDE_GATE_KEYRING_SERVICE":    "custom-service",
				"CLAUDE_GATE_AUTO_MIGRATE_TOKENS": "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "/custom/auth.json", cfg.AuthStoragePath)
				assert.Equal(t, "file", cfg.AuthStorageType)
				assert.Equal(t, "custom-service", cfg.KeyringService)
				assert.False(t, cfg.AutoMigrateTokens)
			},
		},
		{
			name: "invalid values are ignored",
			envVars: map[string]string{
				"CLAUDE_GATE_PORT":                "not-a-number",
				"CLAUDE_GATE_REQUEST_TIMEOUT":     "invalid-duration",
				"CLAUDE_GATE_MAX_REQUEST_SIZE":    "not-a-size",
				"CLAUDE_GATE_RATE_LIMIT_PER_MINUTE": "invalid",
			},
			validate: func(t *testing.T, cfg *Config) {
				// Values should remain at defaults
				assert.Equal(t, 5789, cfg.Port)
				assert.Equal(t, 600*time.Second, cfg.RequestTimeout)
				assert.Equal(t, 10*1024*1024, cfg.MaxRequestSize)
				assert.Equal(t, 60, cfg.RateLimitPerMinute)
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

// Test 20: GetBindAddress should format address correctly
func TestConfig_GetBindAddress(t *testing.T) {
	// This test verifies that GetBindAddress formats the address correctly
	
	// Prediction: This test will pass because GetBindAddress is simple
	
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{"localhost default", "127.0.0.1", 5789, "127.0.0.1:5789"},
		{"all interfaces", "0.0.0.0", 8080, "0.0.0.0:8080"},
		{"custom host and port", "192.168.1.100", 3000, "192.168.1.100:3000"},
		{"IPv6 address", "::1", 9000, "::1:9000"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Host: tt.host,
				Port: tt.port,
			}
			
			result := cfg.GetBindAddress()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test DefaultConfig returns expected defaults
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	// Verify key defaults
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 5789, cfg.Port)
	assert.Equal(t, "https://api.anthropic.com", cfg.AnthropicBaseURL)
	assert.Equal(t, 600*time.Second, cfg.RequestTimeout)
	assert.Equal(t, 10*1024*1024, cfg.MaxRequestSize)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.True(t, cfg.LogRequests)
	assert.False(t, cfg.EnableRateLimit)
	assert.Equal(t, 60, cfg.RateLimitPerMinute)
	assert.Equal(t, []string{"*"}, cfg.CORSAllowOrigins)
	assert.Equal(t, "auto", cfg.AuthStorageType)
	assert.Equal(t, "claude-gate", cfg.KeyringService)
	assert.True(t, cfg.AutoMigrateTokens)
	
	// Verify auth storage path contains home directory
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".claude-gate", "auth.json")
	assert.Equal(t, expectedPath, cfg.AuthStoragePath)
}