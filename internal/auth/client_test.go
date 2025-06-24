package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthTokenExchange(t *testing.T) {
	t.Run("successful token exchange", func(t *testing.T) {
		// Mock OAuth server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			
			// Verify request body
			assert.Equal(t, "authorization_code", body["grant_type"])
			assert.Equal(t, "test-code", body["code"])
			assert.Equal(t, "test-verifier", body["code_verifier"])
			
			// Send response
			resp := map[string]interface{}{
				"access_token":  "test-access-token",
				"refresh_token": "test-refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		
		client := NewOAuthClient()
		client.TokenURL = server.URL
		
		token, err := client.ExchangeCode("test-code", "test-verifier")
		require.NoError(t, err)
		assert.Equal(t, "oauth", token.Type)
		assert.Equal(t, "test-access-token", token.AccessToken)
		assert.Equal(t, "test-refresh-token", token.RefreshToken)
		assert.Greater(t, token.ExpiresAt, time.Now().Unix())
	})
	
	t.Run("handles code with state", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			
			assert.Equal(t, "test-code", body["code"])
			assert.Equal(t, "test-state", body["state"])
			
			resp := map[string]interface{}{
				"access_token":  "test-access-token",
				"refresh_token": "test-refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		
		client := NewOAuthClient()
		client.TokenURL = server.URL
		
		token, err := client.ExchangeCode("test-code#test-state", "test-verifier")
		require.NoError(t, err)
		assert.NotNil(t, token)
	})
	
	t.Run("handles token exchange error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":             "invalid_grant",
				"error_description": "Invalid authorization code",
			})
		}))
		defer server.Close()
		
		client := NewOAuthClient()
		client.TokenURL = server.URL
		
		token, err := client.ExchangeCode("bad-code", "test-verifier")
		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Contains(t, err.Error(), "token request failed")
	})
}

func TestOAuthTokenRefresh(t *testing.T) {
	t.Run("successful token refresh", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			
			assert.Equal(t, "refresh_token", body["grant_type"])
			assert.Equal(t, "old-refresh-token", body["refresh_token"])
			
			resp := map[string]interface{}{
				"access_token":  "new-access-token",
				"refresh_token": "new-refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		
		client := NewOAuthClient()
		client.TokenURL = server.URL
		
		token, err := client.RefreshToken("old-refresh-token")
		require.NoError(t, err)
		assert.Equal(t, "new-access-token", token.AccessToken)
		assert.Equal(t, "new-refresh-token", token.RefreshToken)
	})
}

func TestOAuthTokenProvider(t *testing.T) {
	t.Run("returns valid token", func(t *testing.T) {
		tempDir := t.TempDir()
		storage := NewTokenStorage(tempDir + "/auth.json")
		
		// Store a valid token
		validToken := &TokenInfo{
			Type:         "oauth",
			AccessToken:  "valid-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		}
		err := storage.Set("anthropic", validToken)
		require.NoError(t, err)
		
		provider := NewOAuthTokenProvider(storage)
		token, err := provider.GetAccessToken()
		require.NoError(t, err)
		assert.Equal(t, "valid-token", token)
	})
	
	t.Run("refreshes expired token", func(t *testing.T) {
		// Mock OAuth server for refresh
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]interface{}{
				"access_token":  "refreshed-token",
				"refresh_token": "new-refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		
		tempDir := t.TempDir()
		storage := NewTokenStorage(tempDir + "/auth.json")
		
		// Store an expired token
		expiredToken := &TokenInfo{
			Type:         "oauth",
			AccessToken:  "expired-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    time.Now().Add(-time.Hour).Unix(), // Expired
		}
		err := storage.Set("anthropic", expiredToken)
		require.NoError(t, err)
		
		provider := NewOAuthTokenProvider(storage)
		provider.client.TokenURL = server.URL
		
		token, err := provider.GetAccessToken()
		require.NoError(t, err)
		assert.Equal(t, "refreshed-token", token)
		
		// Verify token was saved
		savedToken, err := storage.Get("anthropic")
		require.NoError(t, err)
		assert.Equal(t, "refreshed-token", savedToken.AccessToken)
	})
	
	t.Run("returns error when no token", func(t *testing.T) {
		tempDir := t.TempDir()
		storage := NewTokenStorage(tempDir + "/auth.json")
		
		provider := NewOAuthTokenProvider(storage)
		token, err := provider.GetAccessToken()
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "no OAuth token found")
	})
}