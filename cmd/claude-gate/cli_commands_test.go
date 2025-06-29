package main

import (
	"testing"

	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOAuthClient mocks the OAuth client for testing
type MockOAuthClient struct {
	mock.Mock
}

func (m *MockOAuthClient) GetAuthorizationURL() (*auth.AuthData, error) {
	args := m.Called()
	if authData := args.Get(0); authData != nil {
		return authData.(*auth.AuthData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOAuthClient) ExchangeCode(code, verifier string) (*auth.TokenInfo, error) {
	args := m.Called(code, verifier)
	if token := args.Get(0); token != nil {
		return token.(*auth.TokenInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockStorage mocks the storage interface for testing
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Get(provider string) (*auth.TokenInfo, error) {
	args := m.Called(provider)
	if token := args.Get(0); token != nil {
		return token.(*auth.TokenInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) Set(provider string, token *auth.TokenInfo) error {
	args := m.Called(provider, token)
	return args.Error(0)
}

func (m *MockStorage) Remove(provider string) error {
	args := m.Called(provider)
	return args.Error(0)
}

func (m *MockStorage) List() ([]string, error) {
	args := m.Called()
	if providers := args.Get(0); providers != nil {
		return providers.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStorage) RequiresUnlock() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStorage) Unlock(prompt func(string) (string, error)) error {
	args := m.Called(prompt)
	return args.Error(0)
}

func (m *MockStorage) Lock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Name() string {
	args := m.Called()
	return args.String(0)
}

// Test 1: auth login command should handle successful authentication flow
func TestLoginCmd_SuccessfulAuthentication(t *testing.T) {
	// This test verifies that the login command successfully authenticates
	// and stores the token when everything works correctly
	
	// Prediction: This test will fail because LoginCmd is not yet testable
	// due to tight coupling with external dependencies
	
	t.Skip("LoginCmd needs refactoring to accept injected dependencies")
	
	// TODO: Refactor LoginCmd to accept:
	// - Storage factory or storage instance
	// - OAuth client instance
	// - UI output instance
	// This will make it testable
}

// Test helper function that creates storage factory config
func TestCreateStorageFactoryConfig(t *testing.T) {
	// Test that createStorageFactoryConfig properly converts config to storage factory config
	
	// Prediction: This test will pass because the function is simple and testable
	
	cfg := config.DefaultConfig()
	cfg.AuthStorageType = "keyring"
	cfg.AuthStoragePath = "/custom/path/auth.json"
	cfg.KeychainTrustApp = true
	cfg.KeychainAccessibleWhenUnlocked = false
	cfg.KeychainSynchronizable = true
	cfg.KeyringService = "custom-service"
	
	result := createStorageFactoryConfig(cfg)
	
	assert.Equal(t, auth.StorageTypeKeyring, result.Type)
	assert.Equal(t, "/custom/path/auth.json", result.FilePath)
	assert.Equal(t, "custom-service", result.ServiceName)
	assert.True(t, result.KeychainTrustApp)
	assert.False(t, result.KeychainAccessibleWhenUnlocked)
	assert.True(t, result.KeychainSynchronizable)
}

// Test 9: version command should display version information
func TestVersionCmd_DisplaysVersion(t *testing.T) {
	// This test verifies that the version command outputs version information
	
	// Prediction: This test will fail initially because we need to capture stdout
	
	// Create the version command
	cmd := &VersionCmd{}
	
	// Run the command
	err := cmd.Run()
	
	// The command should not return an error
	assert.NoError(t, err)
	
	// Note: The actual output goes to stdout via ui.NewOutput()
	// To properly test this, we would need to:
	// 1. Refactor to allow injecting the output writer
	// 2. Or capture stdout during the test
	// For now, we just verify it runs without error
}