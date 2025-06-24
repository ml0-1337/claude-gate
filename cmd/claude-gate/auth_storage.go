package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/config"
)

// AuthStorageCmd handles auth storage management commands
type AuthStorageCmd struct {
	Status  AuthStorageStatusCmd  `cmd:"" help:"Show storage backend status"`
	Migrate AuthStorageMigrateCmd `cmd:"" help:"Migrate tokens between storage backends"`
	Test    AuthStorageTestCmd    `cmd:"" help:"Test storage backend operations"`
	Backup  AuthStorageBackupCmd  `cmd:"" help:"Create manual backup of tokens"`
}

// AuthStorageStatusCmd shows storage backend status
type AuthStorageStatusCmd struct{}

func (cmd *AuthStorageStatusCmd) Run(ctx *kong.Context) error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create storage factory
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	// Get current storage
	storage, err := factory.Create()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	fmt.Println("Storage Backend Status")
	fmt.Println("=====================")
	fmt.Printf("Type: %s\n", cfg.AuthStorageType)
	fmt.Printf("Backend: %s\n", storage.Name())
	fmt.Printf("Available: %v\n", storage.IsAvailable())
	fmt.Printf("Requires Unlock: %v\n", storage.RequiresUnlock())
	
	// List stored providers
	providers, err := storage.List()
	if err != nil {
		fmt.Printf("\nError listing providers: %v\n", err)
	} else {
		fmt.Printf("\nStored Providers: %d\n", len(providers))
		if len(providers) > 0 {
			fmt.Println("\nProviders:")
			for _, provider := range providers {
				token, err := storage.Get(provider)
				if err != nil {
					fmt.Printf("  - %s: Error reading token\n", provider)
					continue
				}
				if token == nil {
					fmt.Printf("  - %s: No token\n", provider)
					continue
				}
				
				status := "Valid"
				if token.IsExpired() {
					status = "Expired"
				} else if token.NeedsRefresh() {
					status = "Needs Refresh"
				}
				
				fmt.Printf("  - %s: %s token (%s)\n", provider, token.Type, status)
			}
		}
	}
	
	// Show configuration
	fmt.Println("\nConfiguration:")
	fmt.Printf("  Storage Path: %s\n", cfg.AuthStoragePath)
	fmt.Printf("  Keyring Service: %s\n", cfg.KeyringService)
	fmt.Printf("  Auto-Migrate: %v\n", cfg.AutoMigrateTokens)
	
	return nil
}

// AuthStorageMigrateCmd migrates tokens between storage backends
type AuthStorageMigrateCmd struct {
	From   string `help:"Source storage type (file/keyring)" default:"file"`
	To     string `help:"Destination storage type (file/keyring)" default:"keyring"`
	DryRun bool   `help:"Show what would be migrated without making changes"`
}

func (cmd *AuthStorageMigrateCmd) Run(ctx *kong.Context) error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create source storage
	sourceCfg := createStorageFactoryConfig(cfg)
	sourceCfg.Type = auth.StorageType(cmd.From)
	sourceFactory := auth.NewStorageFactory(sourceCfg)
	
	source, err := sourceFactory.Create()
	if err != nil {
		return fmt.Errorf("failed to create source storage: %w", err)
	}
	
	// Create destination storage
	destCfg := createStorageFactoryConfig(cfg)
	destCfg.Type = auth.StorageType(cmd.To)
	destCfg.FilePath = cfg.AuthStoragePath + ".migrated"
	destFactory := auth.NewStorageFactory(destCfg)
	
	destination, err := destFactory.Create()
	if err != nil {
		return fmt.Errorf("failed to create destination storage: %w", err)
	}
	
	// List tokens to migrate
	providers, err := source.List()
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}
	
	if len(providers) == 0 {
		fmt.Println("No tokens to migrate")
		return nil
	}
	
	fmt.Printf("Migrating %d tokens from %s to %s\n", len(providers), source.Name(), destination.Name())
	
	if cmd.DryRun {
		fmt.Println("\nDry run - no changes will be made:")
		for _, provider := range providers {
			fmt.Printf("  - Would migrate: %s\n", provider)
		}
		return nil
	}
	
	// Confirm migration
	fmt.Print("\nProceed with migration? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Migration cancelled")
		return nil
	}
	
	// Perform migration
	migrator := auth.NewStorageMigrator(source, destination)
	if err := migrator.Migrate(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	
	// Verify migration
	if err := migrator.VerifyMigration(); err != nil {
		fmt.Printf("Warning: Migration verification failed: %v\n", err)
	} else {
		fmt.Println("Migration completed and verified successfully")
	}
	
	return nil
}

// AuthStorageTestCmd tests storage backend operations
type AuthStorageTestCmd struct{}

func (cmd *AuthStorageTestCmd) Run(ctx *kong.Context) error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create storage
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.Create()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	fmt.Printf("Testing storage backend: %s\n\n", storage.Name())
	
	// Test availability
	fmt.Print("Testing availability... ")
	if storage.IsAvailable() {
		fmt.Println("✓ Available")
	} else {
		fmt.Println("✗ Not available")
		return fmt.Errorf("storage backend not available")
	}
	
	// Test set operation
	fmt.Print("Testing write operation... ")
	testToken := &auth.TokenInfo{
		Type:        "oauth",
		AccessToken: "test-token",
	}
	if err := storage.Set("test-provider", testToken); err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return err
	}
	fmt.Println("✓ Success")
	
	// Test get operation
	fmt.Print("Testing read operation... ")
	retrieved, err := storage.Get("test-provider")
	if err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return err
	}
	if retrieved == nil || retrieved.AccessToken != testToken.AccessToken {
		fmt.Println("✗ Token mismatch")
		return fmt.Errorf("retrieved token does not match")
	}
	fmt.Println("✓ Success")
	
	// Test list operation
	fmt.Print("Testing list operation... ")
	providers, err := storage.List()
	if err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return err
	}
	found := false
	for _, p := range providers {
		if p == "test-provider" {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("✗ Test provider not found in list")
		return fmt.Errorf("test provider not found")
	}
	fmt.Println("✓ Success")
	
	// Test remove operation
	fmt.Print("Testing remove operation... ")
	if err := storage.Remove("test-provider"); err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return err
	}
	
	// Verify removal
	retrieved, err = storage.Get("test-provider")
	if err != nil {
		fmt.Printf("✗ Failed to verify removal: %v\n", err)
		return err
	}
	if retrieved != nil {
		fmt.Println("✗ Token still exists after removal")
		return fmt.Errorf("token not removed")
	}
	fmt.Println("✓ Success")
	
	fmt.Println("\nAll tests passed!")
	return nil
}

// AuthStorageBackupCmd creates manual backup of tokens
type AuthStorageBackupCmd struct{}

func (cmd *AuthStorageBackupCmd) Run(ctx *kong.Context) error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Only backup file storage
	if cfg.AuthStorageType != "file" {
		fmt.Println("Backup is only supported for file storage")
		fmt.Println("Keyring storage is backed up by the operating system")
		return nil
	}
	
	// Check if auth file exists
	if _, err := os.Stat(cfg.AuthStoragePath); os.IsNotExist(err) {
		fmt.Println("No auth file to backup")
		return nil
	}
	
	// Create backup
	fmt.Println("Creating backup...")
	
	// For now, use a simple copy
	backupPath := cfg.AuthStoragePath + ".backup"
	data, err := os.ReadFile(cfg.AuthStoragePath)
	if err != nil {
		return fmt.Errorf("failed to read auth file: %w", err)
	}
	
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}
	
	fmt.Printf("Backup created: %s\n", backupPath)
	return nil
}