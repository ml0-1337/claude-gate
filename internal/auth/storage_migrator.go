package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StorageMigrator handles migration between storage backends
type StorageMigrator struct {
	source      StorageBackend
	destination StorageBackend
	backup      bool
}

// NewStorageMigrator creates a new storage migrator
func NewStorageMigrator(source, destination StorageBackend) *StorageMigrator {
	return &StorageMigrator{
		source:      source,
		destination: destination,
		backup:      true,
	}
}

// Migrate performs the migration from source to destination
func (m *StorageMigrator) Migrate() error {
	// Get all providers from source
	providers, err := m.source.List()
	if err != nil {
		return fmt.Errorf("failed to list providers from source: %w", err)
	}
	
	if len(providers) == 0 {
		return nil // Nothing to migrate
	}
	
	// Create backup if requested
	if m.backup {
		if err := m.createBackup(); err != nil {
			// Log warning but continue
			fmt.Fprintf(os.Stderr, "Warning: Failed to create backup: %v\n", err)
		}
	}
	
	// Track migration results
	migrated := 0
	failed := 0
	
	// Migrate each provider
	for _, provider := range providers {
		// Get token from source
		token, err := m.source.Get(provider)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get token for %s: %v\n", provider, err)
			failed++
			continue
		}
		
		if token == nil {
			continue // Skip empty entries
		}
		
		// Set token in destination
		if err := m.destination.Set(provider, token); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to set token for %s: %v\n", provider, err)
			failed++
			continue
		}
		
		migrated++
	}
	
	// Report results
	if failed > 0 {
		return fmt.Errorf("migration completed with errors: %d migrated, %d failed", migrated, failed)
	}
	
	// Mark source as migrated (don't delete yet)
	if err := m.markSourceMigrated(); err != nil {
		// Non-fatal error
		fmt.Fprintf(os.Stderr, "Warning: Failed to mark source as migrated: %v\n", err)
	}
	
	return nil
}

// MigrateProvider migrates a single provider
func (m *StorageMigrator) MigrateProvider(provider string) error {
	// Get token from source
	token, err := m.source.Get(provider)
	if err != nil {
		return fmt.Errorf("failed to get token from source: %w", err)
	}
	
	if token == nil {
		return fmt.Errorf("no token found for provider: %s", provider)
	}
	
	// Set token in destination
	if err := m.destination.Set(provider, token); err != nil {
		return fmt.Errorf("failed to set token in destination: %w", err)
	}
	
	return nil
}

// Rollback attempts to rollback a failed migration
func (m *StorageMigrator) Rollback() error {
	// Swap source and destination
	rollbackMigrator := &StorageMigrator{
		source:      m.destination,
		destination: m.source,
		backup:      false, // Don't create backup during rollback
	}
	
	return rollbackMigrator.Migrate()
}

// createBackup creates a backup of the source storage
func (m *StorageMigrator) createBackup() error {
	// Only backup file storage
	fileStorage, ok := m.source.(*FileStorage)
	if !ok {
		return nil // Can't backup non-file storage
	}
	
	// Create backup directory
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".claude-gate", "backups")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	
	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("auth-%s.json", timestamp))
	
	// Copy file
	sourceData, err := os.ReadFile(fileStorage.path)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	
	if err := os.WriteFile(backupPath, sourceData, 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	
	fmt.Fprintf(os.Stderr, "Created backup at: %s\n", backupPath)
	return nil
}

// markSourceMigrated marks the source as migrated
func (m *StorageMigrator) markSourceMigrated() error {
	// Only mark file storage
	fileStorage, ok := m.source.(*FileStorage)
	if !ok {
		return nil
	}
	
	// Rename file to indicate migration
	migratedPath := fileStorage.path + ".migrated"
	
	// Check if already exists
	if _, err := os.Stat(migratedPath); err == nil {
		// Already migrated, remove old one
		os.Remove(migratedPath)
	}
	
	return os.Rename(fileStorage.path, migratedPath)
}

// VerifyMigration verifies that all data was migrated correctly
func (m *StorageMigrator) VerifyMigration() error {
	// Get all providers from source
	sourceProviders, err := m.source.List()
	if err != nil {
		return fmt.Errorf("failed to list source providers: %w", err)
	}
	
	// Check each provider in destination
	for _, provider := range sourceProviders {
		sourceToken, err := m.source.Get(provider)
		if err != nil {
			return fmt.Errorf("failed to get source token for %s: %w", provider, err)
		}
		
		destToken, err := m.destination.Get(provider)
		if err != nil {
			return fmt.Errorf("failed to get destination token for %s: %w", provider, err)
		}
		
		// Compare tokens
		if !tokensEqual(sourceToken, destToken) {
			return fmt.Errorf("token mismatch for provider %s", provider)
		}
	}
	
	return nil
}

// tokensEqual compares two tokens for equality
func tokensEqual(a, b *TokenInfo) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	return a.Type == b.Type &&
		a.RefreshToken == b.RefreshToken &&
		a.AccessToken == b.AccessToken &&
		a.ExpiresAt == b.ExpiresAt &&
		a.APIKey == b.APIKey
}