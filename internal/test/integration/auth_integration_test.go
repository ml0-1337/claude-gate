//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/test/helpers"
)

func TestOAuthFlow_CompleteAuthentication(t *testing.T) {
	// Setup mock OAuth server
	mockServer := helpers.CreateMockOAuthServer(t)
	defer mockServer.Close()

	// Create custom OAuth client with mock server URLs
	client := &auth.OAuthClient{
		ClientID:     "test-client-id",
		AuthorizeURL: mockServer.URL + "/oauth/authorize",
		TokenURL:     mockServer.URL + "/oauth/token",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       "read write",
	}

	// Test authorization URL generation
	authData, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	assert.NotEmpty(t, authData.URL)
	assert.NotEmpty(t, authData.Verifier)
	assert.Contains(t, authData.URL, mockServer.URL)
	assert.Contains(t, authData.URL, "code_challenge")
	assert.Contains(t, authData.URL, "code_challenge_method=S256")
}

func TestOAuthFlow_TokenExchange(t *testing.T) {
	// Setup mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"access_token": "test-access-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"refresh_token": "test-refresh-token"
			}`))
		}
	}))
	defer mockServer.Close()

	// Create OAuth client
	client := &auth.OAuthClient{
		ClientID:    "test-client-id",
		TokenURL:    mockServer.URL + "/oauth/token",
		RedirectURI: "http://localhost:8080/callback",
	}

	// Test token exchange
	token, err := client.ExchangeCode("test-code", "test-verifier")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "oauth", token.Type)
	assert.Equal(t, "test-access-token", token.AccessToken)
	assert.Equal(t, "test-refresh-token", token.RefreshToken)
	assert.True(t, token.ExpiresAt > time.Now().Unix())
}

func TestOAuthFlow_TokenRefresh(t *testing.T) {
	// Setup mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			// Parse the JSON body to check grant_type
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			
			if body["grant_type"] == "refresh_token" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{
					"access_token": "new-access-token",
					"token_type": "Bearer",
					"expires_in": 3600,
					"refresh_token": "new-refresh-token"
				}`))
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer mockServer.Close()

	// Create OAuth client
	client := &auth.OAuthClient{
		ClientID: "test-client-id",
		TokenURL: mockServer.URL + "/oauth/token",
	}

	// Test token refresh
	newToken, err := client.RefreshToken("test-refresh-token")
	require.NoError(t, err)
	assert.NotNil(t, newToken)
	assert.Equal(t, "new-access-token", newToken.AccessToken)
	assert.Equal(t, "new-refresh-token", newToken.RefreshToken)
}

func TestOAuthTokenProvider_Integration(t *testing.T) {
	// Create mock storage
	storage := &mockStorage{
		tokens: make(map[string]*auth.TokenInfo),
	}

	// Create token provider
	provider := auth.NewOAuthTokenProvider(storage)

	// Test getting token when none exists
	_, err := provider.GetAccessToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no OAuth token found")

	// Store a valid token
	validToken := &auth.TokenInfo{
		Type:        "oauth",
		AccessToken: "valid-access-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour).Unix(),
	}
	storage.Set("anthropic", validToken)

	// Test getting valid token
	token, err := provider.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, "valid-access-token", token)
}

// mockStorage implements StorageBackend for testing
type mockStorage struct {
	tokens map[string]*auth.TokenInfo
}

func (m *mockStorage) Get(provider string) (*auth.TokenInfo, error) {
	token, ok := m.tokens[provider]
	if !ok {
		return nil, nil
	}
	return token, nil
}

func (m *mockStorage) Set(provider string, token *auth.TokenInfo) error {
	m.tokens[provider] = token
	return nil
}

func (m *mockStorage) Remove(provider string) error {
	delete(m.tokens, provider)
	return nil
}

func (m *mockStorage) List() ([]string, error) {
	var providers []string
	for provider := range m.tokens {
		providers = append(providers, provider)
	}
	return providers, nil
}

func (m *mockStorage) IsAvailable() bool {
	return true
}

func (m *mockStorage) RequiresUnlock() bool {
	return false
}

func (m *mockStorage) Unlock() error {
	return nil
}

func (m *mockStorage) Lock() error {
	return nil
}

func (m *mockStorage) Name() string {
	return "mock"
}