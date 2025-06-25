.PHONY: build test test-unit test-integration test-e2e clean install release snapshot npm-test test-all test-docker test-edge help

# Default target
help:
	@echo "Claude Gate - Available targets:"
	@echo "  make build         - Build for current platform"
	@echo "  make test          - Run unit tests with coverage"
	@echo "  make test-unit     - Run unit tests only (short mode)"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-e2e      - Run end-to-end tests"
	@echo "  make test-all      - Run comprehensive test suite"
	@echo "  make snapshot      - Build snapshot release (all platforms)"
	@echo "  make npm-test      - Test NPM package locally"
	@echo "  make test-docker   - Test in Docker containers"
	@echo "  make test-edge     - Test edge cases"
	@echo "  make install       - Install locally"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make release       - Create a new release (requires version)"

# Build for current platform
build:
	go build -ldflags="-s -w" -o claude-gate ./cmd/claude-gate

# Run unit tests with coverage
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run unit tests only (short mode)
test-unit:
	go test -short -v ./...

# Run integration tests
test-integration:
	go test -tags=integration -v ./internal/test/integration/...

# Run end-to-end tests
test-e2e: build
	go test -tags=e2e -v ./internal/test/e2e/...

# Build snapshot release with GoReleaser
snapshot:
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "Error: goreleaser is not installed. Install with: brew install goreleaser"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean --skip=publish
	@echo ""
	@echo "Snapshot build complete! Binaries in ./dist/"

# Test NPM package locally
npm-test: snapshot
	./scripts/test-npm-local.sh

# Run comprehensive test suite
test-all: test-unit test-integration test-e2e
	./scripts/test-all.sh

# Test in Docker containers
test-docker: snapshot
	./scripts/test-docker.sh

# Test edge cases
test-edge:
	./scripts/test-edge-cases.sh

# Install locally
install: build
	mkdir -p ~/bin
	cp claude-gate ~/bin/
	@echo "Installed to ~/bin/claude-gate"
	@echo "Make sure ~/bin is in your PATH"

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f claude-gate
	rm -f coverage.out coverage.html
	rm -rf npm/node_modules/
	rm -f npm/*.tgz

# Create a new release (requires VERSION parameter)
release:
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=0.1.0)
endif
	@echo "Creating release v$(VERSION)..."
	@echo "This will:"
	@echo "  1. Update version in code"
	@echo "  2. Commit changes"
	@echo "  3. Create tag v$(VERSION)"
	@echo "  4. Push to GitHub"
	@echo ""
	@read -p "Continue? [y/N] " -n 1 -r; \
	echo ""; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		sed -i.bak 's/version = ".*"/version = "$(VERSION)"/' cmd/claude-gate/main.go && rm cmd/claude-gate/main.go.bak; \
		sed -i.bak 's/"version": ".*"/"version": "$(VERSION)"/' npm/package.json && rm npm/package.json.bak; \
		git add -A; \
		git commit -m "chore: release v$(VERSION)"; \
		git tag -a v$(VERSION) -m "Release v$(VERSION)"; \
		echo ""; \
		echo "Release prepared! To publish:"; \
		echo "  git push origin main"; \
		echo "  git push origin v$(VERSION)"; \
	else \
		echo "Release cancelled"; \
	fi