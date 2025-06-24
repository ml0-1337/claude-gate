# Claude Gate Testing Guide

This guide explains how to thoroughly test the cross-platform builds and NPM distribution system.

## Quick Start

Run all tests:
```bash
make test-all
```

## Test Categories

### 1. Unit Tests

**Go Tests:**
```bash
make test
```

**NPM Package Tests:**
```bash
cd npm
node scripts/install.test.js
```

### 2. Build Tests

**Local Build:**
```bash
make build
./claude-gate version
```

**Cross-Platform Build:**
```bash
make snapshot
ls -la dist/
```

### 3. NPM Package Tests

**Local Package Test:**
```bash
make npm-test
```

This will:
- Build all platform binaries
- Create test NPM package
- Simulate installation
- Verify the package works

### 4. Integration Tests

**Comprehensive Test Suite:**
```bash
make test-all
```

This runs:
- Go unit tests
- GoReleaser build tests
- NPM package structure validation
- Platform detection tests
- Binary execution tests
- Error handling tests

### 5. Docker Tests

**Multi-Platform Docker Tests:**
```bash
make test-docker
```

Tests installation on:
- Linux x64
- Linux ARM64 (if supported)
- Different Node.js versions

### 6. Edge Case Tests

**Edge Cases and Error Scenarios:**
```bash
make test-edge
```

Tests:
- Missing binaries
- Permission errors
- Concurrent installations
- Signal handling
- Version mismatches

## Manual Testing

### Test NPM Installation Locally

1. **Create package:**
   ```bash
   cd npm
   npm pack
   ```

2. **Install globally:**
   ```bash
   npm install -g claude-gate-0.1.0.tgz
   ```

3. **Test commands:**
   ```bash
   claude-gate version
   claude-gate --help
   claude-gate auth status
   ```

4. **Uninstall:**
   ```bash
   npm uninstall -g claude-gate
   ```

### Test --ignore-scripts

```bash
npm install -g claude-gate-0.1.0.tgz --ignore-scripts
claude-gate version  # Should trigger fallback installation
```

### Test Platform Detection

```bash
cd npm
node -e "const {getPlatform} = require('./scripts/install.js'); console.log(getPlatform())"
```

## CI/CD Testing

### GitHub Actions

The project includes automated tests that run on:
- Every push to main/develop
- Every pull request

Test matrix includes:
- OS: Ubuntu, macOS, Windows
- Go: 1.22, 1.23
- Node.js: 18, 20, 22

### Local CI Testing

Test GitHub Actions locally using act:
```bash
brew install act
act -j test-go
```

## Platform-Specific Testing

### macOS
```bash
# Test on Intel Mac
GOARCH=amd64 make build

# Test on Apple Silicon
GOARCH=arm64 make build
```

### Linux
Use Docker for testing:
```bash
docker run --rm -it -v $(pwd):/app -w /app golang:1.22 bash
# Inside container:
make test
```

### Windows
Use Windows VM or:
```bash
GOOS=windows make build
# Test .exe in Windows environment
```

## Release Testing

### Pre-Release Checklist

1. **Version Consistency:**
   ```bash
   ./scripts/update-version.sh 0.1.0
   git diff  # Review changes
   ```

2. **Build Test:**
   ```bash
   make snapshot
   ```

3. **NPM Package Test:**
   ```bash
   make npm-test
   ```

4. **Tag and Release (Dry Run):**
   ```bash
   goreleaser release --snapshot --skip-publish --clean
   ```

## Debugging Failed Tests

### Verbose Output
```bash
# Go tests with verbose output
go test -v ./...

# NPM install with debug
DEBUG=* npm install -g claude-gate-0.1.0.tgz
```

### Check Logs
- GitHub Actions: Check workflow run logs
- Local: Check console output and error messages

### Common Issues

1. **GoReleaser not found:**
   ```bash
   brew install goreleaser
   ```

2. **Permission denied:**
   ```bash
   sudo npm install -g claude-gate
   # Or use nvm to avoid sudo
   ```

3. **Binary not found:**
   - Check if platform is supported
   - Verify dist/ directory has binaries
   - Check npm package structure

## Performance Testing

### Binary Size Check
```bash
make snapshot
find dist -name "claude-gate*" -exec ls -lh {} \;
```

Target: <15MB per platform

### Installation Speed
```bash
time npm install -g claude-gate
```

Target: <30 seconds on decent connection

## Security Testing

### Checksum Verification
```bash
cd dist
sha256sum -c checksums.txt
```

### Dependency Audit
```bash
cd npm
npm audit
```

## Reporting Issues

When reporting test failures, include:
1. OS and architecture: `uname -a`
2. Go version: `go version`
3. Node.js version: `node --version`
4. NPM version: `npm --version`
5. Complete error output
6. Steps to reproduce