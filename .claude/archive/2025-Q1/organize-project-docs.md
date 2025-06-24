---
todo_id: organize-project-docs
started: 2025-01-20 09:45:00
completed: 2025-01-20 10:15:00
status: completed
priority: high
---

# Task: Reorganize project documentation structure

## Findings & Research

### Current Documentation Analysis

**Root Level Documentation (5 files)**:
- **README.md**: Main project overview, features, installation methods, quick start guide
- **CONTRIBUTING.md**: Comprehensive contributor guide (345 lines) covering setup, workflow, testing, code style
- **TEST_GUIDE.md**: Testing guide for cross-platform builds and NPM distribution
- **TEST_RESULTS.md**: Test execution results summary dated 2025-06-24
- **TROUBLESHOOTING.md**: Common issues and solutions for installation and runtime

**docs/ Directory (4 files)**:
- **docs/SECURITY.md**: Security policy with OAuth model, token security, threat model
- **docs/api/API.md**: API documentation covering endpoints, transformations, client examples
- **docs/architecture/ARCHITECTURE.md**: System architecture, components, data flow
- **docs/decisions/001-project-structure.md**: ADR explaining project structure choices

**npm/ Directory (1 file)**:
- **npm/README.md**: Simplified README for NPM package distribution

### Major Issues Identified

1. **Duplicate Content**:
   - Root README.md and npm/README.md have overlapping content
   - OAuth/PKCE explanation appears in README, SECURITY, and ARCHITECTURE
   - Installation instructions in README, CONTRIBUTING, and npm/README
   - CLI commands listed in multiple places with slight variations

2. **Inconsistent Documentation Placement**:
   - SECURITY.md is in docs/ while other guides (TEST_GUIDE, TROUBLESHOOTING) are at root
   - No clear pattern for what belongs in docs/ vs root
   - TEST_GUIDE.md and TEST_RESULTS.md could be in a tests/ documentation folder

3. **Missing Documentation**:
   - No CHANGELOG.md despite references in CONTRIBUTING.md
   - No dedicated API reference for the CLI commands
   - No development/debugging guide
   - Missing `.claude/` directory mentioned in multiple docs

4. **Inconsistencies**:
   - Port numbers: README shows port 8000, API.md shows port 9988
   - GitHub URLs use "yourusername" placeholder
   - TEST_RESULTS.md shows future date (2025-06-24)
   - Different Go version requirements mentioned

## Test Strategy

- **Test Framework**: Manual verification (documentation reorganization)
- **Test Types**: Structure validation, link checking, content consistency
- **Coverage Target**: All documentation files moved and properly linked
- **Edge Cases**: 
  - Ensure no broken links after reorganization
  - Verify all code examples still work
  - Check that NPM package docs remain functional

## Test Cases

```bash
# Test 1: Verify all files moved to correct locations
# Input: find . -name "*.md" | grep -E "(README|CONTRIBUTING|TEST_|TROUBLESHOOTING|SECURITY|API|ARCHITECTURE)"
# Expected: All files in new structure under docs/

# Test 2: Check for broken internal links
# Input: grep -r "\[.*\](" docs/ | grep -E "\.md|\.MD"
# Expected: All markdown links resolve correctly

# Test 3: Verify no duplicate content
# Input: Compare key sections across files
# Expected: Single source of truth for each topic

# Test 4: Port consistency check
# Input: grep -r ":[0-9]\{4\}" docs/ README.md
# Expected: Consistent port number (8080) across all files
```

## Maintainability Analysis

- **Readability**: [9/10] Clear directory structure with logical grouping
- **Complexity**: Simple flat structure within categories, easy to navigate
- **Modularity**: Each topic in its own file, easy to update independently
- **Testability**: Documentation can be validated with link checkers and consistency tools
- **Trade-offs**: More files to maintain, but better organization outweighs this

## Implementation Plan

### Phase 1: Create New Directory Structure
```
docs/
├── README.md                # Documentation index
├── getting-started/
│   ├── installation.md
│   ├── quick-start.md
│   └── configuration.md
├── guides/
│   ├── contributing.md
│   ├── development.md
│   └── troubleshooting.md
├── testing/
│   ├── test-guide.md
│   └── test-results/
├── reference/
│   ├── api.md
│   ├── cli.md
│   └── configuration.md
├── architecture/
│   ├── overview.md
│   ├── security.md
│   └── decisions/
└── deployment/
    └── npm-package.md
```

### Phase 2: Content Migration and Consolidation

1. **Create directory structure**
2. **Extract and consolidate installation content**:
   - From README.md → getting-started/installation.md
   - From CONTRIBUTING.md → getting-started/installation.md
   - From npm/README.md → deployment/npm-package.md

3. **Move existing files**:
   - CONTRIBUTING.md → docs/guides/contributing.md
   - TROUBLESHOOTING.md → docs/guides/troubleshooting.md
   - TEST_GUIDE.md → docs/testing/test-guide.md
   - TEST_RESULTS.md → docs/testing/test-results/2025-01-20.md
   - docs/SECURITY.md → docs/architecture/security.md
   - docs/api/API.md → docs/reference/api.md
   - docs/architecture/ARCHITECTURE.md → docs/architecture/overview.md
   - docs/decisions/* → docs/architecture/decisions/

4. **Create new files**:
   - CHANGELOG.md (root)
   - docs/README.md (documentation index)
   - docs/getting-started/quick-start.md
   - docs/getting-started/configuration.md
   - docs/guides/development.md
   - docs/reference/cli.md

### Phase 3: Fix Inconsistencies

1. **Standardize port numbers**: Use 8080 throughout
2. **Fix placeholder URLs**: Replace "yourusername" with "anthropics" or appropriate
3. **Update dates**: Fix future date in test results
4. **Align version requirements**: Go 1.20+ consistently

### Phase 4: Improve Navigation

1. **Add documentation index** (docs/README.md)
2. **Add navigation headers** to each document
3. **Create consistent ToC** for longer documents
4. **Add cross-references** between related documents

## Checklist

- [ ] Create .claude directory if missing
- [ ] Create new directory structure under docs/
- [ ] Create docs/README.md as documentation hub
- [ ] Extract installation content from multiple sources
- [ ] Move CONTRIBUTING.md to docs/guides/
- [ ] Move TROUBLESHOOTING.md to docs/guides/
- [ ] Move TEST_GUIDE.md to docs/testing/
- [ ] Archive TEST_RESULTS.md to docs/testing/test-results/
- [ ] Move SECURITY.md to docs/architecture/
- [ ] Move API.md to docs/reference/
- [ ] Move ARCHITECTURE.md to docs/architecture/overview.md
- [ ] Move decisions to docs/architecture/decisions/
- [ ] Create CHANGELOG.md at root
- [ ] Create getting-started/quick-start.md
- [ ] Create getting-started/configuration.md
- [ ] Create guides/development.md
- [ ] Create reference/cli.md
- [ ] Create deployment/npm-package.md
- [ ] Fix all port number inconsistencies (→ 8080)
- [ ] Fix all placeholder URLs
- [ ] Update all internal links
- [ ] Add navigation to all documents
- [ ] Update root README.md (simplified)
- [ ] Update npm/README.md (minimal, NPM-specific)
- [ ] Verify no broken links
- [ ] Test documentation build/preview if applicable

## Working Scratchpad

### Requirements
- Organize scattered documentation into logical structure
- Eliminate duplicate content
- Fix inconsistencies (ports, URLs, versions)
- Add missing documentation
- Improve navigation and cross-linking

### Approach
1. Create comprehensive directory structure first
2. Move files systematically, updating links as we go
3. Consolidate duplicate content into single sources
4. Create missing documentation files
5. Add navigation and improve discoverability

### Code
N/A - Documentation reorganization task

### Notes
- Need to decide on standard port (8080 seems reasonable)
- Should preserve git history when moving files (use git mv)
- Consider adding a documentation linter/link checker to CI
- May need to update any build scripts that reference doc locations

### Commands & Output

```bash
# Commands will be logged here as executed
```

## Progress Tracking

**Overall Progress**: [24/24] tasks completed ✓

**Phase 1 - Structure**: [2/2] ✓
- [x] Create .claude directory
- [x] Create docs/ subdirectories

**Phase 2 - Migration**: [13/13] ✓
- [x] Create docs/README.md as documentation hub
- [x] Extract installation content from multiple sources
- [x] Move CONTRIBUTING.md to docs/guides/
- [x] Move TROUBLESHOOTING.md to docs/guides/
- [x] Move TEST_GUIDE.md to docs/testing/
- [x] Archive TEST_RESULTS.md to docs/testing/test-results/
- [x] Move SECURITY.md to docs/architecture/
- [x] Move API.md to docs/reference/
- [x] Move ARCHITECTURE.md to docs/architecture/overview.md
- [x] Move decisions to docs/architecture/decisions/
- [x] Update root README.md (simplified)
- [x] Update npm/README.md (minimal, NPM-specific)
- [x] Remove empty directories

**Phase 3 - New Content**: [6/6] ✓
- [x] Create CHANGELOG.md at root
- [x] Create getting-started/installation.md
- [x] Create getting-started/quick-start.md
- [x] Create getting-started/configuration.md
- [x] Create guides/development.md
- [x] Create reference/cli.md
- [x] Create deployment/npm-package.md

**Phase 4 - Polish**: [3/3] ✓
- [x] Fix inconsistencies (ports → 8080, URLs → anthropics)
- [x] Add navigation headers to migrated documents
- [x] Verify all links work correctly (fixed broken links, created missing files)