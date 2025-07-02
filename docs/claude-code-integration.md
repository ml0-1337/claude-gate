# Claude Code Integration

Claude Gate can now use credentials from Claude Code (Anthropic's official CLI) directly, providing a seamless experience for users who have both tools installed.

## How It Works

When you use the `--storage-backend=claude-code` option, Claude Gate reads OAuth credentials directly from Claude Code's keychain storage. This means:

- No need to authenticate twice
- Credentials stay in sync automatically
- Claude Code handles token refresh
- Read-only access (Claude Gate won't modify Claude Code's credentials)

## Usage

### Start the proxy using Claude Code credentials:
```bash
claude-gate start --storage-backend=claude-code
```

### Or use the dashboard:
```bash
claude-gate dashboard --storage-backend=claude-code
```

## Requirements

- Claude Code must be installed and authenticated
- macOS: Keychain access must be allowed
- Linux: Secret Service must be available
- Windows: Windows Credential Manager access

## Storage Backend Options

- `auto` (default): Try keychain first, fall back to file storage
- `keyring`: Use claude-gate's own keychain storage
- `file`: Use JSON file storage
- `claude-code`: Use Claude Code's credentials (read-only)

## Technical Details

### Credential Transformation

Claude Code and claude-gate store credentials in slightly different formats. The adapter automatically handles:

- Field name mapping (e.g., `accessToken` → `AccessToken`)
- Timestamp conversion (milliseconds to seconds)
- Structure flattening (nested to flat JSON)
- Type field addition (always "oauth")

### Security

The Claude Code storage adapter:
- Only reads credentials, never writes
- Requires same user access as Claude Code
- Uses OS-native secure storage (Keychain/Secret Service/Credential Manager)
- No credentials are ever logged or exposed

## Troubleshooting

If Claude Code credentials aren't found:
1. Ensure Claude Code is authenticated: `claude auth login`
2. Check keychain permissions on macOS
3. Try running with `--storage-backend=auto` to fall back to claude-gate's storage
4. Use `claude-gate auth storage status` to debug

## Migration

To migrate from Claude Code to claude-gate's storage:
```bash
# This will copy credentials from Claude Code to claude-gate's keychain
claude-gate auth storage migrate --from=claude-code --to=keyring
```