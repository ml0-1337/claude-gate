package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config holds all configuration for the proxy server
type Config struct {
	// Server settings
	Host string
	Port int
	
	// Anthropic API settings
	AnthropicBaseURL string
	
	// Proxy authentication
	ProxyAuthToken string
	
	// Request settings
	RequestTimeout time.Duration
	MaxRequestSize int
	
	// Logging
	LogLevel     string
	LogRequests  bool
	
	// Rate limiting
	EnableRateLimit     bool
	RateLimitPerMinute  int
	
	// CORS settings
	CORSAllowOrigins []string
	
	// Storage paths
	AuthStoragePath string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	
	return &Config{
		Host:                "127.0.0.1",
		Port:                5789,
		AnthropicBaseURL:    "https://api.anthropic.com",
		RequestTimeout:      600 * time.Second,
		MaxRequestSize:      10 * 1024 * 1024, // 10MB
		LogLevel:            "INFO",
		LogRequests:         true,
		EnableRateLimit:     false,
		RateLimitPerMinute:  60,
		CORSAllowOrigins:    []string{"*"},
		AuthStoragePath:     filepath.Join(homeDir, ".claude-gate", "auth.json"),
	}
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Server settings
	if host := os.Getenv("CLAUDE_GATE_HOST"); host != "" {
		c.Host = host
	}
	if port := os.Getenv("CLAUDE_GATE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Port = p
		}
	}
	
	// Anthropic API
	if url := os.Getenv("CLAUDE_GATE_ANTHROPIC_BASE_URL"); url != "" {
		c.AnthropicBaseURL = url
	}
	
	// Proxy auth
	if token := os.Getenv("CLAUDE_GATE_PROXY_AUTH_TOKEN"); token != "" {
		c.ProxyAuthToken = token
	}
	
	// Request settings
	if timeout := os.Getenv("CLAUDE_GATE_REQUEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			c.RequestTimeout = d
		}
	}
	if size := os.Getenv("CLAUDE_GATE_MAX_REQUEST_SIZE"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			c.MaxRequestSize = s
		}
	}
	
	// Logging
	if level := os.Getenv("CLAUDE_GATE_LOG_LEVEL"); level != "" {
		c.LogLevel = level
	}
	if logReq := os.Getenv("CLAUDE_GATE_LOG_REQUESTS"); logReq != "" {
		c.LogRequests = logReq == "true" || logReq == "1"
	}
	
	// Rate limiting
	if enable := os.Getenv("CLAUDE_GATE_ENABLE_RATE_LIMIT"); enable != "" {
		c.EnableRateLimit = enable == "true" || enable == "1"
	}
	if limit := os.Getenv("CLAUDE_GATE_RATE_LIMIT_PER_MINUTE"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			c.RateLimitPerMinute = l
		}
	}
}

// GetBindAddress returns the server bind address
func (c *Config) GetBindAddress() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}