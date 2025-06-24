# Quick Start Guide

Get up and running with Claude Gate in just a few minutes!

## Prerequisites

- Claude Gate installed ([Installation Guide](./installation.md))
- Claude Pro or Claude Max subscription
- Python, Node.js, or any language with an Anthropic SDK

## Step 1: Authenticate

First, authenticate with your Claude account:

```bash
claude-gate auth login
```

This will:
1. Open your browser to the Claude login page
2. Guide you through the OAuth authentication flow
3. Securely store your authentication token

## Step 2: Start the Proxy Server

Start Claude Gate in the background:

```bash
claude-gate start
```

The proxy server will start on `http://localhost:8080` by default.

To run with a custom port:

```bash
claude-gate start --port 3000
```

## Step 3: Use with Anthropic SDKs

Now you can use any Anthropic SDK with Claude Gate as the proxy:

### Python Example

```python
import anthropic

client = anthropic.Anthropic(
    base_url="http://localhost:8080",
    api_key="sk-dummy"  # Can be any string
)

response = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=300,
    messages=[
        {"role": "user", "content": "Hello, Claude!"}
    ]
)

print(response.content[0].text)
```

### Node.js Example

```javascript
import Anthropic from '@anthropic-ai/sdk';

const client = new Anthropic({
  baseURL: 'http://localhost:8080',
  apiKey: 'sk-dummy', // Can be any string
});

const response = await client.messages.create({
  model: 'claude-3-5-sonnet-20241022',
  max_tokens: 300,
  messages: [
    { role: 'user', content: 'Hello, Claude!' }
  ],
});

console.log(response.content[0].text);
```

### cURL Example

```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-dummy" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-5-sonnet-20241022",
    "max_tokens": 300,
    "messages": [
      {"role": "user", "content": "Hello, Claude!"}
    ]
  }'
```

## Step 4: Monitor Usage (Optional)

View the interactive dashboard to monitor your usage:

```bash
claude-gate start --dashboard
```

This shows:
- Real-time request metrics
- Token usage statistics
- Response times
- Error rates

## Common Commands

### Check Status

```bash
claude-gate status
```

### Stop the Server

```bash
claude-gate stop
```

### View Logs

```bash
claude-gate logs
```

### Re-authenticate

```bash
claude-gate auth refresh
```

## What's Next?

- **Configure Claude Gate**: See the [Configuration Guide](./configuration.md)
- **Integrate with your app**: Check the [API Reference](../reference/api.md)
- **Run into issues?**: Visit our [Troubleshooting Guide](../guides/troubleshooting.md)
- **Want to contribute?**: Read the [Contributing Guide](../guides/contributing.md)

## Tips

1. **Keep it running**: Claude Gate needs to be running for your applications to work
2. **Security**: By default, Claude Gate only accepts connections from localhost
3. **Rate limits**: Claude Gate respects your Claude subscription's rate limits
4. **Token usage**: Monitor your token usage through the dashboard to avoid surprises

---

[← Installation](./installation.md) | [Documentation Home](../README.md) | [Configuration →](./configuration.md)