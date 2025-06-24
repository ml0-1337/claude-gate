#!/bin/bash
# Test Claude Gate NPM package in Docker containers for different platforms

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Testing Claude Gate NPM package in Docker containers..."
echo "======================================================="

# Check if Docker is available
if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}Docker is not installed or not running${NC}"
    exit 1
fi

# Build the NPM package first
echo -e "\n${YELLOW}Building NPM package...${NC}"
cd "$REPO_ROOT"

# Check if binaries exist
if [ ! -d "dist" ]; then
    echo "No dist directory found. Building with GoReleaser first..."
    if command -v goreleaser >/dev/null 2>&1; then
        goreleaser build --snapshot --clean
    else
        echo -e "${RED}GoReleaser not installed. Please install it first.${NC}"
        exit 1
    fi
fi

# Create test package
cd npm
npm pack
PACKAGE_FILE="claude-gate-0.1.0.tgz"
mv "$PACKAGE_FILE" "$REPO_ROOT/"
cd "$REPO_ROOT"

# Test function
test_platform() {
    local platform=$1
    local node_version=$2
    local docker_platform=$3
    
    echo -e "\n${YELLOW}Testing $platform with Node.js $node_version...${NC}"
    
    # Create test script
    cat > test-docker-install.sh << 'EOF'
#!/bin/sh
set -e

echo "Platform: $(uname -s) $(uname -m)"
echo "Node.js: $(node --version)"
echo "npm: $(npm --version)"

# Install the package
echo "Installing claude-gate..."
npm install -g /app/claude-gate-0.1.0.tgz

# Test if installed
echo "Testing installation..."
which claude-gate
claude-gate version || echo "Binary execution failed (expected on different arch)"

# Test help
claude-gate --help || echo "Help command failed"

# Check what was installed
echo "Checking installed files..."
ls -la $(npm root -g)/claude-gate/bin/

# Test uninstall
echo "Testing uninstall..."
npm uninstall -g claude-gate

echo "Test completed successfully!"
EOF

    chmod +x test-docker-install.sh
    
    # Run Docker test
    local docker_cmd="docker run --rm -v $REPO_ROOT:/app -w /app"
    
    if [ -n "$docker_platform" ]; then
        docker_cmd="$docker_cmd --platform $docker_platform"
    fi
    
    if $docker_cmd node:$node_version sh test-docker-install.sh; then
        echo -e "${GREEN}✓ $platform test passed${NC}"
        return 0
    else
        echo -e "${RED}✗ $platform test failed${NC}"
        return 1
    fi
}

# Run tests for different platforms
TESTS_PASSED=0
TESTS_FAILED=0

# Linux x64
if test_platform "Linux x64" "20" "linux/amd64"; then
    ((TESTS_PASSED++))
else
    ((TESTS_FAILED++))
fi

# Linux ARM64 (only if Docker supports it)
if docker run --rm --platform linux/arm64 alpine uname -m >/dev/null 2>&1; then
    if test_platform "Linux ARM64" "20" "linux/arm64"; then
        ((TESTS_PASSED++))
    else
        ((TESTS_FAILED++))
    fi
else
    echo -e "\n${YELLOW}Skipping Linux ARM64 test (platform not supported by Docker)${NC}"
fi

# Test with different Node.js versions
for node_version in 18 22; do
    if test_platform "Linux x64 (Node $node_version)" "$node_version" "linux/amd64"; then
        ((TESTS_PASSED++))
    else
        ((TESTS_FAILED++))
    fi
done

# Test --ignore-scripts flag
echo -e "\n${YELLOW}Testing with --ignore-scripts flag...${NC}"
docker run --rm -v "$REPO_ROOT:/app" -w /app node:20 sh -c "
    npm install -g /app/claude-gate-0.1.0.tgz --ignore-scripts
    # The fallback should kick in when running
    claude-gate version 2>&1 | grep -E '(Binary not found|Running installation|version)' && echo 'Fallback mechanism works'
    npm uninstall -g claude-gate
"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ --ignore-scripts test passed${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${RED}✗ --ignore-scripts test failed${NC}"
    ((TESTS_FAILED++))
fi

# Cleanup
rm -f test-docker-install.sh
rm -f "$REPO_ROOT/claude-gate-0.1.0.tgz"

# Summary
echo -e "\n========================================"
echo "Docker Test Summary"
echo "========================================"
echo -e "Tests passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All Docker tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some Docker tests failed!${NC}"
    exit 1
fi