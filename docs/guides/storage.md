# Token Storage Guide

Claude Gate now supports secure platform-native token storage using your operating system's built-in security features.

## Overview

By default, Claude Gate automatically selects the most secure storage backend available on your system:

- **macOS**: Keychain Services
- **Linux**: Secret Service API (GNOME Keyring, KWallet)
- **Fallback**: Encrypted file storage when OS keychain is unavailable

## Storage Backends

### OS Keychain (Recommended)
Your tokens are stored in the operating system's secure credential storage:
- Hardware-backed security when available
- Encrypted at rest by the OS
- Integrated with system security (Touch ID, etc.)
- Survives application updates

### File Storage (Fallback)
When keychain access is unavailable:
- Tokens stored in `~/.claude-gate/auth.json`
- File permissions restricted to owner only (0600)
- Optional encryption using JOSE (JSON Web Encryption)
- Portable across systems

## Configuration

### Environment Variables

```bash
# Storage backend selection (auto, keyring, file)
export CLAUDE_GATE_AUTH_STORAGE_TYPE=auto

# Keyring service name
export CLAUDE_GATE_KEYRING_SERVICE=claude-gate

# File storage path
export CLAUDE_GATE_AUTH_STORAGE_PATH=~/.claude-gate/auth.json

# Auto-migrate tokens to keyring
export CLAUDE_GATE_AUTO_MIGRATE_TOKENS=true
```

### Storage Commands

```bash
# Check current storage backend
claude-gate auth storage status

# Migrate tokens between backends
claude-gate auth storage migrate --from file --to keyring

# Test storage operations
claude-gate auth storage test

# Create backup (file storage only)
claude-gate auth storage backup
```

## Migration

When you update Claude Gate, existing tokens are automatically migrated to the most secure available backend:

1. **Automatic Migration**: Happens on first use after update
2. **Manual Migration**: Use `claude-gate auth storage migrate`
3. **Backup Creation**: Original tokens are backed up before migration
4. **Verification**: Migration is verified before marking complete

## Troubleshooting

### Keychain Access Issues

**macOS**: If prompted for keychain access, click "Always Allow" to avoid repeated prompts.

**Linux**: Ensure a keyring daemon is running:
```bash
# Check if keyring is available
echo $GNOME_KEYRING_CONTROL
# or
echo $SSH_AUTH_SOCK
```

### Fallback Behavior

If keychain access fails, Claude Gate automatically falls back to file storage with a warning. To force a specific backend:

```bash
# Force file storage
export CLAUDE_GATE_AUTH_STORAGE_TYPE=file

# Force keyring (fails if unavailable)
export CLAUDE_GATE_AUTH_STORAGE_TYPE=keyring
```

### Permission Errors

**Linux**: You may need to unlock your keyring on first use:
```bash
# Unlock GNOME Keyring
gnome-keyring-daemon --unlock
```

**All Platforms**: Ensure you have write permissions to `~/.claude-gate/`:
```bash
mkdir -p ~/.claude-gate
chmod 700 ~/.claude-gate
```

## Security Considerations

1. **Keychain Security**: Your tokens are as secure as your system login
2. **Backup Security**: Backups are stored in `~/.claude-gate/backups/` with restricted permissions
3. **Migration Security**: Tokens are never transmitted during migration
4. **Revocation**: Use `claude-gate auth logout` to remove all stored tokens

## Best Practices

1. **Use OS Keychain**: Let Claude Gate use your OS keychain when available
2. **Regular Backups**: Periodically backup tokens if using file storage
3. **Monitor Access**: Check `claude-gate auth storage status` regularly
4. **Update Promptly**: Keep Claude Gate updated for latest security fixes
5. **Logout When Done**: Remove tokens if not using Claude Gate for extended periods

## Platform-Specific Notes

### macOS
- Tokens stored in login keychain by default
- Synced across devices if iCloud Keychain enabled
- May prompt for password after system updates

### Linux
- Requires D-Bus and keyring daemon
- Works with GNOME Keyring, KDE Wallet, etc.
- May require unlocking after reboot

### Docker/Containers
- Keychain not available in containers
- Always uses file storage
- Mount volume for persistence: `-v ~/.claude-gate:/root/.claude-gate`

## Advanced Usage

### Custom Password for File Backend

When using encrypted file storage, you can provide a custom password:

```go
// In code - implement PasswordPromptFunc
func myPasswordPrompt(prompt string) (string, error) {
    // Return password from secure source
    return os.Getenv("CLAUDE_GATE_MASTER_PASSWORD"), nil
}
```

### Multiple Profiles

Store tokens for different accounts:

```bash
# Use different service names
export CLAUDE_GATE_KEYRING_SERVICE=claude-gate-work
claude-gate auth login

export CLAUDE_GATE_KEYRING_SERVICE=claude-gate-personal  
claude-gate auth login
```

For questions or issues, please see our [Troubleshooting Guide](./troubleshooting.md) or [open an issue](https://github.com/ml0-1337/claude-gate/issues).