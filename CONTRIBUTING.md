# Contributing to Hubless

Thanks for building Hubless! This document summarizes how we collaborate. For the full agent handbook see `AGENTS.md`.

## Prerequisites
- Go 1.22+
- Git 2.30+
- Familiarity with Charmbracelet (Bubbletea, Lipgloss, Bubbles) is helpful.

## Task Workflow
1. Add or pick up a task in `@hubless/issues/tasks.md` and set status to `PLANNED`.
2. When you begin work, flip status to `STARTED` and link the task ID in your branch, PR, and commit messages.
3. Blocked tasks move to `BLOCKED` with context on what’s needed.
4. When finished, mark as `DONE` and ensure it carries the badges **Tested**, **Documented**, **Shipped**.
5. After verifying badges, move the entry to `tasks.archive.md`.

## Git Practices
- **Never amend**: make new commits for follow-up fixes.
- **Never rebase**: merge instead; we keep history messy but truthful.
- **Never force push**: if you think you need to, stop and consult the team.

## Code Quality
- Single Responsibility Principle.
- One file per entity (struct, interface, enum, etc.).
- Test-double friendly design (interfaces over concretions).
- Dependency Injection everywhere practical.
- Hexagonal architecture: inbound adapters (CLI/TUI) → application services → ports → outbound adapters (Git, GitHub).
- Developer and user experience are top priorities.

## Development Flow
1. Read `docs/PRD.md` and `docs/TechSpec.md` to understand the current goals.
2. Update specs/design docs first when changing direction; treat them as living documents.
3. Build features through application services so both CLI and TUI can reuse logic.
4. Write tests alongside features. Unit test domain logic, adapters, and Bubbletea models.
5. Run `go fmt`, `go vet`, and any configured linters before opening a PR.

## Commit & PR Guidance
- Reference task IDs in commit messages (`TASK-123: implement catalog writer`).
- Small, focused commits preferred.
- Ensure docs and changelog entries are updated where relevant.
- PR description should include testing evidence (manual commands, unit test output) and links to updated docs.

## Communication
- Inline documentation lives in `docs/`. Update related files (PRD, TechSpec, design docs) as features evolve.
- Capture architectural decisions either in commit messages or lightweight ADRs under `docs/` if needed.

Welcome to the team—let’s make Git-native planning delightful.
