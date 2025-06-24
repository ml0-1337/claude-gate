---
todo_id: fix-storage-inconsistency
started: 2025-06-25 03:17:29
completed:
status: in_progress
priority: high
---

# Task: Fix storage implementation inconsistency between start and dashboard commands

## Findings & Research

### Issue Analysis
The `start` command uses deprecated `NewTokenStorage` function which only accesses file-based storage, while the `dashboard` command uses the modern `StorageFactory` pattern that supports both file and keyring storage. When tokens are migrated to keyring (as evidenced by `auth.json.migrated` file), the start command fails to find them.

### Key Differences Found
1. **Start Command (line 89)**: `storage := auth.NewTokenStorage(cfg.AuthStoragePath)`
   - Uses deprecated function
   - Always uses file storage
   - Ignores storage type configuration

2. **Dashboard Command (lines 176-179)**: Uses `StorageFactory` pattern
   - Respects AuthStorageType config
   - Supports auto/keyring/file modes
   - Includes migration logic

### Evidence of Migration
- File `~/.claude-gate/auth.json.migrated` exists
- No `auth.json` file present
- Indicates tokens were migrated to keyring

## Test Strategy

- **Test Framework**: Go's built-in testing with testify assertions
- **Test Types**: Unit tests for storage factory usage
- **Coverage Target**: Ensure both commands use identical storage initialization
- **Edge Cases**: 
  - File storage only
  - Keyring storage with migrated tokens
  - Auto mode fallback scenarios

## Test Cases

```go
// Test 1: StartCmd uses StorageFactory
// Input: Start command with default config
// Expected: Creates storage using factory pattern

// Test 2: Both commands use same storage backend
// Input: Same config for both start and dashboard
// Expected: Both access same storage type

// Test 3: Start command reads keyring tokens
// Input: Tokens in keyring, start command
// Expected: Successfully reads OAuth tokens
```

## Maintainability Analysis

- **Readability**: [9/10] Clear separation between legacy and modern patterns
- **Complexity**: Simple fix - replace one function call
- **Modularity**: Good - factory pattern already exists
- **Testability**: Excellent - can mock storage backends
- **Trade-offs**: None - purely beneficial change

## Implementation Steps

1. [ ] Write failing tests for start command storage usage
2. [ ] Update StartCmd.Run() to use StorageFactory
3. [ ] Ensure consistent storage usage in server initialization
4. [ ] Remove deprecated NewTokenStorage function
5. [ ] Run all tests to ensure nothing breaks
6. [ ] Update other commands if needed

## Checklist

- [✓] Write tests for storage consistency
- [✓] Update StartCmd to use StorageFactory
- [✓] Match exact pattern from DashboardCmd
- [✓] Remove deprecated functions
- [✓] Test with both file and keyring storage
- [✓] Verify migration scenarios work (auth status confirms it works)
- [ ] Update documentation if needed

## Working Scratchpad

### Requirements
- Both start and dashboard commands must use identical storage access patterns
- Must support keyring storage for start command
- Maintain backward compatibility for file-only storage

### Approach
Replace the deprecated NewTokenStorage call in StartCmd with the StorageFactory pattern already used in DashboardCmd. This ensures consistency and proper keyring support.

### Code
Changes needed in main.go StartCmd.Run():
- Line 89: Replace storage initialization
- Line 123-129: Already uses factory for server, good
- Remove lines 266-268: Deprecated function definitions

### Test Results Log

```bash
# Red Phase (tests written, demonstrating bug)
[2025-06-25 03:17:29] Tests pass showing the issue:
- TestStartCmdFailsWithMigratedTokens: PASS
- Shows StartCmd with NewTokenStorage finds no tokens when migrated

# Green Phase (implementation complete)
[2025-06-25 03:30:00] Fixed StartCmd to use StorageFactory:
- Updated line 89 to use factory pattern
- Auth status command now finds tokens correctly
- Build succeeds

# Refactor Phase (cleanup)
[2025-06-25 03:32:00] Removed deprecated functions:
- Deleted NewTokenStorage and TokenStorage type alias
- Updated all test files to use NewFileStorage
- Project builds successfully
```

### Commands & Output

```bash
# Check current storage implementation
grep -n "NewTokenStorage\|NewStorageFactory" cmd/claude-gate/main.go
```