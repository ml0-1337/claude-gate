package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// ClaudeCodePrompt is the required system prompt for OAuth authentication
	ClaudeCodePrompt = "You are Claude Code, Anthropic's official CLI for Claude."
)

// ModelAliases maps model aliases to their full names for OAuth compatibility
var ModelAliases = map[string]string{
	"claude-3-5-haiku-latest":  "claude-3-5-haiku-20241022",
	"claude-3-5-sonnet-latest": "claude-3-5-sonnet-20241022",
	"claude-3-7-sonnet-latest": "claude-3-7-sonnet-20250219",
	"claude-3-opus-latest":     "claude-3-opus-20240229",
}

// RequestTransformer handles request body and header transformations
type RequestTransformer struct{}

// NewRequestTransformer creates a new request transformer
func NewRequestTransformer() *RequestTransformer {
	return &RequestTransformer{}
}

// TransformSystemPrompt modifies the system prompt to ensure Claude Code identification comes first
func (t *RequestTransformer) TransformSystemPrompt(body []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body, nil // Return original if not JSON
	}
	
	// Check if request has a system prompt
	systemRaw, hasSystem := data["system"]
	if !hasSystem {
		// No system prompt, inject Claude Code identification
		data["system"] = ClaudeCodePrompt
		return json.Marshal(data)
	}
	
	switch system := systemRaw.(type) {
	case string:
		// Handle string system prompt
		if system == ClaudeCodePrompt {
			// Already correct, leave as-is
			return body, nil
		}
		// Convert to array with Claude Code first
		data["system"] = []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
			map[string]interface{}{"type": "text", "text": system},
		}
		
	case []interface{}:
		// Handle array system prompt
		if len(system) > 0 {
			// Check if first element has correct text
			if first, ok := system[0].(map[string]interface{}); ok {
				if text, ok := first["text"].(string); ok && text == ClaudeCodePrompt {
					// Already has Claude Code first, return as-is
					return body, nil
				}
			}
		}
		// Prepend Claude Code identification
		newSystem := []interface{}{
			map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
		}
		data["system"] = append(newSystem, system...)
	}
	
	// Re-marshal the modified data
	return json.Marshal(data)
}

// MapModelAlias maps model aliases to their full names
func (t *RequestTransformer) MapModelAlias(model string) string {
	if mapped, exists := ModelAliases[model]; exists {
		return mapped
	}
	return model
}

// TransformRequestBody applies all necessary transformations to the request body
func (t *RequestTransformer) TransformRequestBody(body []byte, path string) ([]byte, error) {
	// Handle OpenAI chat completions endpoint
	if path == "/v1/chat/completions" {
		// Convert OpenAI format to Anthropic format
		convertedBody, err := ConvertOpenAIToAnthropic(body)
		if err != nil {
			return nil, fmt.Errorf("failed to convert OpenAI format: %w", err)
		}
		
		// Apply standard transformations to the converted body
		return t.TransformRequestBody(convertedBody, "/v1/messages")
	}
	
	// Only transform messages endpoint
	if path != "/v1/messages" {
		return body, nil
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body, nil // Return original if not JSON
	}
	
	// Transform system prompt
	modifiedBody, err := t.TransformSystemPrompt(body)
	if err != nil {
		return nil, fmt.Errorf("failed to transform system prompt: %w", err)
	}
	
	// Re-unmarshal to apply model mapping
	if err := json.Unmarshal(modifiedBody, &data); err != nil {
		return modifiedBody, nil
	}
	
	// Map model alias if present
	if model, ok := data["model"].(string); ok {
		data["model"] = t.MapModelAlias(model)
	}
	
	return json.Marshal(data)
}

// InjectHeaders creates new headers with OAuth authentication and strips problematic ones
func (t *RequestTransformer) InjectHeaders(headers map[string][]string, accessToken string) http.Header {
	// Create fresh headers with only necessary ones
	newHeaders := http.Header{}
	newHeaders.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	newHeaders.Set("anthropic-beta", "oauth-2025-04-20")
	newHeaders.Set("anthropic-version", "2023-06-01")
	
	// Preserve content headers with defaults
	if contentType := getHeader(headers, "Content-Type"); contentType != "" {
		newHeaders.Set("Content-Type", contentType)
	} else {
		newHeaders.Set("Content-Type", "application/json")
	}
	
	if accept := getHeader(headers, "Accept"); accept != "" {
		newHeaders.Set("Accept", accept)
	} else {
		newHeaders.Set("Accept", "*/*")
	}
	
	// Preserve Connection header for proper SSE handling
	if connection := getHeader(headers, "Connection"); connection != "" {
		newHeaders.Set("Connection", connection)
	}
	
	// Preserve cache control headers for SSE
	if cacheControl := getHeader(headers, "Cache-Control"); cacheControl != "" {
		newHeaders.Set("Cache-Control", cacheControl)
	}
	
	return newHeaders
}

// getHeader performs case-insensitive header lookup
func getHeader(headers map[string][]string, key string) string {
	// Direct lookup
	if values, ok := headers[key]; ok && len(values) > 0 {
		return values[0]
	}
	
	// Case-insensitive lookup
	for k, v := range headers {
		if http.CanonicalHeaderKey(k) == http.CanonicalHeaderKey(key) && len(v) > 0 {
			return v[0]
		}
	}
	
	return ""
}

// TransformResponseBody transforms response body based on the endpoint
func (t *RequestTransformer) TransformResponseBody(body []byte, path string) ([]byte, error) {
	if path == "/v1/chat/completions" {
		// Convert Anthropic response to OpenAI format
		return ConvertAnthropicToOpenAI(body)
	}
	return body, nil
}