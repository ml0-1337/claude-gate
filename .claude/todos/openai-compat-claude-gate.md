---
todo_id: openai-compat-claude-gate
started: 2025-06-25 17:21:15
completed:
status: in_progress
priority: high
---

# Task: Add OpenAI chat/completions compatibility to claude-gate

## Findings & Research

### Problem Analysis
- Cursor IDE sends OpenAI-format requests to `/v1/chat/completions`
- claude-gate only handles Anthropic-format requests to `/v1/messages`
- The transformer only processes `/v1/messages` requests, leaving OpenAI format untransformed
- This causes Anthropic API to reject the requests with "Invalid API Key" error

### Request Flow Investigation
1. llm-router correctly routes to claude-gate with `/v1/chat/completions`
2. claude-gate's mux routes `/v1/` to proxy handler
3. Proxy handler gets OAuth token successfully
4. But transformer skips transformation because path != `/v1/messages`
5. Untransformed OpenAI format sent to Anthropic API fails

### Format Differences
OpenAI format:
```json
{
  "model": "anthropic/claude-opus-4-20250514",
  "messages": [
    {"role": "system", "content": "..."},
    {"role": "user", "content": "..."}
  ],
  "max_tokens": 10,
  "temperature": 1,
  "stream": false
}
```

Anthropic format:
```json
{
  "model": "claude-opus-4-20250514",
  "messages": [
    {"role": "user", "content": "..."}
  ],
  "system": "...",
  "max_tokens": 10,
  "temperature": 1,
  "stream": false
}
```

## Test Strategy

- **Test Framework**: Go's built-in testing + testify
- **Test Types**: Unit tests for converter, Integration tests for full flow
- **Coverage Target**: 100% of new conversion logic
- **Edge Cases**: Missing system messages, multiple system messages, streaming, model name mapping

## Test List

**MANDATORY for TDD**: List all test scenarios BEFORE writing any tests. Focus on behavior, not implementation.

- [x] Test 1: Should convert OpenAI format with system message to Anthropic format
- [x] Test 2: Should handle OpenAI format without system message
- [x] Test 3: Should remove "anthropic/" prefix from model names
- [x] Test 4: Should preserve all other fields (temperature, max_tokens, etc.)
- [x] Test 5: Should handle multiple system messages by concatenating them
- [x] Test 6: Should inject Claude Code prompt into system field
- [ ] Test 7: Should convert streaming requests properly
- [ ] Test 8: Should handle assistant messages in conversation
- [x] Test 9: Should convert error responses back to OpenAI format
- [ ] Test 10: Should handle edge case of empty messages array

**Current Test**: Adding CORS headers and /v1/models endpoint for Cursor IDE compatibility

## Test Cases

```go
// Test 1: Should convert OpenAI format with system message to Anthropic format
// Input: OpenAI request with system and user messages
// Expected: Anthropic format with system field extracted

// Test 2: Should handle OpenAI format without system message
// Input: OpenAI request with only user messages
// Expected: Anthropic format with Claude Code prompt as system
```

## Maintainability Analysis

- **Readability**: [8/10] Clear separation of concerns, well-named functions
- **Complexity**: Low - straightforward JSON transformation
- **Modularity**: High - separate converter from existing code
- **Testability**: High - pure functions for conversion
- **Trade-offs**: Adding complexity for compatibility, but isolated in new module

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

- [x] Create openai_converter.go with conversion logic
- [x] Create openai_converter_test.go with all test cases
- [x] Update transformer.go to handle /v1/chat/completions
- [x] Update handler.go to transform path and response
- [x] Add response format conversion
- [ ] Test with actual Cursor IDE
- [ ] Update documentation
- [ ] Handle streaming responses
- [ ] Add error response conversion

## Working Scratchpad

### Requirements
1. Convert OpenAI chat/completions format to Anthropic messages format
2. Extract system messages to system field
3. Remove "anthropic/" prefix from model names
4. Inject Claude Code prompt
5. Convert responses back to OpenAI format

### Approach
Create a new `openai_converter.go` file in the proxy package with:
- `ConvertOpenAIToAnthropic(body []byte) ([]byte, error)`
- `ConvertAnthropicToOpenAI(body []byte) ([]byte, error)`
- Update transformer to use these for `/v1/chat/completions`

### Code

### Notes
- Need to handle streaming responses differently
- Error format conversion is important for Cursor to understand failures
- Model name mapping already exists in transformer.go

### Commands & Output

```bash

```