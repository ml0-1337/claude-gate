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
- ✅ **Interactive Dashboard** - Real-time monitoring of requests and usage
- ✅ **High Performance** - <50MB memory usage, <5ms request overhead

## Quick Start

### 1. Install

```bash
npm install -g claude-gate
```

### 2. Authenticate

```bash
claude-gate auth login
```

### 3. Start Proxy

```bash
claude-gate start
```

### 4. Use with any SDK

```python
import anthropic

client = anthropic.Anthropic(
    base_url="http://localhost:5789",
    api_key="sk-dummy"  # Can be any string
)

response = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=300,
    messages=[{"role": "user", "content": "Hello, Claude!"}]
)
```

## Documentation

For detailed documentation, see the [docs](./docs) directory:

- **[Getting Started](./docs/getting-started/)** - Installation, configuration, and quick start
- **[User Guides](./docs/guides/)** - Troubleshooting, development, and contributing
- **[API Reference](./docs/reference/)** - CLI commands and HTTP API
- **[Architecture](./docs/architecture/)** - System design and security model

## Contributing

We welcome contributions! Please see our [Contributing Guide](./docs/guides/contributing.md) for details.

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Acknowledgments

- Original Python implementation: [r3ggi/claude-auth-bridge](https://github.com/r3ggi/claude-auth-bridge)
- Anthropic for Claude and the Claude Code CLI

---

⚠️ **Disclaimer**: This project is not affiliated with Anthropic. Use at your own risk and in accordance with Claude's Terms of Service.