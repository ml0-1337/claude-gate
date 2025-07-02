// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/99designs/keyring"
)

func main() {
	fmt.Println("=== Testing Keyring Backends ===")
	fmt.Printf("OS: %s\n", os.Getenv("GOOS"))
	fmt.Printf("User: %s\n", os.Getenv("USER"))
	
	// Test 1: Auto-detect backend
	fmt.Println("\n1. Testing with auto-detect:")
	testKeyring(keyring.Config{
		ServiceName: "Claude Code-credentials",
	})
	
	// Test 2: File backend
	fmt.Println("\n2. Testing with File backend:")
	testKeyring(keyring.Config{
		ServiceName: "Claude Code-credentials",
		AllowedBackends: []keyring.BackendType{
			keyring.FileBackend,
		},
		FileDir: os.Getenv("HOME") + "/.claude-code",
	})
	
	// Test 3: Pass backend (might be available)
	fmt.Println("\n3. Testing with Pass backend:")
	testKeyring(keyring.Config{
		ServiceName: "Claude Code-credentials",
		AllowedBackends: []keyring.BackendType{
			keyring.PassBackend,
		},
	})
	
	// Test 4: Try with different service names
	fmt.Println("\n4. Testing different service names:")
	serviceNames := []string{
		"Claude Code-credentials",
		"claude-code",
		"claude.ai",
		"com.anthropic.claude",
	}
	
	for _, svc := range serviceNames {
		fmt.Printf("\n   Service: %s\n", svc)
		kr, err := keyring.Open(keyring.Config{
			ServiceName: svc,
		})
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			continue
		}
		
		keys, err := kr.Keys()
		if err != nil {
			fmt.Printf("   ❌ Can't list keys: %v\n", err)
		} else if len(keys) == 0 {
			fmt.Printf("   ⚠️  No keys found\n")
		} else {
			fmt.Printf("   ✓ Found %d keys: %v\n", len(keys), keys)
		}
	}
}

func testKeyring(config keyring.Config) {
	kr, err := keyring.Open(config)
	if err != nil {
		fmt.Printf("  ❌ Failed to open: %v\n", err)
		return
	}
	
	fmt.Printf("  ✓ Opened successfully\n")
	
	keys, err := kr.Keys()
	if err != nil {
		fmt.Printf("  ❌ Failed to list keys: %v\n", err)
		return
	}
	
	fmt.Printf("  ✓ Found %d keys: %v\n", len(keys), keys)
	
	// Try to get each key
	for _, key := range keys {
		item, err := kr.Get(key)
		if err != nil {
			fmt.Printf("    ❌ Failed to get %s: %v\n", key, err)
			continue
		}
		
		// Check if it's Claude data
		var data map[string]interface{}
		if err := json.Unmarshal(item.Data, &data); err == nil {
			if _, ok := data["claudeAiOauth"]; ok {
				fmt.Printf("    ✓ %s contains Claude OAuth data\n", key)
			}
		}
	}
}