# Hubless Product Requirements Document

## Document Control

- Version: 0.1 (working draft)
- Last updated: 2025-09-18
- Authors: Hubless Core Team

## 1. Executive Summary

Hubless is a terminal-native work tracker that treats issues, pull requests, and boards as Git-native data. It extends a repository with append-only event streams and presents them through an ergonomic text user interface (TUI). Hubless targets teams that prefer Git workflows, need offline access, and require auditable change history without depending on SaaS task managers. The product must feel as fast and expressive as Magit while eliminating the friction of context switching to web dashboards.

## 2. Product Vision and Goals

- **Vision**: Make collaborative planning feel like a first-class Git primitive.
- **Primary Goal**: Provide conflict-free, offline-capable issue and PR tracking that syncs cleanly across Git remotes and optional GitHub projections.
- **Secondary Goals**:
  - Deliver a TUI that matches developer muscle memory (keyboard-first, Magit-inspired ergonomics).
  - Preserve a verifiable trail for every change without requiring a centralized service.
  - Allow gradual adoption: start with Git-only workflows, add GitHub synchronization when required.

## 3. Target Users and Personas

- **Hands-on Developers**: Live in the terminal, already comfortable with Git plumbing, and want issue tracking that keeps pace with code review workflows.
- **Tech Leads / Maintainers**: Need immediate visibility into status across multiple contributors without manual status reporting.
- **Tooling Enthusiasts**: Evaluate new workflows, expect scripting hooks, and will extend the tool.

## 4. Problem Statement

Traditional issue trackers fragment context between code and planning artifacts and require constant network connectivity. Teams maintaining long-lived repositories lack an auditable, offline-first issue tracker that integrates with Git operations. Hubless solves this by storing planning data directly in the repository while offering an efficient interface and optional projections into GitHub.

## 5. Scope

### 5.1 In Scope

- Append-only event streams for issues, boards, and pull requests stored under `refs/hubless/**`.
- Local-first command-line and TUI workflows to create, view, update, and sync work items.
- Snapshotting and catalog indexes that keep large backlogs responsive.
- Optional GitHub synchronization that treats GitHub as a projection of Hubless state.

### 5.2 Out of Scope (Initial Releases)

- Full parity with GitHub project boards or enterprise integrations (Jira, Linear, etc.).
- Web UI or mobile clients.
- Real-time multi-user collaboration beyond Git’s eventual consistency model.
- Automated analytics dashboards beyond basic activity feeds.

## 6. Product Principles

1. **Git is the source of truth**: All work artifacts originate as Git refs and commits.
2. **No merge conflicts**: Event sourcing and CRDT-friendly data structures avoid write contention.
3. **Offline-first**: The product works without network access; sync happens on demand.
4. **Transparency and auditability**: Every change has a durable author, timestamp, and payload.
5. **Extensible vocabulary**: Event types can evolve without breaking compatibility.

## 7. Use Cases and Functional Requirements

### 7.1 Core Use Cases

- View an overview of open work and drill into issue timelines.
- Create and edit issues using editors developers already use (`$EDITOR`).
- Update status, assignment, and comments from either CLI commands or the TUI.
- Organize work through Kanban-style views driven by board events.
- Synchronize changes between teammates through Git fetch/push.
- Optionally synchronize with GitHub issues, comments, and PRs.

### 7.2 Functional Requirements (MVP)

1. **List issues**: `hubless list` and TUI list view show title, ID, status, priority, assignee.
2. **View issue timeline**: `hubless view <id>` replays events, including comments and status changes.
3. **Create issue**: `hubless create` opens a template, commits an `issue:created` event.
4. **Edit status / assignment**: Commands and TUI actions append corresponding events.
5. **Comment**: `hubless comment <id>` appends an `issue:commented` event.
6. **Kanban navigation**: TUI board view allows drag-and-drop via keyboard, recording move events.
7. **Git sync**: `hubless sync` performs bidirectional set-union with remotes.
8. **GitHub projection** (post-MVP toggle): translate events to GitHub issues/comments and back while preserving event IDs.

### 7.3 Stretch Goals (Phase 2+)

- Promote issue to PR (`hubless pr <branch>`) with event linkage.
- Export state for external systems (e.g., Jira) via CLI commands.
- Provide a feed of recent changes for status reporting.
- Expose an LSP endpoint for IDE integrations.

## 8. Non-Functional Requirements

- **Performance**: Listing 10k issues with 100k total events in under 200 ms via catalog indexing. Viewing a single issue under 100 ms using snapshot + tail replay.
- **Reliability**: Commands are idempotent and safe to retry. Sync detects and de-duplicates events using stable IDs.
- **Security**: Honor repository permissions; reuse existing Git or `gh` authentication flows for GitHub API use.
- **Usability**: Keyboard-first UX, minimal modal dialogs, immediate visual feedback in the TUI.
- **Portability**: Works on macOS, Linux, and WSL out of the box.

## 9. Release Strategy

| Phase | Objectives | Key Deliverables |
|-------|-------------|------------------|
| Phase 0.5 – CLI Proof | Validate event model; basic CLI create/list/view | Event schema, append-only refs, CLI surface |
| Phase 1 – MVP | Ship Magit-speed TUI, snapshots, Git-only sync | Bubbletea-based TUI, catalog index, kanban view |
| Phase 1.5 – Enhancements | Filters, activity feed, saved views | Query language, feed chain, UX polish |
| Phase 2 – GitHub Sync | Round-trip GitHub integration, PR support | Sync adapter, PR events, mapping metadata |
| Phase 3 – IDE Integration | Surface Hubless data in editors | LSP service, integration guides |

## 10. Success Metrics

- Team of 3–10 developers can replace GitHub Issues within one sprint.
- Offline work for two weeks reconciles without conflicts on sync.
- TUI satisfaction scores exceed web-based alternatives in qualitative interviews.
- 90% of status updates originate from Hubless commands/TUI rather than GitHub UI.

## 11. Dependencies and Assumptions

- Developers have Git 2.30+ and can install a Go-based CLI.
- Repositories may already contain Charmbracelet-based tooling; Hubless must coexist without conflicting key bindings.
- GitHub API access tokens available when sync is enabled.
- Repository maintainers permit additional refs under `refs/hubless/**`.

## 12. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| GitHub API mismatch with event model | Data loss or drift | Treat GitHub as projection; store canonical event IDs; run diff checks before publishing |
| Event volume growth | Slow replay and sync | Introduce periodic snapshots and aggregated catalog commits |
| User confusion about boards vs. statuses | Adoption friction | Provide opinionated defaults and education in onboarding tutorials |
| Authentication complexity | Failed sync operations | Reuse `gh` CLI auth and document token scopes |

## 13. Open Questions

- Do we require migration tooling for existing GitHub issues when onboarding a repository?
- What is the default cadence for snapshots (per N events vs. time-based)?
- Should board definitions support custom columns beyond To Do / In Progress / Done in MVP?
- How much configuration should live in repo-level manifests versus CLI config files?

## 14. Related Documents

- [docs/TechSpec.md](docs/TechSpec.md)
- `docs/design/tui.md`
- `docs/reference/implementation-skeleton.md`
