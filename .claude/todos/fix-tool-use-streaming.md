---
todo_id: fix-tool-use-streaming
started: 2025-06-25 18:55:00
completed:
status: in_progress
priority: high
---

# Task: Fix Tool Use Streaming in OpenAI Compatibility Layer

## Findings & Research

### Problem Description
OpenAI-compatible streaming responses stop after the first response when using Cursor. The issue occurs even when sending requests directly to claude-gate (bypassing llm-router).

### Log Analysis
From the debug logs provided:
```
time=2025-06-25T18:53:03.391+09:00 level=DEBUG msg="SSE event received" event=content_block_start
time=2025-06-25T18:53:04.374+09:00 level=DEBUG msg="SSE event received" event=content_block_start
time=2025-06-25T18:53:04.375+09:00 level=DEBUG msg="SSE event received" event=content_block_delta
time=2025-06-25T18:53:05.220+09:00 level=DEBUG msg="SSE event received" event=content_block_delta
[... many more content_block_delta events ...]
```

Key observations:
- Total events received: 20+ content_block_delta events
- Total events converted: Only 7 events (4 content_block_delta events)
- The second `content_block_start` event (at 18:53:04.374) suggests a tool use block
- Many content_block_delta events after this are NOT being converted

### Root Cause
The `ConvertAnthropicSSEToOpenAI` function in `internal/proxy/openai_converter.go` only handles text deltas:

```go
case "content_block_delta":
    if delta, ok := eventData["delta"].(map[string]interface{}); ok {
        if delta["type"] == "text_delta" {  // Only handles text!
            if text, ok := delta["text"].(string); ok {
                // ... conversion logic
            }
        }
    }
```

Tool use deltas with `type: "input_json_delta"` are being ignored, causing the stream to appear incomplete.

### WebSearch Research Results

#### Anthropic Tool Use SSE Format
From Anthropic documentation:
- Tool use blocks have `content_block_start` with type "tool_use"
- Tool inputs are streamed as `input_json_delta` events with `partial_json` field
- Example: `{"type": "input_json_delta", "partial_json": "{\"location\": \"San Fra"}`
- The deltas are partial JSON strings that need to be accumulated
- A `content_block_stop` event signals the end of the tool block

#### OpenAI Streaming Tool Calls Format
From OpenAI documentation:
- Tool calls are part of the delta field: `delta.tool_calls`
- Each chunk contains: `delta.tool_calls[index].function.arguments` with partial JSON
- The tool call has: `id`, `type`, `function.name`, `function.arguments`
- Arguments are streamed incrementally as partial JSON strings

#### Key Differences
1. Anthropic uses `input_json_delta` with `partial_json` field
2. OpenAI uses `delta.tool_calls[].function.arguments` 
3. Both stream partial JSON that needs accumulation
4. Anthropic tracks tool blocks by content block index
5. OpenAI tracks tool calls by array index in delta.tool_calls

## Test Strategy

- **Test Framework**: Go's built-in testing with testify
- **Test Types**: Unit tests for SSE conversion
- **Coverage Target**: All content block delta types (text and tool use)
- **Edge Cases**: 
  - Mixed text and tool use in same stream
  - Multiple tool calls
  - Empty tool inputs
  - Malformed tool data

## Test List

**MANDATORY for TDD**: List all test scenarios BEFORE writing any tests. Focus on behavior, not implementation.

- [x] Test 1: Should convert text_delta events to OpenAI format
- [x] Test 2: Should convert input_json_delta events to OpenAI tool format
- [ ] Test 3: Should handle mixed text and tool content in same stream
- [x] Test 4: Should preserve tool names and arguments correctly
- [x] Test 5: Should handle empty tool input gracefully
- [x] Test 6: Should log skipped/unhandled event types for debugging
- [ ] Test 7: Should maintain message continuity with tool use blocks

**Current Test**: Working on fixing tool index and ID issues in streaming

**Remember**: This is behavioral analysis only. NO implementation details, NO "how", only "what".

## Test Cases

```go
// Test 1: Should convert input_json_delta events
// Input: content_block_delta with type:"input_json_delta" 
// Expected: OpenAI format chunk with tool_calls delta

// Test 2: Should handle complete tool use flow
// Input: content_block_start (tool_use) -> multiple input_json_delta -> content_block_stop
// Expected: Proper OpenAI tool calling format with accumulated JSON
```

## Maintainability Analysis

- **Readability**: [8/10] Current code is clear but missing tool handling
- **Complexity**: Low complexity, just needs additional cases
- **Modularity**: Good separation between event types
- **Testability**: Easy to test with table-driven tests
- **Trade-offs**: Need to balance completeness vs. complexity

## Test Results Log

```bash
# Initial test run (should fail)
[2025-06-25 19:04:00] Red Phase: All new tool use tests failing as expected
- TestConvertAnthropicSSEToOpenAI/should_convert_input_json_delta_events_to_OpenAI_tool_format: FAIL
- TestConvertAnthropicSSEToOpenAI/should_handle_content_block_start_for_tool_use: FAIL
- TestConvertAnthropicSSEToOpenAI/should_handle_empty_tool_input_gracefully: FAIL
- TestConvertAnthropicSSEToOpenAI/should_handle_multiple_tool_deltas_in_sequence: FAIL

# After implementation
[2025-06-25 19:05:20] Green Phase: All tests passing!
- Implemented content_block_start handling for tool_use
- Implemented input_json_delta conversion to OpenAI format
- Added tool state tracking between events
- All test scenarios passing

# After refactoring
[2025-06-25 19:06:50] Refactor Phase: Refactoring complete!
- Added ConvertAnthropicSSEToOpenAIWithLogger for optional logging
- Added debug logging for unhandled event types
- Added content_block_stop handling to clean up tool state
- All tests still passing

# Fix for tool index and ID issues
[2025-06-25 19:29:00] Red Phase: New tests failing as expected
- TestConvertAnthropicSSEToOpenAI/should_include_tool_ID_in_input_json_delta_events: FAIL
- TestConvertAnthropicSSEToOpenAI/should_handle_multiple_tools_with_correct_indices: FAIL

[2025-06-25 19:31:00] Green Phase: All tests passing!
- Fixed tool index mapping (Anthropic block index → OpenAI tool call index)
- Added tool ID to all delta events
- Added ResetSSEConverterState for test isolation
- Reset tool call index on message_start
- All tests passing including new tests for proper tool handling
```

## Checklist

- [x] Research Anthropic's tool use SSE format
- [x] Research OpenAI's streaming tool call format
- [x] Write failing tests for tool use conversion
- [x] Implement tool use delta conversion
- [x] Add debug logging for unhandled events
- [ ] Test with Cursor using tool-heavy prompts
- [ ] Document the streaming format differences

## Working Scratchpad

### Requirements
1. Support tool use deltas in SSE streaming
2. Convert Anthropic tool format to OpenAI function calling format
3. Maintain streaming continuity when switching between text and tools
4. Log any unhandled event types for future debugging

### Approach
1. Extend the switch statement in ConvertAnthropicSSEToOpenAI
2. Add handler for input_json_delta events
3. Track tool use state (name, accumulated JSON)
4. Convert to OpenAI's delta.tool_calls format

### Code
TBD - Write tests first

### Notes
- Anthropic uses content blocks with type "tool_use"
- Tool inputs are streamed as input_json_delta events
- OpenAI uses delta.tool_calls array in streaming
- Need to accumulate partial JSON and track tool indices

### Commands & Output

```bash
# Run existing tests
go test -v ./internal/proxy -run TestConvertAnthropicSSEToOpenAI

# Debug a real streaming response
claude-gate start --log-level DEBUG
# Then send a request that uses tools
```

### Issue Found: Tool Index and ID Problems
From user testing with Cursor:
- Error: "Model provided invalid arguments to terminal tool"
- Problem 1: Always using index: 0 in tool_calls array
- Problem 2: Not including tool ID in input_json_delta events
- Problem 3: Need to map Anthropic content block index to OpenAI tool call index

Fix needed:
- Track OpenAI tool call index (0, 1, 2...) separately from Anthropic block index
- Include tool ID in every delta event
- Properly handle multiple tools in same message

## Implementation Plan

1. **Research Phase**
   - Study Anthropic's tool use SSE format documentation
   - Study OpenAI's streaming function call format
   - Create test fixtures with real tool use events

2. **Test Phase**
   - Write comprehensive tests for tool use conversion
   - Include edge cases and error scenarios
   - Verify tests fail with current implementation

3. **Implementation Phase**
   - Add input_json_delta handling to ConvertAnthropicSSEToOpenAI
   - Implement proper OpenAI tool_calls delta format
   - Add debug logging for unhandled events

4. **Validation Phase**
   - Run tests to ensure all pass
   - Test with Cursor using real tool-calling scenarios
   - Monitor logs to ensure no events are skipped

## Related Issues
- Streaming stops after first response in Cursor
- Only affects OpenAI compatibility endpoint (/v1/chat/completions)
- Native Anthropic streaming (/v1/messages) works correctly