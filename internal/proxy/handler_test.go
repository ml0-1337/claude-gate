package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTokenProvider struct {
	token string
	err   error
}

func (m *mockTokenProvider) GetAccessToken() (string, error) {
	return m.token, m.err
}

func TestProxyHandler(t *testing.T) {
	t.Run("proxies request with transformed body and headers", func(t *testing.T) {
		// Create test server to act as Anthropic API
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify headers
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "oauth-2025-04-20", r.Header.Get("anthropic-beta"))
			assert.Equal(t, "2023-06-01", r.Header.Get("anthropic-version"))
			
			// Verify body transformation
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			require.NoError(t, err)
			
			// Check system prompt was transformed
			system, ok := data["system"].([]interface{})
			require.True(t, ok)
			require.Len(t, system, 2)
			
			first := system[0].(map[string]interface{})
			assert.Equal(t, ClaudeCodePrompt, first["text"])
			
			// Check model alias was mapped
			assert.Equal(t, "claude-3-5-sonnet-20241022", data["model"])
			
			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "msg_123",
				"content": []map[string]interface{}{
					{"type": "text", "text": "Hello from Claude"},
				},
			})
		}))
		defer upstream.Close()
		
		// Create proxy handler
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   upstream.URL,
			TokenProvider: &mockTokenProvider{token: "test-token"},
			Transformer:   NewRequestTransformer(),
		})
		
		// Create test request
		reqBody := map[string]interface{}{
			"model":  "claude-3-5-sonnet-latest",
			"system": "User prompt",
			"messages": []map[string]interface{}{
				{"role": "user", "content": "Hello"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Test/1.0") // Should be stripped
		
		// Execute request
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		
		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "msg_123", response["id"])
	})
	
	t.Run("handles streaming SSE responses", func(t *testing.T) {
		// Create test server with SSE response
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("X-Accel-Buffering", "no")
			
			flusher, ok := w.(http.Flusher)
			require.True(t, ok)
			
			// Send SSE events
			events := []string{
				"event: message_start\ndata: {\"type\":\"message_start\"}\n\n",
				"event: content_block_start\ndata: {\"type\":\"content_block_start\"}\n\n",
				"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\"Hello\"}}\n\n",
				"event: content_block_stop\ndata: {\"type\":\"content_block_stop\"}\n\n",
				"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n",
			}
			
			for _, event := range events {
				_, err := w.Write([]byte(event))
				require.NoError(t, err)
				flusher.Flush()
				time.Sleep(10 * time.Millisecond) // Simulate streaming delay
			}
		}))
		defer upstream.Close()
		
		// Create proxy handler
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   upstream.URL,
			TokenProvider: &mockTokenProvider{token: "test-token"},
			Transformer:   NewRequestTransformer(),
		})
		
		// Create streaming request
		reqBody := map[string]interface{}{
			"model":  "claude-3-5-sonnet-20241022",
			"stream": true,
			"messages": []map[string]interface{}{
				{"role": "user", "content": "Hello"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		
		// Execute request
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		
		// Verify response headers
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
		
		// Verify we received SSE events
		body := w.Body.String()
		assert.Contains(t, body, "event: message_start")
		assert.Contains(t, body, "event: content_block_delta")
		assert.Contains(t, body, "data: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\"Hello\"}}")
	})
	
	t.Run("handles token provider errors", func(t *testing.T) {
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   "http://example.com",
			TokenProvider: &mockTokenProvider{err: assert.AnError},
			Transformer:   NewRequestTransformer(),
		})
		
		req := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		errorObj, ok := response["error"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "OAuth token error", errorObj["type"])
	})
	
	t.Run("passes through non-messages endpoints without transformation", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify no system prompt transformation for non-messages endpoint
			body, _ := io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.Write(body) // Echo back
		}))
		defer upstream.Close()
		
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   upstream.URL,
			TokenProvider: &mockTokenProvider{token: "test-token"},
			Transformer:   NewRequestTransformer(),
		})
		
		reqBody := map[string]interface{}{
			"system": "Should not be transformed",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("GET", "/v1/models", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Should not be transformed", response["system"])
	})
	
	t.Run("handles upstream errors", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"type":    "invalid_request_error",
					"message": "Invalid model",
				},
			})
		}))
		defer upstream.Close()
		
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   upstream.URL,
			TokenProvider: &mockTokenProvider{token: "test-token"},
			Transformer:   NewRequestTransformer(),
		})
		
		req := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid_request_error", response["error"].(map[string]interface{})["type"])
	})
	
	t.Run("injects system prompt for requests without one", func(t *testing.T) {
		// Create test server to verify system prompt injection
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			require.NoError(t, err)
			
			// Verify system prompt was injected
			system, ok := data["system"].(string)
			require.True(t, ok)
			assert.Equal(t, ClaudeCodePrompt, system)
			
			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "msg_456",
				"content": []map[string]interface{}{
					{"type": "text", "text": "Response without system prompt"},
				},
			})
		}))
		defer upstream.Close()
		
		// Create proxy handler
		handler := NewProxyHandler(&ProxyConfig{
			UpstreamURL:   upstream.URL,
			TokenProvider: &mockTokenProvider{token: "test-token"},
			Transformer:   NewRequestTransformer(),
		})
		
		// Create test request WITHOUT system prompt
		reqBody := map[string]interface{}{
			"model": "claude-opus-4-20250514",
			"messages": []map[string]interface{}{
				{"role": "user", "content": "Hello, Claude!"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/v1/messages", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		
		// Execute request
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		
		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "msg_456", response["id"])
	})
}