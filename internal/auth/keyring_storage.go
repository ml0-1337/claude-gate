package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/99designs/keyring"
)

// KeyringStorage implements StorageBackend using OS-native secure storage
type KeyringStorage struct {
	keyring keyring.Keyring
	config  KeyringConfig
	mu      sync.RWMutex
	metrics StorageMetrics
}

// KeyringConfig holds configuration for keyring storage
type KeyringConfig struct {
	ServiceName     string
	AllowedBackends []keyring.BackendType
	FileDir         string
	PasswordPrompt  keyring.PromptFunc
	Debug           bool
}

// NewKeyringStorage creates a new keyring-based storage backend
func NewKeyringStorage(config KeyringConfig) (*KeyringStorage, error) {
	// Set default service name if not provided
	if config.ServiceName == "" {
		config.ServiceName = "claude-gate"
	}

	// Set default file directory for FileBackend fallback
	if config.FileDir == "" {
		homeDir, _ := os.UserHomeDir()
		config.FileDir = homeDir + "/.claude-gate/keyring"
	}

	// Configure keyring
	keyringConfig := keyring.Config{
		ServiceName:     config.ServiceName,
		AllowedBackends: config.AllowedBackends,
		FileDir:         config.FileDir,
		FilePasswordFunc: config.PasswordPrompt,
	}

	// If no backends specified, use platform defaults
	if len(config.AllowedBackends) == 0 {
		keyringConfig.AllowedBackends = getPlatformBackends()
	}

	// Open keyring
	kr, err := keyring.Open(keyringConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &KeyringStorage{
		keyring: kr,
		config:  config,
		metrics: StorageMetrics{
			Operations: make(map[string]int64),
			Errors:     make(map[string]int64),
			Latencies:  make(map[string]time.Duration),
		},
	}, nil
}

// getPlatformBackends returns the recommended backends for the current platform
func getPlatformBackends() []keyring.BackendType {
	switch runtime.GOOS {
	case "darwin":
		return []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.FileBackend,
		}
	case "linux":
		return []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KWalletBackend,
			keyring.FileBackend,
		}
	case "windows":
		return []keyring.BackendType{
			keyring.WinCredBackend,
			keyring.FileBackend,
		}
	default:
		return []keyring.BackendType{
			keyring.FileBackend,
		}
	}
}

// Get retrieves token information for a provider
func (s *KeyringStorage) Get(provider string) (*TokenInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := time.Now()
	s.recordOperation("get")
	
	item, err := s.keyring.Get(s.getKey(provider))
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, nil // Not an error, just not found
		}
		s.recordError("get", err)
		return nil, s.wrapError("get", err)
	}

	var token TokenInfo
	if err := json.Unmarshal(item.Data, &token); err != nil {
		s.recordError("get_unmarshal", err)
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	s.recordLatency("get", time.Since(start))
	return &token, nil
}

// Set stores token information for a provider
func (s *KeyringStorage) Set(provider string, token *TokenInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	s.recordOperation("set")

	data, err := json.Marshal(token)
	if err != nil {
		s.recordError("set_marshal", err)
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	item := keyring.Item{
		Key:         s.getKey(provider),
		Data:        data,
		Label:       fmt.Sprintf("Claude Gate - %s", provider),
		Description: fmt.Sprintf("OAuth token for %s", provider),
	}

	if err := s.keyring.Set(item); err != nil {
		s.recordError("set", err)
		return s.wrapError("set", err)
	}

	s.recordLatency("set", time.Since(start))
	return nil
}

// Remove deletes token information for a provider
func (s *KeyringStorage) Remove(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	s.recordOperation("remove")

	if err := s.keyring.Remove(s.getKey(provider)); err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil // Not an error if already removed
		}
		s.recordError("remove", err)
		return s.wrapError("remove", err)
	}

	s.recordLatency("remove", time.Since(start))
	return nil
}

// List returns all stored provider names
func (s *KeyringStorage) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := time.Now()
	s.recordOperation("list")

	keys, err := s.keyring.Keys()
	if err != nil {
		s.recordError("list", err)
		return nil, s.wrapError("list", err)
	}

	providers := make([]string, 0, len(keys))
	prefix := s.config.ServiceName + "."
	
	for _, key := range keys {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			providers = append(providers, key[len(prefix):])
		}
	}

	s.recordLatency("list", time.Since(start))
	return providers, nil
}

// IsAvailable checks if the backend is available on this system
func (s *KeyringStorage) IsAvailable() bool {
	// Try a simple operation to check availability
	_, err := s.keyring.Keys()
	return err == nil
}

// RequiresUnlock checks if the backend needs to be unlocked
func (s *KeyringStorage) RequiresUnlock() bool {
	// Most keyrings unlock automatically when needed
	// This is mainly for FileBackend with password
	return false
}

// Unlock attempts to unlock the backend
func (s *KeyringStorage) Unlock() error {
	// Most backends handle this automatically
	return nil
}

// Lock locks the backend
func (s *KeyringStorage) Lock() error {
	// Most backends handle this automatically
	return nil
}

// Name returns the backend name for identification
func (s *KeyringStorage) Name() string {
	return fmt.Sprintf("keyring:%s", s.config.ServiceName)
}

// getKey returns the full key name for a provider
func (s *KeyringStorage) getKey(provider string) string {
	return fmt.Sprintf("%s.%s", s.config.ServiceName, provider)
}

// wrapError wraps keyring errors with more context
func (s *KeyringStorage) wrapError(operation string, err error) error {
	switch err {
	case keyring.ErrKeyNotFound:
		return fmt.Errorf("%s: key not found", operation)
	// Note: ErrUnsupportedBackend might not be exported in all versions
	// Handle by checking error message instead
	default:
		// Check for common error patterns
		errStr := err.Error()
		switch {
		case contains(errStr, "locked", "unlock"):
			return ErrKeyringLocked
		case contains(errStr, "denied", "permission"):
			return ErrKeyringAccessDenied
		case contains(errStr, "timeout"):
			return ErrKeyringTimeout
		case contains(errStr, "unsupported", "not supported"):
			return ErrKeyringUnavailable
		default:
			return fmt.Errorf("%s failed: %w", operation, err)
		}
	}
}

// contains checks if any of the substrings are in s (case-insensitive)
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if containsIgnoreCase(s, substr) {
			return true
		}
	}
	return false
}

// containsIgnoreCase checks if substr is in s (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 containsIgnoreCase(s[1:], substr) || 
		 (len(s) > 0 && len(substr) > 0 && 
		  (s[0] == substr[0] || s[0] == substr[0]+32 || s[0] == substr[0]-32) && 
		  s[1:len(substr)] == substr[1:]))
}

// Metrics tracking methods
func (s *KeyringStorage) recordOperation(op string) {
	s.metrics.Operations[op]++
	s.metrics.LastAccess = time.Now()
}

func (s *KeyringStorage) recordError(op string, err error) {
	s.metrics.Errors[op]++
	s.metrics.LastError = err
}

func (s *KeyringStorage) recordLatency(op string, duration time.Duration) {
	s.metrics.Latencies[op] = duration
}