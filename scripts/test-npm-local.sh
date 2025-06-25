#!/bin/bash
set -e

# Script to test NPM package locally without publishing

echo "🧪 Testing NPM package locally..."

# Get current directory
REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$REPO_ROOT"

# Build binaries first
if [ ! -d "dist" ]; then
    echo "No dist directory found. Building binaries first..."
    ./scripts/build-release.sh
fi

# Create temporary directory for testing
TEST_DIR=$(mktemp -d)
echo "Using test directory: $TEST_DIR"

# Function to clean up on exit
cleanup() {
    echo "Cleaning up..."
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Copy NPM package files
echo "Copying NPM package files..."
cp -r npm/* "$TEST_DIR/"

# Create platform packages directory
mkdir -p "$TEST_DIR/node_modules/@claude-gate"

# Extract and copy binaries to simulate platform packages
echo "Preparing platform packages..."

# Find the version from the archive names (e.g., 0.0.1-next)
VERSION=$(ls dist/*.tar.gz 2>/dev/null | head -1 | sed -E 's/.*claude-gate_([^_]+)_.*/\1/' || echo "snapshot")

# macOS Intel
DARWIN_X64_ARCHIVE=$(ls dist/claude-gate_*_Darwin_x86_64.tar.gz 2>/dev/null | head -1)
if [ -f "$DARWIN_X64_ARCHIVE" ]; then
    mkdir -p "$TEST_DIR/node_modules/@claude-gate/darwin-x64"
    tar -xzf "$DARWIN_X64_ARCHIVE" -C "$TEST_DIR/node_modules/@claude-gate/darwin-x64" --strip-components=1
    mv "$TEST_DIR/node_modules/@claude-gate/darwin-x64/claude-gate" "$TEST_DIR/node_modules/@claude-gate/darwin-x64/bin" 2>/dev/null || true
fi

# macOS ARM
DARWIN_ARM64_ARCHIVE=$(ls dist/claude-gate_*_Darwin_arm64.tar.gz 2>/dev/null | head -1)
if [ -f "$DARWIN_ARM64_ARCHIVE" ]; then
    mkdir -p "$TEST_DIR/node_modules/@claude-gate/darwin-arm64"
    tar -xzf "$DARWIN_ARM64_ARCHIVE" -C "$TEST_DIR/node_modules/@claude-gate/darwin-arm64" --strip-components=1
    mv "$TEST_DIR/node_modules/@claude-gate/darwin-arm64/claude-gate" "$TEST_DIR/node_modules/@claude-gate/darwin-arm64/bin" 2>/dev/null || true
fi

# Linux x64
LINUX_X64_ARCHIVE=$(ls dist/claude-gate_*_Linux_x86_64.tar.gz 2>/dev/null | head -1)
if [ -f "$LINUX_X64_ARCHIVE" ]; then
    mkdir -p "$TEST_DIR/node_modules/@claude-gate/linux-x64"
    tar -xzf "$LINUX_X64_ARCHIVE" -C "$TEST_DIR/node_modules/@claude-gate/linux-x64" --strip-components=1
    mv "$TEST_DIR/node_modules/@claude-gate/linux-x64/claude-gate" "$TEST_DIR/node_modules/@claude-gate/linux-x64/bin" 2>/dev/null || true
fi

# Linux ARM64
LINUX_ARM64_ARCHIVE=$(ls dist/claude-gate_*_Linux_arm64.tar.gz 2>/dev/null | head -1)
if [ -f "$LINUX_ARM64_ARCHIVE" ]; then
    mkdir -p "$TEST_DIR/node_modules/@claude-gate/linux-arm64"
    tar -xzf "$LINUX_ARM64_ARCHIVE" -C "$TEST_DIR/node_modules/@claude-gate/linux-arm64" --strip-components=1
    mv "$TEST_DIR/node_modules/@claude-gate/linux-arm64/claude-gate" "$TEST_DIR/node_modules/@claude-gate/linux-arm64/bin" 2>/dev/null || true
fi

# # Windows - commented out as Windows is not supported
# WINDOWS_ARCHIVE=$(ls dist/claude-gate_*_Windows_x86_64.zip 2>/dev/null | head -1)
# if [ -f "$WINDOWS_ARCHIVE" ]; then
#     mkdir -p "$TEST_DIR/node_modules/@claude-gate/win32-x64"
#     unzip -q "$WINDOWS_ARCHIVE" -d "$TEST_DIR/node_modules/@claude-gate/win32-x64"
#     mv "$TEST_DIR/node_modules/@claude-gate/win32-x64/claude-gate.exe" "$TEST_DIR/node_modules/@claude-gate/win32-x64/bin.exe" 2>/dev/null || true
# fi

# Run install script
echo ""
echo "Running install script..."
cd "$TEST_DIR"
node scripts/install.js

# Test the binary
echo ""
echo "Testing installed binary..."
if [ -f "bin/claude-gate" ]; then
    echo "✅ Binary wrapper found"
    
    # Try running it
    echo "Testing command execution..."
    ./bin/claude-gate version || echo "Note: Binary execution test failed (this is expected if built for different platform)"
else
    echo "❌ Binary wrapper not found!"
    ls -la bin/
fi

# Create npm package for testing
echo ""
echo "Creating test package..."
npm pack

echo ""
echo "✅ Local NPM package test complete!"
echo ""
echo "To test global installation:"
echo "  npm install -g $TEST_DIR/claude-gate-0.1.0.tgz"
echo ""
echo "Package file created at: $TEST_DIR/claude-gate-0.1.0.tgz"