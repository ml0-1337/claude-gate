# End-to-End Tests

This directory contains end-to-end (E2E) tests for the Claude Gate proxy. These tests verify complete user workflows from the CLI interface through to API responses.

## Running E2E Tests

E2E tests are tagged with the `e2e` build tag and are not run by default due to their longer execution time.

To run E2E tests:

```bash
# Run only E2E tests
make test-e2e

# Or directly with go test
go test -tags=e2e ./internal/test/e2e/...

# Run all tests including E2E
make test-all
```

## Test Organization

E2E tests should focus on:
- Complete user workflows (auth → start → API call)
- CLI command interactions
- Real OAuth flow testing (with appropriate mocks)
- Cross-platform behavior verification

## Writing E2E Tests

All E2E test files must include the build tag:

```go
//go:build e2e
// +build e2e

package e2e_test

import (
    "testing"
    "os/exec"
    "github.com/stretchr/testify/assert"
)

func TestCompleteWorkflow(t *testing.T) {
    // Test implementation
}
```

## Test Environment

E2E tests may require:
- Built binaries (run `make build` first)
- Mock OAuth server running
- Clean test environment (no existing auth tokens)

## Best Practices

1. Use test fixtures from `internal/test/testdata/`
2. Clean up resources after tests
3. Use timeouts for long-running operations
4. Test both success and failure paths