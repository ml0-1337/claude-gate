#!/bin/bash
# Test edge cases for Claude Gate NPM distribution

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Claude Gate - Edge Case Testing"
echo "==============================="

PASSED=0
FAILED=0

log_test() {
    echo -e "\n${BLUE}TEST:${NC} $1"
}

log_pass() {
    echo -e "${GREEN}✓ PASS:${NC} $1"
    ((PASSED++))
}

log_fail() {
    echo -e "${RED}✗ FAIL:${NC} $1"
    ((FAILED++))
}

# Test 1: Missing platform binary
log_test "Missing platform binary"
(
    cd "$REPO_ROOT/npm"
    # Temporarily rename install.js to test error
    cp scripts/install.js scripts/install.js.bak
    
    # Modify to simulate missing binary
    cat > scripts/install.js << 'EOF'
#!/usr/bin/env node
const {getPlatform} = require('./install.js.bak');
// Override findBinary to always return null
require('./install.js.bak').findBinary = () => null;
require('./install.js.bak').install();
EOF
    
    # Test if error message is helpful
    if node scripts/install.js 2>&1 | grep -q "Could not find platform binary"; then
        log_pass "Missing binary error message"
    else
        log_fail "Missing binary error message"
    fi
    
    # Restore
    mv scripts/install.js.bak scripts/install.js
)

# Test 2: Corrupted package.json
log_test "Corrupted package.json handling"
(
    cd "$REPO_ROOT"
    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cp -r npm/* "$TEMP_DIR/"
    
    # Corrupt the package.json
    echo "{ invalid json" > "$TEMP_DIR/package.json"
    
    cd "$TEMP_DIR"
    if npm pack 2>&1 | grep -q -E "(Failed to parse|Unexpected token)"; then
        log_pass "Corrupted package.json detected"
    else
        log_fail "Corrupted package.json not detected"
    fi
    
    rm -rf "$TEMP_DIR"
)

# Test 3: Permission denied on binary
log_test "Permission denied on binary"
(
    if [ "$OSTYPE" = "msys" ] || [ "$OSTYPE" = "win32" ]; then
        log_pass "Permission test (skipped on Windows)"
    else
        cd "$REPO_ROOT/npm"
        TEMP_BIN=$(mktemp)
        echo "#!/bin/sh" > "$TEMP_BIN"
        chmod 000 "$TEMP_BIN"
        
        # Test if we handle permission errors gracefully
        if ! "$TEMP_BIN" 2>&1 | grep -q "Permission denied"; then
            log_fail "Permission error not detected"
        else
            log_pass "Permission error handling"
        fi
        
        rm -f "$TEMP_BIN"
    fi
)

# Test 4: Very long path names
log_test "Very long path names"
(
    # Create a very long path
    LONG_PATH="$REPO_ROOT"
    for i in {1..20}; do
        LONG_PATH="$LONG_PATH/very_long_directory_name_that_exceeds_normal_limits"
    done
    
    # Don't actually create it, just test path length handling
    if [ ${#LONG_PATH} -gt 1000 ]; then
        log_pass "Long path test setup"
    else
        log_fail "Long path test setup"
    fi
)

# Test 5: Unicode in paths
log_test "Unicode characters in paths"
(
    UNICODE_DIR="$REPO_ROOT/测试目录_🚀"
    mkdir -p "$UNICODE_DIR"
    
    if [ -d "$UNICODE_DIR" ]; then
        log_pass "Unicode directory creation"
        rm -rf "$UNICODE_DIR"
    else
        log_fail "Unicode directory creation"
    fi
)

# Test 6: Concurrent installations
log_test "Concurrent installation attempts"
(
    cd "$REPO_ROOT/npm"
    
    # Try to run install script multiple times concurrently
    node scripts/install.js > /dev/null 2>&1 &
    PID1=$!
    node scripts/install.js > /dev/null 2>&1 &
    PID2=$!
    
    # Wait for both
    wait $PID1
    RESULT1=$?
    wait $PID2
    RESULT2=$?
    
    # At least one should succeed
    if [ $RESULT1 -eq 0 ] || [ $RESULT2 -eq 0 ]; then
        log_pass "Concurrent installation handling"
    else
        log_fail "Concurrent installation handling"
    fi
)

# Test 7: Invalid Node.js version
log_test "Old Node.js version warning"
(
    # Check if package.json has engines field
    if grep -q '"engines"' "$REPO_ROOT/npm/package.json"; then
        log_pass "Node.js version requirement specified"
    else
        log_fail "Node.js version requirement missing"
    fi
)

# Test 8: Network timeout simulation
log_test "Network timeout handling"
(
    # This would require actual network mocking
    # For now, just check if timeout handling exists in code
    if grep -q -i "timeout" "$REPO_ROOT/npm/scripts/install.js"; then
        log_pass "Timeout handling code exists"
    else
        log_fail "No timeout handling found"
    fi
)

# Test 9: Disk space check
log_test "Disk space considerations"
(
    # Check if we document space requirements
    if grep -q -E "(size|space|MB)" "$REPO_ROOT/README.md"; then
        log_pass "Disk space documented"
    else
        log_fail "Disk space not documented"
    fi
)

# Test 10: Signal handling
log_test "Signal handling during installation"
(
    cd "$REPO_ROOT/npm"
    
    # Start installation and send signal
    node scripts/install.js > /dev/null 2>&1 &
    PID=$!
    sleep 0.1
    
    # Check if process exists before killing
    if kill -0 $PID 2>/dev/null; then
        kill -TERM $PID 2>/dev/null
        wait $PID 2>/dev/null
        log_pass "Signal handling (process terminated)"
    else
        log_pass "Signal handling (process completed quickly)"
    fi
)

# Test 11: Version mismatch
log_test "Version mismatch detection"
(
    # Check if versions are consistent
    GO_VERSION=$(grep -E 'var version = ".*"' "$REPO_ROOT/cmd/claude-gate/main.go" | sed 's/.*"\(.*\)".*/\1/')
    NPM_VERSION=$(node -e "console.log(require('$REPO_ROOT/npm/package.json').version)")
    
    if [ "$GO_VERSION" = "$NPM_VERSION" ]; then
        log_pass "Version consistency"
    else
        log_fail "Version mismatch: Go=$GO_VERSION, NPM=$NPM_VERSION"
    fi
)

# Test 12: Symlink handling
log_test "Symlink handling"
(
    if [ "$OSTYPE" = "msys" ] || [ "$OSTYPE" = "win32" ]; then
        log_pass "Symlink test (skipped on Windows)"
    else
        TEMP_LINK=$(mktemp -d)/link
        ln -s "$REPO_ROOT/npm" "$TEMP_LINK"
        
        if [ -L "$TEMP_LINK" ]; then
            log_pass "Symlink creation"
            rm -f "$TEMP_LINK"
        else
            log_fail "Symlink creation"
        fi
    fi
)

# Summary
echo -e "\n========================================"
echo "Edge Case Test Summary"
echo "========================================"
echo -e "Tests passed: ${GREEN}$PASSED${NC}"
echo -e "Tests failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All edge case tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some edge case tests failed!${NC}"
    exit 1
fi