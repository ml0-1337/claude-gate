package proxy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertOpenAIToAnthropic(t *testing.T) {
	t.Run("should convert OpenAI format with system message to Anthropic format", func(t *testing.T) {
		// Arrange
		openAIRequest := map[string]interface{}{
			"model": "anthropic/claude-opus-4-20250514",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				map[string]interface{}{
					"role":    "user",
					"content": "Hello, how are you?",
				},
			},
			"max_tokens":  100,
			"temperature": 0.7,
			"stream":      false,
		}
		
		requestBody, err := json.Marshal(openAIRequest)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertOpenAIToAnthropic(requestBody)
		
		// Assert
		require.NoError(t, err)
		
		var anthropicRequest map[string]interface{}
		err = json.Unmarshal(result, &anthropicRequest)
		require.NoError(t, err)
		
		// Check model name (prefix removed)
		assert.Equal(t, "claude-opus-4-20250514", anthropicRequest["model"])
		
		// Check system field (should include Claude Code prompt + original system)
		expectedSystem := []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
			map[string]interface{}{"type": "text", "text": "You are a helpful assistant."},
		}
		assert.Equal(t, expectedSystem, anthropicRequest["system"])
		
		// Check messages (system message removed)
		messages := anthropicRequest["messages"].([]interface{})
		assert.Len(t, messages, 1)
		assert.Equal(t, "user", messages[0].(map[string]interface{})["role"])
		assert.Equal(t, "Hello, how are you?", messages[0].(map[string]interface{})["content"])
		
		// Check other fields preserved
		assert.Equal(t, float64(100), anthropicRequest["max_tokens"])
		assert.Equal(t, 0.7, anthropicRequest["temperature"])
		assert.Equal(t, false, anthropicRequest["stream"])
	})
	
	t.Run("should handle OpenAI format without system message", func(t *testing.T) {
		// Arrange
		openAIRequest := map[string]interface{}{
			"model": "claude-3-opus-20240229",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "What is the weather today?",
				},
			},
			"max_tokens": 50,
		}
		
		requestBody, err := json.Marshal(openAIRequest)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertOpenAIToAnthropic(requestBody)
		
		// Assert
		require.NoError(t, err)
		
		var anthropicRequest map[string]interface{}
		err = json.Unmarshal(result, &anthropicRequest)
		require.NoError(t, err)
		
		// Check system field (should only have Claude Code prompt)
		expectedSystem := []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
		}
		assert.Equal(t, expectedSystem, anthropicRequest["system"])
		
		// Check messages remain unchanged
		messages := anthropicRequest["messages"].([]interface{})
		assert.Len(t, messages, 1)
		assert.Equal(t, "user", messages[0].(map[string]interface{})["role"])
		assert.Equal(t, "What is the weather today?", messages[0].(map[string]interface{})["content"])
	})
	
	t.Run("should remove anthropic/ prefix from model names", func(t *testing.T) {
		// Arrange
		openAIRequest := map[string]interface{}{
			"model": "anthropic/claude-3-5-sonnet-20241022",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		
		requestBody, err := json.Marshal(openAIRequest)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertOpenAIToAnthropic(requestBody)
		
		// Assert
		require.NoError(t, err)
		
		var anthropicRequest map[string]interface{}
		err = json.Unmarshal(result, &anthropicRequest)
		require.NoError(t, err)
		
		assert.Equal(t, "claude-3-5-sonnet-20241022", anthropicRequest["model"])
	})
	
	t.Run("should handle multiple system messages by concatenating them", func(t *testing.T) {
		// Arrange
		openAIRequest := map[string]interface{}{
			"model": "claude-3-opus-20240229",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				map[string]interface{}{
					"role":    "system",
					"content": "Always be polite and professional.",
				},
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		
		requestBody, err := json.Marshal(openAIRequest)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertOpenAIToAnthropic(requestBody)
		
		// Assert
		require.NoError(t, err)
		
		var anthropicRequest map[string]interface{}
		err = json.Unmarshal(result, &anthropicRequest)
		require.NoError(t, err)
		
		// Check system field has all system messages
		expectedSystem := []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
			map[string]interface{}{"type": "text", "text": "You are a helpful assistant."},
			map[string]interface{}{"type": "text", "text": "Always be polite and professional."},
		}
		assert.Equal(t, expectedSystem, anthropicRequest["system"])
		
		// Check only user message remains in messages
		messages := anthropicRequest["messages"].([]interface{})
		assert.Len(t, messages, 1)
		assert.Equal(t, "user", messages[0].(map[string]interface{})["role"])
	})
	
	t.Run("should handle complex content format from Cursor", func(t *testing.T) {
		// Arrange - mimicking Cursor's actual request format
		openAIRequest := map[string]interface{}{
			"model": "anthropic/claude-3-5-sonnet-20241022",
			"messages": []interface{}{
				map[string]interface{}{
					"role": "system",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "You are an AI programming assistant.",
							"cache_control": map[string]interface{}{
								"type": "ephemeral",
							},
						},
					},
				},
				map[string]interface{}{
					"role": "user",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "Please help me fix this code",
						},
					},
				},
				map[string]interface{}{
					"role": "assistant",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "I'll help you fix the code.",
						},
					},
				},
			},
			"max_tokens": 4096,
			"temperature": 0.2,
			"tools": []interface{}{
				map[string]interface{}{
					"type": "function",
					"function": map[string]interface{}{
						"name": "str_replace_editor",
						"description": "Replace text in a file",
					},
				},
			},
		}
		
		requestBody, err := json.Marshal(openAIRequest)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertOpenAIToAnthropic(requestBody)
		
		// Assert
		require.NoError(t, err)
		
		var anthropicRequest map[string]interface{}
		err = json.Unmarshal(result, &anthropicRequest)
		require.NoError(t, err)
		
		// Check model name (prefix removed)
		assert.Equal(t, "claude-3-5-sonnet-20241022", anthropicRequest["model"])
		
		// Check system field extracted correctly
		expectedSystem := []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
			map[string]interface{}{"type": "text", "text": "You are an AI programming assistant."},
		}
		assert.Equal(t, expectedSystem, anthropicRequest["system"])
		
		// Check messages preserved structured content
		messages := anthropicRequest["messages"].([]interface{})
		assert.Len(t, messages, 2) // user and assistant messages
		
		// Check user message
		userMsg := messages[0].(map[string]interface{})
		assert.Equal(t, "user", userMsg["role"])
		userContent := userMsg["content"].([]interface{})
		assert.Equal(t, "text", userContent[0].(map[string]interface{})["type"])
		assert.Equal(t, "Please help me fix this code", userContent[0].(map[string]interface{})["text"])
		
		// Check assistant message
		assistantMsg := messages[1].(map[string]interface{})
		assert.Equal(t, "assistant", assistantMsg["role"])
		
		// Check tools are preserved
		assert.Equal(t, openAIRequest["tools"], anthropicRequest["tools"])
	})
}

func TestConvertAnthropicToOpenAI(t *testing.T) {
	t.Run("should convert Anthropic response to OpenAI format", func(t *testing.T) {
		// Arrange
		anthropicResponse := map[string]interface{}{
			"id":      "msg_123",
			"type":    "message",
			"role":    "assistant",
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Hello! How can I help you today?",
				},
			},
			"model":         "claude-3-opus-20240229",
			"stop_reason":   "end_turn",
			"stop_sequence": nil,
			"usage": map[string]interface{}{
				"input_tokens":  10,
				"output_tokens": 20,
			},
		}
		
		responseBody, err := json.Marshal(anthropicResponse)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertAnthropicToOpenAI(responseBody)
		
		// Assert
		require.NoError(t, err)
		
		var openAIResponse map[string]interface{}
		err = json.Unmarshal(result, &openAIResponse)
		require.NoError(t, err)
		
		// Check OpenAI format fields
		assert.Equal(t, "msg_123", openAIResponse["id"])
		assert.Equal(t, "chat.completion", openAIResponse["object"])
		assert.Equal(t, "claude-3-opus-20240229", openAIResponse["model"])
		
		// Check choices array
		choices := openAIResponse["choices"].([]interface{})
		assert.Len(t, choices, 1)
		
		choice := choices[0].(map[string]interface{})
		assert.Equal(t, float64(0), choice["index"])
		assert.Equal(t, "stop", choice["finish_reason"])
		
		// Check message in choice
		message := choice["message"].(map[string]interface{})
		assert.Equal(t, "assistant", message["role"])
		assert.Equal(t, "Hello! How can I help you today?", message["content"])
		
		// Check usage
		usage := openAIResponse["usage"].(map[string]interface{})
		assert.Equal(t, float64(10), usage["prompt_tokens"])
		assert.Equal(t, float64(20), usage["completion_tokens"])
		assert.Equal(t, float64(30), usage["total_tokens"])
	})
	
	t.Run("should convert Anthropic error response to OpenAI format", func(t *testing.T) {
		// Arrange
		anthropicError := map[string]interface{}{
			"error": map[string]interface{}{
				"type":    "invalid_request_error",
				"message": "messages: at least one message is required",
			},
		}
		
		errorBody, err := json.Marshal(anthropicError)
		require.NoError(t, err)
		
		// Act
		result, err := ConvertAnthropicToOpenAI(errorBody)
		
		// Assert
		require.NoError(t, err)
		
		var openAIError map[string]interface{}
		err = json.Unmarshal(result, &openAIError)
		require.NoError(t, err)
		
		// Check error format
		errorObj := openAIError["error"].(map[string]interface{})
		assert.Equal(t, "messages: at least one message is required", errorObj["message"])
		assert.Equal(t, "invalid_request_error", errorObj["type"])
		assert.Nil(t, errorObj["param"])
		assert.Nil(t, errorObj["code"])
	})
}

func TestConvertAnthropicSSEToOpenAI(t *testing.T) {
	messageID := "chatcmpl-test123"
	model := "claude-3-opus-20240229"
	created := int64(1719331200)
	
	t.Run("should convert message_start event", func(t *testing.T) {
		// Arrange
		event := "message_start"
		data := `{"type":"message_start","message":{"id":"msg_123","type":"message","role":"assistant","model":"claude-3-opus-20240229","content":[],"stop_reason":null}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"object":"chat.completion.chunk"`)
		assert.Contains(t, result, `"role":"assistant"`)
		assert.Contains(t, result, messageID)
	})
	
	t.Run("should convert content_block_delta event", func(t *testing.T) {
		// Arrange
		event := "content_block_delta"
		data := `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello world"}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"content":"Hello world"`)
		assert.Contains(t, result, `"finish_reason":null`)
	})
	
	t.Run("should convert message_stop event without DONE", func(t *testing.T) {
		// Arrange
		event := "message_stop"
		data := `{"type":"message_stop"}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"finish_reason":"stop"`)
		assert.NotContains(t, result, "[DONE]") // DONE should be sent separately
	})
	
	t.Run("should convert message_delta with stop reason", func(t *testing.T) {
		// Arrange
		event := "message_delta"
		data := `{"type":"message_delta","delta":{"stop_reason":"max_tokens"}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"finish_reason":"length"`)
	})
	
	t.Run("should skip unhandled events", func(t *testing.T) {
		// Arrange
		event := "ping"
		data := `{"type":"ping"}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Empty(t, result)
	})
	
	t.Run("should convert input_json_delta events to OpenAI tool format", func(t *testing.T) {
		// Arrange
		event := "content_block_delta"
		data := `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"location\": \"San Fra"}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"tool_calls"`)
		assert.Contains(t, result, `"function"`)
		assert.Contains(t, result, `"arguments":"{\"location\": \"San Fra"`)
	})
	
	t.Run("should handle content_block_start for tool_use", func(t *testing.T) {
		// Arrange
		event := "content_block_start"
		data := `{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_123","name":"get_weather"}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"tool_calls"`)
		assert.Contains(t, result, `"id":"toolu_123"`)
		assert.Contains(t, result, `"type":"function"`)
		assert.Contains(t, result, `"name":"get_weather"`)
	})
	
	t.Run("should handle empty tool input gracefully", func(t *testing.T) {
		// Arrange
		event := "content_block_delta"
		data := `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":""}}`
		
		// Act
		result, err := ConvertAnthropicSSEToOpenAI(event, data, messageID, model, created)
		
		// Assert
		require.NoError(t, err)
		assert.Contains(t, result, "data: ")
		assert.Contains(t, result, `"arguments":""`)
	})
	
	t.Run("should handle multiple tool deltas in sequence", func(t *testing.T) {
		// Arrange - simulating a sequence of tool use events
		events := []struct {
			event string
			data  string
		}{
			{
				event: "content_block_start",
				data:  `{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_456","name":"calculate"}}`,
			},
			{
				event: "content_block_delta",
				data:  `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"operation\": \"add\","}}`,
			},
			{
				event: "content_block_delta",
				data:  `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":" \"a\": 5, \"b\": 3}"}}`,
			},
		}
		
		// Act & Assert - each event should produce valid output
		for _, tc := range events {
			result, err := ConvertAnthropicSSEToOpenAI(tc.event, tc.data, messageID, model, created)
			require.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "data: ")
		}
	})
}