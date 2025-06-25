package proxy

import (
	"encoding/json"
	"net/http"
	
	"github.com/ml0-1337/claude-gate/internal/auth"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	storage auth.StorageBackend
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(storage auth.StorageBackend) *HealthHandler {
	return &HealthHandler{
		storage: storage,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check OAuth status
	oauthStatus := "not_configured"
	if token, err := h.storage.Get("anthropic"); err == nil && token != nil {
		if token.Type == "oauth" {
			oauthStatus = "ready"
		}
	}
	
	response := map[string]interface{}{
		"status":       "healthy",
		"oauth_status": oauthStatus,
		"proxy_auth":   "disabled", // TODO: get from config
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RootHandler handles the root endpoint
type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":     "Claude OAuth Proxy",
		"description": "Anthropic API proxy with OAuth authentication injection",
		"endpoints": map[string]interface{}{
			"health":       "/health",
			"anthropic_api": "/*",
		},
		"oauth_required": true,
		"proxy_auth": "disabled", // TODO: get from config
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateMux creates the HTTP mux with all routes
func CreateMux(proxyHandler http.Handler, healthHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.Handle("/health", healthHandler)
	
	// Root endpoint
	mux.Handle("/", &RootHandler{})
	
	// Models endpoint for OpenAI compatibility
	mux.Handle("/v1/models", NewModelsHandler())
	
	// All other paths go to the proxy
	mux.Handle("/v1/", proxyHandler)
	
	return mux
}