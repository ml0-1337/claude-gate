# NPM Package Guide

This guide covers the Claude Gate NPM package distribution, including installation, usage, and publishing.

## For Users

### Installation

Install Claude Gate globally via NPM:

```bash
npm install -g claude-gate
```

Or using Yarn:

```bash
yarn global add claude-gate
```

Or using pnpm:

```bash
pnpm add -g claude-gate
```

### What Gets Installed

The NPM package:
1. Downloads the appropriate binary for your platform
2. Installs it in your global NPM bin directory
3. Makes the `claude-gate` command available system-wide

### Supported Platforms

The NPM package supports the following platforms:

| Platform | Architecture | Package Name |
|----------|--------------|--------------|
| macOS | Intel (x64) | `@claude-gate/darwin-x64` |
| macOS | Apple Silicon (arm64) | `@claude-gate/darwin-arm64` |
| Linux | x64 | `@claude-gate/linux-x64` |
| Linux | arm64 | `@claude-gate/linux-arm64` |
| Windows | x64 | `@claude-gate/win32-x64` |

### Verifying Installation

After installation, verify it works:

```bash
claude-gate version
```

### Updating

Update to the latest version:

```bash
npm update -g claude-gate
```

### Uninstalling

Remove Claude Gate:

```bash
npm uninstall -g claude-gate
```

## For Maintainers

### Package Structure

The NPM package uses a multi-package architecture:

```
npm/
├── package.json           # Main package
├── index.js              # Installation script
├── scripts/
│   ├── install.js        # Post-install script
│   └── uninstall.js      # Pre-uninstall script
└── platforms/
    ├── darwin-arm64/
    │   └── package.json  # Platform-specific package
    ├── darwin-x64/
    ├── linux-arm64/
    ├── linux-x64/
    └── win32-x64/
```

### How It Works

1. **Main Package** (`claude-gate`):
   - Contains installation logic
   - Has optional dependencies on platform packages
   - Runs post-install script to set up binary

2. **Platform Packages** (`@claude-gate/platform-arch`):
   - Contains the actual binary for that platform
   - Only the matching platform package is installed

3. **Installation Process**:
   - NPM installs main package
   - NPM installs matching platform package
   - Post-install script extracts and sets up binary

### Building NPM Packages

#### Prerequisites

- Built binaries for all platforms
- Node.js and NPM installed
- NPM authentication token

#### Build Process

1. **Build all binaries**:
   ```bash
   make build-all
   ```

2. **Prepare NPM packages**:
   ```bash
   make build-npm
   ```

3. **Test locally**:
   ```bash
   cd npm
   npm link
   claude-gate version
   ```

### Publishing

#### First-Time Setup

1. **Create NPM organization** (if needed):
   - Go to npmjs.com
   - Create organization: `@claude-gate`

2. **Authenticate with NPM**:
   ```bash
   npm login
   ```

3. **Set up automated publishing** (optional):
   ```bash
   ./scripts/setup-npm-auth.sh
   ```

#### Publishing Process

1. **Update version** in all package.json files:
   ```bash
   ./scripts/update-version.sh 1.2.3
   ```

2. **Build and test packages**:
   ```bash
   make build-npm
   ./scripts/test-npm-local.sh
   ```

3. **Publish packages**:
   ```bash
   # Publish platform packages first
   for platform in npm/platforms/*/; do
     cd "$platform"
     npm publish --access public
     cd -
   done

   # Publish main package
   cd npm
   npm publish --access public
   ```

#### Automated Publishing

The GitHub Actions workflow handles publishing on new releases:

```yaml
# .github/workflows/release.yml
- name: Publish to NPM
  run: |
    echo "//registry.npmjs.org/:_authToken=${{ secrets.NPM_TOKEN }}" > ~/.npmrc
    make publish-npm
```

### Version Management

#### Versioning Strategy

- Main package and platform packages share the same version
- Version follows semantic versioning (semver)
- All packages are published together

#### Updating Versions

Use the update script:

```bash
./scripts/update-version.sh 1.2.3
```

Or manually update:
1. `npm/package.json`
2. `npm/platforms/*/package.json`
3. Binary version in Go code

### Testing NPM Package

#### Local Testing

```bash
# Build packages
make build-npm

# Test installation
cd npm
npm link

# Verify
claude-gate version

# Test uninstall
npm unlink
```

#### Cross-Platform Testing

```bash
# Run comprehensive tests
./scripts/test-npm-local.sh

# Test in Docker containers
./scripts/test-docker.sh
```

#### Integration Tests

Test with different package managers:

```bash
# NPM
npm install -g ./npm

# Yarn
yarn global add file:./npm

# pnpm
pnpm add -g ./npm
```

### Troubleshooting Publishing

#### Authentication Issues

```bash
# Check authentication
npm whoami

# Re-authenticate
npm logout
npm login
```

#### Permission Issues

```bash
# Check package access
npm access ls-packages

# Grant publish access
npm owner add USERNAME @claude-gate/PACKAGE
```

#### Failed Publish

1. Check version conflicts:
   ```bash
   npm view claude-gate versions
   ```

2. Try force publish (careful!):
   ```bash
   npm publish --force
   ```

3. Unpublish broken version (within 72 hours):
   ```bash
   npm unpublish claude-gate@1.2.3
   ```

### NPM Package Best Practices

1. **Security**:
   - Never include sensitive data in packages
   - Use `.npmignore` to exclude unnecessary files
   - Regular security audits: `npm audit`

2. **Size Optimization**:
   - Only include necessary files
   - Compress binaries when possible
   - Check package size: `npm pack --dry-run`

3. **Compatibility**:
   - Test on all supported Node.js versions
   - Handle missing optional dependencies gracefully
   - Provide clear error messages

4. **Documentation**:
   - Keep README.md updated
   - Include platform requirements
   - Document common issues

### Monitoring NPM Package

#### Download Statistics

View package statistics:

```bash
# Command line
npm view claude-gate

# Or visit
# https://www.npmjs.com/package/claude-gate
```

#### User Issues

Monitor:
- GitHub issues with `npm` label
- NPM package page reviews
- Social media mentions

#### Security Monitoring

```bash
# Check for vulnerabilities
npm audit

# Update dependencies
npm update

# Check outdated packages
npm outdated
```

---

[← Deployment](../README.md#deployment) | [Documentation Home](../README.md)