# Claude Gate

OAuth proxy for Anthropic's Claude API - enables FREE Claude usage for Pro/Max subscribers.

## Installation

```bash
npm install -g claude-gate
```

## Quick Start

```bash
# 1. Authenticate
claude-gate auth login

# 2. Start proxy
claude-gate start

# 3. Use with any Anthropic SDK
```

```python
import anthropic

client = anthropic.Anthropic(
    base_url="http://localhost:5789",
    api_key="sk-dummy"
)
```

## Supported Platforms

- macOS (Intel & Apple Silicon)
- Linux (x64 & ARM64)

## Documentation

Full documentation: https://github.com/ml0-1337/claude-gate

## License

MIT