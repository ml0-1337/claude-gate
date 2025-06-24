# Security Policy

## Overview

Claude Gate handles sensitive authentication tokens and proxies API requests. This document outlines our security model, best practices, and vulnerability reporting process.

## Security Model

### OAuth 2.0 Public Client with PKCE

Claude Gate uses OAuth 2.0 with Proof Key for Code Exchange (PKCE) for authentication:

- **Public Client**: The OAuth client ID is intentionally public
- **PKCE Protection**: Security comes from the PKCE challenge/verifier, not client secrecy
- **Industry Standard**: Same approach used by GitHub CLI, Google Cloud SDK, and other major tools

### Token Security

#### Current Implementation
- **OS Keychain Integration** (Implemented)
  - macOS: Keychain Services
  - Linux: Secret Service API (GNOME Keyring, KWallet)
  - Windows: Credential Manager
- **Automatic Fallback**: File storage when keychain unavailable
- **Zero-Configuration**: Auto-detects best storage backend
- **Seamless Migration**: Automatic migration from JSON to keychain
- File permissions set to 0600 (user read/write only)
- Automatic token refresh before expiration
- Tokens never transmitted to clients

#### Storage Backend Selection
1. **Auto Mode** (Default)
   - Automatically selects the most secure available backend
   - Prefers OS keychain/keystore over file storage
   - Transparent fallback on errors
   
2. **Keyring Mode**
   - Forces use of OS-native secure storage
   - Fails if keychain not available
   
3. **File Mode**
   - Traditional JSON file storage
   - Encrypted file storage using JOSE (PBES2_HS256_A128KW + A256GCM)
   - For environments without keychain access

#### Storage Management
- `claude-gate auth storage status`: Check current backend
- `claude-gate auth storage migrate`: Migrate between backends
- `claude-gate auth storage test`: Verify storage operations
- `claude-gate auth storage backup`: Create manual backup

### Network Security

- **TLS Only**: All API communication uses HTTPS
- **Certificate Validation**: Full certificate chain validation
- **No Proxy Bypass**: Security settings cannot be disabled
- **Header Sanitization**: Removes identifying client headers

### Request Security

1. **Authentication Validation**
   - Every request checks token validity
   - Expired tokens trigger automatic refresh
   - Invalid tokens return 401 Unauthorized

2. **System Prompt Injection**
   - Prepends identifier to maintain API compliance
   - No user data logged or stored
   - Transformation is transparent and auditable

3. **Rate Limiting** (Planned)
   - Per-IP rate limiting to prevent abuse
   - Token bucket algorithm with configurable limits
   - Graceful degradation under load

## Security Best Practices

### For Users

1. **Protect Your Token**
   - Never share your `auth.json` file
   - Revoke tokens if system compromised
   - Use `claude-gate auth logout` when done

2. **System Security**
   - Keep your OS and Claude Gate updated
   - Use full-disk encryption
   - Don't run Claude Gate as root/admin

3. **Network Security**
   - Use Claude Gate only on trusted networks
   - Consider VPN for public WiFi
   - Monitor for unusual activity

### For Developers

1. **Code Security**
   - Never log tokens or sensitive data
   - Validate all inputs
   - Use constant-time comparisons for secrets
   - Keep dependencies updated

2. **Build Security**
   - Build with latest Go version
   - Enable all compiler security features
   - Sign releases with GPG
   - Generate reproducible builds

## Threat Model

### In Scope

1. **Token Theft**
   - Mitigation: Keychain storage, file encryption
   - Detection: Token revocation on suspicious activity

2. **Man-in-the-Middle**
   - Mitigation: TLS with certificate pinning
   - Detection: Certificate validation failures

3. **Replay Attacks**
   - Mitigation: OAuth nonce, request signing
   - Detection: Duplicate request detection

4. **Resource Exhaustion**
   - Mitigation: Rate limiting, connection limits
   - Detection: Metrics and monitoring

### Out of Scope

1. **Physical Access**: If attacker has system access, game over
2. **OS Vulnerabilities**: We assume a secure OS
3. **Social Engineering**: User education responsibility
4. **Anthropic API Security**: We trust Anthropic's security

## Vulnerability Reporting

### Reporting Process

1. **Do NOT** create public GitHub issues for security vulnerabilities
2. Email security concerns to: [maintainer email]
3. Include:
   - Description of vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline

- **24 hours**: Initial acknowledgment
- **72 hours**: Preliminary assessment
- **7 days**: Fix development begins
- **30 days**: Fix released (critical issues faster)

### Disclosure Policy

- 90-day disclosure deadline
- Coordinated disclosure preferred
- Credit given to reporters (if desired)
- Security advisories posted to GitHub

## Security Checklist

### For Each Release

- [ ] Update all dependencies
- [ ] Run security scanner (gosec)
- [ ] Test token rotation
- [ ] Verify TLS settings
- [ ] Check file permissions
- [ ] Audit logging statements

### Monthly Reviews

- [ ] Review threat model
- [ ] Check for new CVEs
- [ ] Update security documentation
- [ ] Test incident response plan
- [ ] Review access logs

## Incident Response

### If You Suspect Compromise

1. **Immediately**: Run `claude-gate auth logout`
2. **Then**: Revoke tokens in Anthropic dashboard
3. **Check**: System for other compromises
4. **Report**: To maintainers if Claude Gate issue
5. **Monitor**: For unauthorized API usage

### For Maintainers

1. **Assess**: Severity and scope
2. **Contain**: Patch or disable affected features
3. **Communicate**: Notify affected users
4. **Fix**: Develop and test patch
5. **Release**: Emergency update
6. **Post-mortem**: Document and learn

## Security Features Roadmap

### Q1 2025
- [x] OAuth 2.0 PKCE implementation
- [x] OS keychain integration (Completed)
- [x] Encrypted file storage (via 99designs/keyring FileBackend)
- [ ] Basic rate limiting

### Q2 2025
- [ ] Certificate pinning
- [ ] Advanced rate limiting
- [ ] Security audit
- [ ] Penetration testing

### Q3 2025
- [ ] Hardware token support
- [ ] Multi-factor authentication
- [ ] Audit logging
- [ ] SIEM integration

## Compliance

Claude Gate is designed with security best practices but makes no compliance claims:

- Not evaluated for HIPAA, PCI, or other standards
- Users responsible for their own compliance
- We'll help with security questionnaires
- Open source for full auditability

## Resources

- [OWASP API Security](https://owasp.org/www-project-api-security/)
- [OAuth 2.0 Security Best Practices](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [Go Security Guidelines](https://golang.org/doc/security)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

---

[← Architecture](../README.md#architecture) | [Documentation Home](../README.md) | [Overview →](./overview.md)