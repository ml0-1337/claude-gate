package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// OAuthTokenProvider implements TokenProvider interface for the proxy
type OAuthTokenProvider struct {
	client      *OAuthClient
	storage     StorageBackend
	cachedToken *TokenInfo
	cacheMutex  sync.RWMutex
}

// NewOAuthTokenProvider creates a new OAuth token provider
func NewOAuthTokenProvider(storage StorageBackend) *OAuthTokenProvider {
	return &OAuthTokenProvider{
		client:  NewOAuthClient(),
		storage: storage,
	}
}

// GetAccessToken returns a valid access token, refreshing if necessary
func (p *OAuthTokenProvider) GetAccessToken() (string, error) {
	// First, check if we have a valid cached token
	p.cacheMutex.RLock()
	if p.cachedToken != nil && p.cachedToken.Type == "oauth" && !p.cachedToken.NeedsRefresh() {
		token := p.cachedToken.AccessToken
		p.cacheMutex.RUnlock()
		return token, nil
	}
	p.cacheMutex.RUnlock()
	
	// Need to fetch or refresh token - use write lock
	p.cacheMutex.Lock()
	defer p.cacheMutex.Unlock()
	
	// Double-check after acquiring write lock (another goroutine might have refreshed)
	if p.cachedToken != nil && p.cachedToken.Type == "oauth" && !p.cachedToken.NeedsRefresh() {
		return p.cachedToken.AccessToken, nil
	}
	
	// Fetch token from storage
	token, err := p.storage.Get("anthropic")
	if err != nil {
		return "", fmt.Errorf("failed to get token from storage: %w", err)
	}
	
	if token == nil || token.Type != "oauth" {
		return "", fmt.Errorf("no OAuth token found - please authenticate first")
	}
	
	// Check if token needs refresh
	if token.NeedsRefresh() {
		// Refresh the token
		newToken, err := p.client.RefreshToken(token.RefreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}
		
		// Update storage
		if err := p.storage.Set("anthropic", newToken); err != nil {
			return "", fmt.Errorf("failed to save refreshed token: %w", err)
		}
		
		// Update cache
		p.cachedToken = newToken
		return newToken.AccessToken, nil
	}
	
	// Update cache with the token from storage
	p.cachedToken = token
	return token.AccessToken, nil
}

// ExchangeCode exchanges an authorization code for tokens
func (c *OAuthClient) ExchangeCode(code, verifier string) (*TokenInfo, error) {
	// Parse code and state
	parsedCode, parsedState := c.parseCodeAndState(code)
	
	// Build request body
	reqBody := map[string]interface{}{
		"code":          parsedCode,
		"grant_type":    "authorization_code",
		"client_id":     c.ClientID,
		"redirect_uri":  c.RedirectURI,
		"code_verifier": verifier,
	}
	
	// Include state if present
	if parsedState != "" {
		reqBody["state"] = parsedState
	}
	
	// Make request
	return c.makeTokenRequest(reqBody)
}

// RefreshToken refreshes an access token using a refresh token
func (c *OAuthClient) RefreshToken(refreshToken string) (*TokenInfo, error) {
	reqBody := map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     c.ClientID,
	}
	
	return c.makeTokenRequest(reqBody)
}

// makeTokenRequest makes a token request to the OAuth server
func (c *OAuthClient) makeTokenRequest(body map[string]interface{}) (*TokenInfo, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	req, err := http.NewRequestWithContext(context.Background(), "POST", c.TokenURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("token request failed: %v", errorResp)
		}
		return nil, fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}
	
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	return &TokenInfo{
		Type:         "oauth",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(tokenResp.ExpiresIn),
	}, nil
}