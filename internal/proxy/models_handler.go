package proxy

import (
	"encoding/json"
	"net/http"
	"time"
)

// ModelsHandler handles /v1/models requests for OpenAI compatibility
type ModelsHandler struct{}

// NewModelsHandler creates a new models handler
func NewModelsHandler() *ModelsHandler {
	return &ModelsHandler{}
}

// ServeHTTP handles the models endpoint
func (h *ModelsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle CORS
	if r.Method == "OPTIONS" {
		setCORSHeadersStandalone(w, r)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	setCORSHeadersStandalone(w, r)
	
	// Return available models in OpenAI format
	models := map[string]interface{}{
		"object": "list",
		"data": []interface{}{
			map[string]interface{}{
				"id":       "claude-3-opus-20240229",
				"object":   "model",
				"created":  1706745600,
				"owned_by": "anthropic",
				"permission": []interface{}{
					map[string]interface{}{
						"allow_create_engine":  false,
						"allow_fine_tuning":    false,
						"allow_logprobs":       false,
						"allow_sampling":       true,
						"allow_search_indices": false,
						"allow_view":           true,
						"created":              int(time.Now().Unix()),
						"group":                nil,
						"id":                   "modelperm-" + "claude-3-opus",
						"is_blocking":          false,
						"object":               "model_permission",
						"organization":         "*",
					},
				},
			},
			map[string]interface{}{
				"id":       "claude-3-5-sonnet-20241022",
				"object":   "model",
				"created":  1729555200,
				"owned_by": "anthropic",
				"permission": []interface{}{
					map[string]interface{}{
						"allow_create_engine":  false,
						"allow_fine_tuning":    false,
						"allow_logprobs":       false,
						"allow_sampling":       true,
						"allow_search_indices": false,
						"allow_view":           true,
						"created":              int(time.Now().Unix()),
						"group":                nil,
						"id":                   "modelperm-" + "claude-3-5-sonnet",
						"is_blocking":          false,
						"object":               "model_permission",
						"organization":         "*",
					},
				},
			},
			map[string]interface{}{
				"id":       "claude-3-5-haiku-20241022",
				"object":   "model",
				"created":  1729555200,
				"owned_by": "anthropic",
				"permission": []interface{}{
					map[string]interface{}{
						"allow_create_engine":  false,
						"allow_fine_tuning":    false,
						"allow_logprobs":       false,
						"allow_sampling":       true,
						"allow_search_indices": false,
						"allow_view":           true,
						"created":              int(time.Now().Unix()),
						"group":                nil,
						"id":                   "modelperm-" + "claude-3-5-haiku",
						"is_blocking":          false,
						"object":               "model_permission",
						"organization":         "*",
					},
				},
			},
			map[string]interface{}{
				"id":       "claude-opus-4-20250514",
				"object":   "model",
				"created":  1747353600,
				"owned_by": "anthropic",
				"permission": []interface{}{
					map[string]interface{}{
						"allow_create_engine":  false,
						"allow_fine_tuning":    false,
						"allow_logprobs":       false,
						"allow_sampling":       true,
						"allow_search_indices": false,
						"allow_view":           true,
						"created":              int(time.Now().Unix()),
						"group":                nil,
						"id":                   "modelperm-" + "claude-opus-4",
						"is_blocking":          false,
						"object":               "model_permission",
						"organization":         "*",
					},
				},
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

// setCORSHeadersStandalone is a standalone CORS header setter
func setCORSHeadersStandalone(w http.ResponseWriter, r *http.Request) {
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