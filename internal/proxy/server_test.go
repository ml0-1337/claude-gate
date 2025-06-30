package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock storage backend for testing
type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Get(provider string) (*auth.TokenInfo, error) {
	args := m.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenInfo), args.Error(1)
}

func (m *mockStorage) Set(provider string, token *auth.TokenInfo) error {
	args := m.Called(provider, token)
	return args.Error(0)
}

func (m *mockStorage) Remove(provider string) error {
	args := m.Called(provider)
	return args.Error(0)
}

func (m *mockStorage) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockStorage) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockStorage) RequiresUnlock() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockStorage) Unlock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockStorage) Lock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockStorage) Name() string {
	args := m.Called()
	return args.String(0)
}

// Test 6: HealthHandler should return 200 OK
func TestHealthHandler_Returns200OK(t *testing.T) {
	// Prediction: This test will pass - handler returns 200
	
	mockStorage := new(mockStorage)
	mockStorage.On("Get", "anthropic").Return(nil, nil)
	
	handler := NewHealthHandler(mockStorage)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test 7: HealthHandler should return proper JSON response
func TestHealthHandler_ReturnsProperJSON(t *testing.T) {
	// Prediction: This test will pass - handler returns JSON
	
	t.Run("without OAuth token", func(t *testing.T) {
		mockStorage := new(mockStorage)
		mockStorage.On("Get", "anthropic").Return(nil, nil)
		
		handler := NewHealthHandler(mockStorage)
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Check Content-Type
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Check response structure
		assert.Equal(t, "healthy", response["status"])
		assert.Equal(t, "not_configured", response["oauth_status"])
		assert.Equal(t, "disabled", response["proxy_auth"])
	})
	
	t.Run("with OAuth token", func(t *testing.T) {
		mockStorage := new(mockStorage)
		token := &auth.TokenInfo{Type: "oauth", AccessToken: "test"}
		mockStorage.On("Get", "anthropic").Return(token, nil)
		
		handler := NewHealthHandler(mockStorage)
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Check OAuth status
		assert.Equal(t, "ready", response["oauth_status"])
	})
	
	t.Run("with non-OAuth token", func(t *testing.T) {
		mockStorage := new(mockStorage)
		token := &auth.TokenInfo{Type: "api_key", AccessToken: "test"}
		mockStorage.On("Get", "anthropic").Return(token, nil)
		
		handler := NewHealthHandler(mockStorage)
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)
		
		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Should still be not_configured for non-OAuth tokens
		assert.Equal(t, "not_configured", response["oauth_status"])
	})
}

// Test 8: RootHandler should return 200 OK with JSON content type
func TestRootHandler_Returns200WithJSON(t *testing.T) {
	// Prediction: This test will pass - handler returns 200 with JSON
	
	handler := &RootHandler{}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

// Test 9: RootHandler should return correct proxy info structure
func TestRootHandler_ReturnsCorrectStructure(t *testing.T) {
	// Prediction: This test will pass - handler returns expected JSON
	
	handler := &RootHandler{}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Check response structure
	assert.Equal(t, "Claude OAuth Proxy", response["service"])
	assert.Contains(t, response, "description")
	assert.Contains(t, response, "endpoints")
	assert.Equal(t, true, response["oauth_required"])
	assert.Equal(t, "disabled", response["proxy_auth"])
	
	// Check endpoints
	endpoints := response["endpoints"].(map[string]interface{})
	assert.Equal(t, "/health", endpoints["health"])
	assert.Equal(t, "/*", endpoints["anthropic_api"])
}

// Test 10: CreateMux should register all expected routes
func TestCreateMux_RegistersAllRoutes(t *testing.T) {
	// Prediction: This test will pass - mux registers all routes
	
	// Create mock handlers
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("proxy"))
	})
	
	mockStorage := new(mockStorage)
	mockStorage.On("Get", "anthropic").Return(nil, nil)
	healthHandler := NewHealthHandler(mockStorage)
	
	// Create mux
	mux := CreateMux(proxyHandler, healthHandler)
	
	// Test health endpoint
	t.Run("health endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "healthy")
	})
	
	// Test root endpoint
	t.Run("root endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Claude OAuth Proxy")
	})
	
	// Test models endpoint
	t.Run("models endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/models", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "list", response["object"])
	})
	
	// Test proxy endpoint
	t.Run("proxy endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/chat/completions", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "proxy", w.Body.String())
	})
}