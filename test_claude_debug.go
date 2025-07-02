// +build ignore

package main

import (
	"fmt"
	"os"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
)

func main() {
	os.Setenv("CLAUDE_GATE_DEBUG", "true")
	
	fmt.Println("=== Testing Claude Code Storage with Debug ===")
	
	// Create storage with debug
	storage, err := auth.NewClaudeCodeStorageDebug()
	if err != nil {
		fmt.Printf("Failed to create storage: %v\n", err)
		return
	}
	
	fmt.Println("\nTrying to get token...")
	
	// Call GetDebug directly
	token, err := storage.GetDebug("anthropic")
	
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}
	
	if token == nil {
		fmt.Println("No token found")
		return
	}
	
	fmt.Printf("\nSuccess! Found token:\n")
	fmt.Printf("  Type: %s\n", token.Type)
	fmt.Printf("  Has Access Token: %v\n", token.AccessToken != "")
	fmt.Printf("  Has Refresh Token: %v\n", token.RefreshToken != "")
	fmt.Printf("  Expires At: %d\n", token.ExpiresAt)
}