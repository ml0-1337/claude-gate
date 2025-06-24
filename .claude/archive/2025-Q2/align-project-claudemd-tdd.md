---
todo_id: align-project-claudemd-tdd
started: 2025-06-24 19:56:22
completed: 2025-06-24 20:00:51
status: completed
priority: high
---

# Task: Update project CLAUDE.md to explicitly reinforce TDD workflow from system CLAUDE.md

## Findings & Research

### Issue Identified
- System CLAUDE.md enforces strict TDD with mandatory Red-Green-Refactor-Commit cycle
- Project CLAUDE.md mentions TDD but lacks the rigor and mandatory nature
- This creates ambiguity about whether TDD is required or optional

### System CLAUDE.md Requirements
- Law 2: Plan-First Workflow
- Law 4: Test-First for Features & Fixes (MANDATORY)
- RGRC Cycle: Red → Green → Refactor → Commit
- WebSearch mandatory for technical decisions
- Todo creation for every task

### Current Project CLAUDE.md Gaps
- Only mentions "Always write tests before implementing features" once
- Testing strategy section is descriptive, not prescriptive
- No explicit reference to system CLAUDE.md requirements
- Missing Go-specific TDD patterns and examples

## Test Strategy

Since this is a documentation update, traditional tests don't apply. However, we'll ensure:
- The updated documentation is clear and unambiguous
- TDD requirements are prominently featured
- Go-specific examples are accurate
- No conflicts with system CLAUDE.md

## Implementation Plan

1. Add new "TDD Requirements" section after "Testing Strategy"
2. Update "Development Workflow" to make TDD mandatory
3. Add Go-specific TDD examples
4. Include reference to system CLAUDE.md laws
5. Add quick reference for RGRC cycle

## Checklist

- [✓] Create todo and document research
- [✓] Update todo status to in_progress
- [✓] Read current project CLAUDE.md
- [✓] Draft TDD sections to add
- [✓] Update project CLAUDE.md
- [✓] Verify no conflicts with system requirements
- [✓] Commit changes
- [✓] Complete and archive todo

## Working Scratchpad

### Requirements
- Make TDD mandatory and explicit in project CLAUDE.md
- Reference system CLAUDE.md laws
- Add Go-specific TDD guidance
- Include RGRC cycle reference

### Approach
1. Insert new sections that complement existing content
2. Keep project-specific information while adding TDD rigor
3. Ensure consistency with system CLAUDE.md

### Notes
- User wants to maintain TDD rigor across all projects
- Project CLAUDE.md should reinforce, not dilute, system standards
- Go projects have specific testing patterns (table-driven tests, testify)