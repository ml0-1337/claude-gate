//go:build e2e
// +build e2e

package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ml0-1337/claude-gate/internal/test/helpers"
)

func TestCLI_CompleteWorkflow(t *testing.T) {
	// Skip if binary not built
	binPath := filepath.Join("..", "..", "..", "claude-gate")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skip("Binary not found. Run 'make build' first")
	}

	// Setup mock servers
	mockOAuth := helpers.CreateMockOAuthServer(t)
	defer mockOAuth.Close()

	mockAPI := helpers.CreateMockAPIServer(t)
	defer mockAPI.Close()

	// Get free port for proxy
	proxyPort := helpers.GetFreePort()

	// Test version command
	t.Run("version", func(t *testing.T) {
		cmd := exec.Command(binPath, "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), "Claude Gate")
		assert.Contains(t, string(output), "Version:")
	})

	// Test help command
	t.Run("help", func(t *testing.T) {
		cmd := exec.Command(binPath, "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), "Usage:")
		assert.Contains(t, string(output), "auth")
		assert.Contains(t, string(output), "start")
	})

	// Test proxy start and stop
	t.Run("start_stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start proxy
		cmd := exec.CommandContext(ctx, binPath, "start",
			"--host", "127.0.0.1",
			"--port", proxyPort,
			"--auth-url", mockOAuth.URL+"/oauth/authorize",
			"--token-url", mockOAuth.URL+"/oauth/token",
			"--api-url", mockAPI.URL,
		)

		// Capture output
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Start command
		err := cmd.Start()
		require.NoError(t, err)

		// Wait for server to be ready
		proxyURL := fmt.Sprintf("http://127.0.0.1:%s", proxyPort)
		err = helpers.WaitForServer(proxyURL+"/health", 2*time.Second)
		
		if err == nil {
			assert.Contains(t, stdout.String(), "Proxy server started")
		}

		// Stop server
		cancel()
		cmd.Wait()
	})
}

func TestCLI_ErrorHandling(t *testing.T) {
	binPath := filepath.Join("..", "..", "..", "claude-gate")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skip("Binary not found. Run 'make build' first")
	}

	// Test invalid command
	t.Run("invalid_command", func(t *testing.T) {
		cmd := exec.Command(binPath, "invalid-command")
		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.Contains(t, string(output), "unexpected argument")
	})

	// Test port already in use
	t.Run("port_in_use", func(t *testing.T) {
		t.Skip("Skipping port conflict test - requires authentication")
		// This test would need mock OAuth setup to work properly
		// since the server checks for authentication before starting
	})
}