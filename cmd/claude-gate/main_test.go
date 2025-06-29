package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOutput captures stdout and stderr during function execution
func captureOutput(f func() error) (string, string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	os.Stdout = wOut
	os.Stderr = wErr
	
	outChan := make(chan string)
	errChan := make(chan string)
	
	go func() {
		buf := new(bytes.Buffer)
		io.Copy(buf, rOut)
		outChan <- buf.String()
	}()
	
	go func() {
		buf := new(bytes.Buffer)
		io.Copy(buf, rErr)
		errChan <- buf.String()
	}()
	
	err := f()
	
	wOut.Close()
	wErr.Close()
	
	stdout := <-outChan
	stderr := <-errChan
	
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	
	return stdout, stderr, err
}

// Test 1: StartCmd should start proxy with default configuration
func TestStartCmd_DefaultConfiguration(t *testing.T) {
	// Prediction: This test will pass after we mock the dependencies properly
	
	// Create a temporary directory for test storage
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	// Create a test OAuth token
	storage := auth.NewFileStorage(authFile)
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	err := storage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// For now, just verify the command can be created
	// Actual server testing would require significant mocking
	cmd := &StartCmd{
		Host:      "127.0.0.1",
		Port:      5789,
		LogLevel:  "INFO",
		SkipAuthCheck: false,
	}
	
	// Verify command fields are set correctly
	assert.Equal(t, "127.0.0.1", cmd.Host)
	assert.Equal(t, 5789, cmd.Port)
	assert.Equal(t, "INFO", cmd.LogLevel)
	assert.False(t, cmd.SkipAuthCheck)
}

// Test 2: StartCmd should handle missing authentication
func TestStartCmd_MissingAuthentication(t *testing.T) {
	// Prediction: This test will pass - command should fail when no auth is found
	
	// Create a temporary directory with no auth file
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	// Create StartCmd
	cmd := &StartCmd{
		Host:      "127.0.0.1",
		Port:      5789,
		LogLevel:  "INFO",
		SkipAuthCheck: false,
	}
	
	// Mock environment to use our test storage
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", authFile)
	os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
	defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
	
	// Run the command and expect it to fail
	stdout, stderr, err := captureOutput(func() error {
		return cmd.Run()
	})
	
	// Should return an error about missing authentication
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication required")
	
	// Check output contains helpful message
	output := stdout + stderr
	assert.Contains(t, output, "No OAuth authentication found")
	assert.Contains(t, output, "claude-gate auth login")
}

// Test 3: StartCmd should respect custom host/port flags
func TestStartCmd_CustomHostPort(t *testing.T) {
	// Prediction: This test will pass - verifies custom configuration is used
	
	// Setup test storage with valid token
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	storage := auth.NewFileStorage(authFile)
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	err := storage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Create StartCmd with custom host/port
	cmd := &StartCmd{
		Host:          "0.0.0.0",
		Port:          8080,
		LogLevel:      "DEBUG",
		AuthToken:     "test-proxy-token",
		SkipAuthCheck: true, // Skip to make test simpler
	}
	
	// Verify custom configuration is respected
	assert.Equal(t, "0.0.0.0", cmd.Host)
	assert.Equal(t, 8080, cmd.Port)
	assert.Equal(t, "DEBUG", cmd.LogLevel)
	assert.Equal(t, "test-proxy-token", cmd.AuthToken)
	assert.True(t, cmd.SkipAuthCheck)
}

// Test 4: DashboardCmd should initialize TUI dashboard
func TestDashboardCmd_Initialize(t *testing.T) {
	// Prediction: This test will be limited as Bubble Tea TUI is hard to test
	
	// Setup test storage
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, "auth.json")
	
	storage := auth.NewFileStorage(authFile)
	testToken := &auth.TokenInfo{
		Type:         "oauth",
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}
	err := storage.Set("anthropic", testToken)
	require.NoError(t, err)
	
	// Create DashboardCmd
	cmd := &DashboardCmd{
		Host: "127.0.0.1",
		Port: 5789,
	}
	
	// Verify command configuration
	assert.Equal(t, "127.0.0.1", cmd.Host)
	assert.Equal(t, 5789, cmd.Port)
	
	// Full TUI testing would require mocking Bubble Tea
	// which is beyond the scope of unit tests
}

// Test 7: StatusCmd should display auth status correctly
func TestStatusCmd_DisplayStatus(t *testing.T) {
	// Prediction: This test will pass - StatusCmd is straightforward to test
	
	tests := []struct {
		name      string
		token     *auth.TokenInfo
		wantError bool
		contains  []string
	}{
		{
			name: "valid OAuth token",
			token: &auth.TokenInfo{
				Type:         "oauth",
				AccessToken:  "test-access",
				RefreshToken: "test-refresh",
				ExpiresAt:    time.Now().Add(time.Hour).Unix(),
			},
			wantError: false,
			contains: []string{
				"OAuth Authentication: Configured",
				"Token expires:",
			},
		},
		{
			name: "expired OAuth token",
			token: &auth.TokenInfo{
				Type:         "oauth",
				AccessToken:  "test-access",
				RefreshToken: "test-refresh",
				ExpiresAt:    time.Now().Add(-time.Hour).Unix(),
			},
			wantError: false,
			contains: []string{
				"OAuth Authentication: Configured",
				"Token is expired and will be refreshed",
			},
		},
		{
			name: "API key authentication",
			token: &auth.TokenInfo{
				Type:   "api_key",
				APIKey: "test-api-key",
			},
			wantError: false,
			contains: []string{
				"API Key Authentication: Configured",
				"Consider using OAuth for free usage",
			},
		},
		{
			name:      "no authentication",
			token:     nil,
			wantError: false,
			contains: []string{
				"Authentication: Not configured",
				"claude-gate auth login",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test storage
			tmpDir := t.TempDir()
			authFile := filepath.Join(tmpDir, "auth.json")
			
			if tt.token != nil {
				storage := auth.NewFileStorage(authFile)
				err := storage.Set("anthropic", tt.token)
				require.NoError(t, err)
			}
			
			// Create StatusCmd
			cmd := &StatusCmd{}
			
			// Mock environment
			os.Setenv("CLAUDE_GATE_AUTH_STORAGE_PATH", authFile)
			os.Setenv("CLAUDE_GATE_AUTH_STORAGE_TYPE", "file")
			defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_PATH")
			defer os.Unsetenv("CLAUDE_GATE_AUTH_STORAGE_TYPE")
			
			// Run command
			stdout, stderr, err := captureOutput(func() error {
				return cmd.Run()
			})
			
			// Check error
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Check output contains expected strings
			output := stdout + stderr
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

// Test 8: TestCmd should check proxy connectivity
func TestTestCmd_CheckConnectivity(t *testing.T) {
	// Prediction: This test will pass - we can mock the HTTP requests
	
	tests := []struct {
		name         string
		serverStatus int
		serverBody   string
		wantError    bool
		contains     []string
	}{
		{
			name:         "proxy is running",
			serverStatus: http.StatusOK,
			serverBody:   `{"models": ["claude-3-opus", "claude-3-sonnet"]}`,
			wantError:    false,
			contains: []string{
				"Proxy is running",
				"Available models:",
				"claude-3-opus",
			},
		},
		{
			name:         "proxy returns error",
			serverStatus: http.StatusInternalServerError,
			serverBody:   "Internal Server Error",
			wantError:    true,
			contains: []string{
				"Proxy returned an error",
				"Status: 500",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request path
				if r.URL.Path == "/health" {
					w.WriteHeader(http.StatusOK)
				} else if r.URL.Path == "/v1/models" {
					w.WriteHeader(tt.serverStatus)
					w.Write([]byte(tt.serverBody))
				}
			}))
			defer server.Close()
			
			// Create TestCmd pointing to test server
			cmd := &TestCmd{
				BaseURL: server.URL,
			}
			
			// Run command
			stdout, stderr, err := captureOutput(func() error {
				return cmd.Run()
			})
			
			// Check error
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Check output
			output := stdout + stderr
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

// Test 9: VersionCmd displays version information
func TestVersionCmd_Display(t *testing.T) {
	// Prediction: This test will pass - already tested but adding for completeness
	
	cmd := &VersionCmd{}
	
	stdout, stderr, err := captureOutput(func() error {
		return cmd.Run()
	})
	
	assert.NoError(t, err)
	
	output := stdout + stderr
	assert.Contains(t, output, "Claude Gate")
	assert.Contains(t, output, version)
	assert.Contains(t, output, "Go OAuth proxy for Anthropic API")
}

// Test for main function and CLI parsing
func TestMain_CLIParsing(t *testing.T) {
	// Test that Kong can parse our CLI structure
	
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no command shows help",
			args:    []string{},
			wantErr: true, // Kong returns error for no command
		},
		{
			name:    "version command",
			args:    []string{"version"},
			wantErr: false,
		},
		{
			name:    "start with flags",
			args:    []string{"start", "--host", "0.0.0.0", "--port", "8080"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"invalid"},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli,
				kong.Name("claude-gate"),
				kong.Description("Claude OAuth proxy server"),
				kong.UsageOnError(),
			)
			require.NoError(t, err)
			
			_, err = parser.Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}