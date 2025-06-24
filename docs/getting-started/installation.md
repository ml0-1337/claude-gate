# Installation Guide

This guide covers all the ways to install Claude Gate on your system.

## Prerequisites

- **Claude Pro or Claude Max subscription** - Required for authentication
- **Supported Operating System**:
  - macOS (Intel & Apple Silicon)
  - Linux (x64 & ARM64)
  - Windows (x64)

## Installation Methods

### Via NPM (Recommended)

The easiest way to install Claude Gate is through NPM:

```bash
npm install -g claude-gate
```

This will automatically download the correct binary for your platform.

### Via Homebrew (macOS/Linux)

If you prefer Homebrew:

```bash
brew tap ml0-1337/tap
brew install claude-gate
```

### From Source

If you have Go 1.22+ installed:

```bash
go install github.com/ml0-1337/claude-gate/cmd/claude-gate@latest
```

### Direct Download

Download pre-built binaries from the [releases page](https://github.com/ml0-1337/claude-gate/releases).

1. Download the appropriate binary for your platform
2. Extract the archive
3. Move the binary to a location in your PATH (e.g., `/usr/local/bin`)
4. Make it executable: `chmod +x claude-gate` (macOS/Linux)

## Verify Installation

After installation, verify Claude Gate is working:

```bash
claude-gate version
```

## Development Installation

If you're planning to contribute to Claude Gate, see our [Development Guide](../guides/development.md) for setting up a development environment.

### Development Prerequisites

- Go 1.22 or later
- Node.js 18 or later (for NPM package testing)
- Git
- Make
- GoReleaser (for release builds)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/ml0-1337/claude-gate.git
cd claude-gate

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

## Configuration

Claude Gate uses the following default settings:

- **Host**: `127.0.0.1` (localhost only)
- **Port**: `5789`

You can customize these using environment variables:

```bash
export CLAUDE_GATE_HOST=127.0.0.1
export CLAUDE_GATE_PORT=5789
export CLAUDE_GATE_PROXY_AUTH_TOKEN=your-token  # Optional: Enable proxy authentication
export CLAUDE_GATE_LOG_LEVEL=INFO               # Options: DEBUG, INFO, WARNING, ERROR
```

For more configuration options, see the [Configuration Guide](./configuration.md).

## Troubleshooting Installation

Having issues? Check our [Troubleshooting Guide](../guides/troubleshooting.md) or common solutions below:

### NPM Installation Fails

If NPM installation fails, try:

1. Clear NPM cache: `npm cache clean --force`
2. Use a different registry: `npm install -g claude-gate --registry https://registry.npmjs.org/`
3. Install with verbose logging: `npm install -g claude-gate --verbose`

### Permission Errors

On macOS/Linux, you might need to use sudo:

```bash
sudo npm install -g claude-gate
```

Or configure NPM to use a different directory for global packages.

### Binary Not Found

If the `claude-gate` command is not found after installation:

1. Check if it's in your PATH: `which claude-gate`
2. Find where NPM installs global packages: `npm root -g`
3. Add the NPM bin directory to your PATH

## Next Steps

Once installed, proceed to the [Quick Start Guide](./quick-start.md) to authenticate and start using Claude Gate.

---

[← Back to Documentation](../README.md) | [Quick Start →](./quick-start.md)