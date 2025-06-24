#!/bin/bash
set -e

# Script to update version across all files

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 0.2.0"
    exit 1
fi

VERSION=$1

# Validate version format
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.\-]+)?$ ]]; then
    echo "Error: Invalid version format. Use semantic versioning (e.g., 1.0.0, 1.0.0-beta.1)"
    exit 1
fi

echo "Updating version to $VERSION..."

# Update Go source
echo "Updating cmd/claude-gate/main.go..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/var version = \".*\"/var version = \"$VERSION\"/" cmd/claude-gate/main.go
else
    # Linux
    sed -i "s/var version = \".*\"/var version = \"$VERSION\"/" cmd/claude-gate/main.go
fi

# Update main NPM package
echo "Updating npm/package.json..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" npm/package.json
    sed -i '' "s/\"@claude-gate\/darwin-x64\": \".*\"/\"@claude-gate\/darwin-x64\": \"$VERSION\"/" npm/package.json
    sed -i '' "s/\"@claude-gate\/darwin-arm64\": \".*\"/\"@claude-gate\/darwin-arm64\": \"$VERSION\"/" npm/package.json
    sed -i '' "s/\"@claude-gate\/linux-x64\": \".*\"/\"@claude-gate\/linux-x64\": \"$VERSION\"/" npm/package.json
    sed -i '' "s/\"@claude-gate\/linux-arm64\": \".*\"/\"@claude-gate\/linux-arm64\": \"$VERSION\"/" npm/package.json
    sed -i '' "s/\"@claude-gate\/win32-x64\": \".*\"/\"@claude-gate\/win32-x64\": \"$VERSION\"/" npm/package.json
else
    sed -i "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" npm/package.json
    sed -i "s/\"@claude-gate\/darwin-x64\": \".*\"/\"@claude-gate\/darwin-x64\": \"$VERSION\"/" npm/package.json
    sed -i "s/\"@claude-gate\/darwin-arm64\": \".*\"/\"@claude-gate\/darwin-arm64\": \"$VERSION\"/" npm/package.json
    sed -i "s/\"@claude-gate\/linux-x64\": \".*\"/\"@claude-gate\/linux-x64\": \"$VERSION\"/" npm/package.json
    sed -i "s/\"@claude-gate\/linux-arm64\": \".*\"/\"@claude-gate\/linux-arm64\": \"$VERSION\"/" npm/package.json
    sed -i "s/\"@claude-gate\/win32-x64\": \".*\"/\"@claude-gate\/win32-x64\": \"$VERSION\"/" npm/package.json
fi

# Update platform packages
for platform in darwin-x64 darwin-arm64 linux-x64 linux-arm64 win32-x64; do
    echo "Updating npm/platforms/$platform/package.json..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" "npm/platforms/$platform/package.json"
    else
        sed -i "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" "npm/platforms/$platform/package.json"
    fi
done

# Update npm/index.js
echo "Updating npm/index.js..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/version: '.*'/version: '$VERSION'/" npm/index.js
else
    sed -i "s/version: '.*'/version: '$VERSION'/" npm/index.js
fi

echo ""
echo "✅ Version updated to $VERSION in all files!"
echo ""
echo "Next steps:"
echo "1. Review changes: git diff"
echo "2. Run tests: make test"
echo "3. Build snapshot: make snapshot"
echo "4. Test NPM package: make npm-test"
echo "5. Commit: git add -A && git commit -m \"chore: bump version to $VERSION\""
echo "6. Tag: git tag -a v$VERSION -m \"Release v$VERSION\""
echo "7. Push: git push origin main && git push origin v$VERSION"