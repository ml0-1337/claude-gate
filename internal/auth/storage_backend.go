package auth

import (
	"errors"
	"time"
)

// Common storage errors
var (
	ErrKeyringLocked       = errors.New("keyring is locked - unlock required")
	ErrKeyringUnavailable  = errors.New("keyring backend not available")
	ErrKeyringAccessDenied = errors.New("keyring access denied")
	ErrKeyringCorrupted    = errors.New("keyring data corrupted")
	ErrKeyringTimeout      = errors.New("keyring operation timed out")
)

// StorageBackend defines the interface for token storage implementations
type StorageBackend interface {
	// Get retrieves token information for a provider
	Get(provider string) (*TokenInfo, error)
	
	// Set stores token information for a provider
	Set(provider string, token *TokenInfo) error
	
	// Remove deletes token information for a provider
	Remove(provider string) error
	
	// List returns all stored provider names
	List() ([]string, error)
	
	// IsAvailable checks if the backend is available on this system
	IsAvailable() bool
	
	// RequiresUnlock checks if the backend needs to be unlocked
	RequiresUnlock() bool
	
	// Unlock attempts to unlock the backend (if supported)
	Unlock() error
	
	// Lock locks the backend (if supported)
	Lock() error
	
	// Name returns the backend name for identification
	Name() string
}

// StorageMetrics tracks storage operation metrics
type StorageMetrics struct {
	Operations map[string]int64
	Errors     map[string]int64
	Latencies  map[string]time.Duration
	LastError  error
	LastAccess time.Time
}

// BackendHealth represents the health status of a storage backend
type BackendHealth struct {
	Available   bool
	Locked      bool
	LastChecked time.Time
	Error       error
}