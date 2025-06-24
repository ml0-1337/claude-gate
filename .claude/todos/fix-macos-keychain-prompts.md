---
todo_id: fix-macos-keychain-prompts
started: 2025-06-24 22:49:42
completed:
status: in_progress
priority: high
---

# Task: Fix macOS keychain password prompts by adding KeychainTrustApplication configuration

## Findings & Research

### Problem Analysis
- The claude-gate proxy triggers password prompts on macOS every time it accesses the keychain
- Root cause: `KeychainTrustApplication` flag is not being set to `true` in the keyring configuration
- Current implementation creates keyring.Config but doesn't set macOS-specific trust fields

### Library Investigation
- Using `github.com/99designs/keyring v1.2.2`
- The keyring.Config struct has these macOS-specific fields available:
  - `KeychainTrustApplication bool` - whether the calling application should be trusted by default
  - `KeychainAccessibleWhenUnlocked bool` - whether the item is accessible when device is locked
  - `KeychainSynchronizable bool` - whether the item can be synchronized to iCloud
  - `KeychainName string` - name of the macOS keychain to use

### WebSearch Findings
- No programmatic way to bypass keychain prompts entirely (security by design)
- Setting `KeychainTrustApplication: true` allows the app to be trusted after first "Always Allow" click
- This is the standard approach used by tools like aws-vault
- Users will see one prompt where they can choose "Always Allow" and won't be prompted again

## Test Strategy

- **Test Framework**: Go's built-in testing with testify
- **Test Types**: Unit tests for configuration, integration tests for keyring behavior
- **Coverage Target**: 100% coverage of new configuration code
- **Edge Cases**: 
  - Non-macOS platforms should ignore macOS-specific settings
  - Default values should be applied when not explicitly set
  - Environment variable overrides should work correctly

## Test Cases

```go
// Test 1: macOS-specific config is applied on Darwin
// Input: runtime.GOOS = "darwin", create KeyringStorage
// Expected: keyring.Config has KeychainTrustApplication = true

// Test 2: macOS-specific config is ignored on Linux
// Input: runtime.GOOS = "linux", create KeyringStorage  
// Expected: keyring.Config fields are unmodified

// Test 3: Environment variables override defaults
// Input: CLAUDE_GATE_KEYCHAIN_TRUST_APP=false
// Expected: KeychainTrustApplication = false

// Test 4: KeyringConfig struct preserves settings
// Input: KeyringConfig with custom macOS settings
// Expected: Settings passed through to keyring.Config
```

## Maintainability Analysis

- **Readability**: [9/10] Clear field names, follows existing patterns
- **Complexity**: [10/10] Simple boolean flags, no complex logic
- **Modularity**: [9/10] Platform-specific code isolated, easy to extend
- **Testability**: [10/10] Simple to test with mocks and environment variables
- **Trade-offs**: None - this is a straightforward configuration enhancement

## Test Results Log

```bash
# Initial test run - all passing
[2025-06-24 22:59:29] All tests passing:
- Config tests: PASS (TestConfig_macOSKeychainDefaults, TestConfig_LoadFromEnv_macOSKeychain)
- Auth tests: PASS (TestKeyringStorage_macOSConfig, TestStorageFactory_macOSDefaults)
- Coverage: 74.8% for auth package
```

## Checklist

- [✓] Update KeyringConfig struct with macOS fields
- [✓] Modify NewKeyringStorage to set macOS config
- [✓] Update storage_factory.go for platform detection
- [✓] Add environment variable support in config.go
- [✓] Write unit tests for configuration
- [✓] Write integration tests for keyring behavior
- [✓] Run tests and verify they pass
- [ ] Test manually on macOS
- [✓] Update troubleshooting documentation
- [ ] Commit with tests and implementation

## Working Scratchpad

### Requirements
- Eliminate repeated password prompts on macOS
- Maintain security (one-time "Always Allow" is acceptable)
- Don't affect other platforms
- Allow environment variable overrides

### Approach
1. Add macOS-specific fields to KeyringConfig
2. Set sensible defaults in NewKeyringStorage
3. Apply settings only on macOS platform
4. Support environment variable overrides

### Code

Implementation will go here as I work on it.

### Notes

- KeychainTrustApplication allows the app to be trusted after first authorization
- This matches behavior of other CLI tools like aws-vault
- Setting KeychainSynchronizable to false prevents tokens from syncing to iCloud

### Commands & Output

```bash
# Test commands will be logged here
```