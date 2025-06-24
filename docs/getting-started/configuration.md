# Configuration Guide

Claude Gate can be configured through command-line flags, environment variables, or configuration files.

## Configuration Methods

### 1. Command-Line Flags

Command-line flags take the highest precedence:

```bash
claude-gate start --port 3000 --host 0.0.0.0 --log-level debug
```

### 2. Environment Variables

Set environment variables before starting Claude Gate:

```bash
export CLAUDE_GATE_PORT=3000
export CLAUDE_GATE_HOST=0.0.0.0
export CLAUDE_GATE_LOG_LEVEL=debug
claude-gate start
```

### 3. Configuration File

Create a configuration file at `~/.claude-gate/config.yaml`:

```yaml
# ~/.claude-gate/config.yaml
host: 127.0.0.1
port: 5789
log_level: info
proxy_auth_token: your-secret-token
dashboard:
  enabled: false
  refresh_interval: 1s
```

## Configuration Options

### Server Settings

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--host` | `CLAUDE_GATE_HOST` | `127.0.0.1` | IP address to bind to |
| `--port` | `CLAUDE_GATE_PORT` | `5789` | Port number to listen on |
| `--proxy-auth-token` | `CLAUDE_GATE_PROXY_AUTH_TOKEN` | (none) | Enable proxy authentication |

### Logging

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--log-level` | `CLAUDE_GATE_LOG_LEVEL` | `INFO` | Log level: DEBUG, INFO, WARNING, ERROR |
| `--log-file` | `CLAUDE_GATE_LOG_FILE` | (stdout) | Log to file instead of console |
| `--log-format` | `CLAUDE_GATE_LOG_FORMAT` | `text` | Log format: text, json |

### Dashboard

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--dashboard` | `CLAUDE_GATE_DASHBOARD` | `false` | Enable interactive dashboard |
| `--dashboard-refresh` | `CLAUDE_GATE_DASHBOARD_REFRESH` | `1s` | Dashboard refresh interval |

### Security

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--allowed-origins` | `CLAUDE_GATE_ALLOWED_ORIGINS` | `*` | CORS allowed origins |
| `--tls-cert` | `CLAUDE_GATE_TLS_CERT` | (none) | TLS certificate file |
| `--tls-key` | `CLAUDE_GATE_TLS_KEY` | (none) | TLS key file |

## Common Configurations

### Development Setup

For local development with maximum visibility:

```bash
claude-gate start \
  --host 127.0.0.1 \
  --port 5789 \
  --log-level debug \
  --dashboard
```

### Production Setup

For production use with security:

```bash
export CLAUDE_GATE_HOST=127.0.0.1
export CLAUDE_GATE_PORT=5789
export CLAUDE_GATE_PROXY_AUTH_TOKEN=$(openssl rand -hex 32)
export CLAUDE_GATE_LOG_LEVEL=info
export CLAUDE_GATE_LOG_FILE=/var/log/claude-gate.log

claude-gate start
```

### Docker Setup

Using environment variables with Docker:

```dockerfile
FROM golang:1.22-alpine
# ... build steps ...

ENV CLAUDE_GATE_HOST=0.0.0.0
ENV CLAUDE_GATE_PORT=5789
ENV CLAUDE_GATE_LOG_LEVEL=info

EXPOSE 5789
CMD ["claude-gate", "start"]
```

### Systemd Service

Create `/etc/systemd/system/claude-gate.service`:

```ini
[Unit]
Description=Claude Gate Proxy
After=network.target

[Service]
Type=simple
User=claude-gate
Environment="CLAUDE_GATE_PORT=5789"
Environment="CLAUDE_GATE_HOST=127.0.0.1"
Environment="CLAUDE_GATE_LOG_LEVEL=info"
ExecStart=/usr/local/bin/claude-gate start
Restart=always

[Install]
WantedBy=multi-user.target
```

## Authentication Storage

Claude Gate stores authentication tokens in:
- **macOS**: `~/Library/Application Support/claude-gate/`
- **Linux**: `~/.config/claude-gate/`
- **Windows**: `%APPDATA%\claude-gate\`

The token file permissions are set to 600 (user read/write only) for security.

## Proxy Authentication

To require authentication for proxy access:

1. Generate a secure token:
   ```bash
   openssl rand -hex 32
   ```

2. Set the environment variable:
   ```bash
   export CLAUDE_GATE_PROXY_AUTH_TOKEN=your-generated-token
   ```

3. Include the token in your API requests:
   ```python
   client = anthropic.Anthropic(
       base_url="http://localhost:5789",
       api_key="your-generated-token"  # Use the proxy auth token
   )
   ```

## TLS/HTTPS Support

To enable HTTPS:

1. Generate or obtain TLS certificates
2. Start Claude Gate with TLS:
   ```bash
   claude-gate start \
     --tls-cert /path/to/cert.pem \
     --tls-key /path/to/key.pem
   ```

3. Update your client to use HTTPS:
   ```python
   client = anthropic.Anthropic(
       base_url="https://localhost:5789",
       api_key="sk-dummy"
   )
   ```

## Configuration Best Practices

1. **Security First**
   - Always use `127.0.0.1` unless you need external access
   - Enable proxy authentication in production
   - Use TLS for external access

2. **Logging**
   - Use `INFO` level for production
   - Enable `DEBUG` only when troubleshooting
   - Rotate log files regularly

3. **Performance**
   - Keep the dashboard disabled in production
   - Use appropriate log levels to reduce I/O

4. **Monitoring**
   - Set up log aggregation for production use
   - Monitor disk space for log files
   - Track proxy usage patterns

## Troubleshooting Configuration

### Port Already in Use

If you get "address already in use" error:

1. Check what's using the port:
   ```bash
   lsof -i :5789  # macOS/Linux
   netstat -ano | findstr :5789  # Windows
   ```

2. Use a different port:
   ```bash
   claude-gate start --port 5790
   ```

### Configuration Not Loading

1. Check file location:
   ```bash
   claude-gate config --show-path
   ```

2. Validate YAML syntax:
   ```bash
   claude-gate config --validate
   ```

3. Check environment variables:
   ```bash
   env | grep CLAUDE_GATE
   ```

---

[← Quick Start](./quick-start.md) | [Documentation Home](../README.md) | [Troubleshooting →](../guides/troubleshooting.md)