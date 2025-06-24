# CLI Reference

Complete reference for all Claude Gate command-line interface commands and options.

## Global Options

These options can be used with any command:

```bash
claude-gate [global options] command [command options] [arguments...]
```

| Option | Description | Default |
|--------|-------------|---------|
| `--help`, `-h` | Show help | - |
| `--version`, `-v` | Show version | - |
| `--config FILE` | Load configuration from FILE | `~/.claude-gate/config.yaml` |
| `--log-level LEVEL` | Set log level (DEBUG, INFO, WARNING, ERROR) | `INFO` |

## Commands

### `auth` - Authentication Management

Manage authentication with Claude Pro/Max accounts.

#### `auth login`

Authenticate with your Claude account:

```bash
claude-gate auth login [options]
```

**Options:**
- `--browser` - Force browser authentication (default: auto-detect)
- `--no-browser` - Use terminal-only authentication
- `--timeout DURATION` - Authentication timeout (default: `5m`)

**Example:**
```bash
claude-gate auth login --timeout 10m
```

#### `auth logout`

Remove stored authentication:

```bash
claude-gate auth logout
```

#### `auth status`

Check authentication status:

```bash
claude-gate auth status [options]
```

**Options:**
- `--json` - Output in JSON format
- `--verbose` - Show detailed token information

**Example:**
```bash
claude-gate auth status --json
```

#### `auth refresh`

Refresh authentication token:

```bash
claude-gate auth refresh
```

### `start` - Start Proxy Server

Start the Claude Gate proxy server:

```bash
claude-gate start [options]
```

**Options:**
| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--host` | `CLAUDE_GATE_HOST` | `127.0.0.1` | Host to bind to |
| `--port` | `CLAUDE_GATE_PORT` | `5789` | Port to listen on |
| `--dashboard` | - | `false` | Enable interactive dashboard |
| `--daemon` | - | `false` | Run in background |
| `--proxy-auth-token` | `CLAUDE_GATE_PROXY_AUTH_TOKEN` | - | Require authentication |
| `--tls-cert` | - | - | TLS certificate file |
| `--tls-key` | - | - | TLS key file |

**Examples:**
```bash
# Start on default port
claude-gate start

# Start with dashboard
claude-gate start --dashboard

# Start on custom port with auth
claude-gate start --port 3000 --proxy-auth-token "secret-token"

# Start with TLS
claude-gate start --tls-cert cert.pem --tls-key key.pem
```

### `stop` - Stop Proxy Server

Stop the running Claude Gate proxy:

```bash
claude-gate stop [options]
```

**Options:**
- `--force` - Force stop without graceful shutdown
- `--timeout DURATION` - Shutdown timeout (default: `30s`)

**Example:**
```bash
claude-gate stop --timeout 10s
```

### `status` - Check Server Status

Check if Claude Gate is running:

```bash
claude-gate status [options]
```

**Options:**
- `--json` - Output in JSON format
- `--verbose` - Show detailed server information

**Output includes:**
- Running status
- PID
- Port number
- Uptime
- Request statistics

**Example:**
```bash
claude-gate status --json
```

### `logs` - View Server Logs

Display Claude Gate server logs:

```bash
claude-gate logs [options]
```

**Options:**
- `--follow`, `-f` - Follow log output
- `--tail N` - Show last N lines (default: `100`)
- `--since DURATION` - Show logs since duration ago
- `--level LEVEL` - Filter by log level

**Examples:**
```bash
# Follow logs in real-time
claude-gate logs -f

# Show last 50 error logs
claude-gate logs --tail 50 --level ERROR

# Show logs from last hour
claude-gate logs --since 1h
```

### `config` - Configuration Management

Manage Claude Gate configuration:

#### `config show`

Display current configuration:

```bash
claude-gate config show [options]
```

**Options:**
- `--json` - Output in JSON format
- `--show-defaults` - Include default values

#### `config validate`

Validate configuration file:

```bash
claude-gate config validate [FILE]
```

#### `config init`

Create initial configuration file:

```bash
claude-gate config init [options]
```

**Options:**
- `--force` - Overwrite existing configuration

### `version` - Show Version Information

Display Claude Gate version:

```bash
claude-gate version [options]
```

**Options:**
- `--json` - Output in JSON format
- `--verbose` - Show build information

**Example output:**
```
Claude Gate version 1.0.0
Built: 2024-01-20T10:00:00Z
Go version: go1.22.0
Platform: darwin/arm64
```

## Environment Variables

Claude Gate respects the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `CLAUDE_GATE_HOST` | Default host for server | `127.0.0.1` |
| `CLAUDE_GATE_PORT` | Default port for server | `5789` |
| `CLAUDE_GATE_CONFIG` | Configuration file path | `~/.claude-gate/config.yaml` |
| `CLAUDE_GATE_LOG_LEVEL` | Default log level | `INFO` |
| `CLAUDE_GATE_LOG_FILE` | Log file path | - |
| `CLAUDE_GATE_PROXY_AUTH_TOKEN` | Proxy authentication token | - |
| `CLAUDE_GATE_DASHBOARD` | Enable dashboard by default | `false` |
| `CLAUDE_GATE_ALLOWED_ORIGINS` | CORS allowed origins | `*` |
| `NO_COLOR` | Disable colored output | - |

## Exit Codes

Claude Gate uses the following exit codes:

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Authentication error |
| 4 | Server already running |
| 5 | Server not running |
| 127 | Command not found |

## Examples

### Basic Usage

```bash
# Authenticate
claude-gate auth login

# Start server
claude-gate start

# Check status
claude-gate status

# View logs
claude-gate logs -f

# Stop server
claude-gate stop
```

### Production Setup

```bash
# Set up with authentication
export CLAUDE_GATE_PROXY_AUTH_TOKEN=$(openssl rand -hex 32)
export CLAUDE_GATE_LOG_LEVEL=INFO
export CLAUDE_GATE_LOG_FILE=/var/log/claude-gate.log

# Start with specific configuration
claude-gate start --host 127.0.0.1 --port 5789

# Monitor
claude-gate status --json | jq
```

### Development Setup

```bash
# Start with debug logging and dashboard
claude-gate start --log-level DEBUG --dashboard

# Follow logs in another terminal
claude-gate logs -f --level DEBUG
```

## Tips and Tricks

1. **Auto-completion**: Enable shell completion:
   ```bash
   # Bash
   claude-gate completion bash > /etc/bash_completion.d/claude-gate
   
   # Zsh
   claude-gate completion zsh > "${fpath[1]}/_claude-gate"
   ```

2. **Aliases**: Add useful aliases:
   ```bash
   alias cg='claude-gate'
   alias cgs='claude-gate start --dashboard'
   alias cgl='claude-gate logs -f'
   ```

3. **Configuration Profiles**: Use different configs:
   ```bash
   claude-gate --config ~/.claude-gate/dev.yaml start
   claude-gate --config ~/.claude-gate/prod.yaml start
   ```

---

[← Reference](../README.md#reference) | [Documentation Home](../README.md) | [API Reference →](./api.md)