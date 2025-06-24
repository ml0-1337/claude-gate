---
todo_id: fix-keychain-trust-settings
started: 2025-06-24 23:51:40
completed:
status: in_progress
priority: high
---

# Task: Fix keychain item-level trust settings to prevent password prompts

## Findings & Research

### Problem Analysis
- KeychainTrustApplication is set at Config level, but items need individual trust settings
- The keyring.Item struct has `KeychainNotTrustApplication` field (note the "Not")
- When `KeychainNotTrustApplication = false`, it means the app SHOULD be trusted (double negative)
- Current implementation doesn't set this field, so it defaults to false (which is actually what we want)
- The issue might be that existing items were created without proper trust, or the Config-level setting isn't being applied

### Root Cause
- Code signing changes (rebuilds) invalidate previous trust settings
- Existing keychain items may have been created before the trust config was added
- Each rebuild creates a new binary with different checksum, requiring new trust approval

### Solution
1. Explicitly set item-level trust settings when creating keychain items
2. Add a reset command to clear and recreate keychain items with proper trust
3. Ensure both Config-level and Item-level trust settings are applied

## Test Strategy

- **Test Framework**: Go's built-in testing with testify
- **Test Types**: Unit tests for item creation, integration tests for trust behavior
- **Coverage Target**: 100% coverage of modified code
- **Edge Cases**: 
  - New items created with trust settings
  - Existing items without trust settings
  - Items after binary rebuild
  - Concurrent item creation

## Test Cases

```go
// Test 1: Item created with explicit trust settings
// Input: Set() called with token
// Expected: Item has KeychainNotTrustApplication = false

// Test 2: Trust settings preserved on Get
// Input: Get() called after Set() with trust
// Expected: No password prompt (manual test)

// Test 3: Reset command clears and recreates items
// Input: reset-keychain command
// Expected: All items recreated with trust settings
```

## Maintainability Analysis

- **Readability**: [10/10] Explicit field setting is clear
- **Complexity**: [10/10] Simple boolean field addition
- **Modularity**: [10/10] Isolated to Set() method
- **Testability**: [9/10] Easy to test field values, manual test for prompts
- **Trade-offs**: None - this is a necessary fix

## Test Results Log

```bash
# Test run completed
[2025-06-24 23:59:17] All tests passing:
- TestKeyringStorage_SetWithTrustSettings: PASS
- All auth package tests: PASS
```

## Checklist

- [✓] Update Set() method to explicitly set KeychainNotTrustApplication = false
- [✓] Add KeychainNotSynchronizable = true for security
- [✓] Write unit tests for item creation
- [✓] Add reset-keychain command
- [ ] Test manually on macOS (pending user verification)
- [✓] Update documentation
- [ ] Commit with tests and implementation

## Working Scratchpad

### Requirements
- Eliminate password prompts by setting item-level trust
- Ensure items don't sync to iCloud
- Provide way to reset existing items

### Approach
1. Modify Set() in keyring_storage.go to set item fields
2. Add auth storage reset-keychain command
3. Test with fresh keychain items

### Code

Implementation will go here as I work on it.

### Notes

- KeychainNotTrustApplication = false means app IS trusted (double negative)
- Need to be careful with the logic
- Users may need to re-authenticate after reset

### Commands & Output

```bash
# Test commands will be logged here
```