# Hubless

Hubless is a terminal-native, Git-backed work tracker. It treats issues, pull requests, and boards as append-only event streams under `refs/hubless/**`, then presents them through a Charmbracelet-powered TUI and CLI. This repository houses the specs, tooling, and implementation that turn Git repositories into fully auditable planning systems.

## Getting Started

> **Status:** Early development. Specs are in place; implementation is in progress.

### Prerequisites

- Go 1.22+
- Git 2.30+
- Optional: [gh](https://github.com/cli/cli) for GitHub integration experiments

### Clone

```bash
git clone https://github.com/flyingrobots/hubless.git
cd hubless
```

### Build the utilities

The Go module is initialized but the primary CLI is still under construction. A helper binary for progress updates exists today:

```bash
go build ./cmd/update-progress
```

### Run the progress updater

The legacy Python script has been replaced with the Go implementation (spec documented in `docs/reference/update-progress-algorithm.md`). Point the tool at your `git-mind` checkout once the Go port is finished.

```bash
./update-progress --root ../git-mind
```

## Project Docs

- `docs/PRD.md` – Product requirements and roadmap.
- `docs/TechSpec.md` – Architecture, data model, sync contracts.
- `docs/design/tui.md` – Bubbletea TUI views, interactions, styling.
- `docs/reference/implementation-skeleton.md` – Hexagonal layout and scaffolding.
- `docs/reference/update-progress-algorithm.md` – Transcription of the ledger updater logic.
- `AGENTS.md` – Workflow rules, coding standards, collaboration notes.
- `@hubless/` – Structured planning data (tasks, stories, features, milestones schemas).

## Development Principles

- Git is the source of truth; no central server required.
- Conflict-free, append-only event streams for issues, boards, and PRs.
- Hexagonal architecture with Go application services, Git adapters, and Charmbracelet UI layers.
- CLI command surface will use Charmbracelet Fang/Cobra to keep styling consistent with the TUI.

## Contributing

See `CONTRIBUTING.md` for task workflow, branching rules, and code quality expectations.

## License

This project is licensed under the MIT License – see `LICENSE` for details.
