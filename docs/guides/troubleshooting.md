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
  git clone https://github.com/ml0-1337/claude-gate
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
4. Check storage backend status:
   ```bash
   claude-gate auth storage status
   ```
5. For keychain issues, see [Token Storage Guide](./storage.md)

### Proxy Connection Refused

**Problem:** Cannot connect to the proxy server.

**Solution:**
1. Check if the proxy is running:
   ```bash
   claude-gate test
   ```
2. Verify the port isn't already in use:
   ```bash
   lsof -i :5789  # macOS/Linux
   netstat -an | findstr :5789  # Windows
   ```
3. Try a different port:
   ```bash
   claude-gate start --port 5789
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

1. Check existing issues: https://github.com/ml0-1337/claude-gate/issues
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

## Token Storage Issues

### Keychain Access Denied

**Problem:** "keyring access denied" or repeated password prompts on macOS.

**Solution:**
1. When prompted, click "Always Allow" instead of "Allow"
2. Check Keychain Access app for any denied entries
3. Reset keychain permissions:
   ```bash
   security unlock-keychain ~/Library/Keychains/login.keychain
   ```

### macOS Repeated Password Prompts (Fixed)

**Problem:** macOS shows password prompts every time Claude Gate accesses the keychain.

**Solution:**
As of version 0.2.0+, Claude Gate automatically configures the `KeychainTrustApplication` setting to prevent repeated password prompts. You should only see one prompt on first use where you can click "Always Allow".

If you still experience repeated prompts:

1. **Update to the latest version:**
   ```bash
   npm update -g claude-gate
   ```

2. **Clear existing keychain entries and re-authenticate:**
   ```bash
   claude-gate auth logout
   claude-gate auth login
   ```

3. **Override default settings (if needed):**
   ```bash
   # Disable application trust (not recommended)
   export CLAUDE_GATE_KEYCHAIN_TRUST_APP=false
   
   # Change accessibility settings
   export CLAUDE_GATE_KEYCHAIN_ACCESSIBLE_WHEN_UNLOCKED=false
   
   # Enable iCloud sync (security risk - not recommended)
   export CLAUDE_GATE_KEYCHAIN_SYNCHRONIZABLE=true
   ```

4. **Check keychain permissions in Keychain Access app:**
   - Open Keychain Access (Applications > Utilities)
   - Search for "claude-gate"
   - Double-click the entry
   - Go to Access Control tab
   - Ensure claude-gate is in the allowed applications list

### Linux Keyring Not Available

**Problem:** "keyring backend not available" on Linux.

**Solution:**
1. Install required packages:
   ```bash
   # Debian/Ubuntu
   sudo apt-get install gnome-keyring libsecret-1-0
   
   # Fedora/RHEL
   sudo dnf install gnome-keyring libsecret
   ```
2. Ensure D-Bus is running:
   ```bash
   echo $DBUS_SESSION_BUS_ADDRESS
   ```
3. Start keyring daemon:
   ```bash
   gnome-keyring-daemon --start --daemonize
   ```

### Automatic Migration Fails

**Problem:** Token migration from file to keychain fails.

**Solution:**
1. Check current storage status:
   ```bash
   claude-gate auth storage status
   ```
2. Manually migrate:
   ```bash
   claude-gate auth storage migrate --from file --to keyring
   ```
3. Force file storage if keychain issues persist:
   ```bash
   export CLAUDE_GATE_AUTH_STORAGE_TYPE=file
   ```

### Tokens Lost After Update

**Problem:** Authentication required after updating Claude Gate.

**Solution:**
1. Check for migrated file:
   ```bash
   ls ~/.claude-gate/auth.json.migrated
   ```
2. Restore from backup:
   ```bash
   cp ~/.claude-gate/backups/auth-*.json ~/.claude-gate/auth.json
   ```
3. Re-authenticate if necessary:
   ```bash
   claude-gate auth login
   ```

For more storage-related issues, see the [Token Storage Guide](./storage.md).

---

[← Guides](../README.md#guides) | [Documentation Home](../README.md) | [Contributing →](./contributing.md)