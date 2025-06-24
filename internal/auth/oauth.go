package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// OAuthClient handles OAuth authentication with Anthropic
type OAuthClient struct {
	ClientID     string
	AuthorizeURL string
	TokenURL     string
	RedirectURI  string
	Scopes       string
}

// AuthData contains authorization URL and PKCE verifier
type AuthData struct {
	URL      string
	Verifier string
}

// NewOAuthClient creates a new OAuth client with Anthropic configuration
func NewOAuthClient() *OAuthClient {
	return &OAuthClient{
		ClientID:     "9d1c250a-e61b-44d9-88ed-5944d1962f5e",
		AuthorizeURL: "https://claude.ai/oauth/authorize",
		TokenURL:     "https://console.anthropic.com/v1/oauth/token",
		RedirectURI:  "https://console.anthropic.com/oauth/code/callback",
		Scopes:       "org:create_api_key user:profile user:inference",
	}
}

// GeneratePKCE generates PKCE verifier and challenge for OAuth flow
func GeneratePKCE() (verifier, challenge string, err error) {
	// Generate 32 bytes of random data
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Encode verifier as base64url without padding
	verifier = base64.RawURLEncoding.EncodeToString(verifierBytes)
	
	// Generate challenge by SHA256 hashing the verifier
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	
	return verifier, challenge, nil
}

// GetAuthorizationURL generates the authorization URL with PKCE parameters
func (c *OAuthClient) GetAuthorizationURL() (*AuthData, error) {
	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE: %w", err)
	}
	
	params := url.Values{
		"code":                  {"true"},
		"client_id":             {c.ClientID},
		"response_type":         {"code"},
		"redirect_uri":          {c.RedirectURI},
		"scope":                 {c.Scopes},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"state":                 {verifier}, // Using verifier as state (following Python impl)
	}
	
	authURL := fmt.Sprintf("%s?%s", c.AuthorizeURL, params.Encode())
	
	return &AuthData{
		URL:      authURL,
		Verifier: verifier,
	}, nil
}

// parseCodeAndState parses the authorization code and state from the callback
func (c *OAuthClient) parseCodeAndState(code string) (parsedCode, parsedState string) {
	splits := strings.Split(code, "#")
	parsedCode = splits[0]
	if len(splits) > 1 {
		parsedState = splits[1]
	}
	return
}