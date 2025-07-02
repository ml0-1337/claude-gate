package auth

import (
	"encoding/json"
	"fmt"
	"os/user"

	"github.com/99designs/keyring"
)

// ClaudeCodeStorage is an adapter that reads credentials from Claude Code's keychain storage
type ClaudeCodeStorage struct {
	keyring keyring.Keyring
}

// NewClaudeCodeStorage creates a new Claude Code storage adapter
func NewClaudeCodeStorage() (*ClaudeCodeStorage, error) {
	// Configure keyring to access Claude Code's service
	config := keyring.Config{
		ServiceName: "Claude Code-credentials",
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.FileBackend,
		},
	}
	
	kr, err := keyring.Open(config)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}
	
	return &ClaudeCodeStorage{
		keyring: kr,
	}, nil
}

// claudeCodeCredentials represents the structure of Claude Code's stored credentials
type claudeCodeCredentials struct {
	ClaudeAiOauth struct {
		AccessToken      string   `json:"accessToken"`
		RefreshToken     string   `json:"refreshToken"`
		ExpiresAt        int64    `json:"expiresAt"`
		Scopes           []string `json:"scopes"`
		SubscriptionType string   `json:"subscriptionType"`
	} `json:"claudeAiOauth"`
}

// Get retrieves token information for a provider (always returns anthropic tokens from Claude Code)
func (s *ClaudeCodeStorage) Get(provider string) (*TokenInfo, error) {
	// Claude Code only stores anthropic credentials
	if provider != "anthropic" {
		return nil, nil
	}
	
	// Get current username for account name
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	
	// Retrieve item from keyring
	item, err := s.keyring.Get(currentUser.Username)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, nil // Not an error, just not found
		}
		return nil, fmt.Errorf("failed to get credentials from keyring: %w", err)
	}
	
	// Parse Claude Code's JSON structure
	var creds claudeCodeCredentials
	if err := json.Unmarshal(item.Data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Claude Code credentials: %w", err)
	}
	
	// Transform to claude-gate's TokenInfo format
	token := &TokenInfo{
		Type:         "oauth",
		AccessToken:  creds.ClaudeAiOauth.AccessToken,
		RefreshToken: creds.ClaudeAiOauth.RefreshToken,
		ExpiresAt:    creds.ClaudeAiOauth.ExpiresAt / 1000, // Convert milliseconds to seconds
	}
	
	return token, nil
}

// Set is a no-op for the read-only adapter
func (s *ClaudeCodeStorage) Set(provider string, token *TokenInfo) error {
	// Read-only adapter - do nothing
	return nil
}

// Remove is a no-op for the read-only adapter
func (s *ClaudeCodeStorage) Remove(provider string) error {
	// Read-only adapter - do nothing
	return nil
}

// List returns the providers available (always just anthropic for Claude Code)
func (s *ClaudeCodeStorage) List() ([]string, error) {
	// Check if we can access the credentials
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	
	_, err = s.keyring.Get(currentUser.Username)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return []string{}, nil // No providers
		}
		return nil, fmt.Errorf("failed to check credentials: %w", err)
	}
	
	return []string{"anthropic"}, nil
}

// IsAvailable checks if the backend is available on this system
func (s *ClaudeCodeStorage) IsAvailable() bool {
	// Try to list keys to check availability
	_, err := s.keyring.Keys()
	return err == nil
}

// RequiresUnlock checks if the backend needs to be unlocked
func (s *ClaudeCodeStorage) RequiresUnlock() bool {
	// Keychain typically unlocks automatically
	return false
}

// Unlock attempts to unlock the backend
func (s *ClaudeCodeStorage) Unlock() error {
	// Most backends handle this automatically
	return nil
}

// Lock locks the backend
func (s *ClaudeCodeStorage) Lock() error {
	// Most backends handle this automatically
	return nil
}

// Name returns the backend name for identification
func (s *ClaudeCodeStorage) Name() string {
	return "claude-code-adapter"
}