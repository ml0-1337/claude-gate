package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/yourusername/claude-gate/internal/auth"
	"github.com/yourusername/claude-gate/internal/config"
	"github.com/yourusername/claude-gate/internal/proxy"
)

var version = "0.1.0"

type CLI struct {
	Start StartCmd `cmd:"" help:"Start the Claude OAuth proxy server"`
	Auth  AuthCmd  `cmd:"" help:"Authentication management commands"`
	Test  TestCmd  `cmd:"" help:"Test the proxy connection"`
	Version VersionCmd `cmd:"" help:"Show version information"`
}

type StartCmd struct {
	Host      string `help:"Host to bind the proxy server" default:"127.0.0.1"`
	Port      int    `help:"Port to bind the proxy server" default:"8000"`
	AuthToken string `help:"Enable proxy authentication with this token" env:"CLAUDE_GATE_PROXY_AUTH_TOKEN"`
	LogLevel  string `help:"Logging level (DEBUG, INFO, WARNING, ERROR)" default:"INFO"`
	SkipAuthCheck bool `help:"Skip OAuth authentication check"`
}

type AuthCmd struct {
	Login  LoginCmd  `cmd:"" help:"Authenticate with Claude Pro/Max using OAuth"`
	Logout LogoutCmd `cmd:"" help:"Clear stored authentication credentials"`
	Status StatusCmd `cmd:"" help:"Check authentication status"`
}

type LoginCmd struct{}
type LogoutCmd struct{}
type StatusCmd struct{}

type TestCmd struct {
	BaseURL string `help:"Proxy server URL" default:"http://localhost:8000"`
}

type VersionCmd struct{}

func (s *StartCmd) Run() error {
	cfg := config.DefaultConfig()
	cfg.Host = s.Host
	cfg.Port = s.Port
	cfg.ProxyAuthToken = s.AuthToken
	cfg.LogLevel = s.LogLevel
	cfg.LoadFromEnv()
	
	// Check authentication unless skipped
	if !s.SkipAuthCheck {
		storage := auth.NewTokenStorage(cfg.AuthStoragePath)
		token, err := storage.Get("anthropic")
		if err != nil || token == nil || token.Type != "oauth" {
			fmt.Println("❌ No OAuth authentication found!")
			fmt.Println("Please run 'claude-gate auth login' first to set up OAuth.")
			return fmt.Errorf("authentication required")
		}
		fmt.Println("✅ OAuth authentication configured and ready")
	}
	
	// Print startup banner
	fmt.Printf("\n🚀 Claude OAuth Proxy\n")
	fmt.Printf("Server URL:          http://%s\n", cfg.GetBindAddress())
	fmt.Printf("Anthropic API:       %s\n", cfg.AnthropicBaseURL)
	fmt.Printf("Proxy Auth:          %s\n", func() string {
		if cfg.ProxyAuthToken != "" {
			return "Enabled"
		}
		return "Disabled"
	}())
	fmt.Printf("OpenAI Compatible:   http://%s/v1\n\n", cfg.GetBindAddress())
	
	if cfg.ProxyAuthToken == "" {
		fmt.Println("⚠️  Proxy authentication disabled - anyone can use this proxy")
	}
	
	fmt.Println("\nPress CTRL+C to stop the server\n")
	
	// Create proxy server
	storage := auth.NewTokenStorage(cfg.AuthStoragePath)
	tokenProvider := auth.NewOAuthTokenProvider(storage)
	transformer := proxy.NewRequestTransformer()
	
	proxyConfig := &proxy.ProxyConfig{
		UpstreamURL:   cfg.AnthropicBaseURL,
		TokenProvider: tokenProvider,
		Transformer:   transformer,
		Timeout:       cfg.RequestTimeout,
	}
	
	server := proxy.NewProxyServer(proxyConfig, cfg.GetBindAddress(), storage)
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		fmt.Println("\n\n✋ Shutting down proxy server...")
		if err := server.Stop(30 * time.Second); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()
	
	// Start server
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	
	fmt.Println("✅ Proxy server stopped")
	return nil
}

func (l *LoginCmd) Run() error {
	cfg := config.DefaultConfig()
	storage := auth.NewTokenStorage(cfg.AuthStoragePath)
	client := auth.NewOAuthClient()
	
	// Check if already authenticated
	existing, _ := storage.Get("anthropic")
	if existing != nil && existing.Type == "oauth" {
		fmt.Println("Already authenticated!")
		fmt.Print("Do you want to re-authenticate? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return nil
		}
		storage.Remove("anthropic")
	}
	
	fmt.Println("\n🔐 Starting Claude Pro/Max OAuth authentication...")
	fmt.Println(strings.Repeat("-", 50))
	
	// Get authorization URL
	authData, err := client.GetAuthorizationURL()
	if err != nil {
		return fmt.Errorf("failed to generate authorization URL: %w", err)
	}
	
	fmt.Println("Please visit this URL to authorize:")
	fmt.Printf("\n%s\n\n", authData.URL)
	fmt.Println("After authorizing, you'll receive an authorization code.")
	fmt.Println(strings.Repeat("-", 50))
	
	// Get authorization code from user
	fmt.Print("Enter the authorization code: ")
	var code string
	fmt.Scanln(&code)
	code = strings.TrimSpace(code)
	
	// Exchange code for tokens
	fmt.Println("\nExchanging code for tokens...")
	token, err := client.ExchangeCode(code, authData.Verifier)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	
	// Save tokens
	if err := storage.Set("anthropic", token); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}
	
	fmt.Println("\n✅ Authentication successful!")
	fmt.Println("Your Claude Pro/Max account is now connected.")
	fmt.Println("Tokens are securely stored for future use.")
	
	return nil
}

func (l *LogoutCmd) Run() error {
	cfg := config.DefaultConfig()
	storage := auth.NewTokenStorage(cfg.AuthStoragePath)
	
	fmt.Print("Are you sure you want to logout? (y/N): ")
	var response string
	fmt.Scanln(&response)
	
	if response == "y" || response == "Y" {
		if err := storage.Remove("anthropic"); err != nil {
			return fmt.Errorf("failed to remove authentication: %w", err)
		}
		fmt.Println("✅ Logged out successfully")
	}
	
	return nil
}

func (s *StatusCmd) Run() error {
	cfg := config.DefaultConfig()
	storage := auth.NewTokenStorage(cfg.AuthStoragePath)
	
	fmt.Println("\nClaude Auth Bridge Status\n")
	
	// Check authentication
	token, err := storage.Get("anthropic")
	if err != nil || token == nil {
		fmt.Println("❌ Authentication: Not configured")
		fmt.Println("Run 'claude-gate auth login' to authenticate")
		return nil
	}
	
	if token.Type == "oauth" {
		fmt.Println("✅ OAuth Authentication: Configured")
		expires := time.Unix(token.ExpiresAt, 0)
		fmt.Printf("   Token expires: %s\n", expires.Format("2006-01-02 15:04:05"))
		if token.IsExpired() {
			fmt.Println("   ⚠️  Token is expired and will be refreshed on next use")
		} else if token.NeedsRefresh() {
			fmt.Println("   ⚠️  Token expires soon and will be refreshed on next use")
		}
	} else {
		fmt.Println("⚠️  API Key Authentication: Configured")
		fmt.Println("   Consider using OAuth for free usage")
	}
	
	// Show proxy configuration
	fmt.Println("\nProxy Configuration:")
	fmt.Printf("  Default host: %s\n", cfg.Host)
	fmt.Printf("  Default port: %d\n", cfg.Port)
	fmt.Printf("  Auth required: %s\n", func() string {
		if cfg.ProxyAuthToken != "" {
			return "Yes"
		}
		return "No"
	}())
	fmt.Printf("  Log level: %s\n", cfg.LogLevel)
	
	return nil
}

func (t *TestCmd) Run() error {
	fmt.Printf("Testing proxy at %s...\n\n", t.BaseURL)
	
	// Test root endpoint
	resp, err := http.Get(t.BaseURL + "/health")
	if err != nil {
		fmt.Printf("❌ Could not connect to proxy at %s\n", t.BaseURL)
		fmt.Println("Is the proxy server running?")
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ Proxy server is running")
		
		// Parse response
		var health map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&health); err == nil {
			if status, ok := health["oauth_status"].(string); ok {
				fmt.Printf("   OAuth status: %s\n", status)
			}
			if proxyAuth, ok := health["proxy_auth"].(string); ok {
				fmt.Printf("   Proxy auth: %s\n", proxyAuth)
			}
		}
	} else {
		fmt.Printf("❌ Unexpected status code: %d\n", resp.StatusCode)
	}
	
	return nil
}

func (v *VersionCmd) Run() error {
	fmt.Printf("claude-gate version %s\n", version)
	fmt.Println("Go OAuth proxy for Anthropic API")
	fmt.Println("https://github.com/yourusername/claude-gate")
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