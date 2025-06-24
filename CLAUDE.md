# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Claude Gate is a high-performance Go OAuth proxy for Anthropic's Claude API that enables FREE Claude usage for Pro/Max subscribers by identifying as "Claude Code" (Anthropic's official CLI). This is the official Go port of the original Python implementation.

**Key Features:**
- OAuth 2.0 PKCE authentication flow
- Secure token storage with OS keychain integration (planned)
- Interactive TUI dashboard for monitoring
- Cross-platform support (macOS, Linux, Windows)
- Multiple distribution methods (NPM, Homebrew, direct binary)

## Common Commands

### Building
```bash
make build        # Build for current platform
make snapshot     # Build all platforms (uses GoReleaser)
make install      # Install to ~/bin
```

### Testing
```bash
make test         # Run Go tests with race detection
make test-all     # Comprehensive test suite (unit, integration, edge cases)
make npm-test     # Test NPM package locally
go test -v ./...  # Quick test during development
```

### Running
```bash
claude-gate start --host 127.0.0.1 --port 8000  # Start proxy server
claude-gate dashboard                            # Start with interactive dashboard
claude-gate auth login                           # Authenticate with Claude
```

### Releasing
```bash
make release VERSION=0.2.0  # Create new release
./scripts/update-version.sh # Update version in all files
```

## Architecture

The codebase follows clean architecture principles with clear separation of concerns:

### Core Components

1. **CLI Layer** (`cmd/claude-gate/`)
   - Uses Kong framework for command parsing
   - Entry point for all operations

2. **Auth Package** (`internal/auth/`)
   - OAuth 2.0 PKCE implementation
   - Token storage and management
   - Browser automation for login flow

3. **Proxy Package** (`internal/proxy/`)
   - HTTP proxy server implementation
   - Request/response transformation
   - Enhanced server with monitoring capabilities

4. **UI Package** (`internal/ui/`)
   - Bubble Tea-based TUI components
   - Interactive dashboard for monitoring
   - Reusable components (spinner, progress, styles)

### Request Flow
1. Client connects to local proxy (default: 127.0.0.1:8000)
2. Proxy validates authentication token
3. Request transformed to identify as "Claude Code"
4. Forwarded to Claude API with OAuth credentials
5. Response streamed back to client

### Security Model
- OAuth 2.0 PKCE flow for authentication
- Tokens stored securely (keychain integration planned)
- Local-only proxy binding by default
- Optional proxy authentication token for additional security

## Testing Strategy

The project uses Test-Driven Development (TDD) with comprehensive test coverage:

- **Unit Tests**: Alongside source files (`*_test.go`)
- **Integration Tests**: `internal/test/integration/`
- **E2E Tests**: `internal/test/e2e/`
- **Cross-Platform Tests**: Via Docker containers
- **NPM Package Tests**: Validates installation and binary selection

Always write tests before implementing features. Use testify for assertions.

## Development Workflow

1. **Feature Development**
   - Create feature branch
   - Write tests first (TDD)
   - Implement feature
   - Run `make test-all`
   - Submit PR with tests passing

2. **Bug Fixes**
   - Write test to reproduce bug
   - Fix the issue
   - Ensure all tests pass
   - Document in TROUBLESHOOTING.md if user-facing

3. **Documentation Updates**
   - Update relevant docs in `docs/`
   - Keep README.md concise
   - Add troubleshooting entries as needed

## NPM Package Management

The project distributes platform-specific binaries via NPM:

- Main package: `npm/package.json`
- Platform packages: `npm/platforms/*/package.json`
- Installation scripts: `npm/scripts/`
- Binary wrappers: `npm/bin/`

Test NPM changes with: `make npm-test`

## Key Dependencies

- **Kong**: CLI framework for command parsing
- **Bubble Tea**: Terminal UI framework
- **Lipgloss**: Terminal styling
- **Testify**: Testing assertions and mocks
- **GoReleaser**: Multi-platform release automation

## Release Process

1. Update version: `./scripts/update-version.sh`
2. Create release: `make release VERSION=x.y.z`
3. Push tags: `git push origin main && git push origin vx.y.z`
4. GitHub Actions automatically:
   - Builds binaries for all platforms
   - Publishes to NPM registry
   - Creates GitHub release

## Important Patterns

- Use `internal/` for private packages (Go convention)
- Follow clean architecture: separate concerns between packages
- Use interfaces for testability and flexibility
- Implement context propagation for cancellation
- Handle errors explicitly, never ignore them
- Use structured logging with clear messages
- Write self-documenting code with meaningful names

## Contribution Guidelines

See CONTRIBUTING.md for detailed guidelines. Key points:
- All code must have tests
- Follow Go best practices and idioms
- Use `go fmt` and `go vet`
- Document public APIs
- Sign commits with GPG key
- Squash commits before merging