# ADR-001: Project Structure

## Status
Accepted

## Context
The Claude Gate project needs a clear, maintainable structure that balances Go community conventions with practical needs. The golang-standards/project-layout is popular but controversial, and we need to decide what structure best serves our project.

## Decision
We will use a simplified structure based on Go conventions:
- `/cmd/claude-gate/` for the CLI entry point
- `/internal/` for private packages (auth, proxy, config)
- No `/pkg/` directory unless we develop public APIs
- `/docs/` for documentation
- `/scripts/` for build automation
- `/.claude/` for project management

We will NOT use:
- `/pkg/` - No public API currently planned
- `/api/` - No API definitions needed
- `/web/` - No web assets
- Deep nesting - Keep structure flat

## Consequences

### Positive
- Simple and easy to navigate
- Follows Go compiler conventions (internal)
- Clear separation of concerns
- Easy for new contributors to understand
- Avoids over-engineering

### Negative
- May need restructuring if public API added
- Less "standard" than full project-layout
- Some Go developers expect /pkg/

## Alternatives Considered

1. **Full golang-standards/project-layout**
   - Pro: Familiar to many Go developers
   - Con: Overkill for our current needs
   - Con: Includes many unnecessary directories

2. **Flat structure (everything in root)**
   - Pro: Maximum simplicity
   - Con: No privacy enforcement
   - Con: Harder to organize as project grows

3. **Domain-driven structure**
   - Pro: Business logic focused
   - Con: Less conventional in Go
   - Con: Overhead for small project

## References
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [Go Project Structure Best Practices](https://tutorialedge.net/golang/go-project-structure-best-practices/)
- [Organizing a Go module](https://go.dev/doc/modules/layout)