---
completed: ""
current_test: 'Test 1: Valid credentials transformation'
priority: high
started: "2025-07-02 17:38:38"
status: in_progress
todo_id: implement-claude-code-credential-adapter-for
type: feature
---

# Task: Implement Claude Code credential adapter for keychain integration

## Findings & Research
## Token Structure Analysis

### Claude Code Keychain Storage
- Service name: `"Claude Code-credentials"`
- Account name: `"macbook"` (username)
- Token format:
```json
{
  "claudeAiOauth": {
    "accessToken": "sk-ant-oat01-...",
    "refreshToken": "sk-ant-ort01-...",
    "expiresAt": 1751458199105,  // Unix timestamp in milliseconds
    "scopes": ["user:inference", "user:profile"],
    "subscriptionType": "max"
  }
}
```

### claude-gate Keychain Storage
- Service name: `"claude-gate"`
- Account name: `"claude-gate.anthropic"` (format: `serviceName.provider`)
- Token format:
```json
{
  "type": "oauth",
  "refresh": "sk-ant-ort01-...",
  "access": "sk-ant-oat01-...",
  "expires": 1751309205  // Unix timestamp in seconds
}
```

### Key Differences
1. Service names differ
2. Account naming conventions differ
3. Nested vs flat JSON structure
4. Different field names (accessToken vs access, etc.)
5. Timestamp format (milliseconds vs seconds)
6. Claude Code includes scopes and subscriptionType
7. claude-gate includes type field

### Architecture Findings
- claude-gate has clean StorageBackend interface
- KeyringStorage implements keychain operations
- StorageFactory creates backends based on config
- No existing adapter patterns found
## Web Searches

## Test Strategy
- **Test Framework**: Go standard testing with testify assertions
- **Test Types**: Unit tests for transformation logic, integration tests with mock keychain
- **Coverage Target**: 90% for the adapter implementation
- **Edge Cases**: 
  - Missing Claude Code credentials
  - Invalid JSON format
  - Expired tokens
  - Empty/null fields
  - Timestamp conversion edge cases
## Test List
- [ ] Test 1: Should retrieve valid Claude Code credentials and return in claude-gate format
- [ ] Test 2: Should return nil when Claude Code credentials don't exist in keychain
- [ ] Test 3: Should handle invalid JSON format in Claude Code credentials gracefully
- [ ] Test 4: Should convert expiration timestamp from milliseconds to seconds correctly
- [ ] Test 5: Should map all token fields correctly (accessToken→access, refreshToken→refresh)
- [ ] Test 6: Should add "oauth" type field to transformed tokens
- [ ] Test 7: Should return appropriate provider name ("anthropic") in List operation
- [ ] Test 8: Should report as available when keychain is accessible
- [ ] Test 9: Should handle Set operation as no-op (read-only adapter)
- [ ] Test 10: Should handle Remove operation as no-op (read-only adapter)
- [ ] Test 11: Should return correct adapter name for identification
- [ ] Test 12: Should handle missing required fields in Claude Code JSON
## Test Cases

## Maintainability Analysis

## Test Results Log

[2025-07-02 17:39:43] ```bash
[2025-07-02 17:39:43] # Initial test run - Test 1 (should fail)
[2025-07-02 06:42:01] Red Phase:
[2025-07-02 17:39:43] internal/auth/claude_code_storage_test.go:15:3: unknown field items in struct literal of type mockKeyring
[2025-07-02 17:39:43] internal/auth/claude_code_storage_test.go:29:14: undefined: ClaudeCodeStorage
[2025-07-02 17:39:43] FAIL	github.com/ml0-1337/claude-gate/internal/auth [build failed]
[2025-07-02 17:39:43] ```
[2025-07-02 17:40:58] # After implementation
[2025-07-02 06:44:15] Green Phase:
[2025-07-02 17:40:58] === RUN   TestClaudeCodeStorage_Get_ValidCredentials
[2025-07-02 17:40:58] --- PASS: TestClaudeCodeStorage_Get_ValidCredentials (0.00s)
[2025-07-02 17:40:58] PASS
[2025-07-02 17:40:58] ok  	github.com/ml0-1337/claude-gate/internal/auth	0.278s
## Checklist

## Working Scratchpad
