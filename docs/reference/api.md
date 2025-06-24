# Claude Gate API Documentation

## Overview

Claude Gate acts as a transparent proxy to the Anthropic Claude API, adding OAuth authentication and system prompt injection. It maintains full API compatibility while identifying as "Claude Code".

## Base URL

```
http://localhost:8080
```

The proxy listens on port 8080 by default. This can be configured via environment variables.

## Authentication

Claude Gate handles OAuth authentication transparently. Before making API requests:

1. Authenticate using the CLI: `claude-gate auth login`
2. Start the proxy server: `claude-gate start`
3. Configure your client to use `http://localhost:8080` instead of `https://api.anthropic.com`

## Proxied Endpoints

All Anthropic Claude API endpoints are proxied transparently:

### Messages API
```
POST /v1/messages
```

Creates a message with the Claude model. The proxy automatically:
- Adds OAuth authentication headers
- Injects "Claude Code" system prompt
- Maps model aliases (e.g., "latest" to specific versions)

### Models API
```
GET /v1/models
```

Lists available models (proxied directly).

### Other Endpoints

All other Anthropic API endpoints are proxied without modification, with only authentication headers added.

## Request Transformation

### System Prompt Injection

The proxy ensures all requests identify as "Claude Code":

**Original Request:**
```json
{
  "model": "claude-3-5-sonnet-20241022",
  "system": "You are a helpful assistant",
  "messages": [{"role": "user", "content": "Hello"}]
}
```

**Transformed Request:**
```json
{
  "model": "claude-3-5-sonnet-20241022",
  "system": [
    {"type": "text", "text": "You are Claude Code, Anthropic's official CLI for Claude."},
    {"type": "text", "text": "You are a helpful assistant"}
  ],
  "messages": [{"role": "user", "content": "Hello"}]
}
```

### Model Alias Mapping

The proxy automatically maps model aliases:

| Alias | Actual Model |
|-------|--------------|
| claude-3-5-haiku-latest | claude-3-5-haiku-20241022 |
| claude-3-5-sonnet-latest | claude-3-5-sonnet-20241022 |
| claude-3-7-sonnet-latest | claude-3-7-sonnet-20250219 |
| claude-3-opus-latest | claude-3-opus-20240229 |

### Header Transformation

**Added Headers:**
- `Authorization: Bearer <oauth-token>`
- `anthropic-beta: oauth-2025-04-20`
- `anthropic-version: 2023-06-01`

**Removed Headers:**
- `User-Agent` (identifies client application)
- Custom headers not in allowlist

## Response Handling

### Streaming Responses

Server-Sent Events (SSE) are fully supported with immediate flushing for real-time streaming.

### Error Responses

Errors maintain Anthropic's format:
```json
{
  "type": "error",
  "error": {
    "type": "authentication_error",
    "message": "Authentication required. Please run: claude-gate auth login"
  }
}
```

## Client Configuration Examples

### Python (anthropic)
```python
from anthropic import Anthropic

client = Anthropic(
    base_url="http://localhost:8080/v1",
    api_key="dummy"  # OAuth handled by proxy
)
```

### Node.js (@anthropic-ai/sdk)
```javascript
import Anthropic from '@anthropic-ai/sdk';

const client = new Anthropic({
  baseURL: 'http://localhost:8080/v1',
  apiKey: 'dummy'  // OAuth handled by proxy
});
```

### cURL
```bash
curl http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

## Rate Limiting

Currently, rate limiting is enforced by Anthropic's API. Local rate limiting is planned for future releases.

## Health Check

```
GET /health
```

Returns 200 OK when the proxy is running and authenticated.

## Metrics (Planned)

```
GET /metrics
```

Prometheus-compatible metrics endpoint (future release).

## Security Considerations

1. The proxy runs locally and should not be exposed to the internet
2. OAuth tokens are never exposed to clients
3. All communication with Anthropic uses TLS
4. See [Security Policy](../SECURITY.md) for details

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Run `claude-gate auth login`
   - Check token expiration with `claude-gate auth status`

2. **Connection Refused**
   - Ensure proxy is running: `claude-gate start`
   - Check port 8080 is not in use

3. **Model Not Found**
   - Use valid model names or aliases
   - Check Anthropic API documentation

### Debug Mode

Set environment variable for verbose logging:
```bash
CLAUDE_GATE_DEBUG=true claude-gate start
```

## Version Compatibility

Claude Gate maintains compatibility with:
- Anthropic API version: 2023-06-01
- OAuth Beta: oauth-2025-04-20

Check your version:
```bash
claude-gate version
```

---

[← Reference](../README.md#reference) | [Documentation Home](../README.md) | [CLI Reference →](./cli.md)