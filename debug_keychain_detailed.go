// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/99designs/keyring"
)

func main() {
	fmt.Println("=== Detailed Keychain Debug ===")
	fmt.Printf("Current user: %s\n", os.Getenv("USER"))
	fmt.Printf("Home directory: %s\n", os.Getenv("HOME"))
	
	// Test 1: Try Claude Code service name
	fmt.Println("\n1. Testing 'Claude Code-credentials' service:")
	testKeyring("Claude Code-credentials")
	
	// Test 2: Try claude.ai service name
	fmt.Println("\n2. Testing 'claude.ai' service:")
	testKeyring("claude.ai")
	
	// Test 3: Try com.anthropic.claude-code
	fmt.Println("\n3. Testing 'com.anthropic.claude-code' service:")
	testKeyring("com.anthropic.claude-code")
	
	// Test 4: Try Claude Code (with space)
	fmt.Println("\n4. Testing 'Claude Code' service:")
	testKeyring("Claude Code")
	
	// Test 5: List all available keychains using security command
	fmt.Println("\n5. Searching keychain with security command:")
	searchKeychain()
}

func testKeyring(serviceName string) {
	kr, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		fmt.Printf("  ❌ Failed to open keyring: %v\n", err)
		return
	}
	
	keys, err := kr.Keys()
	if err != nil {
		fmt.Printf("  ❌ Failed to list keys: %v\n", err)
		return
	}
	
	if len(keys) == 0 {
		fmt.Printf("  ⚠️  No keys found\n")
		return
	}
	
	fmt.Printf("  ✓ Found %d keys:\n", len(keys))
	for _, key := range keys {
		fmt.Printf("    - %s\n", key)
		
		// Try to get the item
		item, err := kr.Get(key)
		if err != nil {
			fmt.Printf("      ❌ Failed to get item: %v\n", err)
			continue
		}
		
		// Try to parse as JSON
		var data map[string]interface{}
		if err := json.Unmarshal(item.Data, &data); err != nil {
			fmt.Printf("      ⚠️  Not JSON data\n")
		} else {
			// Check if it looks like Claude credentials
			if claudeOauth, ok := data["claudeAiOauth"]; ok {
				fmt.Printf("      ✓ Found Claude OAuth data!\n")
				if oauthMap, ok := claudeOauth.(map[string]interface{}); ok {
					if token, ok := oauthMap["accessToken"].(string); ok && len(token) > 20 {
						fmt.Printf("      ✓ Has valid access token (length: %d)\n", len(token))
					}
				}
			}
		}
	}
}

func searchKeychain() {
	// Use the security command to search for claude-related items
	fmt.Println("  Running: security dump-keychain | grep -i claude")
	fmt.Println("  (This may prompt for your password)")
	
	// Note: We can't actually run this from Go easily due to password prompts
	// But we can suggest the command
	fmt.Println("\n  Please run this command manually to search for Claude entries:")
	fmt.Println("  security find-generic-password -s 'Claude' 2>&1 | grep -E 'svce|acct'")
	fmt.Println("  security find-generic-password -s 'claude' 2>&1 | grep -E 'svce|acct'")
	fmt.Println("  security find-internet-password -s 'claude.ai' 2>&1 | grep -E 'srvr|acct'")
}