package auth

import (
	"fmt"
	"os"
	"runtime"

	"github.com/99designs/keyring"
)

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeAuto    StorageType = "auto"    // Automatically select best available
	StorageTypeKeyring StorageType = "keyring" // Force keyring storage
	StorageTypeFile    StorageType = "file"    // Force file storage
)

// StorageFactory creates storage backends based on configuration
type StorageFactory struct {
	storageType    StorageType
	filePath       string
	keyringConfig  KeyringConfig
	passwordPrompt keyring.PromptFunc
}

// StorageFactoryConfig holds configuration for the storage factory
type StorageFactoryConfig struct {
	Type           StorageType
	FilePath       string
	ServiceName    string
	PasswordPrompt keyring.PromptFunc
	
	// macOS-specific settings
	KeychainTrustApp               bool
	KeychainAccessibleWhenUnlocked bool
	KeychainSynchronizable         bool
}

// NewStorageFactory creates a new storage factory
func NewStorageFactory(config StorageFactoryConfig) *StorageFactory {
	// Set defaults
	if config.Type == "" {
		config.Type = StorageTypeAuto
	}
	
	if config.ServiceName == "" {
		config.ServiceName = "claude-gate"
	}
	
	if config.FilePath == "" {
		homeDir, _ := os.UserHomeDir()
		config.FilePath = homeDir + "/.claude-gate/auth.json"
	}
	
	// Default password prompt if not provided
	if config.PasswordPrompt == nil {
		config.PasswordPrompt = defaultPasswordPrompt
	}
	
	// Set keyring config
	keyringCfg := KeyringConfig{
		ServiceName:    config.ServiceName,
		PasswordPrompt: config.PasswordPrompt,
		Debug:          false,
	}
	
	// Apply macOS settings from config (with defaults if not set)
	if runtime.GOOS == "darwin" {
		// Use config values, but default to true/true/false if not explicitly set
		keyringCfg.KeychainTrustApplication = config.KeychainTrustApp
		keyringCfg.KeychainAccessibleWhenUnlocked = config.KeychainAccessibleWhenUnlocked
		keyringCfg.KeychainSynchronizable = config.KeychainSynchronizable
		
		// If the struct fields are zero values and we're on macOS, apply sensible defaults
		// This handles backward compatibility when the fields aren't explicitly set
		if !config.KeychainTrustApp && !config.KeychainAccessibleWhenUnlocked && !config.KeychainSynchronizable {
			keyringCfg.KeychainTrustApplication = true       // Trust app by default
			keyringCfg.KeychainAccessibleWhenUnlocked = true // Accessible when unlocked
			keyringCfg.KeychainSynchronizable = false        // Don't sync to iCloud
		}
	}
	
	return &StorageFactory{
		storageType:    config.Type,
		filePath:       config.FilePath,
		keyringConfig:  keyringCfg,
		passwordPrompt: config.PasswordPrompt,
	}
}

// Create creates a storage backend based on configuration
func (f *StorageFactory) Create() (StorageBackend, error) {
	switch f.storageType {
	case StorageTypeFile:
		return NewFileStorage(f.filePath), nil
		
	case StorageTypeKeyring:
		ks, err := NewKeyringStorage(f.keyringConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create keyring storage: %w", err)
		}
		return ks, nil
		
	case StorageTypeAuto:
		// Try keyring first, fall back to file
		if isKeyringAvailable() {
			ks, err := NewKeyringStorage(f.keyringConfig)
			if err == nil && ks.IsAvailable() {
				return ks, nil
			}
			// Log warning and fall back
			fmt.Fprintf(os.Stderr, "Warning: Keyring storage unavailable, falling back to file storage: %v\n", err)
		}
		return NewFileStorage(f.filePath), nil
		
	default:
		return nil, fmt.Errorf("unknown storage type: %s", f.storageType)
	}
}

// CreateWithMigration creates a storage backend and migrates data if needed
func (f *StorageFactory) CreateWithMigration() (StorageBackend, error) {
	// Create the target storage
	storage, err := f.Create()
	if err != nil {
		return nil, err
	}
	
	// Check if we need to migrate from file storage
	if f.storageType != StorageTypeFile {
		fileStorage := NewFileStorage(f.filePath)
		
		// Check if file storage has data
		providers, err := fileStorage.List()
		if err == nil && len(providers) > 0 {
			// Migrate data
			migrator := NewStorageMigrator(fileStorage, storage)
			if err := migrator.Migrate(); err != nil {
				return nil, fmt.Errorf("failed to migrate storage: %w", err)
			}
			
			fmt.Fprintf(os.Stderr, "Successfully migrated %d tokens to %s\n", len(providers), storage.Name())
		}
	}
	
	return storage, nil
}

// isKeyringAvailable checks if keyring functionality is available on this system
func isKeyringAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		// macOS always has Keychain
		return true
	case "linux":
		// Check for Secret Service or KWallet
		// This is a simplified check - in reality we'd check D-Bus
		return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	case "windows":
		// Windows always has Credential Manager
		return true
	default:
		return false
	}
}

// defaultPasswordPrompt provides a default password prompt function
func defaultPasswordPrompt(prompt string) (string, error) {
	// In a real implementation, this would use terminal input
	// For now, return an error to force non-interactive mode
	return "", fmt.Errorf("interactive password prompt not implemented - use environment variable or config file")
}