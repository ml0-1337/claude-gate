#!/bin/bash
# Comprehensive test script for Claude Gate cross-platform builds and NPM distribution

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
PASSED=0
FAILED=0
SKIPPED=0

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Test result tracking (use arrays for compatibility)
TEST_NAMES=()
TEST_STATUSES=()

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED++))
    TEST_NAMES+=("$1")
    TEST_STATUSES+=("PASSED")
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED++))
    TEST_NAMES+=("$1")
    TEST_STATUSES+=("FAILED")
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_skip() {
    echo -e "${YELLOW}[SKIP]${NC} $1"
    ((SKIPPED++))
    TEST_NAMES+=("$1")
    TEST_STATUSES+=("SKIPPED")
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    local prereqs_met=true
    
    # Check Go
    if command -v go >/dev/null 2>&1; then
        log_success "Go installed: $(go version)"
    else
        log_error "Go not installed"
        prereqs_met=false
    fi
    
    # Check Node.js
    if command -v node >/dev/null 2>&1; then
        log_success "Node.js installed: $(node --version)"
    else
        log_error "Node.js not installed"
        prereqs_met=false
    fi
    
    # Check npm
    if command -v npm >/dev/null 2>&1; then
        log_success "npm installed: $(npm --version)"
    else
        log_error "npm not installed"
        prereqs_met=false
    fi
    
    # Check GoReleaser
    if command -v goreleaser >/dev/null 2>&1; then
        log_success "GoReleaser installed: $(goreleaser --version | head -1)"
    else
        log_warning "GoReleaser not installed - some tests will be skipped"
    fi
    
    if [ "$prereqs_met" = false ]; then
        log_error "Missing prerequisites. Please install required tools."
        exit 1
    fi
}

# Test 1: Go build and tests
test_go_build() {
    log_info "Testing Go build and unit tests..."
    
    cd "$REPO_ROOT"
    
    # Run Go tests
    if go test -v -race ./... >/dev/null 2>&1; then
        log_success "Go unit tests"
    else
        log_error "Go unit tests"
        return 1
    fi
    
    # Test basic build
    if go build -o claude-gate-test ./cmd/claude-gate; then
        log_success "Go build"
        rm -f claude-gate-test
    else
        log_error "Go build"
        return 1
    fi
}

# Test 2: GoReleaser snapshot build
test_goreleaser_build() {
    log_info "Testing GoReleaser snapshot build..."
    
    if ! command -v goreleaser >/dev/null 2>&1; then
        log_skip "GoReleaser build (not installed)"
        return 0
    fi
    
    cd "$REPO_ROOT"
    
    # Clean previous builds
    rm -rf dist/
    
    # Run snapshot build
    if goreleaser build --snapshot --clean >/dev/null 2>&1; then
        log_success "GoReleaser snapshot build"
        
        # Verify expected files exist
        local expected_files=(
            "dist/claude-gate_darwin_amd64_v1/claude-gate"
            "dist/claude-gate_darwin_arm64/claude-gate"
            "dist/claude-gate_linux_amd64_v1/claude-gate"
            "dist/claude-gate_linux_arm64/claude-gate"
            # "dist/claude-gate_windows_amd64_v1/claude-gate.exe"  # Windows not supported
        )
        
        for file in "${expected_files[@]}"; do
            if [ -f "$file" ]; then
                log_success "Binary created: $file"
            else
                log_error "Binary missing: $file"
            fi
        done
        
        # Check binary sizes
        log_info "Binary sizes:"
        find dist -name "claude-gate*" -type f -exec ls -lh {} \; | awk '{print "  " $9 ": " $5}'
        
    else
        log_error "GoReleaser snapshot build"
        return 1
    fi
}

# Test 3: NPM package structure
test_npm_package_structure() {
    log_info "Testing NPM package structure..."
    
    cd "$REPO_ROOT"
    
    # Check required files
    local required_files=(
        "npm/package.json"
        "npm/index.js"
        "npm/scripts/install.js"
        "npm/scripts/uninstall.js"
        "npm/bin/claude-gate"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            log_success "File exists: $file"
        else
            log_error "File missing: $file"
        fi
    done
    
    # Check platform packages
    local platforms=(
        "darwin-x64"
        "darwin-arm64"
        "linux-x64"
        "linux-arm64"
        "win32-x64"
    )
    
    for platform in "${platforms[@]}"; do
        if [ -f "npm/platforms/$platform/package.json" ]; then
            log_success "Platform package: $platform"
        else
            log_error "Platform package missing: $platform"
        fi
    done
}

# Test 4: NPM install script
test_npm_install_script() {
    log_info "Testing NPM install script..."
    
    cd "$REPO_ROOT/npm"
    
    # Test platform detection
    local platform_output=$(node -e "const {getPlatform} = require('./scripts/install.js'); console.log(JSON.stringify(getPlatform()))" 2>&1)
    
    if [[ $? -eq 0 ]]; then
        log_success "Platform detection works: $platform_output"
    else
        log_error "Platform detection failed"
        return 1
    fi
    
    # Test install script syntax
    if node -c scripts/install.js 2>/dev/null; then
        log_success "Install script syntax valid"
    else
        log_error "Install script syntax error"
        return 1
    fi
    
    # Test uninstall script syntax
    if node -c scripts/uninstall.js 2>/dev/null; then
        log_success "Uninstall script syntax valid"
    else
        log_error "Uninstall script syntax error"
        return 1
    fi
}

# Test 5: Local NPM package
test_npm_local_package() {
    log_info "Testing local NPM package creation..."
    
    if [ ! -d "$REPO_ROOT/dist" ]; then
        log_skip "NPM local package test (no binaries built)"
        return 0
    fi
    
    cd "$REPO_ROOT"
    
    # Run the test script
    if ./scripts/test-npm-local.sh >/dev/null 2>&1; then
        log_success "Local NPM package test"
    else
        log_error "Local NPM package test"
        return 1
    fi
}

# Test 6: Package.json validation
test_package_json_validation() {
    log_info "Testing package.json files..."
    
    cd "$REPO_ROOT"
    
    # Main package.json
    if node -e "JSON.parse(require('fs').readFileSync('npm/package.json'))" 2>/dev/null; then
        log_success "Main package.json valid"
    else
        log_error "Main package.json invalid"
    fi
    
    # Check version consistency
    local main_version=$(node -e "console.log(require('./npm/package.json').version)")
    local go_version=$(grep -E 'var version = ".*"' cmd/claude-gate/main.go | sed 's/.*"\(.*\)".*/\1/')
    
    if [ "$main_version" = "$go_version" ]; then
        log_success "Version consistency: $main_version"
    else
        log_error "Version mismatch: npm=$main_version, go=$go_version"
    fi
    
    # Platform packages
    for platform in darwin-x64 darwin-arm64 linux-x64 linux-arm64 win32-x64; do
        if node -e "JSON.parse(require('fs').readFileSync('npm/platforms/$platform/package.json'))" 2>/dev/null; then
            log_success "Platform package.json valid: $platform"
        else
            log_error "Platform package.json invalid: $platform"
        fi
    done
}

# Test 7: GitHub Actions workflow
test_github_actions() {
    log_info "Testing GitHub Actions workflow..."
    
    cd "$REPO_ROOT"
    
    # Check workflow file exists
    if [ -f ".github/workflows/release.yml" ]; then
        log_success "Release workflow exists"
        
        # Basic YAML validation
        if command -v python3 >/dev/null 2>&1; then
            # Check if PyYAML is installed
            if python3 -c "import yaml" 2>/dev/null; then
                if python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))" 2>/dev/null; then
                    log_success "Release workflow YAML valid"
                else
                    log_error "Release workflow YAML invalid"
                fi
            else
                log_skip "YAML validation (PyYAML module not installed)"
            fi
        elif command -v ruby >/dev/null 2>&1; then
            # Try Ruby as fallback
            if ruby -e "require 'yaml'; YAML.load_file('.github/workflows/release.yml')" 2>/dev/null; then
                log_success "Release workflow YAML valid"
            else
                log_error "Release workflow YAML invalid"
            fi
        else
            log_skip "YAML validation (no parser available)"
        fi
    else
        log_error "Release workflow missing"
    fi
}

# Test 8: Binary execution (current platform)
test_binary_execution() {
    log_info "Testing binary execution on current platform..."
    
    if [ ! -d "$REPO_ROOT/dist" ]; then
        log_skip "Binary execution test (no binaries built)"
        return 0
    fi
    
    cd "$REPO_ROOT"
    
    # Detect current platform
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # Map to GoReleaser naming
    case "$os" in
        darwin) os="darwin" ;;
        linux) os="linux" ;;
        *) log_skip "Binary execution test (unsupported OS: $os)"; return 0 ;;
    esac
    
    case "$arch" in
        x86_64) arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *) log_skip "Binary execution test (unsupported arch: $arch)"; return 0 ;;
    esac
    
    # Find the binary
    local binary_path
    for dir in dist/*; do
        if [[ -d "$dir" && "$dir" == *"${os}_${arch}"* ]]; then
            binary_path="$dir/claude-gate"
            break
        fi
    done
    
    if [ -z "$binary_path" ] || [ ! -f "$binary_path" ]; then
        log_error "Binary not found for $os/$arch"
        return 1
    fi
    
    # Test execution
    if "$binary_path" version >/dev/null 2>&1; then
        log_success "Binary executes: $binary_path"
        
        # Test help command
        if "$binary_path" --help >/dev/null 2>&1; then
            log_success "Binary help command works"
        else
            log_error "Binary help command failed"
        fi
    else
        log_error "Binary execution failed: $binary_path"
    fi
}

# Test 9: Error handling
test_error_handling() {
    log_info "Testing error handling..."
    
    cd "$REPO_ROOT/npm"
    
    # Test unsupported platform error
    local error_output=$(node -e "
        const originalPlatform = process.platform;
        Object.defineProperty(process, 'platform', {value: 'freebsd', configurable: true});
        const {getPlatform} = require('./scripts/install.js');
        try { getPlatform(); } catch(e) { console.log('ERROR_CAUGHT'); }
    " 2>&1)
    
    if [[ "$error_output" == *"ERROR_CAUGHT"* ]]; then
        log_success "Unsupported platform error handling"
    else
        log_error "Unsupported platform error not caught"
    fi
}

# Test 10: Scripts permissions
test_script_permissions() {
    log_info "Testing script permissions..."
    
    cd "$REPO_ROOT"
    
    local scripts=(
        "scripts/build-release.sh"
        "scripts/test-npm-local.sh"
        "scripts/setup-npm-auth.sh"
        "scripts/update-version.sh"
    )
    
    for script in "${scripts[@]}"; do
        if [ -x "$script" ]; then
            log_success "Script executable: $script"
        else
            log_error "Script not executable: $script"
        fi
    done
}

# Generate test report
generate_report() {
    echo ""
    echo "========================================"
    echo "Test Report"
    echo "========================================"
    echo "Total Tests: $((PASSED + FAILED + SKIPPED))"
    echo -e "Passed: ${GREEN}$PASSED${NC}"
    echo -e "Failed: ${RED}$FAILED${NC}"
    echo -e "Skipped: ${YELLOW}$SKIPPED${NC}"
    echo ""
    
    if [ $FAILED -gt 0 ]; then
        echo "Failed Tests:"
        for i in "${!TEST_NAMES[@]}"; do
            if [ "${TEST_STATUSES[$i]}" = "FAILED" ]; then
                echo -e "  ${RED}✗${NC} ${TEST_NAMES[$i]}"
            fi
        done
        echo ""
    fi
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}Some tests failed!${NC}"
        return 1
    fi
}

# Main test execution
main() {
    echo "Claude Gate - Cross-Platform Build & NPM Distribution Tests"
    echo "==========================================================="
    echo ""
    
    check_prerequisites
    echo ""
    
    # Run all tests
    test_go_build
    test_goreleaser_build
    test_npm_package_structure
    test_npm_install_script
    test_npm_local_package
    test_package_json_validation
    test_github_actions
    test_binary_execution
    test_error_handling
    test_script_permissions
    
    # Generate report
    generate_report
}

# Run tests
main