---
todo_id: future-claude-gate-enhancements
started: 2025-06-24 04:56:03
completed:
status: pending
priority: high
---

# Task: Claude Gate Future Enhancements - Distribution & Security

## Findings & Research

### Cross-Platform Builds & NPM Distribution
- Use GoReleaser for automated cross-platform builds
- Create npm packages with platform-specific binaries
- Package structure:
  ```
  @claude-gate/cli (main package)
  ├── @claude-gate/darwin-x64
  ├── @claude-gate/darwin-arm64
  ├── @claude-gate/linux-x64
  ├── @claude-gate/linux-arm64
  └── @claude-gate/win32-x64
  ```
- Use postinstall scripts to download correct binary
- Binary naming: claude-gate-darwin-arm64, claude-gate-linux-x64, etc.

### OS Keychain Integration
- 99designs/keyring library supports:
  - macOS Keychain
  - Windows Credential Manager
  - Linux Secret Service (GNOME Keyring, KWallet)
  - Encrypted file fallback
- Replace file-based storage with keyring
- Maintain backward compatibility with existing tokens

### Integration Tests
- Mock Anthropic API server for testing
- Test scenarios:
  - OAuth flow with mock endpoints
  - Token refresh cycle
  - SSE streaming responses
  - Error handling and retries
  - System prompt transformation
  - Model alias mapping

## Test Strategy

- **Test Framework**: Go standard testing + testify + httptest for mocks
- **Test Types**: Integration tests with mock server, E2E tests with real binary
- **Coverage Target**: 90% for critical paths
- **Edge Cases**: 
  - Network failures during OAuth
  - Expired tokens during streaming
  - Malformed SSE responses
  - Keychain access failures

## Test Cases

```go
// Test 1: Cross-platform binary execution
// Input: Built binaries for each platform
// Expected: Each binary runs without missing dependencies

// Test 2: NPM package installation
// Input: npm install @claude-gate/cli
// Expected: Correct platform binary downloaded and executable

// Test 3: Keychain token storage
// Input: OAuth tokens
// Expected: Stored in OS keychain, retrieved successfully

// Test 4: Keychain fallback
// Input: Keychain unavailable
// Expected: Falls back to encrypted file storage

// Test 5: Mock OAuth flow
// Input: Full OAuth login with mock server
// Expected: Complete flow without external dependencies

// Test 6: Mock SSE streaming
// Input: Streaming request to mock server
// Expected: Proper chunked response handling
```

## Maintainability Analysis

- **Readability**: [9/10] Clear separation of platform-specific code
- **Complexity**: Keep platform abstractions simple
- **Modularity**: Separate packages for each platform
- **Testability**: Mock servers enable offline testing
- **Trade-offs**: Multiple npm packages vs single binary

## Test Results Log

```bash
# Future test runs will be logged here
```

## Checklist

### Cross-Platform Builds & NPM Distribution
- [ ] Set up GoReleaser configuration
- [ ] Create build scripts for all platforms
- [ ] Test binaries on each platform
- [ ] Create npm package structure
- [ ] Write npm postinstall scripts
- [ ] Test npm installation on all platforms
- [ ] Set up GitHub Actions for automated releases
- [ ] Document installation methods

### OS Keychain Integration
- [ ] Add 99designs/keyring dependency
- [ ] Create keyring adapter interface
- [ ] Implement platform-specific storage
- [ ] Add migration from file storage
- [ ] Test on macOS Keychain
- [ ] Test on Windows Credential Manager
- [ ] Test on Linux Secret Service
- [ ] Implement encrypted file fallback
- [ ] Update documentation

### Integration Tests
- [ ] Create mock Anthropic API server
- [ ] Mock OAuth endpoints
- [ ] Mock messages endpoint with SSE
- [ ] Mock error responses
- [ ] Write OAuth flow integration tests
- [ ] Write proxy request integration tests
- [ ] Write token refresh integration tests
- [ ] Set up CI test matrix for all platforms

## Working Scratchpad

### Requirements
1. Easy installation via npm for all platforms
2. Secure token storage using OS keychains
3. Comprehensive testing without external dependencies
4. Maintain <15MB binary size for npm

### Approach
1. Start with GoReleaser setup
2. Build platform binaries
3. Create npm packages
4. Implement keychain integration
5. Build mock server for testing
6. Set up CI/CD pipeline

### Code

### Notes
- GoReleaser can handle code signing for macOS
- NPM package should detect platform automatically
- Keychain integration needs user permission prompts
- Mock server should closely mimic Anthropic's behavior

### Commands & Output

```bash
# Install GoReleaser
brew install goreleaser

# Test build
goreleaser build --snapshot --clean

# NPM publish (after setup)
npm publish --access public
```