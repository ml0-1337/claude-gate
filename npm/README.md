# Claude Gate

High-performance OAuth proxy for Anthropic's Claude API - FREE usage for Pro/Max subscribers.

## Installation

```bash
npm install -g claude-gate
```

## Quick Start

1. **Authenticate with your Claude Pro/Max account:**
   ```bash
   claude-gate auth login
   ```

2. **Start the proxy server:**
   ```bash
   claude-gate start
   ```

3. **Use with any Anthropic SDK:**
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

## Features

- ✅ **OAuth Authentication** - Login with your Claude Pro/Max account
- ✅ **System Prompt Injection** - Automatic Claude Code identification 
- ✅ **Model Alias Support** - Use `latest` model versions
- ✅ **SSE Streaming** - Full support for streaming responses
- ✅ **Cross-Platform** - Works on macOS, Linux, and Windows
- ✅ **High Performance** - Minimal overhead and memory usage

## Commands

- `claude-gate start` - Start the proxy server
- `claude-gate auth login` - Authenticate with Claude
- `claude-gate auth logout` - Clear credentials
- `claude-gate auth status` - Check authentication
- `claude-gate test` - Test proxy connection
- `claude-gate version` - Show version info
- `claude-gate --help` - Show all commands

## Supported Platforms

- macOS (Intel & Apple Silicon)
- Linux (x64 & ARM64)
- Windows (x64)

## Documentation

Full documentation available at: https://github.com/yourusername/claude-gate

## License

MIT