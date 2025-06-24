# Contributing to Claude Gate

Thank you for your interest in contributing to Claude Gate! This guide will help you get started with development and ensure your contributions align with our project standards.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

- Go 1.22 or later
- Node.js 18 or later (for NPM package testing)
- Git
- Make
- GoReleaser (for release builds)

### Initial Setup

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

## Development Setup

### Environment Setup

1. **Go Environment**
   ```bash
   # Ensure Go is properly installed
   go version
   
   # Set up Go workspace (if needed)
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```

2. **Development Tools**
   ```bash
   # Install GoReleaser
   brew install goreleaser/tap/goreleaser
   
   # Install linting tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **Editor Configuration**
   - VS Code: Install Go extension
   - GoLand: Built-in support
   - Vim/Neovim: Use vim-go or nvim-lspconfig

## Project Structure

```
claude-gate/
├── cmd/claude-gate/        # CLI entry point
│   └── main.go            # Main application
├── internal/              # Private packages
│   ├── auth/             # OAuth authentication
│   ├── config/           # Configuration management
│   └── proxy/            # Core proxy logic
├── docs/                 # Documentation
│   ├── architecture/     # System design docs
│   └── decisions/        # ADRs
├── scripts/              # Build and test scripts
├── npm/                  # NPM distribution
└── .claude/              # Project management
    ├── todos/            # Active tasks
    └── knowledge/        # Research docs
```

## Development Workflow

### 1. Creating a New Feature

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make changes and test
make test

# Build and test locally
make build && ./claude-gate test
```

### 2. Task Management

We use `.claude/todos/` for task tracking:

```bash
# Create a task file
touch .claude/todos/your-task.md

# Use the template format (see existing todos)
```

### 3. Making Changes

1. **Write tests first** (TDD approach)
2. **Implement the feature**
3. **Update documentation**
4. **Run full test suite**

### 4. Code Organization

- Put reusable code in appropriate `internal/` packages
- Keep `main.go` minimal - delegate to packages
- Use meaningful package names
- Follow single responsibility principle

## Testing

### Running Tests

```bash
# Run all tests with coverage
make test

# Run specific package tests
go test ./internal/auth/...

# Run with race detection
go test -race ./...

# Run integration tests
make test-all
```

### Writing Tests

1. **Unit Tests**
   - Test individual functions
   - Use table-driven tests
   - Mock external dependencies
   - Achieve >80% coverage

2. **Integration Tests**
   - Test component interactions
   - Use real implementations
   - Test error scenarios

3. **Test Organization**
   ```go
   func TestFunctionName(t *testing.T) {
       t.Run("successful case", func(t *testing.T) {
           // Test implementation
       })
       
       t.Run("error case", func(t *testing.T) {
           // Test implementation
       })
   }
   ```

## Code Style

### Go Code Guidelines

1. **Formatting**
   - Use `gofmt` (automatic with most editors)
   - Keep lines under 120 characters
   - Use meaningful variable names

2. **Error Handling**
   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to process request: %w", err)
   }
   
   // Bad
   if err != nil {
       return err
   }
   ```

3. **Comments**
   - Document all exported functions
   - Explain "why", not "what"
   - Keep comments up-to-date

4. **Package Design**
   - Small, focused packages
   - Clear interfaces
   - Minimize dependencies

### Commit Messages

Follow conventional commits:

```
type(scope): subject

body

footer
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Testing
- `refactor`: Code refactoring
- `chore`: Maintenance

Example:
```
feat(auth): add keychain support for token storage

Implement OS-specific keychain integration for secure token storage.
Supports macOS Keychain, Linux Secret Service, and Windows Credential Manager.

Closes #123
```

## Submitting Changes

### Pull Request Process

1. **Before Submitting**
   - [ ] Tests pass (`make test`)
   - [ ] Code is formatted (`gofmt`)
   - [ ] Documentation updated
   - [ ] Commit messages follow convention
   - [ ] Branch is up-to-date with main

2. **PR Description Template**
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Manual testing completed
   
   ## Checklist
   - [ ] My code follows the project style
   - [ ] I have added tests for my changes
   - [ ] Documentation has been updated
   ```

3. **Review Process**
   - Maintainers will review within 48 hours
   - Address feedback promptly
   - Keep PR focused and small
   - Be patient and respectful

### Code Review Guidelines

**For Authors:**
- Keep PRs small and focused
- Respond to feedback constructively
- Update based on suggestions
- Test thoroughly before requesting review

**For Reviewers:**
- Be constructive and specific
- Suggest improvements, don't just criticize
- Approve when ready, not perfect
- Consider the bigger picture

## Release Process

### Version Numbering

We use semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

### Creating a Release

```bash
# Update version
make release VERSION=1.2.3

# This will:
# 1. Update version in all files
# 2. Create git commit
# 3. Create git tag
# 4. Show push instructions
```

### Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG updated
- [ ] Version bumped
- [ ] Tag created
- [ ] Release notes written

## Getting Help

### Resources

- [Project Documentation](../README.md)
- [Architecture Overview](../architecture/overview.md)
- [Security Policy](../architecture/security.md)
- [Issue Tracker](https://github.com/ml0-1337/claude-gate/issues)

### Communication

- **Issues**: Bug reports and feature requests
- **Discussions**: General questions and ideas
- **Pull Requests**: Code contributions

### First-Time Contributors

Look for issues labeled:
- `good first issue`
- `help wanted`
- `documentation`

## Thank You!

Your contributions make Claude Gate better for everyone. We appreciate your time and effort in improving this project!

---

[← Guides](../README.md#guides) | [Documentation Home](../README.md) | [Development →](./development.md)