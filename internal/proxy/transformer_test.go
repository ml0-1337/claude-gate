package proxy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const claudeCodePrompt = "You are Claude Code, Anthropic's official CLI for Claude."

func TestSystemPromptTransformation(t *testing.T) {
	transformer := NewRequestTransformer()
	
	t.Run("converts string system prompt to array", func(t *testing.T) {
		input := map[string]interface{}{
			"system": "Custom user prompt",
			"model":  "claude-3-5-sonnet-20241022",
		}
		
		body, err := json.Marshal(input)
		require.NoError(t, err)
		
		result, err := transformer.TransformSystemPrompt(body)
		require.NoError(t, err)
		
		var output map[string]interface{}
		err = json.Unmarshal(result, &output)
		require.NoError(t, err)
		
		system, ok := output["system"].([]interface{})
		require.True(t, ok)
		require.Len(t, system, 2)
		
		// First element should be Claude Code
		first := system[0].(map[string]interface{})
		assert.Equal(t, "text", first["type"])
		assert.Equal(t, claudeCodePrompt, first["text"])
		
		// Second element should be original prompt
		second := system[1].(map[string]interface{})
		assert.Equal(t, "text", second["type"])
		assert.Equal(t, "Custom user prompt", second["text"])
	})
	
	t.Run("leaves string prompt alone if already Claude Code", func(t *testing.T) {
		input := map[string]interface{}{
			"system": claudeCodePrompt,
			"model":  "claude-3-5-sonnet-20241022",
		}
		
		body, err := json.Marshal(input)
		require.NoError(t, err)
		
		result, err := transformer.TransformSystemPrompt(body)
		require.NoError(t, err)
		
		var output map[string]interface{}
		err = json.Unmarshal(result, &output)
		require.NoError(t, err)
		
		// Should remain as string
		system, ok := output["system"].(string)
		require.True(t, ok)
		assert.Equal(t, claudeCodePrompt, system)
	})
	
	t.Run("prepends to array if Claude Code not first", func(t *testing.T) {
		input := map[string]interface{}{
			"system": []interface{}{
				map[string]interface{}{"type": "text", "text": "User prompt 1"},
				map[string]interface{}{"type": "text", "text": "User prompt 2"},
			},
			"model": "claude-3-5-sonnet-20241022",
		}
		
		body, err := json.Marshal(input)
		require.NoError(t, err)
		
		result, err := transformer.TransformSystemPrompt(body)
		require.NoError(t, err)
		
		var output map[string]interface{}
		err = json.Unmarshal(result, &output)
		require.NoError(t, err)
		
		system, ok := output["system"].([]interface{})
		require.True(t, ok)
		require.Len(t, system, 3)
		
		// First element should be Claude Code
		first := system[0].(map[string]interface{})
		assert.Equal(t, "text", first["type"])
		assert.Equal(t, claudeCodePrompt, first["text"])
		
		// Original prompts should follow
		second := system[1].(map[string]interface{})
		assert.Equal(t, "User prompt 1", second["text"])
		
		third := system[2].(map[string]interface{})
		assert.Equal(t, "User prompt 2", third["text"])
	})
	
	t.Run("leaves array alone if Claude Code already first", func(t *testing.T) {
		input := map[string]interface{}{
			"system": []interface{}{
				map[string]interface{}{"type": "text", "text": claudeCodePrompt},
				map[string]interface{}{"type": "text", "text": "User prompt"},
			},
			"model": "claude-3-5-sonnet-20241022",
		}
		
		body, err := json.Marshal(input)
		require.NoError(t, err)
		
		result, err := transformer.TransformSystemPrompt(body)
		require.NoError(t, err)
		
		var output map[string]interface{}
		err = json.Unmarshal(result, &output)
		require.NoError(t, err)
		
		system, ok := output["system"].([]interface{})
		require.True(t, ok)
		require.Len(t, system, 2)
		
		// Should remain unchanged
		first := system[0].(map[string]interface{})
		assert.Equal(t, claudeCodePrompt, first["text"])
	})
	
	t.Run("handles request without system prompt", func(t *testing.T) {
		input := map[string]interface{}{
			"model": "claude-3-5-sonnet-20241022",
			"messages": []interface{}{
				map[string]interface{}{"role": "user", "content": "Hello"},
			},
		}
		
		body, err := json.Marshal(input)
		require.NoError(t, err)
		
		result, err := transformer.TransformSystemPrompt(body)
		require.NoError(t, err)
		
		// Should remain unchanged
		assert.Equal(t, body, result)
	})
}

func TestModelAliasMapping(t *testing.T) {
	transformer := NewRequestTransformer()
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"claude-3-5-haiku-latest", "claude-3-5-haiku-20241022"},
		{"claude-3-5-sonnet-latest", "claude-3-5-sonnet-20241022"},
		{"claude-3-7-sonnet-latest", "claude-3-7-sonnet-20250219"},
		{"claude-3-opus-latest", "claude-3-opus-20240229"},
		{"claude-3-5-sonnet-20241022", "claude-3-5-sonnet-20241022"}, // No change
		{"gpt-4", "gpt-4"}, // Unknown model, no change
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := transformer.MapModelAlias(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHeaderInjection(t *testing.T) {
	transformer := NewRequestTransformer()
	
	t.Run("injects OAuth headers and strips problematic ones", func(t *testing.T) {
		headers := map[string][]string{
			"User-Agent":   {"Zed/0.191.7 (macos; aarch64)"}, // Should be stripped
			"Content-Type": {"application/json"},
			"Accept":       {"application/json"},
			"X-Custom":     {"value"},
		}
		
		token := "test-access-token"
		result := transformer.InjectHeaders(headers, token)
		
		// Check OAuth headers are added
		assert.Equal(t, "Bearer test-access-token", result.Get("Authorization"))
		assert.Equal(t, "oauth-2025-04-20", result.Get("anthropic-beta"))
		assert.Equal(t, "2023-06-01", result.Get("anthropic-version"))
		
		// Check content headers are preserved
		assert.Equal(t, "application/json", result.Get("Content-Type"))
		assert.Equal(t, "application/json", result.Get("Accept"))
		
		// Check User-Agent is stripped
		assert.Empty(t, result.Get("User-Agent"))
		
		// Check other headers are stripped
		assert.Empty(t, result.Get("X-Custom"))
	})
	
	t.Run("uses defaults for missing headers", func(t *testing.T) {
		headers := map[string][]string{}
		
		token := "test-access-token"
		result := transformer.InjectHeaders(headers, token)
		
		assert.Equal(t, "application/json", result.Get("Content-Type"))
		assert.Equal(t, "*/*", result.Get("Accept"))
	})
}