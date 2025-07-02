package auth

import (
	"encoding/json"
	"fmt"

	"github.com/99designs/keyring"
)

// ClaudeCodeStorageV2 implements StorageBackend by trying multiple possible Claude Code keyring names
type ClaudeCodeStorageV2 struct {
	keyrings []keyring.Keyring
}

// NewClaudeCodeStorageV2 creates a new Claude Code storage adapter that tries multiple service names
func NewClaudeCodeStorageV2() (*ClaudeCodeStorageV2, error) {
	// Try multiple possible service names that Claude Code might use
	serviceNames := []string{
		"Claude Code-credentials",  // Original expected name
		"claude.ai",               // Possible web-based name
		"Claude Code",             // Without -credentials suffix
		"com.anthropic.claude",    // Reverse domain notation
		"com.anthropic.claude-code", // Full reverse domain
	}
	
	var keyrings []keyring.Keyring
	var lastErr error
	
	for _, serviceName := range serviceNames {
		kr, err := keyring.Open(keyring.Config{
			ServiceName: serviceName,
		})
		if err != nil {
			lastErr = err
			continue
		}
		keyrings = append(keyrings, kr)
	}
	
	if len(keyrings) == 0 {
		return nil, fmt.Errorf("failed to open any Claude Code keyring: %w", lastErr)
	}
	
	return &ClaudeCodeStorageV2{
		keyrings: keyrings,
	}, nil
}

// Get retrieves and transforms Claude Code credentials from any available keyring
func (s *ClaudeCodeStorageV2) Get(provider string) (*TokenInfo, error) {
	// Claude Code only supports anthropic
	if provider != "anthropic" {
		return nil, nil
	}

	// Try each keyring
	for _, kr := range s.keyrings {
		// List all keys in this keyring
		keys, err := kr.Keys()
		if err != nil {
			continue // Try next keyring
		}

		// Try each key
		for _, key := range keys {
			item, err := kr.Get(key)
			if err != nil {
				continue
			}

			// Try to parse as Claude Code format
			var creds ClaudeCodeCredentials
			if err := json.Unmarshal(item.Data, &creds); err != nil {
				continue
			}

			// Check if it has the expected structure
			if creds.ClaudeAiOauth.AccessToken == "" || creds.ClaudeAiOauth.RefreshToken == "" {
				continue
			}

			// Found valid credentials!
			token := &TokenInfo{
				Type:         "oauth",
				AccessToken:  creds.ClaudeAiOauth.AccessToken,
				RefreshToken: creds.ClaudeAiOauth.RefreshToken,
				ExpiresAt:    creds.ClaudeAiOauth.ExpiresAt / 1000, // Convert milliseconds to seconds
			}

			return token, nil
		}
	}

	// No valid credentials found in any keyring
	return nil, nil
}

// Set is not supported for read-only adapter
func (s *ClaudeCodeStorageV2) Set(provider string, token *TokenInfo) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// Remove is not supported for read-only adapter
func (s *ClaudeCodeStorageV2) Remove(provider string) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// List returns available providers
func (s *ClaudeCodeStorageV2) List() ([]string, error) {
	// Check if any keyring has credentials
	for _, kr := range s.keyrings {
		keys, err := kr.Keys()
		if err != nil {
			continue
		}
		
		if len(keys) > 0 {
			// Try to verify at least one key has valid Claude credentials
			for _, key := range keys {
				item, err := kr.Get(key)
				if err != nil {
					continue
				}
				
				var creds ClaudeCodeCredentials
				if err := json.Unmarshal(item.Data, &creds); err == nil {
					if creds.ClaudeAiOauth.AccessToken != "" {
						return []string{"anthropic"}, nil
					}
				}
			}
		}
	}
	
	return []string{}, nil
}

// IsAvailable checks if any backend is available
func (s *ClaudeCodeStorageV2) IsAvailable() bool {
	for _, kr := range s.keyrings {
		if _, err := kr.Keys(); err == nil {
			return true
		}
	}
	return false
}

// RequiresUnlock checks if the backend needs to be unlocked
func (s *ClaudeCodeStorageV2) RequiresUnlock() bool {
	return false
}

// Unlock is a no-op for Claude Code storage
func (s *ClaudeCodeStorageV2) Unlock() error {
	return nil
}

// Lock is a no-op for Claude Code storage
func (s *ClaudeCodeStorageV2) Lock() error {
	return nil
}

// Name returns the backend name
func (s *ClaudeCodeStorageV2) Name() string {
	return "claude-code-adapter-v2"
}