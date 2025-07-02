// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/99designs/keyring"
)

func main() {
	fmt.Println("=== Testing Keyring Isolation ===")
	
	// Test 1: Open Claude Code keyring
	fmt.Println("\n1. Claude Code keyring (service: Claude Code-credentials):")
	kr1, err := keyring.Open(keyring.Config{
		ServiceName: "Claude Code-credentials",
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		keys, err := kr1.Keys()
		if err != nil {
			fmt.Printf("   Error listing keys: %v\n", err)
		} else {
			fmt.Printf("   Found %d keys: %v\n", len(keys), keys)
			for _, key := range keys {
				item, err := kr1.Get(key)
				if err == nil {
					var data map[string]interface{}
					if json.Unmarshal(item.Data, &data) == nil {
						if _, ok := data["claudeAiOauth"]; ok {
							fmt.Printf("   ✓ Key '%s' has Claude OAuth data\n", key)
						}
					}
				}
			}
		}
	}
	
	// Test 2: Open claude-gate keyring
	fmt.Println("\n2. claude-gate keyring (service: claude-gate):")
	kr2, err := keyring.Open(keyring.Config{
		ServiceName: "claude-gate",
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		keys, err := kr2.Keys()
		if err != nil {
			fmt.Printf("   Error listing keys: %v\n", err)
		} else {
			fmt.Printf("   Found %d keys: %v\n", len(keys), keys)
			for _, key := range keys {
				fmt.Printf("   - Key: %s\n", key)
			}
		}
	}
	
	fmt.Println("\n3. Checking isolation:")
	fmt.Println("   Claude Code stores under: service='Claude Code-credentials', account=username")
	fmt.Println("   claude-gate stores under: service='claude-gate', account='claude-gate.anthropic'")
	fmt.Println("   These are completely separate keychain items.")
}