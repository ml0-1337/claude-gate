// +build ignore

package main

import (
	"fmt"
	"os"
	"runtime"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
)

func main() {
	fmt.Println("=== Debug Storage Factory ===")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("User: %s\n", os.Getenv("USER"))
	
	// Test 1: Try to create ClaudeCodeStorage directly
	fmt.Println("\n1. Creating ClaudeCodeStorage directly:")
	ccs, err := auth.NewClaudeCodeStorage()
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		fmt.Printf("   Error contains 'keyring unavailable': %v\n", err.Error() == "keyring unavailable on macOS")
	} else {
		fmt.Printf("   Success! Storage: %s\n", ccs.Name())
	}
	
	// Test 2: Create through factory
	fmt.Println("\n2. Creating through factory:")
	factory := auth.NewStorageFactory(auth.StorageFactoryConfig{
		Type: auth.StorageTypeClaudeCode,
	})
	
	storage, err := factory.Create()
	if err != nil {
		fmt.Printf("   Factory error: %v\n", err)
	} else {
		fmt.Printf("   Factory success! Storage: %s\n", storage.Name())
		
		// Try to get token
		token, err := storage.Get("anthropic")
		if err != nil {
			fmt.Printf("   Get error: %v\n", err)
		} else if token == nil {
			fmt.Printf("   No token found\n")
		} else {
			fmt.Printf("   Token found! Type: %s\n", token.Type)
		}
	}
	
	// Test 3: Direct macOS adapter
	fmt.Println("\n3. Creating macOS adapter directly:")
	macosStorage := auth.NewClaudeCodeStorageMacOS()
	fmt.Printf("   Storage: %s\n", macosStorage.Name())
	fmt.Printf("   Available: %v\n", macosStorage.IsAvailable())
	
	token, err := macosStorage.Get("anthropic")
	if err != nil {
		fmt.Printf("   Get error: %v\n", err)
	} else if token == nil {
		fmt.Printf("   No token found\n")
	} else {
		fmt.Printf("   Token found! Type: %s\n", token.Type)
	}
}