# Contributing to Hubless

Thanks for building Hubless! This document summarizes how we collaborate. For the full agent handbook see `AGENTS.md`.

## Prerequisites

- Go 1.26.3
- Git 2.30+
- Familiarity with Charmbracelet (Bubbletea, Lipgloss, Bubbles) is helpful.

## Task Workflow

1. Add or pick up a task in `@hubless/issues/tasks/*.json` (the table in `@hubless/issues/generated/tasks.md` is generated) and set status to `PLANNED`.
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

## Documentation Automation

- Structured data under `@hubless/` now feeds Markdown via reusable snippets in `docs/components/`.
- Install or clone [`markdown-transclusion`](https://github.com/flyingrobots/markdown-transclusion) (Node ≥20). Set `MARKDOWN_TRANSCLUSION_BIN` to the executable (`markdown-transclusion` if installed globally, or `node`) and `MARKDOWN_TRANSCLUSION_SCRIPT` to the CLI script path when using a local clone (e.g., `/path/to/markdown-transclusion/dist/cli.js`). Use `MARKDOWN_TRANSCLUSION_ARGS` only for additional CLI flags.
- Run `make docs` (or `./scripts/render-docs.sh`) after editing JSON records or templates. This regenerates shared snippets and rewrites `@hubless/roadmap/generated/README.md` and `@hubless/issues/generated/tasks.md` from their templates.
- Run `make docs-test` to execute generator unit tests and ensure snippets format as expected.
- Run `make docs-verify` to confirm all generated Markdown is fully transcluded (no `![[…]]` placeholders).
- For custom dependency graph styling, pass `--graph-direction`, `--graph-clusters`, or `--graph-palette` to `cmd/docs-components` (see `README.md`).
- `CHANGELOG.md` is generated from `CHANGELOG.template.md` and `docs/reference/release-notes.*`; edit the template or JSON, not the generated file.
- Palette overrides live in `docs/reference/palettes.json` (validated by `docs/reference/palettes.schema.json`); point `--palette-file` elsewhere if you keep custom palettes in another location.
- Install [`golangci-lint`](https://golangci-lint.run/) locally so `make lint` and the pre-commit hook can execute successfully. Recommended install command:

  ```bash
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
  ```

## Tooling & Hooks

- Run `make fmt`, `make lint`, and `make test` before opening a PR. CI expects them to pass.
- Install local git hooks via `make hooks` (wraps `scripts/install-git-hooks.sh`) so pre-commit checks run automatically.
- `.golangci.yml` houses lint configuration; feel free to propose tweaks but keep the suite running clean.
- `.editorconfig` defines base formatting (tabs for Makefiles, spaces elsewhere).

## Commit & PR Guidance

- Reference task IDs in commit messages (`TASK-123: implement catalog writer`).
- Small, focused commits preferred.
- Ensure docs and changelog entries are updated where relevant.
- PR description should include testing evidence (manual commands, unit test output) and links to updated docs.

## Communication

- Inline documentation lives in `docs/`. Update related files (PRD, TechSpec, design docs) as features evolve.
- Capture architectural decisions either in commit messages or lightweight ADRs under `docs/` if needed.

Welcome to the team—let’s make Git-native planning delightful.
