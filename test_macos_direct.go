// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
)

func main() {
	fmt.Println("=== Testing Direct macOS Security Command ===")
	fmt.Printf("User: %s\n", os.Getenv("USER"))
	
	// Test 1: Direct security command
	fmt.Println("\n1. Testing security command directly:")
	cmd := exec.Command("security", "find-generic-password",
		"-s", "Claude Code-credentials",
		"-a", os.Getenv("USER"),
		"-w")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
	} else {
		fmt.Printf("Success! Password data length: %d bytes\n", len(output))
		// Don't print the actual password
		if len(output) > 100 {
			fmt.Println("Password looks like OAuth data (long JSON)")
		}
	}
	
	// Test 2: Use the macOS storage adapter
	fmt.Println("\n2. Testing macOS storage adapter:")
	storage := auth.NewClaudeCodeStorageMacOS()
	
	fmt.Printf("Storage available: %v\n", storage.IsAvailable())
	fmt.Printf("Storage name: %s\n", storage.Name())
	
	// Try to get token
	token, err := storage.Get("anthropic")
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
	} else if token == nil {
		fmt.Println("No token found")
	} else {
		fmt.Println("Successfully got token!")
		fmt.Printf("  Type: %s\n", token.Type)
		fmt.Printf("  Has access token: %v\n", token.AccessToken != "")
		fmt.Printf("  Has refresh token: %v\n", token.RefreshToken != "")
	}
	
	// Test 3: List providers
	providers, err := storage.List()
	if err != nil {
		fmt.Printf("Error listing: %v\n", err)
	} else {
		fmt.Printf("Providers: %v\n", providers)
	}
}