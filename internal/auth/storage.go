package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TokenInfo represents stored authentication information
type TokenInfo struct {
	Type         string `json:"type"`          // "oauth" or "api"
	RefreshToken string `json:"refresh,omitempty"`
	AccessToken  string `json:"access,omitempty"`
	ExpiresAt    int64  `json:"expires,omitempty"`
	APIKey       string `json:"key,omitempty"`
}

// IsExpired checks if the token is expired
func (t *TokenInfo) IsExpired() bool {
	if t.Type != "oauth" || t.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() >= t.ExpiresAt
}

// NeedsRefresh checks if the token needs refresh (5 minutes before expiry)
func (t *TokenInfo) NeedsRefresh() bool {
	if t.Type != "oauth" || t.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() >= (t.ExpiresAt - 300) // 5 minutes buffer
}

// TokenStorage handles secure storage of authentication tokens
type TokenStorage struct {
	path string
	mu   sync.RWMutex
}

// NewTokenStorage creates a new token storage instance
func NewTokenStorage(path string) *TokenStorage {
	return &TokenStorage{
		path: path,
	}
}

// Get retrieves token information for a provider
func (s *TokenStorage) Get(provider string) (*TokenInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, err := s.loadData()
	if err != nil {
		return nil, err
	}
	
	if tokenData, exists := data[provider]; exists {
		var token TokenInfo
		// Re-marshal and unmarshal to convert map to struct
		jsonData, err := json.Marshal(tokenData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal token data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &token); err != nil {
			return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
		}
		return &token, nil
	}
	
	return nil, nil
}

// Set stores token information for a provider
func (s *TokenStorage) Set(provider string, token *TokenInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Ensure directory exists
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	data, err := s.loadData()
	if err != nil {
		return err
	}
	
	data[provider] = token
	
	// Write data to file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := os.WriteFile(s.path, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// Remove deletes token information for a provider
func (s *TokenStorage) Remove(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, err := s.loadData()
	if err != nil {
		return err
	}
	
	delete(data, provider)
	
	if len(data) == 0 {
		// Remove file if no data left
		return os.Remove(s.path)
	}
	
	// Write updated data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := os.WriteFile(s.path, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// loadData loads the stored data from file
func (s *TokenStorage) loadData() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	
	fileData, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	if len(fileData) > 0 {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}
	
	return data, nil
}