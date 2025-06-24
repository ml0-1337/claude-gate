#!/bin/bash
set -e

# Build script for testing GoReleaser locally

echo "🔨 Building claude-gate with GoReleaser..."

# Clean previous builds
rm -rf dist/

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "❌ GoReleaser is not installed!"
    echo "Install with: brew install goreleaser"
    exit 1
fi

# Build snapshot (without publishing)
echo "Building snapshot release..."
goreleaser build --snapshot --clean

echo ""
echo "✅ Build complete! Binaries are in ./dist/"
echo ""
echo "📦 Built archives:"
ls -la dist/*.tar.gz dist/*.zip 2>/dev/null | awk '{print "  - " $9 " (" $5 " bytes)"}'

echo ""
echo "🔍 Binary sizes:"
for binary in dist/*/claude-gate*; do
    if [[ -f "$binary" && -x "$binary" ]]; then
        size=$(ls -lh "$binary" | awk '{print $5}')
        platform=$(echo "$binary" | grep -oE '(darwin|linux|windows)_[^/]+')
        echo "  - $platform: $size"
    fi
done

echo ""
echo "To test a specific binary:"
echo "  ./dist/claude-gate_darwin_arm64/claude-gate version"
echo ""
echo "To create a real release:"
echo "  1. git tag v0.1.0"
echo "  2. git push origin v0.1.0"
echo "  3. GitHub Actions will handle the rest"