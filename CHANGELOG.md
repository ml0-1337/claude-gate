# Changelog

All notable changes to Claude Gate will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Interactive server monitoring dashboard with real-time metrics
- OAuth flow with interactive TUI and browser automation
- Bubble Tea UI foundation for enhanced CLI experience
- Comprehensive documentation structure

### Changed
- Reorganized documentation into logical categories
- Standardized port configuration to 8080
- Improved authentication flow with better error handling

### Fixed
- Dashboard requests/sec metric showing 0.0
- Various documentation inconsistencies

## [0.1.0] - 2024-01-01

### Added
- Initial release
- OAuth/PKCE authentication with Claude Pro/Max accounts
- HTTP proxy server for Anthropic API
- Cross-platform support (macOS, Linux, Windows)
- NPM package distribution
- Basic CLI commands (auth, start, stop, status)
- Environment variable configuration
- Secure token storage

[Unreleased]: https://github.com/anthropics/claude-gate/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/anthropics/claude-gate/releases/tag/v0.1.0