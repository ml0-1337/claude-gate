package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test 1: ModelsHandler should return 200 OK status
func TestModelsHandler_Returns200OK(t *testing.T) {
	// Prediction: This test will fail - handler not tested yet
	
	handler := NewModelsHandler()
	req := httptest.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test 2: ModelsHandler should return correct OpenAI models JSON structure
func TestModelsHandler_ReturnsCorrectJSONStructure(t *testing.T) {
	// Prediction: This test will fail - JSON structure not verified
	
	handler := NewModelsHandler()
	req := httptest.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Check top-level structure
	assert.Equal(t, "list", response["object"])
	assert.Contains(t, response, "data")
	
	// Check data is an array
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(data), 4) // At least 4 models
	
	// Check first model structure
	firstModel := data[0].(map[string]interface{})
	assert.Contains(t, firstModel, "id")
	assert.Equal(t, "model", firstModel["object"])
	assert.Contains(t, firstModel, "created")
	assert.Equal(t, "anthropic", firstModel["owned_by"])
	assert.Contains(t, firstModel, "permission")
}

// Test 3: ModelsHandler should set proper Content-Type header
func TestModelsHandler_SetsContentTypeHeader(t *testing.T) {
	// Prediction: This test will pass - Content-Type should be set
	
	handler := NewModelsHandler()
	req := httptest.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

// Test 4: ModelsHandler should handle CORS headers correctly
func TestModelsHandler_SetsCORSHeaders(t *testing.T) {
	// Prediction: This test will pass - CORS headers should be set
	
	handler := NewModelsHandler()
	req := httptest.NewRequest("GET", "/v1/models", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	// Check CORS headers
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "3600", w.Header().Get("Access-Control-Max-Age"))
}

// Test 5: setCORSHeadersStandalone should set all required CORS headers
func TestSetCORSHeadersStandalone(t *testing.T) {
	// Prediction: This test will pass - testing CORS function directly
	
	t.Run("with origin header", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/v1/models", nil)
		req.Header.Set("Origin", "https://app.example.com")
		w := httptest.NewRecorder()
		
		setCORSHeadersStandalone(w, req)
		
		assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "3600", w.Header().Get("Access-Control-Max-Age"))
	})
	
	t.Run("without origin header", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/v1/models", nil)
		// No Origin header set
		w := httptest.NewRecorder()
		
		setCORSHeadersStandalone(w, req)
		
		// Should default to "*"
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// Test OPTIONS request handling
func TestModelsHandler_HandlesOPTIONSRequest(t *testing.T) {
	// Prediction: This test will pass - OPTIONS should return 204
	
	handler := NewModelsHandler()
	req := httptest.NewRequest("OPTIONS", "/v1/models", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String()) // No content for OPTIONS
}