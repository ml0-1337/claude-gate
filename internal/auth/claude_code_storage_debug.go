package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/99designs/keyring"
)

// NewClaudeCodeStorageDebug creates a Claude Code storage adapter with debug output
func NewClaudeCodeStorageDebug() (*ClaudeCodeStorage, error) {
	debug := os.Getenv("CLAUDE_GATE_DEBUG") == "true"
	
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Creating Claude Code storage adapter\n")
		fmt.Fprintf(os.Stderr, "[DEBUG] Service name: Claude Code-credentials\n")
	}
	
	// Open keyring with Claude Code's service name
	kr, err := keyring.Open(keyring.Config{
		ServiceName: "Claude Code-credentials",
		// Let it auto-detect the backend
	})
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Failed to open keyring: %v\n", err)
		}
		return nil, fmt.Errorf("failed to open Claude Code keyring: %w", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Keyring opened successfully\n")
		
		// Try to list keys
		keys, err := kr.Keys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] Failed to list keys: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] Found %d keys: %v\n", len(keys), keys)
		}
	}

	return &ClaudeCodeStorage{
		keyring: kr,
	}, nil
}

// GetDebug is a debug version of Get that prints what it's doing
func (s *ClaudeCodeStorage) GetDebug(provider string) (*TokenInfo, error) {
	debug := os.Getenv("CLAUDE_GATE_DEBUG") == "true"
	
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Get called with provider: %s\n", provider)
	}
	
	// Claude Code only supports anthropic
	if provider != "anthropic" {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Provider is not 'anthropic', returning nil\n")
		}
		return nil, nil
	}

	// List all keys
	keys, err := s.keyring.Keys()
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Failed to list keys: %v\n", err)
		}
		return nil, nil
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Found %d keys in keyring\n", len(keys))
	}

	// If no keys, no credentials
	if len(keys) == 0 {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] No keys found, returning nil\n")
		}
		return nil, nil
	}

	// Try each key
	var item keyring.Item
	var foundItem bool
	for _, key := range keys {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Trying key: %s\n", key)
		}
		
		item, err = s.keyring.Get(key)
		if err != nil {
			if debug {
				fmt.Fprintf(os.Stderr, "[DEBUG] Failed to get key %s: %v\n", key, err)
			}
			continue
		}
		
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Successfully got key %s, data length: %d\n", key, len(item.Data))
		}
		
		foundItem = true
		break
	}

	if !foundItem {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Could not retrieve any items\n")
		}
		return nil, nil
	}

	// Parse Claude Code format
	var creds ClaudeCodeCredentials
	if err := json.Unmarshal(item.Data, &creds); err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Failed to parse JSON: %v\n", err)
			fmt.Fprintf(os.Stderr, "[DEBUG] Raw data: %s\n", string(item.Data))
		}
		return nil, fmt.Errorf("failed to parse Claude Code credentials: %w", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Successfully parsed credentials\n")
		fmt.Fprintf(os.Stderr, "[DEBUG] Has access token: %v\n", creds.ClaudeAiOauth.AccessToken != "")
		fmt.Fprintf(os.Stderr, "[DEBUG] Has refresh token: %v\n", creds.ClaudeAiOauth.RefreshToken != "")
		fmt.Fprintf(os.Stderr, "[DEBUG] Expires at: %d\n", creds.ClaudeAiOauth.ExpiresAt)
	}

	// Transform to our format
	token := &TokenInfo{
		Type:         "oauth",
		AccessToken:  creds.ClaudeAiOauth.AccessToken,
		RefreshToken: creds.ClaudeAiOauth.RefreshToken,
		ExpiresAt:    creds.ClaudeAiOauth.ExpiresAt / 1000, // Convert milliseconds to seconds
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Returning token with type: %s\n", token.Type)
	}

	return token, nil
}