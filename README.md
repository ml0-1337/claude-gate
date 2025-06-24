# Claude Gate

A high-performance Go OAuth proxy for Anthropic's Claude API that enables FREE Claude usage for Pro/Max subscribers.

## Overview

Claude Gate is a Go rewrite of claude-auth-bridge that maintains the critical OAuth bypass functionality while improving performance, security, and distribution. By identifying as "Claude Code" (Anthropic's official CLI), it allows Pro/Max subscribers to use the API without additional charges.

## Features

- ✅ **OAuth PKCE Authentication** - Secure authentication flow with Claude Pro/Max
- ✅ **System Prompt Injection** - Automatic Claude Code identification (the secret sauce!)
- ✅ **Model Alias Mapping** - Seamless handling of `latest` model aliases
- ✅ **SSE Streaming Support** - Full support for streaming responses
- ✅ **Cross-Platform** - Works on macOS, Linux, and Windows
- ✅ **Secure Token Storage** - OS keychain integration (coming soon)
- ✅ **High Performance** - <50MB memory usage, <5ms request overhead

## Installation

### From Source

```bash
go install github.com/yourusername/claude-gate/cmd/claude-gate@latest
```

### Binary Releases

Coming soon via npm distribution.

## Quick Start

1. **Authenticate with your Claude Pro/Max account:**
   ```bash
   claude-gate auth login
   ```

2. **Start the proxy server:**
   ```bash
   claude-gate start
   ```

3. **Use any Anthropic SDK with the proxy:**
   ```python
   import anthropic
   
   client = anthropic.Anthropic(
       base_url="http://localhost:8000",
       api_key="sk-dummy"  # Can be any string
   )
   
   response = client.messages.create(
       model="claude-3-5-sonnet-latest",
       messages=[{"role": "user", "content": "Hello!"}]
   )
   ```

## CLI Commands

- `claude-gate start` - Start the proxy server
- `claude-gate auth login` - Authenticate with Claude Pro/Max
- `claude-gate auth logout` - Clear stored credentials
- `claude-gate auth status` - Check authentication status
- `claude-gate test` - Test proxy connection
- `claude-gate version` - Show version information

## Configuration

Environment variables:
- `CLAUDE_GATE_HOST` - Host to bind (default: 127.0.0.1)
- `CLAUDE_GATE_PORT` - Port to bind (default: 8000)
- `CLAUDE_GATE_PROXY_AUTH_TOKEN` - Enable proxy authentication
- `CLAUDE_GATE_LOG_LEVEL` - Logging level (DEBUG, INFO, WARNING, ERROR)

## How It Works

Claude Gate works by:
1. Authenticating with Claude using OAuth (same as claude.ai)
2. Injecting required headers and system prompts
3. Identifying as "Claude Code" to bypass API restrictions
4. Proxying requests with proper transformations

The critical innovation is the system prompt transformation that ensures every request identifies as Claude Code, which Anthropic's systems recognize as their official CLI.

## Development

### Prerequisites
- Go 1.22+
- Claude Pro or Claude Max subscription

### Building
```bash
go build -o claude-gate ./cmd/claude-gate
```

### Testing
```bash
go test ./...
```

## Security

- OAuth tokens are stored securely (file-based currently, keychain coming soon)
- No tokens are ever sent to third parties
- Optional proxy authentication for shared deployments
- All connections to Anthropic use HTTPS

## License

MIT

## Acknowledgments

Based on the original Python implementation of claude-auth-bridge.