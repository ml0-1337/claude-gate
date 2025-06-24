# Development Guide

This guide covers everything you need to know to set up a development environment for Claude Gate.

## Prerequisites

- **Go 1.22 or later** - [Install Go](https://golang.org/dl/)
- **Node.js 18 or later** - For NPM package testing
- **Git** - Version control
- **Make** - Build automation
- **GoReleaser** - For release builds (optional)

### macOS Setup

```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install prerequisites
brew install go node git
brew install goreleaser/tap/goreleaser
```

### Linux Setup

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang nodejs npm git make

# Fedora/RHEL
sudo dnf install golang nodejs npm git make
```

### Windows Setup

1. Install [Go](https://golang.org/dl/) from the official website
2. Install [Node.js](https://nodejs.org/) from the official website
3. Install [Git for Windows](https://git-scm.com/download/win)
4. Install [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

## Setting Up the Development Environment

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/claude-gate.git
cd claude-gate

# Add upstream remote
git remote add upstream https://github.com/anthropics/claude-gate.git
```

### 2. Install Dependencies

```bash
# Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
```

### 3. Build and Test

```bash
# Build the binary
make build

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow the project structure:

```
claude-gate/
├── cmd/claude-gate/      # Main application entry point
├── internal/             # Internal packages
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration management
│   ├── proxy/           # Proxy server implementation
│   └── ui/              # Terminal UI components
├── pkg/                  # Public packages
└── scripts/             # Build and utility scripts
```

### 3. Write Tests

All new features should include tests:

```go
// internal/auth/oauth_test.go
func TestOAuthFlow(t *testing.T) {
    // Your test implementation
}
```

Run tests for a specific package:

```bash
go test ./internal/auth/...
```

### 4. Run Local Development Server

```bash
# Build and run with debug logging
go run cmd/claude-gate/main.go start --log-level debug
```

### 5. Test NPM Package Locally

```bash
# Build for all platforms
make build-all

# Test NPM package installation
cd npm
npm link
cd ..
claude-gate version
```

## Code Style Guide

### Go Code Style

We follow standard Go conventions with some additions:

1. **Format code** with `gofumpt` (stricter than `gofmt`):
   ```bash
   gofumpt -w .
   ```

2. **Import ordering** with `goimports`:
   ```bash
   goimports -w .
   ```

3. **Linting** with `golangci-lint`:
   ```bash
   golangci-lint run
   ```

### Code Organization

- Keep functions small and focused
- Use meaningful variable and function names
- Add comments for exported functions
- Group related functionality in packages

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to authenticate: %w", err)
}

// Bad
if err != nil {
    return err
}
```

## Debugging

### Enable Debug Logging

```bash
# Via environment variable
export CLAUDE_GATE_LOG_LEVEL=debug
claude-gate start

# Or via flag
claude-gate start --log-level debug
```

### Using Delve Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug cmd/claude-gate/main.go -- start --log-level debug
```

### Common Debugging Commands

```bash
# Check binary location
which claude-gate

# View configuration
claude-gate config --show

# Test authentication
claude-gate auth status

# View stored tokens (macOS)
ls -la ~/Library/Application\ Support/claude-gate/
```

## Testing

### Unit Tests

```bash
# Run all unit tests
make test

# Run tests for specific package
go test ./internal/proxy/...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

### Integration Tests

```bash
# Run integration tests
make test-integration

# Run specific integration test
go test -tags=integration ./internal/test/integration/...
```

### End-to-End Tests

```bash
# Run e2e tests
make test-e2e

# Test cross-platform builds
./scripts/test-all.sh
```

## Building

### Local Build

```bash
# Build for current platform
make build

# Build with version info
make build VERSION=v1.2.3
```

### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Build for specific platform
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build
GOOS=windows GOARCH=amd64 make build
```

### Release Builds

```bash
# Create a release with GoReleaser
goreleaser release --snapshot --clean

# Build NPM packages
make build-npm
```

## Troubleshooting Development Issues

### Go Module Issues

```bash
# Clean module cache
go clean -modcache

# Update dependencies
go mod tidy
go mod vendor
```

### Build Failures

```bash
# Clean build artifacts
make clean

# Rebuild everything
make clean build
```

### Test Failures

```bash
# Run tests with more output
go test -v -count=1 ./...

# Skip test cache
go test -count=1 ./...
```

## Contributing Your Changes

1. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

2. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**:
   - Go to GitHub and create a PR from your fork
   - Fill out the PR template
   - Wait for CI to pass

See our [Contributing Guide](./contributing.md) for more details on the PR process.

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Project Architecture](../architecture/overview.md)

---

[← Guides](../README.md#guides) | [Documentation Home](../README.md) | [Contributing →](./contributing.md)