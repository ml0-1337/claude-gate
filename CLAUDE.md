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

## Quick Navigation

### Core Implementation Files
- Main entry: `cmd/claude-gate/main.go`
- OAuth client: `internal/auth/client.go`
- Proxy handler: `internal/proxy/handler.go`
- Dashboard: `internal/ui/dashboard/dashboard.go`

### Configuration Files
- Go module: `go.mod`
- Build config: `.goreleaser.yml`
- CI/CD: `.github/workflows/release.yml`
- NPM package: `npm/package.json`

### Test Files
- Unit tests: `*_test.go` (alongside source)
- Integration: `internal/test/integration/*_test.go`
- E2E: `internal/test/e2e/*_test.go`
- Test helpers: `internal/test/helpers/helpers.go`

### Scripts
- Version update: `scripts/update-version.sh`
- NPM test: `scripts/test-npm-local.sh`
- Add co-author: `.claude/scripts/add-claude-coauthor.sh`

## Common Commands

### Building
```bash
make build        # Build for current platform
make snapshot     # Build all platforms (uses GoReleaser)
make install      # Install to ~/bin
```

### Testing
```bash
make test              # Run unit tests with coverage
make test-unit         # Run unit tests only (short mode)
make test-integration  # Run integration tests
make test-e2e          # Run end-to-end tests
make test-all          # Run all test types
make npm-test          # Test NPM package locally
go test -v ./...       # Quick test during development
```

### Running
```bash
claude-gate start --host 127.0.0.1 --port 5789  # Start proxy server
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
1. Client connects to local proxy (default: 127.0.0.1:5789)
2. Proxy validates authentication token
3. Request transformed to identify as "Claude Code"
4. Forwarded to Claude API with OAuth credentials
5. Response streamed back to client

### Security Model
- OAuth 2.0 PKCE flow for authentication
- Tokens stored securely (keychain integration planned)
- Local-only proxy binding by default
- Optional proxy authentication token for additional security

## Project Structure

```
claude-gate/
├── cmd/
│   └── claude-gate/          # CLI application entry point
│       ├── main.go          # Main entry point with Kong CLI setup
│       └── auth_storage.go  # Storage-related CLI commands
├── internal/                 # Private application code (Go convention)
│   ├── auth/                # Authentication & token management
│   │   ├── client.go        # OAuth client implementation
│   │   ├── oauth.go         # OAuth flow logic
│   │   ├── storage_*.go     # Various storage backends
│   │   └── *_test.go        # Unit tests
│   ├── config/              # Configuration management
│   ├── proxy/               # Proxy server implementation
│   │   ├── handler.go       # Main proxy handler
│   │   ├── server.go        # Enhanced server with monitoring
│   │   └── transformer.go   # Request/response transformation
│   ├── test/                # Test infrastructure
│   │   ├── integration/     # Integration tests (build tag: integration)
│   │   ├── e2e/            # End-to-end tests (build tag: e2e)
│   │   ├── helpers/        # Shared test utilities
│   │   └── testdata/       # Test fixtures
│   └── ui/                  # Terminal UI components
│       ├── components/      # Reusable UI components
│       ├── dashboard/       # Interactive dashboard
│       ├── styles/         # Terminal styling (Lipgloss)
│       └── utils/          # UI utilities
├── docs/                    # Project documentation
│   ├── architecture/       # Architecture decisions and diagrams
│   ├── deployment/         # Deployment guides
│   ├── getting-started/    # User guides
│   └── testing/           # Testing documentation
├── npm/                    # NPM package distribution
│   ├── package.json       # Main NPM package
│   ├── platforms/         # Platform-specific packages
│   └── scripts/           # Installation scripts
├── scripts/               # Build and utility scripts
├── .github/              # GitHub Actions workflows
│   └── workflows/        # CI/CD pipelines
├── .claude/              # Claude-specific files
│   ├── todos/           # Active todo files
│   ├── archive/         # Completed todos
│   └── scripts/         # Claude utility scripts
├── Makefile             # Build automation
├── go.mod              # Go module definition
└── .goreleaser.yml     # GoReleaser configuration
```

### File Naming Conventions
- Test files: `*_test.go` alongside source files
- Integration tests: `*_integration_test.go` with build tags
- E2E tests: `*_e2e_test.go` with build tags
- Mock implementations: `mock_*.go`
- Interfaces: Often in the same file as primary implementation

### Key Files
- Entry point: `cmd/claude-gate/main.go`
- OAuth config: `internal/auth/oauth.go` (contains client ID)
- Proxy handler: `internal/proxy/handler.go`
- Dashboard UI: `internal/ui/dashboard/dashboard.go`
- Storage factory: `internal/auth/storage_factory.go`

## Testing Strategy

The project uses Test-Driven Development (TDD) with comprehensive test coverage:

- **Unit Tests**: Alongside source files (`*_test.go`)
- **Integration Tests**: `internal/test/integration/`
- **E2E Tests**: `internal/test/e2e/`
- **Cross-Platform Tests**: Via Docker containers
- **NPM Package Tests**: Validates installation and binary selection

Always write tests before implementing features. Use testify for assertions.

## TDD Requirements (MANDATORY)

**IMPORTANT**: This project follows the strict TDD workflow defined in the system CLAUDE.md. Test-first development is MANDATORY, not optional.

### The Red-Green-Refactor-Commit (RGRC) Cycle

1. **Red Phase**: Write failing tests FIRST
   ```bash
   # Write test in *_test.go file
   go test -v ./path/to/package -run TestNewFeature
   # Verify test fails before proceeding
   ```

2. **Green Phase**: Write MINIMUM code to pass
   ```bash
   # Implement just enough to make tests pass
   go test -v ./path/to/package -run TestNewFeature
   # All tests should now pass
   ```

3. **Refactor Phase**: Improve code quality
   ```bash
   # Clean up implementation while keeping tests green
   make test  # Run full test suite with race detection
   ```

4. **Commit Phase**: Save progress
   ```bash
   git add -A
   git commit -m "test: Add tests for [feature]"
   git commit -m "feat: Implement [feature] to pass tests"
   ```

### When TDD is MANDATORY

- **New Features**: Write acceptance tests first
- **Bug Fixes**: Write test that reproduces bug first
- **Refactoring**: Ensure tests exist before changing
- **API Changes**: Contract tests before implementation

### Go-Specific TDD Patterns

```go
// Example: Table-driven test (write BEFORE implementation)
func TestAuthenticateRequest(t *testing.T) {
    tests := []struct {
        name    string
        token   string
        want    bool
        wantErr bool
    }{
        {"valid token", "Bearer valid-token", true, false},
        {"invalid token", "Bearer invalid", false, true},
        {"missing token", "", false, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := AuthenticateRequest(tt.token)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## Development Workflow

**REMINDER**: All development MUST follow the TDD workflow. No exceptions.

1. **Feature Development** (TDD MANDATORY)
   - Create feature branch
   - **RED**: Write failing tests first
   - Verify tests fail (run `go test -v`)
   - **GREEN**: Implement minimum code to pass
   - Verify all tests pass
   - **REFACTOR**: Clean up implementation
   - Run `make test-all`
   - **COMMIT**: Tests and implementation together
   - Submit PR with tests passing

2. **Bug Fixes** (TDD MANDATORY)
   - **RED**: Write test that reproduces bug
   - Verify test fails with current code
   - **GREEN**: Fix the issue
   - Ensure all tests pass
   - **REFACTOR**: Improve fix if needed
   - **COMMIT**: Test and fix together
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

- **TEST FIRST**: Never write implementation before tests
- Use `internal/` for private packages (Go convention)
- Follow clean architecture: separate concerns between packages
- Use interfaces for testability and flexibility
- Implement context propagation for cancellation
- Handle errors explicitly, never ignore them
- Use structured logging with clear messages
- Write self-documenting code with meaningful names

## System CLAUDE.md Compliance

This project adheres to all laws and workflows defined in the system CLAUDE.md, including:
- **Law 1**: NEVER Code Without Approval
- **Law 2**: Plan-First Workflow
- **Law 3**: Document Everything
- **Law 4**: Test-First for Features & Fixes (MANDATORY)

Refer to system CLAUDE.md for complete workflow requirements.

## Contribution Guidelines

See CONTRIBUTING.md for detailed guidelines. Key points:
- All code must have tests
- Follow Go best practices and idioms
- Use `go fmt` and `go vet`
- Document public APIs
- Sign commits with GPG key
- Squash commits before merging

## Utility Scripts

- `.claude/scripts/add-claude-coauthor.sh [num_commits]` - Add Claude as co-author to recent commits (default: 10)
  - Creates backup branch before making changes
  - Use when commits are missing Claude co-author attribution

## Important Instruction Reminders

The following instructions are inherited from system CLAUDE.md and apply to ALL work in this project:

- **ALWAYS** use Test-Driven Development (Red-Green-Refactor-Commit)
- **NEVER** write code without tests failing first
- **ALWAYS** get approval before implementing
- **ALWAYS** create todos for task tracking
- **ALWAYS** document all decisions and research

When in doubt, refer to system CLAUDE.md for the complete workflow.