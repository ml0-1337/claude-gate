package proxy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
)

// TokenProvider interface for OAuth token management
type TokenProvider interface {
	GetAccessToken() (string, error)
}

// ProxyConfig holds configuration for the proxy handler
type ProxyConfig struct {
	UpstreamURL   string
	TokenProvider TokenProvider
	Transformer   *RequestTransformer
	Timeout       time.Duration
}

// ProxyHandler handles HTTP requests and proxies them to Anthropic API
type ProxyHandler struct {
	config     *ProxyConfig
	httpClient *http.Client
}

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(config *ProxyConfig) *ProxyHandler {
	if config.Timeout == 0 {
		config.Timeout = 600 * time.Second // 10 minutes default
	}
	
	// Create HTTP client with custom transport for better streaming support
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true, // Important for SSE
	}
	
	return &ProxyHandler{
		config: config,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		},
	}
}

// ServeHTTP implements http.Handler interface
func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight requests
	if r.Method == "OPTIONS" {
		h.handleCORS(w, r)
		return
	}
	
	// Set CORS headers for all requests
	h.setCORSHeaders(w, r)
	
	// Get OAuth token
	token, err := h.config.TokenProvider.GetAccessToken()
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "OAuth token error", err.Error())
		return
	}
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}
	defer r.Body.Close()
	
	// Transform request body if needed
	path := r.URL.Path
	transformedBody, err := h.config.Transformer.TransformRequestBody(body, path)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to transform request", err.Error())
		return
	}
	
	// Transform path for OpenAI endpoints
	upstreamPath := path
	if path == "/v1/chat/completions" {
		upstreamPath = "/v1/messages"
	}
	
	// Build upstream URL
	upstreamURL, err := url.Parse(h.config.UpstreamURL)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Invalid upstream URL", err.Error())
		return
	}
	upstreamURL.Path = upstreamPath
	upstreamURL.RawQuery = r.URL.RawQuery
	
	// Create upstream request
	upstreamReq, err := http.NewRequest(r.Method, upstreamURL.String(), bytes.NewReader(transformedBody))
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create upstream request", err.Error())
		return
	}
	
	// Inject OAuth headers
	upstreamReq.Header = h.config.Transformer.InjectHeaders(r.Header, token)
	
	// Make upstream request
	resp, err := h.httpClient.Do(upstreamReq)
	if err != nil {
		h.writeError(w, http.StatusBadGateway, "Upstream request failed", err.Error())
		return
	}
	defer resp.Body.Close()
	
	// Check if this is a streaming response
	isStreaming := strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") ||
		strings.Contains(r.URL.RawQuery, "stream=true")
	
	// Also check the request body for stream parameter
	if !isStreaming && len(body) > 0 {
		var reqData map[string]interface{}
		if err := json.Unmarshal(body, &reqData); err == nil {
			if stream, ok := reqData["stream"].(bool); ok && stream {
				isStreaming = true
			}
		}
	}
	
	// Handle response body
	if isStreaming {
		// Copy response headers for streaming
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		
		// Write status code
		w.WriteHeader(resp.StatusCode)
		
		// For OpenAI endpoints, convert SSE format
		if path == "/v1/chat/completions" {
			// Log that we're converting for OpenAI
			if h.config.Transformer != nil {
				// Using a simple log for debugging
				json.NewEncoder(io.Discard).Encode(map[string]string{
					"debug": "Converting Anthropic SSE to OpenAI format for streaming",
					"path":  path,
				})
			}
			h.streamOpenAIResponse(w, resp, path)
		} else {
			// For SSE, we need to flush after each write
			h.streamResponse(w, resp)
		}
	} else {
		// For OpenAI endpoints, transform response back
		if path == "/v1/chat/completions" {
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				h.writeError(w, http.StatusInternalServerError, "Failed to read response", err.Error())
				return
			}
			
			// Transform Anthropic response to OpenAI format
			transformedResp, err := h.config.Transformer.TransformResponseBody(respBody, path)
			if err != nil {
				// If transformation fails, return original
				// Copy headers excluding Content-Length
				for key, values := range resp.Header {
					if strings.ToLower(key) != "content-length" {
						for _, value := range values {
							w.Header().Add(key, value)
						}
					}
				}
				w.WriteHeader(resp.StatusCode)
				w.Write(respBody)
				return
			}
			
			// Copy headers excluding Content-Length and Content-Encoding
			for key, values := range resp.Header {
				if strings.ToLower(key) != "content-length" && strings.ToLower(key) != "content-encoding" {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
			}
			
			// Write status code
			w.WriteHeader(resp.StatusCode)
			
			// Write transformed response (Go will set correct Content-Length)
			w.Write(transformedResp)
		} else {
			// Regular response - copy headers and body
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			
			// Write status code
			w.WriteHeader(resp.StatusCode)
			
			// Just copy
			io.Copy(w, resp.Body)
		}
	}
}

// streamResponse handles Server-Sent Events streaming
func (h *ProxyHandler) streamResponse(w http.ResponseWriter, resp *http.Response) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		// Fallback to regular copy if flusher not available
		io.Copy(w, resp.Body)
		return
	}
	
	// Create a custom writer that flushes after each write
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			flusher.Flush()
		}
		if err != nil {
			if err != io.EOF {
				// Log error but don't write it to response
				// as we're already in the middle of streaming
			}
			return
		}
	}
}

// streamOpenAIResponse converts Anthropic SSE to OpenAI SSE format
func (h *ProxyHandler) streamOpenAIResponse(w http.ResponseWriter, resp *http.Response, path string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		// Fallback to regular streaming if flusher not available
		h.streamResponse(w, resp)
		return
	}
	
	// Generate message ID and timestamp for consistency
	messageID := "chatcmpl-" + generateRandomID()
	created := time.Now().Unix()
	model := "claude-3-5-sonnet-20241022" // Default model
	
	scanner := bufio.NewScanner(resp.Body)
	var currentEvent string
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// Extract model from message_start if available
			if currentEvent == "message_start" {
				var msgData map[string]interface{}
				if err := json.Unmarshal([]byte(data), &msgData); err == nil {
					if msg, ok := msgData["message"].(map[string]interface{}); ok {
						if m, ok := msg["model"].(string); ok {
							model = m
						}
					}
				}
			}
			
			// Convert the SSE event
			converted, err := ConvertAnthropicSSEToOpenAI(currentEvent, data, messageID, model, created)
			if err == nil && converted != "" {
				w.Write([]byte(converted))
				flusher.Flush()
			}
		}
	}
}

// generateRandomID generates a random ID for OpenAI format
func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 29)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// writeError writes an error response in Anthropic's error format
func (h *ProxyHandler) writeError(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResp := map[string]interface{}{
		"error": map[string]interface{}{
			"type":    errorType,
			"message": message,
		},
	}
	
	json.NewEncoder(w).Encode(errorResp)
}

// setCORSHeaders sets CORS headers for all responses
func (h *ProxyHandler) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")
}

// handleCORS handles CORS preflight requests
func (h *ProxyHandler) handleCORS(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w, r)
	w.WriteHeader(http.StatusNoContent)
}

// ProxyServer wraps the handler with additional server functionality
type ProxyServer struct {
	handler *ProxyHandler
	server  *http.Server
}

// NewProxyServer creates a new proxy server with health endpoints
func NewProxyServer(config *ProxyConfig, addr string, storage auth.StorageBackend) *ProxyServer {
	proxyHandler := NewProxyHandler(config)
	healthHandler := NewHealthHandler(storage)
	mux := CreateMux(proxyHandler, healthHandler)
	
	return &ProxyServer{
		handler: proxyHandler,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: config.Timeout + 10*time.Second, // Slightly more than request timeout
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Start starts the proxy server
func (s *ProxyServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop gracefully stops the proxy server
func (s *ProxyServer) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}