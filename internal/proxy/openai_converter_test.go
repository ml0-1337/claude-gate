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
}