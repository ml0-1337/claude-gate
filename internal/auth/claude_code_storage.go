package auth

import (
	"encoding/json"
	"fmt"

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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open Claude Code keyring: %w", err)
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

	// Get item from keyring
	item, err := s.keyring.Get("Claude Code-credentials")
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
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
	// Check if Claude Code credentials exist
	_, err := s.keyring.Get("Claude Code-credentials")
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return []string{}, nil
		}
		return nil, err
	}
	return []string{"anthropic"}, nil
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