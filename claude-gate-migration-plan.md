# Claude Gate - Migration Plan

## Project Overview
Migrate Claude Auth Bridge from Python to Go with npm distribution, creating a high-performance OAuth proxy server named "claude-gate".

## Phase 1: Core Go Implementation (Week 1)

### Project Setup
- [ ] Create new GitHub repository: `claude-gate`
- [ ] Initialize Go module: `go mod init github.com/yourusername/claude-gate`
- [ ] Set up project structure:
  ```
  claude-gate/
  ├── cmd/claude-gate/          # CLI entry point
  ├── internal/
  │   ├── auth/                 # OAuth implementation
  │   │   ├── oauth.go         # PKCE flow
  │   │   └── storage.go       # Token encryption
  │   ├── proxy/
  │   │   ├── server.go        # HTTP server
  │   │   ├── handler.go       # Request processing
  │   │   └── middleware.go    # Auth injection
  │   └── config/
  │       └── config.go        # Settings management
  ├── go.mod
  └── go.sum
  ```

### OAuth Authentication
- [ ] Implement PKCE generation (base64url encoded 32-byte verifier)
- [ ] Create OAuth flow with Claude's client ID: `9d1c250a-e61b-44d9-88ed-5944d1962f5e`
- [ ] Build authorization URL with required parameters
- [ ] Implement token exchange endpoint
- [ ] Add token refresh mechanism
- [ ] Create secure token storage (~/.claude-gate/auth.json)
- [ ] Implement token encryption using system keychain

### HTTP Proxy Server
- [ ] Create base HTTP server using standard library
- [ ] Implement request forwarding to Anthropic API
- [ ] Add OAuth header injection:
  - `Authorization: Bearer {token}`
  - `anthropic-beta: oauth-2025-04-20`
  - `anthropic-version: 2023-06-01`
- [ ] Handle streaming responses (SSE)
- [ ] Implement request/response logging
- [ ] Add error handling and retry logic

### System Prompt Handler
- [ ] Detect system prompt in request body
- [ ] Convert string format to array format when needed
- [ ] Ensure Claude Code identification comes first:
  ```json
  "system": [
    {"type": "text", "text": "You are Claude Code, Anthropic's official CLI for Claude."},
    {"type": "text", "text": "original custom prompt"}
  ]
  ```
- [ ] Handle existing array formats
- [ ] Map model aliases (e.g., `claude-3-5-sonnet-latest` → `claude-3-5-sonnet-20241022`)

### CLI Interface
- [ ] Implement CLI using cobra or urfave/cli
- [ ] Create commands:
  - [ ] `claude-gate start` - Start proxy server
  - [ ] `claude-gate auth login` - OAuth authentication
  - [ ] `claude-gate auth logout` - Clear credentials
  - [ ] `claude-gate auth status` - Check auth status
  - [ ] `claude-gate test` - Test proxy connection
  - [ ] `claude-gate status` - Show configuration
- [ ] Add flags:
  - [ ] `--host` (default: 127.0.0.1)
  - [ ] `--port` (default: 8000)
  - [ ] `--auth-token` (optional proxy auth)
  - [ ] `--log-level` (DEBUG, INFO, WARNING, ERROR)

### Configuration
- [ ] Support environment variables (CLAUDE_GATE_*)
- [ ] Load from .env file
- [ ] Command line flag precedence
- [ ] Create default configuration
- [ ] Add CORS settings
- [ ] Implement rate limiting options

## Phase 2: NPM Distribution Setup (Week 2)

### Package Structure
- [ ] Create npm package directory structure:
  ```
  npm/
  ├── package.json
  ├── cli.js
  ├── README.md
  └── packages/
      ├── claude-gate-darwin-x64/
      ├── claude-gate-darwin-arm64/
      ├── claude-gate-linux-x64/
      ├── claude-gate-linux-arm64/
      └── claude-gate-windows-x64/
  ```

### Main Package Configuration
- [ ] Create main package.json with optionalDependencies
- [ ] Write cli.js wrapper to detect platform and run binary
- [ ] Add binary detection logic
- [ ] Implement fallback download mechanism
- [ ] Create installation instructions

### Platform Packages
- [ ] Create package.json for each platform with os/cpu fields
- [ ] Set up binary packaging scripts
- [ ] Configure npm publishing for each platform
- [ ] Test installation on different platforms

### Build Automation
- [ ] Create cross-platform build script
- [ ] Set up GitHub Actions for automated builds
- [ ] Configure release workflow
- [ ] Add version bumping automation
- [ ] Create npm publish workflow

## Phase 3: Feature Parity (Week 3)

### Core Features
- [ ] Port all OAuth endpoints from Python
- [ ] Implement health check endpoint
- [ ] Add proxy authentication middleware
- [ ] Create request size limits
- [ ] Implement timeout handling
- [ ] Add graceful shutdown

### Advanced Features
- [ ] WebSocket support for future streaming
- [ ] Metrics collection (request count, latency)
- [ ] Multi-account support
- [ ] Configuration hot-reload
- [ ] Request/response transformers

### Testing
- [ ] Unit tests for OAuth flow
- [ ] Integration tests for proxy
- [ ] End-to-end tests with real API
- [ ] Cross-platform installation tests
- [ ] Performance benchmarks
- [ ] Load testing

### Documentation
- [ ] API documentation
- [ ] Installation guide
- [ ] Migration guide from Python version
- [ ] Configuration reference
- [ ] Troubleshooting guide
- [ ] Contributing guidelines

## Phase 4: Release & Migration (Week 4)

### Pre-release
- [ ] Security audit
- [ ] Performance profiling
- [ ] Memory leak testing
- [ ] Cross-platform testing
- [ ] Documentation review

### Release Process
- [ ] Tag v1.0.0-beta
- [ ] Publish to npm as @next
- [ ] Create GitHub release
- [ ] Announce in relevant communities
- [ ] Gather feedback

### Migration Support
- [ ] Create migration script from Python version
- [ ] Transfer existing auth tokens
- [ ] Update documentation references
- [ ] Create comparison guide
- [ ] Support both versions temporarily

### Post-release
- [ ] Monitor issue reports
- [ ] Address critical bugs
- [ ] Plan feature roadmap
- [ ] Deprecate Python version
- [ ] Archive old repository

## Technical Specifications

### Dependencies
```go
// go.mod
require (
    golang.org/x/oauth2 v0.27.0
    github.com/spf13/cobra v1.8.0
    github.com/joho/godotenv v1.5.1
    github.com/zalando/go-keyring v0.2.3
)
```

### API Compatibility
- Base URL: `http://localhost:8000`
- All Anthropic API endpoints supported
- Request/response format unchanged
- Compatible with any Anthropic SDK

### Performance Targets
- Binary size: < 10MB
- Memory usage: < 50MB idle
- Startup time: < 100ms
- Request latency: < 5ms overhead
- Concurrent connections: 1000+

### Security Requirements
- No postinstall scripts
- Encrypted token storage
- Optional proxy authentication
- No external dependencies at runtime
- Regular security updates

## Success Metrics

- [ ] All Python features ported
- [ ] Performance improvement verified (>50% memory reduction)
- [ ] npm installation works on all platforms
- [ ] Zero security vulnerabilities
- [ ] User migration completed
- [ ] Documentation comprehensive
- [ ] Community adoption growing

## Notes

- Prioritize security and performance
- Maintain API compatibility
- Focus on developer experience
- Keep binary size minimal
- Ensure cross-platform support
- Document everything thoroughly