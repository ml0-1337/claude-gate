// +build darwin

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// ClaudeCodeStorageMacOS implements StorageBackend using direct macOS security command
type ClaudeCodeStorageMacOS struct{}

// NewClaudeCodeStorageMacOS creates a new macOS-specific Claude Code storage adapter
func NewClaudeCodeStorageMacOS() *ClaudeCodeStorageMacOS {
	return &ClaudeCodeStorageMacOS{}
}

// Get retrieves and transforms Claude Code credentials using macOS security command
func (s *ClaudeCodeStorageMacOS) Get(provider string) (*TokenInfo, error) {
	if provider != "anthropic" {
		return nil, nil
	}

	// Try to get the password using the security command
	username := os.Getenv("USER")
	if username == "" {
		return nil, fmt.Errorf("could not determine username")
	}

	// Use the security command to find the password
	cmd := exec.Command("security", "find-generic-password",
		"-s", "Claude Code-credentials",
		"-a", username,
		"-w") // -w returns just the password

	output, err := cmd.Output()
	if err != nil {
		// No credentials found
		return nil, nil
	}

	// Parse the password data
	var creds ClaudeCodeCredentials
	if err := json.Unmarshal(output, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse Claude Code credentials: %w", err)
	}

	// Check if it has the expected structure
	if creds.ClaudeAiOauth.AccessToken == "" || creds.ClaudeAiOauth.RefreshToken == "" {
		return nil, nil
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
func (s *ClaudeCodeStorageMacOS) Set(provider string, token *TokenInfo) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// Remove is not supported for read-only adapter
func (s *ClaudeCodeStorageMacOS) Remove(provider string) error {
	return fmt.Errorf("Claude Code storage is read-only")
}

// List returns available providers
func (s *ClaudeCodeStorageMacOS) List() ([]string, error) {
	// Check if credentials exist
	username := os.Getenv("USER")
	if username == "" {
		return []string{}, nil
	}

	cmd := exec.Command("security", "find-generic-password",
		"-s", "Claude Code-credentials",
		"-a", username)

	if err := cmd.Run(); err != nil {
		// No credentials found
		return []string{}, nil
	}

	return []string{"anthropic"}, nil
}

// IsAvailable checks if the backend is available
func (s *ClaudeCodeStorageMacOS) IsAvailable() bool {
	// Check if security command exists
	_, err := exec.LookPath("security")
	return err == nil
}

// RequiresUnlock checks if the backend needs to be unlocked
func (s *ClaudeCodeStorageMacOS) RequiresUnlock() bool {
	return false
}

// Unlock is a no-op for Claude Code storage
func (s *ClaudeCodeStorageMacOS) Unlock() error {
	return nil
}

// Lock is a no-op for Claude Code storage
func (s *ClaudeCodeStorageMacOS) Lock() error {
	return nil
}

// Name returns the backend name
func (s *ClaudeCodeStorageMacOS) Name() string {
	return "claude-code-macos"
}