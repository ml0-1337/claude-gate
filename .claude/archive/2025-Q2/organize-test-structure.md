---
todo_id: organize-test-structure
started: 2025-06-24 20:01:35
completed: 2025-06-24 20:07:45
status: completed
priority: high
---

# Task: Organize test structure with dedicated integration and E2E test directories

## Findings & Research

### Current Test Organization
The project currently follows standard Go conventions with unit tests alongside source files (*_test.go). Documentation mentions integration and E2E tests in `internal/test/integration/` and `internal/test/e2e/`, but these directories don't exist yet.

### Go Testing Best Practices (2024)
Based on web research:
1. **Build Tags Approach**: Use build tags like `//go:build integration` to separate test types
2. **Directory Structure**: Common patterns include:
   - Unit tests: alongside source files (*_test.go)
   - Integration tests: separate directory (often `integration/` or `it/`)
   - E2E tests: separate directory (often `e2e/`)
   - Test data: `testdata/` directory (ignored by Go tools)
3. **Test Helpers**: Shared utilities in dedicated helper packages
4. **Running Tests**: Use `-short` flag or build tags to control which tests run

### Makefile Analysis
Current Makefile has targets for:
- `test`: Basic Go test with race detection
- `test-all`: Runs comprehensive test suite via script
- `test-docker`: Tests in Docker containers
- `test-edge`: Edge case testing

## Test Strategy

- **Test Framework**: Go standard testing + testify (already in use)
- **Test Types**: Unit (existing), Integration (new), E2E (new)
- **Coverage Target**: Maintain current coverage, expand with integration tests
- **Edge Cases**: OAuth flows, proxy handling, error conditions

## Test Cases

```go
// Integration Test Example: OAuth Flow
func TestOAuthFlow_Integration(t *testing.T) {
    // Test 1: Complete OAuth flow with mock server
    // Input: Valid OAuth credentials
    // Expected: Successfully authenticated, token stored

    // Test 2: OAuth flow with expired token
    // Input: Expired OAuth token
    // Expected: Token refresh successful

    // Test 3: OAuth flow with invalid credentials
    // Input: Invalid client credentials
    // Expected: Authentication error
}

// E2E Test Example: CLI Commands
func TestCLI_E2E(t *testing.T) {
    // Test 1: Full authentication and proxy start
    // Input: claude-gate auth login && claude-gate start
    // Expected: Proxy running, accepting requests

    // Test 2: Dashboard interaction
    // Input: claude-gate dashboard
    // Expected: TUI displays correctly, responds to input
}
```

## Maintainability Analysis

- **Readability**: [9/10] Clear test organization improves understanding
- **Complexity**: Simple directory structure, build tags are standard Go practice
- **Modularity**: Separate test types allow focused testing
- **Testability**: Already using interfaces, good foundation
- **Trade-offs**: Slightly more complex build commands, but better organization

## Test Results Log

```bash
# Initial test run (should pass existing tests)
[2025-06-24 20:01:35] Red Phase: All existing tests pass, new directories don't exist yet

# After implementation
[2025-06-24 20:05:30] Green Phase: All tests passing:
- Unit tests: PASS
- Integration tests: PASS (4 tests)
- E2E tests: PASS (5 tests, 1 skipped)

# After refactoring
[2025-06-24 20:07:30] Refactor Phase: Cleaned up test imports and fixed expected outputs
```

## Checklist

- [x] Research Go test organization best practices
- [x] Analyze current test structure
- [x] Design new test organization
- [x] Create internal/test directory structure
- [x] Add integration test directory with README
- [x] Add E2E test directory with README
- [x] Create testdata directory for fixtures
- [x] Create helpers directory for shared utilities
- [x] Write example integration test
- [x] Write example E2E test
- [x] Update Makefile with new test targets
- [x] Test all make targets work correctly
- [x] Update documentation

## Working Scratchpad

### Requirements
- Maintain Go best practices
- Keep unit tests alongside source files
- Add dedicated directories for integration/E2E tests
- Use build tags to control test execution
- Update Makefile for easy test running

### Approach
1. Create directory structure under internal/test/
2. Add README files explaining each test type
3. Create example tests with proper build tags
4. Update Makefile with new targets
5. Ensure all existing tests still pass

### Code

Directory structure to create:
```
internal/
└── test/
    ├── integration/
    │   ├── README.md
    │   └── auth_integration_test.go
    ├── e2e/
    │   ├── README.md
    │   └── cli_e2e_test.go
    ├── testdata/
    │   └── .gitkeep
    └── helpers/
        ├── helpers.go
        └── mock_server.go
```

### Notes
- The project already mentions these directories in CLAUDE.md but they don't exist
- Need to ensure backward compatibility with existing test commands
- Build tags will prevent integration/E2E tests from running by default

### Commands & Output

```bash

```