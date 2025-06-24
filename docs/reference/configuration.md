# Configuration Reference

Complete reference for all Claude Gate configuration options.

## Configuration Hierarchy

Claude Gate uses a hierarchical configuration system where settings can be specified at different levels, with higher-priority sources overriding lower ones:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration file**
4. **Default values** (lowest priority)

## Configuration File

Claude Gate looks for configuration in the following locations:

- `~/.claude-gate/config.yaml` (default)
- Path specified by `--config` flag
- Path specified by `CLAUDE_GATE_CONFIG` environment variable

### Configuration File Format

```yaml
# ~/.claude-gate/config.yaml
host: 127.0.0.1
port: 5789
log_level: info
log_format: text
proxy_auth_token: your-secret-token
allowed_origins:
  - http://localhost:3000
  - https://myapp.com
dashboard:
  enabled: false
  refresh_interval: 1s
tls:
  cert: /path/to/cert.pem
  key: /path/to/key.pem
```

## All Configuration Options

### Server Configuration

| Option | CLI Flag | Environment Variable | Config Key | Default | Description |
|--------|----------|---------------------|------------|---------|-------------|
| Host | `--host` | `CLAUDE_GATE_HOST` | `host` | `127.0.0.1` | IP address to bind to |
| Port | `--port` | `CLAUDE_GATE_PORT` | `port` | `5789` | Port number for the server |
| Proxy Auth Token | `--proxy-auth-token` | `CLAUDE_GATE_PROXY_AUTH_TOKEN` | `proxy_auth_token` | (none) | Token for proxy authentication |

### Logging Configuration

| Option | CLI Flag | Environment Variable | Config Key | Default | Description |
|--------|----------|---------------------|------------|---------|-------------|
| Log Level | `--log-level` | `CLAUDE_GATE_LOG_LEVEL` | `log_level` | `info` | Logging level: debug, info, warning, error |
| Log File | `--log-file` | `CLAUDE_GATE_LOG_FILE` | `log_file` | (stdout) | Path to log file |
| Log Format | `--log-format` | `CLAUDE_GATE_LOG_FORMAT` | `log_format` | `text` | Log format: text, json |

### Security Configuration

| Option | CLI Flag | Environment Variable | Config Key | Default | Description |
|--------|----------|---------------------|------------|---------|-------------|
| Allowed Origins | `--allowed-origins` | `CLAUDE_GATE_ALLOWED_ORIGINS` | `allowed_origins` | `["*"]` | CORS allowed origins |
| TLS Certificate | `--tls-cert` | `CLAUDE_GATE_TLS_CERT` | `tls.cert` | (none) | Path to TLS certificate |
| TLS Key | `--tls-key` | `CLAUDE_GATE_TLS_KEY` | `tls.key` | (none) | Path to TLS private key |

### Dashboard Configuration

| Option | CLI Flag | Environment Variable | Config Key | Default | Description |
|--------|----------|---------------------|------------|---------|-------------|
| Dashboard Enabled | `--dashboard` | `CLAUDE_GATE_DASHBOARD` | `dashboard.enabled` | `false` | Enable interactive dashboard |
| Refresh Interval | `--dashboard-refresh` | `CLAUDE_GATE_DASHBOARD_REFRESH` | `dashboard.refresh_interval` | `1s` | Dashboard update frequency |

### Authentication Configuration

| Option | Environment Variable | Config Key | Default | Description |
|--------|---------------------|------------|---------|-------------|
| Token Storage Path | `CLAUDE_GATE_TOKEN_PATH` | `auth.token_path` | (platform-specific) | Where to store auth tokens |
| Token Encryption | `CLAUDE_GATE_TOKEN_ENCRYPT` | `auth.encrypt_tokens` | `true` | Encrypt stored tokens |

## Platform-Specific Defaults

### Token Storage Locations

- **macOS**: `~/Library/Application Support/claude-gate/`
- **Linux**: `~/.config/claude-gate/`
- **Windows**: `%APPDATA%\claude-gate\`

### Log File Locations (when specified)

- **macOS/Linux**: `/var/log/claude-gate.log` (requires permissions)
- **Windows**: `%LOCALAPPDATA%\claude-gate\logs\claude-gate.log`

## Environment Variable Format

### Simple Values

```bash
export CLAUDE_GATE_HOST=0.0.0.0
export CLAUDE_GATE_PORT=3000
export CLAUDE_GATE_LOG_LEVEL=debug
```

### List Values

For array configurations, use comma-separated values:

```bash
export CLAUDE_GATE_ALLOWED_ORIGINS="http://localhost:3000,https://myapp.com"
```

### Boolean Values

Use `true`, `false`, `1`, or `0`:

```bash
export CLAUDE_GATE_DASHBOARD=true
export CLAUDE_GATE_TOKEN_ENCRYPT=1
```

## Configuration Examples

### Development Configuration

```yaml
# ~/.claude-gate/config.dev.yaml
host: 127.0.0.1
port: 5789
log_level: debug
log_format: text
dashboard:
  enabled: true
  refresh_interval: 500ms
```

### Production Configuration

```yaml
# ~/.claude-gate/config.prod.yaml
host: 127.0.0.1
port: 5789
log_level: info
log_format: json
log_file: /var/log/claude-gate.log
proxy_auth_token: ${PROXY_AUTH_TOKEN}  # Can use env var substitution
allowed_origins:
  - https://app.example.com
  - https://www.example.com
tls:
  cert: /etc/claude-gate/cert.pem
  key: /etc/claude-gate/key.pem
```

### Docker Configuration

```yaml
# config.docker.yaml
host: 0.0.0.0  # Listen on all interfaces in container
port: 5789
log_level: info
log_format: json  # Structured logs for container logging
```

## Configuration Validation

Claude Gate validates configuration on startup:

1. **Type validation** - Ensures values are correct types
2. **Range validation** - Ports must be 1-65535
3. **File validation** - TLS cert/key files must exist if specified
4. **Permission validation** - Log file must be writable

To validate a configuration file without starting the server:

```bash
claude-gate config validate ~/.claude-gate/config.yaml
```

## Configuration Best Practices

1. **Security**
   - Never commit `proxy_auth_token` to version control
   - Use environment variables for sensitive values
   - Restrict file permissions on config files with secrets

2. **Performance**
   - Disable dashboard in production
   - Use appropriate log levels (info or warning for production)
   - Consider log rotation for file logging

3. **Maintainability**
   - Use separate config files for different environments
   - Document custom configuration in your project
   - Keep configuration minimal - rely on defaults where possible

4. **Monitoring**
   - Enable JSON logging for log aggregation systems
   - Include instance identifiers in multi-instance deployments
   - Monitor configuration changes

---

[← Reference](../README.md#reference) | [Documentation Home](../README.md) | [CLI Reference →](./cli.md)