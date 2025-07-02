#!/bin/bash

echo "=== Claude Code Authentication Test ==="
echo "Username: $USER"
echo

echo "1. Checking for Claude entries in keychain..."
security find-generic-password -a "$USER" 2>&1 | grep -B5 -A5 -i claude || echo "No Claude entries found for user $USER"

echo
echo "2. Checking all generic passwords with 'Claude' in service name..."
security dump-keychain | grep -B2 -A2 "Claude" 2>/dev/null || echo "No entries found (may need keychain password)"

echo
echo "3. Testing claude-gate with debug output..."
export CLAUDE_GATE_LOG_LEVEL=DEBUG
./claude-gate-intel start --storage-backend=claude-code --skip-auth-check 2>&1 | head -20

echo
echo "4. Trying to read with Go test program..."
cat > test_read.go << 'EOF'
package main
import (
    "fmt"
    "github.com/99designs/keyring"
)
func main() {
    kr, err := keyring.Open(keyring.Config{ServiceName: "Claude Code-credentials"})
    if err != nil { fmt.Printf("Error: %v\n", err); return }
    keys, _ := kr.Keys()
    fmt.Printf("Found %d keys: %v\n", len(keys), keys)
}
EOF
go run test_read.go
rm test_read.go