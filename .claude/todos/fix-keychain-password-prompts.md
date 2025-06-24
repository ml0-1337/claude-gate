---
todo_id: fix-keychain-password-prompts
started: 2025-06-24 23:22:59
completed:
status: in_progress
priority: high
---

# Task: Fix macOS keychain password prompts by implementing token caching in OAuthTokenProvider

## Findings & Research

### Problem Analysis
1. **Root Cause**: The `OAuthTokenProvider.GetAccessToken()` method calls `p.storage.Get("anthropic")` on EVERY proxy request, triggering keychain access each time.

2. **Flow Analysis**:
   - Every HTTP request → ProxyHandler.ServeHTTP() → TokenProvider.GetAccessToken() → storage.Get() → keyring.Get() → macOS password prompt
   - Despite `KeychainTrustApp: true` configuration, macOS still prompts because each keyring operation is treated as a new access

3. **Configuration Check**:
   - Config correctly sets `KeychainTrustApp: true` by default (config.go:68)
   - StorageFactory properly passes this to KeyringStorage (storage_factory.go:72-82)
   - KeyringConfig applies it to the keyring library (keyring_storage.go:60-62)

### WebSearch Findings
1. Known issue with 99designs/keyring library (Issue #118: "Is there a way to prevent continuous popups for password on keychain?")
2. Common workarounds:
   - Setting KeychainTrustApplication: true (already implemented)
   - Clicking "Always Allow" instead of "Allow" (user action required)
   - Implementing token caching to reduce keychain access frequency

## Test Strategy

- **Test Framework**: Go's built-in testing package (already in use)
- **Test Types**: Unit tests for the caching logic
- **Coverage Target**: 100% of new caching code
- **Edge Cases**: 
  - Cache miss on first access
  - Token expiration and refresh
  - Concurrent access to cache
  - Cache invalidation

## Test Cases

```go
// Test 1: Cache hit returns cached token without storage access
// Input: Second call to GetAccessToken
// Expected: Returns cached token, no storage.Get() call

// Test 2: Cache miss triggers storage access
// Input: First call to GetAccessToken
// Expected: Calls storage.Get(), caches result

// Test 3: Expired cached token triggers refresh
// Input: Cached token with past expiration
// Expected: Refreshes token, updates cache

// Test 4: Concurrent access is thread-safe
// Input: Multiple goroutines calling GetAccessToken
// Expected: No race conditions, consistent results
```

## Maintainability Analysis

- **Readability**: [9/10] Simple caching pattern, clear variable names
- **Complexity**: Low - adds single cache layer with clear expiration logic
- **Modularity**: High - caching logic isolated to OAuthTokenProvider
- **Testability**: High - easy to mock storage and test cache behavior
- **Trade-offs**: Small memory overhead for significant UX improvement

## Test Results Log

```bash
# Initial test run (should fail)
[2025-06-24 23:24:29] Red Phase: TestOAuthTokenProvider/caches_token_to_avoid_repeated_storage_access FAILED
- Expected 1 storage call, got 2 (second access)
- Expected 1 storage call, got 12 (after 10 more accesses)
- Current implementation calls storage.Get() on every GetAccessToken() call

# After implementation
[2025-06-24 23:26:14] Green Phase: All tests PASS
- Caching test: ✓ Storage only accessed once, subsequent calls use cache
- Refresh test: ✓ Cache updated when token refreshed
- Concurrent test: ✓ Thread-safe with no race conditions
- All existing tests: ✓ Still passing

# After refactoring
[timestamp] Refactor Phase: [pending]
```

## Checklist

- [ ] Write failing tests for token caching
- [ ] Add cached token fields to OAuthTokenProvider struct
- [ ] Implement GetAccessToken with caching logic
- [ ] Add mutex for thread-safe cache access
- [ ] Handle token expiration and refresh in cache
- [ ] Run tests and ensure all pass
- [ ] Test manually on macOS to verify password prompts reduced
- [ ] Update documentation if needed

## Working Scratchpad

### Requirements
- Cache tokens in memory to avoid repeated keychain access
- Refresh cache when token expires or needs refresh
- Thread-safe implementation for concurrent requests
- Maintain existing security (tokens still stored in keychain)

### Approach
1. Add fields to OAuthTokenProvider:
   - cachedToken *TokenInfo
   - cacheTime time.Time
   - cacheMutex sync.RWMutex
2. Modify GetAccessToken to check cache first
3. Only access storage on cache miss or expiration
4. Update cache after token refresh

### Code
```go
type OAuthTokenProvider struct {
    client       *OAuthClient
    storage      StorageBackend
    cachedToken  *TokenInfo
    cacheTime    time.Time
    cacheMutex   sync.RWMutex
}
```

### Notes
- This aligns with standard OAuth practices
- Similar pattern used in aws-vault and other tools
- Will reduce keychain access from every request to once per session

### Commands & Output

```bash
# Will update with test output
```