package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePKCE(t *testing.T) {
	t.Run("generates valid PKCE verifier and challenge", func(t *testing.T) {
		verifier, challenge, err := GeneratePKCE()
		require.NoError(t, err)
		
		// Verifier should be base64url encoded without padding
		assert.NotEmpty(t, verifier)
		assert.NotContains(t, verifier, "=")
		assert.NotContains(t, verifier, "+")
		assert.NotContains(t, verifier, "/")
		
		// Decode verifier to check it's 32 bytes
		decoded, err := base64.RawURLEncoding.DecodeString(verifier)
		require.NoError(t, err)
		assert.Len(t, decoded, 32)
		
		// Challenge should be base64url encoded without padding
		assert.NotEmpty(t, challenge)
		assert.NotContains(t, challenge, "=")
		assert.NotContains(t, challenge, "+")
		assert.NotContains(t, challenge, "/")
		
		// Challenge should be SHA256 of verifier
		challengeDecoded, err := base64.RawURLEncoding.DecodeString(challenge)
		require.NoError(t, err)
		assert.Len(t, challengeDecoded, 32) // SHA256 is 32 bytes
	})
	
	t.Run("generates different values each time", func(t *testing.T) {
		v1, c1, err := GeneratePKCE()
		require.NoError(t, err)
		
		v2, c2, err := GeneratePKCE()
		require.NoError(t, err)
		
		assert.NotEqual(t, v1, v2)
		assert.NotEqual(t, c1, c2)
	})
}

func TestOAuthClient(t *testing.T) {
	client := NewOAuthClient()
	
	t.Run("has correct configuration", func(t *testing.T) {
		assert.Equal(t, "9d1c250a-e61b-44d9-88ed-5944d1962f5e", client.ClientID)
		assert.Equal(t, "https://claude.ai/oauth/authorize", client.AuthorizeURL)
		assert.Equal(t, "https://console.anthropic.com/v1/oauth/token", client.TokenURL)
		assert.Equal(t, "https://console.anthropic.com/oauth/code/callback", client.RedirectURI)
		assert.Equal(t, "org:create_api_key user:profile user:inference", client.Scopes)
	})
	
	t.Run("generates authorization URL", func(t *testing.T) {
		authData, err := client.GetAuthorizationURL()
		require.NoError(t, err)
		
		assert.NotEmpty(t, authData.URL)
		assert.Contains(t, authData.URL, client.AuthorizeURL)
		assert.Contains(t, authData.URL, "client_id="+client.ClientID)
		assert.Contains(t, authData.URL, "response_type=code")
		assert.Contains(t, authData.URL, "redirect_uri=")
		assert.Contains(t, authData.URL, "scope=")
		assert.Contains(t, authData.URL, "code_challenge=")
		assert.Contains(t, authData.URL, "code_challenge_method=S256")
		assert.Contains(t, authData.URL, "state=")
		assert.NotEmpty(t, authData.Verifier)
	})
}

func TestTokenExchange(t *testing.T) {
	client := NewOAuthClient()
	
	t.Run("handles code with state", func(t *testing.T) {
		// This will be tested with integration tests
		// For now, we test the code parsing logic
		code := "test-code#test-state"
		parsedCode, parsedState := client.parseCodeAndState(code)
		
		assert.Equal(t, "test-code", parsedCode)
		assert.Equal(t, "test-state", parsedState)
	})
	
	t.Run("handles code without state", func(t *testing.T) {
		code := "test-code"
		parsedCode, parsedState := client.parseCodeAndState(code)
		
		assert.Equal(t, "test-code", parsedCode)
		assert.Empty(t, parsedState)
	})
}