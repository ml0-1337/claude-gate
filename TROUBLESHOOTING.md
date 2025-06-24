# Claude Gate Troubleshooting Guide

## Installation Issues

### NPM Installation Fails

**Problem:** `npm install -g claude-gate` fails with permission errors.

**Solution:**
```bash
# Option 1: Use a Node version manager (recommended)
# Install nvm: https://github.com/nvm-sh/nvm
nvm install node
npm install -g claude-gate

# Option 2: Change npm's default directory
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
npm install -g claude-gate
```

### Binary Not Found After Installation

**Problem:** `claude-gate: command not found` after successful installation.

**Solution:**
1. Check where npm installs global packages:
   ```bash
   npm config get prefix
   ```
2. Ensure the bin directory is in your PATH:
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   export PATH="$(npm config get prefix)/bin:$PATH"
   ```

### Platform Not Supported Error

**Problem:** "Unsupported platform" error during installation.

**Solution:**
- Check supported platforms: darwin-x64, darwin-arm64, linux-x64, linux-arm64, win32-x64
- For other platforms, build from source:
  ```bash
  git clone https://github.com/yourusername/claude-gate
  cd claude-gate
  go build -o claude-gate ./cmd/claude-gate
  ```

### Installation with --ignore-scripts

**Problem:** Binary doesn't work when installed with `npm install -g claude-gate --ignore-scripts`.

**Solution:**
The package includes a fallback mechanism that should handle this automatically. If it doesn't work:
```bash
# Manually run the install script
cd $(npm root -g)/claude-gate
node scripts/install.js
```

## Runtime Issues

### OAuth Authentication Fails

**Problem:** "OAuth token error" when starting the proxy.

**Solution:**
1. Re-authenticate:
   ```bash
   claude-gate auth logout
   claude-gate auth login
   ```
2. Check if you have an active Claude Pro/Max subscription
3. Ensure you're using the correct authorization code

### Proxy Connection Refused

**Problem:** Cannot connect to the proxy server.

**Solution:**
1. Check if the proxy is running:
   ```bash
   claude-gate test
   ```
2. Verify the port isn't already in use:
   ```bash
   lsof -i :8000  # macOS/Linux
   netstat -an | findstr :8000  # Windows
   ```
3. Try a different port:
   ```bash
   claude-gate start --port 8080
   ```

### SSE Streaming Not Working

**Problem:** Streaming responses appear to hang or timeout.

**Solution:**
1. Ensure your HTTP client supports SSE
2. Check if a proxy or firewall is interfering
3. Verify your client isn't buffering the response

### Token Expiration Issues

**Problem:** "Token is expired" errors during use.

**Solution:**
The proxy should automatically refresh tokens. If not:
```bash
claude-gate auth status  # Check token status
claude-gate auth login   # Re-authenticate if needed
```

## Build Issues

### GoReleaser Build Fails

**Problem:** `goreleaser build` fails with errors.

**Solution:**
1. Ensure Go 1.22+ is installed:
   ```bash
   go version
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Check for syntax errors:
   ```bash
   go build ./...
   ```

### NPM Package Testing Fails

**Problem:** `make npm-test` doesn't work correctly.

**Solution:**
1. Ensure binaries are built first:
   ```bash
   make snapshot
   ```
2. Check Node.js version (16+ required):
   ```bash
   node --version
   ```

## Common Error Messages

### "OAuth token error"
- **Cause:** No valid authentication or token expired
- **Fix:** Run `claude-gate auth login`

### "Failed to transform request"
- **Cause:** Invalid request format
- **Fix:** Ensure you're using a compatible Anthropic SDK

### "Upstream request failed"
- **Cause:** Network issues or Anthropic API down
- **Fix:** Check internet connection and Anthropic status

### "Unsupported platform"
- **Cause:** Running on an unsupported OS/architecture
- **Fix:** Build from source for your platform

## Getting Help

If you're still experiencing issues:

1. Check existing issues: https://github.com/yourusername/claude-gate/issues
2. Create a new issue with:
   - Your OS and architecture (`uname -a` on Unix)
   - Node.js version (`node --version`)
   - NPM version (`npm --version`)
   - Complete error message
   - Steps to reproduce

## Debug Mode

For more detailed logging:
```bash
CLAUDE_GATE_LOG_LEVEL=DEBUG claude-gate start
```

This will show detailed information about:
- OAuth token operations
- Request/response transformations
- Proxy operations
- Error details