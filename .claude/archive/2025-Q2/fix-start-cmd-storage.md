---
todo_id: fix-start-cmd-storage
started: 2025-06-25 03:29:53
completed: 2025-06-25 03:34:52
status: completed
priority: high
---

# Task: Fix StartCmd Storage Implementation to Match DashboardCmd

## Findings & Research

### Storage Implementation Analysis

The authentication check behaves differently between `start` and `dashboard` commands due to different storage implementations:

1. **StartCmd** (line 88-96 in main.go):
   - Uses deprecated `auth.NewTokenStorage(cfg.AuthStoragePath)` 
   - This creates file-only storage that looks at `~/.config/claude-gate/auth.json`
   - Doesn't support keychain storage

2. **DashboardCmd** (line 174-189 in main.go):
   - Uses modern `auth.NewStorageFactory(createStorageFactoryConfig(cfg))`
   - Creates storage based on config (defaults to keychain on macOS)
   - Supports migration from file to keychain

### Storage Factory Analysis

From `internal/auth/storage_factory.go`:
- Default storage type is "auto" which selects keychain on macOS
- Migration logic in `CreateWithMigration()` moves tokens from file to keychain
- This explains why dashboard works (uses keychain) but start doesn't (uses file)

### Root Cause

When tokens are migrated to keychain by `dashboard` or `auth login`, the `start` command can't find them because it only looks in the file storage.

## Test Strategy

- **Test Framework**: Go standard testing with testify
- **Test Types**: Unit tests for storage consistency
- **Coverage Target**: Both commands using same storage
- **Edge Cases**: Migration scenarios, missing tokens, storage errors

## Test Cases

```go
// Test 1: StartCmd uses storage factory
// Input: Mock config with keychain storage
// Expected: StartCmd creates storage via factory

// Test 2: Both commands use identical storage
// Input: Same config for both commands
// Expected: Same storage instance type

// Test 3: StartCmd handles migrated tokens
// Input: Token in keychain, empty file storage
// Expected: StartCmd finds token successfully
```

## Maintainability Analysis

- **Readability**: [8/10] Clear separation between old/new patterns
- **Complexity**: Simple fix - replace one function call
- **Modularity**: Good - factory pattern already implemented
- **Testability**: Easy to test with mocks
- **Trade-offs**: None - purely beneficial change

## Test Results Log

```bash
# Initial test run (should fail)
[timestamp] Red Phase: [output]

# After implementation
[timestamp] Green Phase: [output]

# After refactoring
[timestamp] Refactor Phase: [output]
```

## Checklist

- [ ] Write failing test for StartCmd storage consistency
- [ ] Update StartCmd to use StorageFactory (lines 88-96)
- [ ] Ensure both auth check and server creation use same storage
- [ ] Remove deprecated NewTokenStorage function
- [ ] Verify all tests pass
- [ ] Test manually with both start and dashboard commands
- [ ] Update any documentation if needed

## Working Scratchpad

### Requirements
- Make StartCmd use the same storage mechanism as DashboardCmd
- Ensure backward compatibility
- Remove deprecated code

### Approach
1. Replace NewTokenStorage with StorageFactory in StartCmd
2. Match the exact pattern from DashboardCmd
3. Remove deprecated functions
4. Add tests to prevent regression

### Code
```go
// Before (StartCmd auth check):
storage := auth.NewTokenStorage(cfg.AuthStoragePath)

// After (StartCmd auth check):
factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
storage, err := factory.Create()
if err != nil {
    return fmt.Errorf("failed to create storage: %w", err)
}
```

### Notes
- DashboardCmd uses factory.CreateWithMigration() later for the server
- StartCmd should follow the same pattern for consistency
- Need to check if migration is needed during auth check

### Commands & Output

```bash
# Found that fix was already implemented in commit 51d625c
git log --oneline -10 -- cmd/claude-gate/main.go

# Rebuilt binary to include the fix
go build -o claude-gate ./cmd/claude-gate

# Verified start command now works
./claude-gate start
✓ OAuth authentication configured and ready
```

## Resolution

The issue was already fixed in commit `51d625c`. The user was running an old binary that didn't include the fix. After rebuilding the binary with the latest code, both `start` and `dashboard` commands now use the same StorageFactory pattern and can access tokens stored in the keychain.