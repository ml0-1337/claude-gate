// +build ignore

package main

import (
	"fmt"
	"os"
	"runtime"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/config"
)

func main() {
	fmt.Println("=== Testing Storage Creation ===")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("User: %s\n", os.Getenv("USER"))
	
	// Create config
	cfg := config.DefaultConfig()
	cfg.AuthStorageType = "claude-code"
	
	// Create factory
	factory := auth.NewStorageFactory(auth.StorageFactoryConfig{
		Type: auth.StorageTypeClaudeCode,
	})
	
	fmt.Println("\nCreating storage...")
	
	// Try to create storage
	storage, err := factory.Create()
	if err != nil {
		fmt.Printf("Error creating storage: %v\n", err)
		return
	}
	
	fmt.Printf("Storage created: %s\n", storage.Name())
	fmt.Printf("Storage available: %v\n", storage.IsAvailable())
	
	// Try to get token
	fmt.Println("\nGetting token...")
	token, err := storage.Get("anthropic")
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}
	
	if token == nil {
		fmt.Println("No token found")
		return
	}
	
	fmt.Println("Token found!")
	fmt.Printf("  Type: %s\n", token.Type)
	fmt.Printf("  Has credentials: %v\n", token.AccessToken != "")
}