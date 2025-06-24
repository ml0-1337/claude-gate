# Claude Gate Architecture

## Overview

Claude Gate is a high-performance OAuth proxy for Anthropic's Claude API that enables free Claude usage for Pro/Max subscribers by identifying as "Claude Code" (Anthropic's official CLI). This document describes the system architecture, design decisions, and component interactions.

## System Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│                 │     │                  │     │                 │
│  User's Editor  │────▶│   Claude Gate   │────▶│  Anthropic API  │
│   (Zed, etc)    │     │   OAuth Proxy    │     │                 │
│                 │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └─────────────────┘
         │                       │                         │
         │                       │                         │
         ▼                       ▼                         ▼
    API Requests          System Prompt              OAuth 2.0
                          Injection +                Authentication
                          Token Addition

```

## Core Components

### 1. CLI Layer (`cmd/claude-gate/`)

The command-line interface provides user interaction and server management:

- **Framework**: Kong CLI framework for command parsing
- **Commands**:
  - `start` - Launches the proxy server
  - `auth login` - Initiates OAuth flow
  - `auth logout` - Revokes tokens
  - `auth status` - Shows authentication status
  - `test` - Validates proxy functionality
  - `version` - Displays version information

### 2. Authentication Layer (`internal/auth/`)

Handles OAuth 2.0 PKCE flow with Anthropic:

- **OAuth Client** (`client.go`) - HTTP client for OAuth operations
- **OAuth Flow** (`oauth.go`) - PKCE implementation with automatic token refresh
- **Token Storage** (`storage.go`) - Secure token persistence (file-based, keychain planned)

Key features:
- Public client with PKCE for security
- Automatic token refresh before expiration
- Graceful degradation for storage backends

### 3. Proxy Layer (`internal/proxy/`)

Core proxy functionality for API translation:

- **Server** (`server.go`) - HTTP server management and lifecycle
- **Handler** (`handler.go`) - Request routing and response handling
- **Transformer** (`transformer.go`) - Request/response transformation

Key transformations:
1. System prompt injection (prepends "Claude Code" identifier)
2. OAuth header injection (adds authentication)
3. Model alias mapping (e.g., "latest" → specific versions)
4. SSE stream handling with proper flushing

### 4. Configuration (`internal/config/`)

Centralized configuration management:
- Default port: 5789
- Configurable via environment variables
- Runtime configuration validation

## Data Flow

### 1. Authentication Flow

```
User ──login──▶ OAuth Provider ──code──▶ Claude Gate ──exchange──▶ Access Token
                                                           │
                                                           ▼
                                                     Token Storage
```

### 2. Request Flow

```
1. Client Request → Claude Gate
2. Validate authentication
3. Transform request:
   - Inject "Claude Code" system prompt
   - Map model aliases
   - Add OAuth headers
4. Forward to Anthropic API
5. Stream response back to client
```

### 3. System Prompt Transformation

The proxy ensures all requests identify as "Claude Code":

- String prompts: Converted to array with Claude Code first
- Array prompts: Claude Code prepended if not present
- Empty prompts: Left unchanged

## Security Architecture

### Authentication Security

- **OAuth 2.0 PKCE**: Proof Key for Code Exchange prevents authorization code interception
- **Public Client**: Client ID is public by design (security via PKCE, not secrecy)
- **Token Rotation**: Automatic refresh before expiration
- **Secure Storage**: Tokens stored in OS keychain (planned) or encrypted file

### Network Security

- **TLS Only**: All communication with Anthropic over HTTPS
- **Header Stripping**: Removes identifying headers from client
- **No Logging**: Sensitive data never logged
- **Rate Limiting**: Planned to prevent abuse

### Defense in Depth

1. Authentication layer validates all requests
2. Proxy never exposes raw tokens to clients
3. Automatic token refresh prevents expiration attacks
4. File-based storage fallback with encryption (planned)

## Deployment Architecture

### Binary Distribution

```
GitHub Release ──▶ Platform Binaries ──▶ Direct Download
       │                                        │
       └──▶ NPM Packages ──▶ Auto-install ──────┘
```

### NPM Package Structure

```
claude-gate (main package)
    ├── @claude-gate/darwin-x64
    ├── @claude-gate/darwin-arm64
    ├── @claude-gate/linux-x64
    ├── @claude-gate/linux-arm64
    └── @claude-gate/win32-x64
```

### Runtime Requirements

- **Go Binary**: Single static binary, no runtime dependencies
- **Port**: 9988 (configurable)
- **Storage**: ~/.claude-gate/ for tokens and config
- **Permissions**: User-level only, no elevated privileges

## Performance Characteristics

### Optimizations

1. **Zero-copy streaming**: Direct pipe between client and API
2. **Minimal allocations**: Reuses buffers where possible
3. **No caching**: Stateless proxy for simplicity
4. **Connection pooling**: Reuses HTTP connections

### Bottlenecks

1. **Token refresh**: Synchronous operation (minimal impact)
2. **System prompt parsing**: JSON unmarshal/marshal overhead
3. **SSE buffering**: Requires immediate flushing

## Monitoring and Observability

### Current State

- Basic console logging
- Error responses in Anthropic format
- No metrics collection

### Planned Enhancements

1. Structured logging with levels
2. Prometheus metrics endpoint
3. Request ID tracking
4. Performance profiling endpoints

## Future Architecture Considerations

### Planned Features

1. **OS Keychain Integration**: Secure token storage
2. **Rate Limiting**: Token bucket algorithm
3. **Circuit Breaker**: Fault tolerance for API outages
4. **WebSocket Support**: Future API compatibility

### Scalability Path

Currently designed for single-user desktop use. For multi-user:

1. External token storage (Redis/PostgreSQL)
2. Horizontal scaling with load balancer
3. Distributed rate limiting
4. Centralized logging

## Architecture Decisions Records (ADRs)

See `/docs/decisions/` for detailed rationale on:

1. Why PKCE over client credentials
2. Why file-based storage with keychain upgrade path
3. Why system prompt injection over API translation
4. Why Go over Python (performance, distribution)
5. Why NPM distribution alongside binaries

---

[← Architecture](../README.md#architecture) | [Documentation Home](../README.md) | [Security →](./security.md)