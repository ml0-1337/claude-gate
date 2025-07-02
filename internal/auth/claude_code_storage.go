package auth

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/99designs/keyring"
)

// ClaudeCodeCredentials represents the structure of Claude Code's keychain data
type ClaudeCodeCredentials struct {
	ClaudeAiOauth struct {
		AccessToken      string   `json:"accessToken"`
		RefreshToken     string   `json:"refreshToken"`
		ExpiresAt        int64    `json:"expiresAt"` // milliseconds
		Scopes           []string `json:"scopes"`
		SubscriptionType string   `json:"subscriptionType"`
	} `json:"claudeAiOauth"`
}

// ClaudeCodeStorage implements StorageBackend by reading from Claude Code's keychain
type ClaudeCodeStorage struct {
	keyring keyring.Keyring
}

// NewClaudeCodeStorage creates a new Claude Code storage adapter
func NewClaudeCodeStorage() (*ClaudeCodeStorage, error) {
	// Open keyring with Claude Code's service name
	kr, err := keyring.Open(keyring.Config{
		ServiceName: "Claude Code-credentials",
		// Let it auto-detect the best backend
	})
	if err != nil {
		// On macOS, if keyring fails, we can fall back to the direct implementation
		if runtime.GOOS == "darwin" {
			// Return a special error that the factory can handle
			return nil, fmt.Errorf("keyring unavailable on macOS: %w", err)
		}
		return nil, fmt.Errorf("failed to open Claude Code keyring: %w", err)
	}

	// On macOS, test if the keyring actually works
	if runtime.GOOS == "darwin" {
		// Try to list keys to see if it's actually functional
		_, err := kr.Keys()
		if err != nil {
			// Keyring opened but doesn't work properly
			return nil, fmt.Errorf("keyring unavailable on macOS: %w", err)
		}
	}

	return &ClaudeCodeStorage{
		keyring: kr,
	}, nil
}

// NewClaudeCodeStorageWithKeyring creates a new Claude Code storage adapter with a custom keyring
// This is primarily for testing purposes
func NewClaudeCodeStorageWithKeyring(kr keyring.Keyring) *ClaudeCodeStorage {
	return &ClaudeCodeStorage{
		keyring: kr,
	}
}

// Get retrieves and transforms Claude Code credentials
func (s *ClaudeCodeStorage) Get(provider string) (*TokenInfo, error) {
	// Claude Code only supports anthropic
	if provider != "anthropic" {
		return nil, nil
	}

	// Claude Code stores credentials under the username, not a fixed key
	// First try to list all keys to find the right one
	keys, err := s.keyring.Keys()
	if err != nil {
		// If we can't list keys, assume no credentials
		return nil, nil
	}

	// If no keys, no credentials
	if len(keys) == 0 {
		return nil, nil
	}

	// Try each key (should typically only be one - the username)
	var item keyring.Item
	var foundItem bool
	for _, key := range keys {
		item, err = s.keyring.Get(key)
		if err == nil {
			foundItem = true
			break
		}
	}

	if !foundItem {
		// Couldn't get any items
		return nil, nil
	}

	// Parse Claude Code format
	var creds ClaudeCodeCredentials
	if err := json.Unmarshal(item.Data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse Claude Code credentials: %w", err)
	}

	// Transform to our format
	token := &TokenInfo{
		Type:         "oauth",
		AccessToken:  creds.ClaudeAiOauth.AccessToken,
		RefreshToken: creds.ClaudeAiOauth.RefreshToken,
		ExpiresAt:    creds.ClaudeAiOauth.ExpiresAt / 1000, // Convert milliseconds to seconds
	}

	return token, nil
}

// Set is not supported for read-only adapter
func (s *ClaudeCodeStorage) Set(provider string, token *TokenInfo) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// Remove is not supported for read-only adapter
func (s *ClaudeCodeStorage) Remove(provider string) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// List returns available providers
func (s *ClaudeCodeStorage) List() ([]string, error) {
	// Check if any Claude Code credentials exist
	keys, err := s.keyring.Keys()
	if err != nil {
		return nil, err
	}
	
	// If there are any keys, we have credentials
	if len(keys) > 0 {
		return []string{"anthropic"}, nil
	}
	
	return []string{}, nil
}

// IsAvailable checks if the backend is available
func (s *ClaudeCodeStorage) IsAvailable() bool {
	_, err := s.keyring.Keys()
	return err == nil
}

// RequiresUnlock checks if the backend needs to be unlocked
func (s *ClaudeCodeStorage) RequiresUnlock() bool {
	return false
}

// Unlock is a no-op for Claude Code storage
func (s *ClaudeCodeStorage) Unlock() error {
	return nil
}

// Lock is a no-op for Claude Code storage
func (s *ClaudeCodeStorage) Lock() error {
	return nil
}

// Name returns the backend name
func (s *ClaudeCodeStorage) Name() string {
	return "claude-code-adapter"
}