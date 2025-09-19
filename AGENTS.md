## Workflow

- Tasks live in @hubless/issues/tasks.md
- Try to associate all work with a task id
- Maintain task dependencies so we have an accurate DAG; update `@hubless/issues/tasks/` whenever prerequisites change. This DAG will drive the rolling frontier worker pool.
  - Tasks are defined as JSON files under `@hubless/issues/tasks/` matching the schema in `@hubless/schema/task.schema.json`; IDs follow `{project}/{milestone}/{type}/{number}` (e.g., `hubless/m0/task/0001`).
- Follow the Task lifecycle:

1. Task added to @hubless/issues/tasks.md; status = "PLANNED"
2. Start task? status = "STARTED"
3. Task blocked? status = "BLOCKED"
4. Task finished? status = "DONE"
5. Once status = "DONE", Tasks need the following badges: (i) Tested (ii) Documented (iii) Shipped; only then shall we remove them from the tasks.md file, and add it to the "tasks.archive.md" file

- NEVER GIT AMEND; just make a new commit.
- NEVER REBASE; just git merge. Embrace the messy history–the truth shall set you free.
- NEVER EVER FORCE PUSH!!! If you feel like you must halt and seek permission from the user.

## Code Quality

- SRP
- One file per entity (class, struct, enum, object, whatever)
- Test-double friendliness
- Dependency Injection
- Hexagonal Architecture
- DX and UX are paramount

## Project Overview

- Hubless is a terminal-native, Git-backed work tracker; see docs/PRD.md for the product mandate and milestones.
- Architecture and data model live in docs/TechSpec.md, covering refs/hubless/** namespaces, event vocabulary, snapshots, catalog, and sync.
- TUI experience (Bubbletea, Bubbles, Lipgloss, Glamour, optional Huh/Wish) is documented in docs/design/tui.md.
- Implementation scaffolding, including hexagonal layout and Go module expectations, sits in docs/reference/implementation-skeleton.md.
- Structured planning data lives under `@hubless/` (schemas, milestones, features, stories, tasks).
- Progress ledger algorithm is captured in docs/reference/update-progress-algorithm.md for parity with the retired Python script.
- Go module: github.com/flyingrobots/hubless; current CLI code lives under cmd/update-progress (to be replaced by full Fang-powered hubless CLI).

## Technical Stack & Practices

- Language: Go 1.22; CLI stack targets Charmbracelet ecosystem (Fang/Cobra for commands, Bubbletea suite for TUI).
- Persistence: Git plumbing (mktree, commit-tree, update-ref) writing refs/hubless/**; catalog + snapshots for fast reads.
- Sync roadmap: Git remotes first, GitHub projection later with stable event IDs and refs/hubless/meta/github-map mapping.
- Testing strategy: unit-test domain replay, adapter plumbing, and Bubbletea models (see docs/design/tui.md and reference skeleton).
- Packaging/ops: plan for hubless doctor, structured logs via HUBLESS_LOG, and repo-friendly refspec configuration.

## Collaboration Notes

- Default branch currently main; initial commit message: "Document hubless specs and add progress updater".
- Before implementing new features, update corresponding spec/design docs; keep docs as living sources of truth.
- Prefer adding new CLI/TUI commands via thin adapters that call internal/application services to preserve hexagonal boundaries.
- When extending progress tooling, consult docs/reference/update-progress-algorithm.md to keep output deterministic.
