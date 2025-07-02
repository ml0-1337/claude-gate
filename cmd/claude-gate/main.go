package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ml0-1337/claude-gate/internal/auth"
	"github.com/ml0-1337/claude-gate/internal/config"
	"github.com/ml0-1337/claude-gate/internal/logger"
	"github.com/ml0-1337/claude-gate/internal/proxy"
	"github.com/ml0-1337/claude-gate/internal/ui"
	"github.com/ml0-1337/claude-gate/internal/ui/components"
)

var version = "0.1.0"

// createStorageFactoryConfig creates a StorageFactoryConfig from the main Config
func createStorageFactoryConfig(cfg *config.Config) auth.StorageFactoryConfig {
	return auth.StorageFactoryConfig{
		Type:                           auth.StorageType(cfg.AuthStorageType),
		FilePath:                       cfg.AuthStoragePath,
		ServiceName:                    cfg.KeyringService,
		KeychainTrustApp:               cfg.KeychainTrustApp,
		KeychainAccessibleWhenUnlocked: cfg.KeychainAccessibleWhenUnlocked,
		KeychainSynchronizable:         cfg.KeychainSynchronizable,
	}
}

type CLI struct {
	Start     StartCmd     `cmd:"" help:"Start the Claude OAuth proxy server"`
	Dashboard DashboardCmd `cmd:"" help:"Start server with interactive dashboard"`
	Auth      AuthCmd      `cmd:"" help:"Authentication management commands"`
	Test      TestCmd      `cmd:"" help:"Test the proxy connection"`
	Version   VersionCmd   `cmd:"" help:"Show version information"`
}

type StartCmd struct {
	Host          string `help:"Host to bind the proxy server" default:"127.0.0.1"`
	Port          int    `help:"Port to bind the proxy server" default:"5789"`
	AuthToken     string `help:"Enable proxy authentication with this token" env:"CLAUDE_GATE_PROXY_AUTH_TOKEN"`
	LogLevel      string `help:"Logging level (DEBUG, INFO, WARNING, ERROR)" default:"INFO"`
	StorageBackend string `help:"Storage backend (auto, keyring, file, claude-code)" default:"auto" enum:"auto,keyring,file,claude-code"`
	SkipAuthCheck bool   `help:"Skip OAuth authentication check"`
}

type DashboardCmd struct {
	Host          string `help:"Host to bind the proxy server" default:"127.0.0.1"`
	Port          int    `help:"Port to bind the proxy server" default:"5789"`
	AuthToken     string `help:"Enable proxy authentication with this token" env:"CLAUDE_GATE_PROXY_AUTH_TOKEN"`
	LogLevel      string `help:"Logging level (DEBUG, INFO, WARNING, ERROR)" default:"INFO"`
	StorageBackend string `help:"Storage backend (auto, keyring, file, claude-code)" default:"auto" enum:"auto,keyring,file,claude-code"`
	SkipAuthCheck bool   `help:"Skip OAuth authentication check"`
}

type AuthCmd struct {
	Login   LoginCmd         `cmd:"" help:"Authenticate with Claude Pro/Max using OAuth"`
	Logout  LogoutCmd        `cmd:"" help:"Clear stored authentication credentials"`
	Status  StatusCmd        `cmd:"" help:"Check authentication status"`
	Storage AuthStorageCmd   `cmd:"" help:"Manage token storage backends"`
}

type LoginCmd struct{}
type LogoutCmd struct{}
type StatusCmd struct{}

type TestCmd struct {
	BaseURL string `help:"Proxy server URL" default:"http://localhost:5789"`
}

type VersionCmd struct{}

func (s *StartCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.Host = s.Host
	cfg.Port = s.Port
	cfg.ProxyAuthToken = s.AuthToken
	cfg.LogLevel = s.LogLevel
	cfg.AuthStorageType = s.StorageBackend
	cfg.LoadFromEnv()
	
	out := ui.NewOutput()
	
	// Check authentication unless skipped
	if !s.SkipAuthCheck {
		// Create storage using factory
		factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
		
		storage, err := factory.Create()
		if err != nil {
			return fmt.Errorf("failed to create storage: %w", err)
		}
		
		token, err := storage.Get("anthropic")
		if err != nil || token == nil || token.Type != "oauth" {
			out.Error("No OAuth authentication found!")
			out.Info("Please run 'claude-gate auth login' first to set up OAuth.")
			return fmt.Errorf("authentication required")
		}
		out.Success("OAuth authentication configured and ready")
	}
	
	// Print startup banner
	out.Title("🚀 Claude OAuth Proxy")
	
	headers := []string{"Configuration", "Value"}
	rows := [][]string{
		{"Server URL", fmt.Sprintf("http://%s", cfg.GetBindAddress())},
		{"Anthropic API", cfg.AnthropicBaseURL},
		{"Proxy Auth", func() string {
			if cfg.ProxyAuthToken != "" {
				return "Enabled"
			}
			return "Disabled"
		}()},
		{"OpenAI Compatible", fmt.Sprintf("http://%s/v1", cfg.GetBindAddress())},
	}
	out.Table(headers, rows)
	
	if cfg.ProxyAuthToken == "" {
		out.Warning("Proxy authentication disabled - anyone can use this proxy")
	}
	
	out.Info("\nPress CTRL+C to stop the server")
	
	// Create proxy server with storage
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.CreateWithMigration()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	tokenProvider := auth.NewOAuthTokenProvider(storage)
	transformer := proxy.NewRequestTransformer()
	
	// Create logger
	log := logger.New(logger.ParseLevel(cfg.LogLevel))
	
	proxyConfig := &proxy.ProxyConfig{
		UpstreamURL:   cfg.AnthropicBaseURL,
		TokenProvider: tokenProvider,
		Transformer:   transformer,
		Timeout:       cfg.RequestTimeout,
		Logger:        log,
	}
	
	server := proxy.NewProxyServer(proxyConfig, cfg.GetBindAddress(), storage)
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		out.Info("\n\nShutting down proxy server...")
		if err := server.Stop(30 * time.Second); err != nil {
			out.Error("Error during shutdown: %v", err)
		}
	}()
	
	// Start server
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	
	out.Success("Proxy server stopped")
	return nil
}

func (d *DashboardCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.Host = d.Host
	cfg.Port = d.Port
	cfg.ProxyAuthToken = d.AuthToken
	cfg.LogLevel = d.LogLevel
	cfg.AuthStorageType = d.StorageBackend
	cfg.LoadFromEnv()
	
	out := ui.NewOutput()
	
	// Check authentication unless skipped
	if !d.SkipAuthCheck {
		// Create storage using factory
		factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
		
		storage, err := factory.Create()
		if err != nil {
			return fmt.Errorf("failed to create storage: %w", err)
		}
		
		token, err := storage.Get("anthropic")
		if err != nil || token == nil || token.Type != "oauth" {
			out.Error("No OAuth authentication found!")
			out.Info("Please run 'claude-gate auth login' first to set up OAuth.")
			return fmt.Errorf("authentication required")
		}
	}
	
	// Create proxy server with storage
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.CreateWithMigration()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	tokenProvider := auth.NewOAuthTokenProvider(storage)
	transformer := proxy.NewRequestTransformer()
	
	// Create logger
	log := logger.New(logger.ParseLevel(cfg.LogLevel))
	
	proxyConfig := &proxy.ProxyConfig{
		UpstreamURL:   cfg.AnthropicBaseURL,
		TokenProvider: tokenProvider,
		Transformer:   transformer,
		Timeout:       cfg.RequestTimeout,
		Logger:        log,
	}
	
	server := proxy.NewEnhancedProxyServer(proxyConfig, cfg.GetBindAddress(), storage)
	
	// Get dashboard model
	dashboardModel := server.GetDashboard()
	
	// Start server in background
	serverErrChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()
	
	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Check if server started successfully
	select {
	case err := <-serverErrChan:
		return fmt.Errorf("failed to start server: %w", err)
	default:
		// Server started successfully
	}
	
	// Run the dashboard UI
	p := tea.NewProgram(dashboardModel, tea.WithAltScreen())
	
	// Handle shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		p.Quit()
	}()
	
	// Run the dashboard
	if _, err := p.Run(); err != nil {
		out.Error("Dashboard error: %v", err)
	}
	
	// Shutdown server
	out.Info("\nShutting down proxy server...")
	if err := server.Stop(30 * time.Second); err != nil {
		out.Error("Error during shutdown: %v", err)
	}
	
	out.Success("Dashboard stopped")
	return nil
}

func (l *LoginCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create storage using factory
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.Create()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	client := auth.NewOAuthClient()
	out := ui.NewOutput()
	
	// Check if already authenticated
	existing, _ := storage.Get("anthropic")
	if existing != nil && existing.Type == "oauth" {
		out.Warning("Already authenticated!")
		if !components.Confirm("Do you want to re-authenticate?") {
			return nil
		}
		err := components.RunSpinner("Removing existing authentication...", func() error {
			return storage.Remove("anthropic")
		})
		if err != nil {
			return err
		}
	}
	
	// Get authorization URL
	var authData *auth.AuthData
	var authErr error
	
	// Generate URL in background
	go func() {
		authData, authErr = client.GetAuthorizationURL()
	}()
	
	// Show spinner while generating
	err = components.RunSpinner("Generating authorization URL...", func() error {
		// Wait for URL generation
		for authData == nil && authErr == nil {
			time.Sleep(100 * time.Millisecond)
		}
		return authErr
	})
	if err != nil {
		return fmt.Errorf("failed to generate authorization URL: %w", err)
	}
	
	// Run interactive OAuth flow
	code, err := ui.RunOAuthFlow(authData.URL)
	if err != nil {
		return fmt.Errorf("authentication canceled: %w", err)
	}
	code = strings.TrimSpace(code)
	
	// Exchange code for tokens
	var token *auth.TokenInfo
	err = components.RunSpinner("Exchanging code for tokens...", func() error {
		var err error
		token, err = client.ExchangeCode(code, authData.Verifier)
		if err != nil {
			return err
		}
		// Save tokens
		return storage.Set("anthropic", token)
	})
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	
	out.Success("\nAuthentication successful!")
	out.Success("Your Claude Pro/Max account is now connected.")
	out.Info("Tokens are securely stored for future use.")
	
	return nil
}

func (l *LogoutCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create storage using factory
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.Create()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	out := ui.NewOutput()
	
	if !components.Confirm("Are you sure you want to logout?") {
		return nil
	}
	
	err = components.RunSpinner("Removing authentication...", func() error {
		return storage.Remove("anthropic")
	})
	if err != nil {
		return fmt.Errorf("failed to remove authentication: %w", err)
	}
	
	out.Success("Logged out successfully")
	return nil
}

func (s *StatusCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	// Create storage using factory
	factory := auth.NewStorageFactory(createStorageFactoryConfig(cfg))
	
	storage, err := factory.Create()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	
	out := ui.NewOutput()
	
	out.Title("Claude Gate Status")
	
	// Check authentication
	token, err := storage.Get("anthropic")
	if err != nil || token == nil {
		out.Error("Authentication: Not configured")
		out.Info("Run 'claude-gate auth login' to authenticate")
		return nil
	}
	
	if token.Type == "oauth" {
		out.Success("OAuth Authentication: Configured")
		expires := time.Unix(token.ExpiresAt, 0)
		if token.IsExpired() {
			out.Warning("Token is expired and will be refreshed on next use")
		} else if token.NeedsRefresh() {
			out.Warning("Token expires soon and will be refreshed on next use")
		} else {
			out.Info("Token expires: %s", expires.Format("2006-01-02 15:04:05"))
		}
	} else {
		out.Warning("API Key Authentication: Configured")
		out.Info("Consider using OAuth for free usage")
	}
	
	// Show proxy configuration
	out.Subtitle("\nProxy Configuration")
	headers := []string{"Setting", "Value"}
	rows := [][]string{
		{"Default host", cfg.Host},
		{"Default port", fmt.Sprintf("%d", cfg.Port)},
		{"Auth required", func() string {
			if cfg.ProxyAuthToken != "" {
				return "Yes"
			}
			return "No"
		}()},
		{"Log level", cfg.LogLevel},
	}
	out.Table(headers, rows)
	
	return nil
}

func (t *TestCmd) Run() error {
	out := ui.NewOutput()
	out.Title("Testing Claude Gate Proxy")
	out.Info("Testing proxy at %s...", t.BaseURL)
	
	var resp *http.Response
	err := components.RunSpinner("Connecting to proxy...", func() error {
		var err error
		resp, err = http.Get(t.BaseURL + "/health")
		return err
	})
	if err != nil {
		out.Error("Could not connect to proxy at %s", t.BaseURL)
		out.Info("Is the proxy server running?")
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		out.Success("Proxy server is running")
		
		// Parse response
		var health map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&health); err == nil {
			headers := []string{"Status", "Value"}
			rows := [][]string{}
			
			if status, ok := health["oauth_status"].(string); ok {
				rows = append(rows, []string{"OAuth status", status})
			}
			if proxyAuth, ok := health["proxy_auth"].(string); ok {
				rows = append(rows, []string{"Proxy auth", proxyAuth})
			}
			
			if len(rows) > 0 {
				out.Table(headers, rows)
			}
		}
	} else {
		out.Error("Unexpected status code: %d", resp.StatusCode)
	}
	
	return nil
}

func (v *VersionCmd) Run() error {
	out := ui.NewOutput()
	out.Title("Claude Gate")
	out.Info("Version: %s", version)
	out.Info("Go OAuth proxy for Anthropic API")
	out.Info("https://github.com/ml0-1337/claude-gate")
	return nil
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("claude-gate"),
		kong.Description("Claude OAuth proxy server - FREE Claude usage for Pro/Max subscribers"),
		kong.UsageOnError(),
	)
	
	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}