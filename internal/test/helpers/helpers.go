package helpers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TokenResponse represents an OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// CreateMockOAuthServer creates a mock OAuth server for testing
func CreateMockOAuthServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/authorize":
			// Mock authorization endpoint
			code := r.URL.Query().Get("code")
			if code == "" {
				code = "test-auth-code"
			}
			redirectURI := r.URL.Query().Get("redirect_uri")
			http.Redirect(w, r, redirectURI+"?code="+code, http.StatusFound)

		case "/oauth/token":
			// Mock token endpoint
			w.Header().Set("Content-Type", "application/json")
			
			// Check for refresh token
			if r.FormValue("grant_type") == "refresh_token" {
				refreshToken := r.FormValue("refresh_token")
				if refreshToken == "invalid-refresh-token" {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "invalid_grant",
					})
					return
				}
			}

			// Return successful token response
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:  "test-access-token-" + time.Now().Format("20060102150405"),
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				RefreshToken: "test-refresh-token",
			})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// CreateMockAPIServer creates a mock Claude API server for testing
func CreateMockAPIServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || authHeader == "Bearer invalid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}

		// Check user agent for Claude Code identification
		userAgent := r.Header.Get("User-Agent")
		if userAgent != "Claude-Code/1.0" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid_client",
			})
			return
		}

		// Mock successful API response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      "msg_test123",
			"type":    "message",
			"content": []map[string]string{
				{"type": "text", "text": "Hello from mock Claude API"},
			},
		})
	}))
}

// WaitForServer waits for a server to be ready
func WaitForServer(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusServiceUnavailable {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return http.ErrServerClosed
}

// GetFreePort returns a free port for testing
func GetFreePort() string {
	// Using port 0 lets the OS assign a free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return fmt.Sprintf("%d", port)
}

// CleanupTestTokens removes any test tokens from storage
func CleanupTestTokens(t *testing.T) {
	t.Helper()
	// Implementation depends on storage backend
	// This is a placeholder that should be implemented based on actual storage
}