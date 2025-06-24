# Integration Tests

This directory contains integration tests for the Claude Gate proxy. These tests verify that different components work correctly together.

## Running Integration Tests

Integration tests are tagged with the `integration` build tag and are not run by default.

To run integration tests:

```bash
# Run only integration tests
make test-integration

# Or directly with go test
go test -tags=integration ./internal/test/integration/...

# Run all tests including integration
make test-all
```

## Test Organization

Integration tests should focus on:
- Component interactions (e.g., auth + proxy)
- Database/storage operations
- External service mocking (OAuth providers)
- Multi-step workflows

## Writing Integration Tests

All integration test files must include the build tag:

```go
//go:build integration
// +build integration

package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestComponentIntegration(t *testing.T) {
    // Test implementation
}
```

## Test Helpers

Use the shared test helpers from `internal/test/helpers/` for:
- Mock OAuth servers
- Test token providers
- Common test utilities