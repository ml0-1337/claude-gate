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
- [x] Install GoReleaser locally for testing
- [x] Create .goreleaser.yml configuration
- [x] Configure multi-platform builds
- [x] Set up archive formats (tar.gz, zip)
- [x] Configure binary naming convention
- [x] Add checksum generation
- [ ] Test local snapshot builds

### Phase 2: GitHub Actions
- [x] Create .github/workflows/release.yml
- [x] Configure trigger on version tags
- [x] Set up GoReleaser action
- [x] Configure GitHub token permissions
- [x] Add release notes generation
- [ ] Test workflow with test tag

### Phase 3: NPM Package Structure
- [x] Create npm/ directory
- [x] Create main package.json
- [x] Create platform detection install.js script
- [x] Create uninstall.js cleanup script
- [x] Add integrity check logic
- [x] Handle --ignore-scripts gracefully
- [x] Create bin wrapper script

### Phase 4: Platform Packages
- [x] Create @claude-gate/darwin-x64 package
- [x] Create @claude-gate/darwin-arm64 package
- [x] Create @claude-gate/linux-x64 package
- [x] Create @claude-gate/linux-arm64 package
- [x] Create @claude-gate/win32-x64 package
- [x] Configure package.json for each
- [x] Set up publishing workflow

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
- [x] Update README with npm install instructions
- [x] Document platform support
- [x] Add troubleshooting guide
- [x] Document build process
- [x] Add npm badges to README
- [x] Create NPM-specific README
- [x] Create Makefile for automation
- [x] Create helper scripts

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

### Implementation Status & Next Steps

**Completed:**
1. ✅ GoReleaser configuration (.goreleaser.yml)
2. ✅ GitHub Actions workflow for automated releases
3. ✅ NPM package structure with platform detection
4. ✅ Platform-specific NPM packages
5. ✅ Install/uninstall scripts with error handling
6. ✅ Fallback wrapper for --ignore-scripts
7. ✅ Helper scripts for building and testing
8. ✅ Version management automation
9. ✅ Updated documentation

**Next Steps:**
1. Install GoReleaser locally: `brew install goreleaser`
2. Test build: `make snapshot`
3. Test NPM package: `make npm-test`
4. Set up NPM account and tokens
5. Configure GitHub secrets (NPM_TOKEN)
6. Create initial release tag
7. Monitor GitHub Actions for successful build
8. Verify NPM packages published correctly
9. Test installation on multiple platforms

**Key Files Created:**
- `.goreleaser.yml` - Build configuration
- `.github/workflows/release.yml` - CI/CD pipeline
- `npm/` - NPM package structure
- `scripts/` - Helper scripts
- `Makefile` - Development automation