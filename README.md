# Claude Gate

A high-performance Go OAuth proxy for Anthropic's Claude API that enables FREE Claude usage for Pro/Max subscribers.

## Overview

Claude Gate is a Go rewrite of claude-auth-bridge that maintains the critical OAuth bypass functionality while improving performance, security, and distribution. By identifying as "Claude Code" (Anthropic's official CLI), it allows Pro/Max subscribers to use the API without additional charges.

## Features

- ✅ **OAuth PKCE Authentication** - Secure authentication flow with Claude Pro/Max
- ✅ **System Prompt Injection** - Automatic Claude Code identification (the secret sauce!)
- ✅ **Model Alias Mapping** - Seamless handling of `latest` model aliases
- ✅ **SSE Streaming Support** - Full support for streaming responses
- ✅ **Cross-Platform** - Works on macOS and Linux
- ✅ **Interactive Dashboard** - Real-time monitoring of requests and usage
- ✅ **High Performance** - <50MB memory usage, <5ms request overhead

## Quick Start

### 1. Install

**Option A: Build from source** (Currently available)
```bash
git clone https://github.com/ml0-1337/claude-gate.git
cd claude-gate
make build
sudo mv claude-gate /usr/local/bin/
```

**Option B: NPM** (Coming soon - will be available after first release)
```bash
# npm install -g claude-gate
```

### 2. Authenticate

```bash
claude-gate auth login
```

### 3. Start Proxy

```bash
claude-gate start
```

### 4. Use with SDK

#### Using Anthropic SDK
```python
import anthropic

client = anthropic.Anthropic(
    base_url="http://localhost:5789",
    api_key="sk-dummy"  # Can be any string
)

response = client.messages.create(
    model="claude-opus-4-20250514",  # Latest Claude 4 Opus
    max_tokens=300,
    messages=[{"role": "user", "content": "Hello, Claude!"}]
)

print(response.content[0].text)
```

#### Using OpenAI SDK (compatibility mode)
```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-dummy",  # Can be any string
    base_url="http://localhost:5789/v1/"  # Note the /v1/ suffix
)

response = client.chat.completions.create(
    model="claude-opus-4-20250514",  # Latest Claude 4 Opus
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello, Claude!"}
    ]
)

print(response.choices[0].message.content)
```

**Note**: OpenAI SDK compatibility has some limitations. System messages are concatenated to the conversation start, and some OpenAI-specific parameters are ignored.

### 5. Using with Zed Editor

Configure Zed to use Claude Gate by adding this to your `settings.json`:

```json
{
  "language_models": {
    "anthropic": {
      "api_url": "http://127.0.0.1:5789"
    }
  }
}
```

You can find your Zed settings at:
- macOS: `~/.config/zed/settings.json`
- Linux: `~/.config/zed/settings.json`

This configuration redirects all Anthropic API calls from Zed to your local Claude Gate proxy, allowing you to use Claude in Zed for FREE with your Pro/Max subscription.

### 6. Using with Cursor IDE

Cursor requires a public HTTPS endpoint, so you'll need to create a tunnel to your local Claude Gate instance.

#### Step 1: Start Claude Gate
```bash
claude-gate start
```

#### Step 2: Create a tunnel (choose one option)

**Option A: Using Cloudflared**
```bash
cloudflared tunnel --url localhost:5789
```

**Option B: Using ngrok**
```bash
ngrok http 5789
```

Take note of the HTTPS URL provided (e.g., `https://xxxx.trycloudflare.com` or `https://xxxx.ngrok-free.app`)

#### Step 3: Configure Cursor

1. Open Cursor's settings and go to the "Models" section
2. Enter any API key in the "OpenAI API Key" field (e.g., `sk-dummy`)
3. Click the dropdown beneath the API key field labeled "Override OpenAI Base URL"
4. Enter your tunnel URL with `/v1` suffix (e.g., `https://xxxx.trycloudflare.com/v1`)
5. Click "Save" next to the URL field

#### Step 4: Configure your models

In Cursor, use the `anthropic/` prefix for Claude models:
- `anthropic/claude-opus-4-20250514` (recommended - latest Claude 4 Opus)
- `anthropic/claude-sonnet-4-20250514` (Claude 4 Sonnet)
- `anthropic/claude-3-5-sonnet-20241022` (Claude 3.5 Sonnet)
- `anthropic/claude-3-5-haiku-20241022` (Claude 3.5 Haiku)

⚠️ **Important**: When clicking "Verify" in Cursor, make sure to disable any models in Cursor's model list that aren't Claude models. Cursor randomly selects a model to test, and verification will fail if it tries a non-Claude model.

## Documentation

For detailed documentation, see the [docs](./docs) directory:

- **[Getting Started](./docs/getting-started/)** - Installation, configuration, and quick start
- **[User Guides](./docs/guides/)** - Troubleshooting, development, and contributing
- **[API Reference](./docs/reference/)** - CLI commands and HTTP API
- **[Architecture](./docs/architecture/)** - System design and security model

## Development

For development and testing:

```bash
# Prerequisites: Go 1.22+, Node.js 18+, GoReleaser
make build         # Build for current platform
make test          # Run tests
make npm-test      # Build all platforms and test NPM package
```

See our [Development Guide](./docs/guides/development.md) for detailed instructions.

## Contributing

We welcome contributions! Please see our [Contributing Guide](./docs/guides/contributing.md) for details.

## License

MIT License - see [LICENSE](./LICENSE) for details.

---

⚠️ **Disclaimer**: This project is not affiliated with Anthropic. Use at your own risk and in accordance with Claude's Terms of Service.
