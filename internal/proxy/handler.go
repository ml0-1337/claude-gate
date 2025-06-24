package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/yourusername/claude-gate/internal/auth"
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
	
	// Build upstream URL
	upstreamURL, err := url.Parse(h.config.UpstreamURL)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Invalid upstream URL", err.Error())
		return
	}
	upstreamURL.Path = path
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
	
	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	
	// Write status code
	w.WriteHeader(resp.StatusCode)
	
	// Handle response body
	if isStreaming {
		// For SSE, we need to flush after each write
		h.streamResponse(w, resp)
	} else {
		// Regular response - just copy
		io.Copy(w, resp.Body)
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

// ProxyServer wraps the handler with additional server functionality
type ProxyServer struct {
	handler *ProxyHandler
	server  *http.Server
}

// NewProxyServer creates a new proxy server with health endpoints
func NewProxyServer(config *ProxyConfig, addr string, storage *auth.TokenStorage) *ProxyServer {
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