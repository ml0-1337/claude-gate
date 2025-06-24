---
todo_id: implement-token-caching
started: 2025-06-24 23:34:56
completed:
status: in_progress
priority: high
---

# Task: Implement token caching to prevent repeated macOS keychain prompts

## Findings & Research

### Root Cause Analysis
- `OAuthTokenProvider.GetAccessToken()` calls `p.storage.Get("anthropic")` on EVERY proxy request
- Each HTTP request to the proxy triggers a keychain access, causing password prompts
- Despite `KeychainTrustApp: true`, macOS prompts because of the frequency of access
- The configuration is correct, but the access pattern is the problem

### Current Flow
1. Client makes request to proxy
2. ProxyHandler calls tokenProvider.GetAccessToken()
3. GetAccessToken() calls storage.Get() - **triggers keychain access**
4. Keychain prompts for password
5. Token is returned and used
6. Process repeats for EVERY request

## Test Strategy

- **Test Framework**: Go's built-in testing with testify
- **Test Types**: Unit tests for cache behavior, integration tests for token provider
- **Coverage Target**: 100% coverage of new caching code
- **Edge Cases**: 
  - Cache miss on first access
  - Cache hit within TTL
  - Cache invalidation on expiry
  - Token refresh clears cache
  - Concurrent access to cache

## Test Cases

```go
// Test 1: Cache miss triggers storage access
// Input: Empty cache, call GetAccessToken
// Expected: Storage.Get called once, token cached

// Test 2: Cache hit avoids storage access
// Input: Cached valid token, call GetAccessToken
// Expected: Storage.Get not called, cached token returned

// Test 3: Expired cache triggers refresh
// Input: Cached expired token, call GetAccessToken
// Expected: Storage.Get called, token refreshed and cached

// Test 4: Concurrent access is thread-safe
// Input: Multiple goroutines calling GetAccessToken
// Expected: No race conditions, correct token returned
```

## Maintainability Analysis

- **Readability**: [9/10] Simple caching pattern, clear TTL logic
- **Complexity**: [9/10] Minimal added complexity, standard caching approach
- **Modularity**: [10/10] Cache isolated within OAuthTokenProvider
- **Testability**: [10/10] Easy to test with time mocking
- **Trade-offs**: Small memory usage for massive UX improvement

## Test Results Log

```bash
# Will be updated as tests are run
```

## Checklist

- [ ] Add cached token field to OAuthTokenProvider struct
- [ ] Add mutex for thread-safe access
- [ ] Implement GetAccessToken with cache check
- [ ] Add cache invalidation on token refresh
- [ ] Write unit tests for caching behavior
- [ ] Write concurrent access tests
- [ ] Run all tests
- [ ] Test manually on macOS
- [ ] Update documentation if needed
- [ ] Commit with tests and implementation

## Working Scratchpad

### Requirements
- Cache tokens in memory to avoid repeated keychain access
- Thread-safe implementation
- Automatic cache invalidation on expiry
- Clear cache on token refresh

### Approach
1. Add `cachedToken *TokenInfo` and `cacheMutex sync.RWMutex` to OAuthTokenProvider
2. Check cache before accessing storage
3. Cache tokens with appropriate TTL
4. Invalidate on refresh

### Code

Implementation will go here as I work on it.

### Notes

- This maintains security (tokens stay in keychain)
- Reduces keychain access from every request to once per session
- Standard OAuth token handling practice

### Commands & Output

```bash
# Test commands will be logged here
```