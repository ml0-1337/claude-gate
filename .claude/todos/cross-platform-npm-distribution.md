---
todo_id: cross-platform-npm-distribution
started: 2025-06-24 11:51:47
completed:
status: in_progress
priority: high
---

# Task: Cross-Platform Builds & NPM Distribution for Claude Gate

## Findings & Research

### GoReleaser Best Practices (2025)
- GoReleaser is the de facto standard for Go release automation
- Supports multiple platforms and architectures out of the box
- Integrates with GitHub Actions for CI/CD
- Can generate checksums and sign binaries with cosign
- Supports NPM publishing as a distribution method

### NPM Distribution Strategy
- Use optionalDependencies for platform-specific binaries (more secure than postinstall)
- Create separate packages for each platform to minimize download size
- Main package acts as installer/selector for correct platform binary
- Postinstall scripts can be security risk - use with caution and integrity checks

### go-npm Alternative
- Established pattern for Go binary distribution via NPM
- Uses postinstall to download binaries from GitHub Releases
- Simpler setup but requires internet access during install
- Less secure than bundling binaries in platform packages

### Platform Matrix
- darwin/amd64 (macOS Intel)
- darwin/arm64 (macOS Apple Silicon) 
- linux/amd64 (most servers)
- linux/arm64 (ARM servers, Raspberry Pi)
- windows/amd64 (Windows 64-bit)

### Binary Size Optimization
- Use ldflags="-s -w" to strip debug info
- UPX compression optional but may trigger antivirus
- Target <15MB per platform for fast npm installs

## Test Strategy

- **Test Framework**: Go standard testing + GitHub Actions matrix builds
- **Test Types**: Build tests, Installation tests, Integration tests
- **Coverage Target**: All platforms must build and run successfully
- **Edge Cases**: 
  - Installation with --ignore-scripts
  - Upgrade from previous versions
  - Permission issues on Unix systems
  - PATH configuration on Windows

## Test Cases

```bash
# Test 1: Cross-platform builds
# Input: goreleaser build --snapshot
# Expected: Binaries for all platforms generated

# Test 2: NPM package installation - macOS
# Input: npm install -g claude-gate (on macOS)
# Expected: Correct darwin binary installed and executable

# Test 3: NPM package installation - Windows
# Input: npm install -g claude-gate (on Windows)
# Expected: Windows binary + .cmd wrapper installed

# Test 4: Binary execution
# Input: claude-gate version (after npm install)
# Expected: Version output without errors

# Test 5: Upgrade scenario
# Input: npm update -g claude-gate
# Expected: New version installed, old version removed

# Test 6: Uninstall cleanup
# Input: npm uninstall -g claude-gate
# Expected: All binaries and wrappers removed
```

## Maintainability Analysis

- **Readability**: [9/10] Clear separation of build configs and platform code
- **Complexity**: GoReleaser abstracts most complexity
- **Modularity**: Separate npm packages for each platform
- **Testability**: GitHub Actions matrix enables comprehensive testing
- **Trade-offs**: More npm packages to maintain vs single package with all binaries

## Test Results Log

```bash
# Test runs will be logged here
```

## Checklist

### Phase 1: GoReleaser Setup
- [ ] Install GoReleaser locally for testing
- [ ] Create .goreleaser.yml configuration
- [ ] Configure multi-platform builds
- [ ] Set up archive formats (tar.gz, zip)
- [ ] Configure binary naming convention
- [ ] Add checksum generation
- [ ] Test local snapshot builds

### Phase 2: GitHub Actions
- [ ] Create .github/workflows/release.yml
- [ ] Configure trigger on version tags
- [ ] Set up GoReleaser action
- [ ] Configure GitHub token permissions
- [ ] Add release notes generation
- [ ] Test workflow with test tag

### Phase 3: NPM Package Structure
- [ ] Create npm/ directory
- [ ] Create main package.json
- [ ] Create platform detection install.js script
- [ ] Create uninstall.js cleanup script
- [ ] Add integrity check logic
- [ ] Handle --ignore-scripts gracefully
- [ ] Create bin wrapper script

### Phase 4: Platform Packages
- [ ] Create @claude-gate/darwin-x64 package
- [ ] Create @claude-gate/darwin-arm64 package
- [ ] Create @claude-gate/linux-x64 package
- [ ] Create @claude-gate/linux-arm64 package
- [ ] Create @claude-gate/win32-x64 package
- [ ] Configure package.json for each
- [ ] Set up publishing workflow

### Phase 5: Testing
- [ ] Test macOS Intel installation
- [ ] Test macOS Apple Silicon installation
- [ ] Test Linux amd64 installation
- [ ] Test Linux arm64 installation
- [ ] Test Windows installation
- [ ] Test upgrade scenarios
- [ ] Test with yarn and pnpm
- [ ] Verify PATH setup

### Phase 6: Documentation
- [ ] Update README with npm install instructions
- [ ] Document platform support
- [ ] Add troubleshooting guide
- [ ] Document build process
- [ ] Add npm badges to README

### Phase 7: Publishing
- [ ] Set up npm account/org
- [ ] Configure npm publish tokens
- [ ] Test npm publish --dry-run
- [ ] Publish beta version
- [ ] Test beta installation
- [ ] Publish stable version
- [ ] Announce release

## Working Scratchpad

### Requirements
1. Simple installation: `npm install -g claude-gate`
2. Support all major platforms (Windows, macOS, Linux)
3. Binary size <15MB per platform
4. Automatic PATH configuration
5. Clean uninstall process
6. Secure distribution with checksums

### Approach
1. GoReleaser for automated multi-platform builds
2. GitHub Actions for CI/CD
3. NPM with optionalDependencies for platform packages
4. Platform detection in main package
5. Binary integrity verification
6. Comprehensive testing matrix

### Code

### Notes
- Consider using esbuild's approach with optionalDependencies
- NPM scoped packages allow better organization
- GoReleaser can handle NPM publishing directly
- Windows needs special .cmd wrapper for global bins
- Consider Homebrew tap for macOS as alternative

### Commands & Output

```bash
# Install GoReleaser
brew install goreleaser

# Test build locally
goreleaser build --snapshot --clean

# Check binary sizes
ls -lh dist/

# Test npm package locally
cd npm && npm pack

# Dry run npm publish
npm publish --dry-run
```

### Current Binary Size
- claude-gate binary: ~8-10MB (estimated)
- After stripping: ~6-8MB (estimated)
- Compressed in archive: ~3-4MB (estimated)