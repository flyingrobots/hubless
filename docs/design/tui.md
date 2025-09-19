# Hubless TUI Design Specification

## Document Control
- Version: 0.1
- Last updated: 2025-09-18
- Owner: Developer Experience Team

## 1. Overview
The Hubless text user interface (TUI) delivers the primary interactive experience for browsing and managing work items. It must feel immediate, keyboard-centric, and familiar to developers who rely on Magit, Tig, or other terminal-first tools. This document outlines framework choices, view layouts, navigation flows, and stylistic guidelines.

## 2. Framework Stack
| Library | Purpose | Notes |
|---------|---------|-------|
| [Bubbletea](https://github.com/charmbracelet/bubbletea) | Elm-inspired TUI state machine | Provides update loop, message handling, and window resizing support |
| [Bubbles](https://github.com/charmbracelet/bubbles) | Reusable UI widgets (list, viewport, text input, progress) | Accelerates prototyping; customize delegates for domain data |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Styling primitives (colors, borders, spacing) | Use for consistent theming and status-specific highlights |
| [Glamour](https://github.com/charmbracelet/glamour) | Markdown rendering | Applied to issue descriptions and comments |
| [Huh](https://github.com/charmbracelet/huh) (optional) | Form builder for interactive dialogs | Useful for guided issue creation beyond `$EDITOR` |
| [Wish](https://github.com/charmbracelet/wish) (stretch) | SSH-hosted multi-user sessions | Enables collaborative board sessions in future phases |

## 3. Primary Views
### 3.1 Issue List View
- **Layout**: Split-pane; left column uses `bubbles/list` for issues, right column uses `bubbles/viewport` for a preview of the selected issue.
- **Displayed Fields**: Issue ID, title, status, priority indicator, assignee.
- **Interactions**:
  - `↑/↓` or `k/j` to change selection.
  - `enter` to open detail view.
  - `c` to create an issue via editor or inline form.
  - `/` to open the filter bar.

### 3.2 Issue Detail View
- **Layout**: Full-width viewport rendering Glamour-formatted markdown.
- **Content**: Entire event timeline (creation, edits, comments, status changes, assignments).
- **Actions**:
  - `c` add comment (opens `$EDITOR` or inline form).
  - `s` change status via quick-select menu.
  - `a` assign/unassign.
  - `b` return to list.

### 3.3 Kanban View
- **Layout**: Three default columns (Open, In Progress, Closed) rendered with Lipgloss borders. Additional columns supported once board configuration events exist.
- **Navigation**: Arrow keys or `h/l` to switch columns; `k/j` to move within a column.
- **Actions**: `space` or `enter` moves the selected issue to the adjacent column based on column ordering, emitting `board:moved` and `issue:status_changed` events.

### 3.4 Filter & Search View
- **Trigger**: `/` opens a `bubbles/textinput` component at the bottom of the screen.
- **Query Language**: `status:open assignee:me priority:high`. Saved filters stored in CLI config for quick access.
- **Persistence**: Named filters appear in a palette (e.g., `1`, `2`, `3` shortcuts).

### 3.5 Sync Progress Overlay
- **Trigger**: `hubless sync` or pressing `S` inside the TUI.
- **Display**: `bubbles/progress` bar with stages (fetch, apply, project) and counters for events pushed/pulled.
- **Feedback**: Success and error notifications rendered as temporary toast panels.

## 4. Interaction Model
- **Startup Flow**: Load catalog and most recent snapshots concurrently; show spinner until data available (<200 ms target).
- **Message Handling**: Each view implements Bubbletea `Model`, with a root `Model` coordinating child updates.
- **Keyboard Defaults**: Provide both Vim-style (`h/j/k/l`) and arrow key mappings. Avoid hidden multi-key chords.
- **Undo**: Not supported initially; users rely on Git history for rollbacks. Provide `git` command hints when destructive actions occur.

## 5. Visual Styling Guidelines
- **Color Palette**:
  - Status Open: blue (`#1E90FF` equivalent).
  - Status In Progress: amber (`#FFBE00`).
  - Status Closed: green (`#32CD32`).
  - Snapshots or derived data: muted gray italics.
- **Priority Indicators**: 🔥 (high), ● (medium), · (low) prefixed to issue titles.
- **Borders**: Use Lipgloss rounded borders sparingly to delineate columns without overwhelming text density.
- **Typography**: Monospaced fonts by default; rely on padding rather than heavy separators.

## 6. Performance Considerations
- Avoid recomputing derived views in the render loop; memoize catalog-derived summaries.
- Limit Glamour rendering to visible content; paginate long comment threads.
- Defer loading of non-critical resources (e.g., feed view) until requested to maintain startup performance.

## 7. Accessibility and Internationalization
- Provide configuration to disable emoji indicators in terminals lacking emoji support.
- Future enhancement: support alternate keymaps and color themes for accessibility.
- Markdown rendering must respect ANSI-stripping when copying to clipboard.

## 8. Testing Strategy
- Use Bubbletea’s `tea.Msg` simulation to unit test navigation and state transitions.
- Snapshot test Lipgloss output for key views to catch regressions in layout.
- Provide integration tests that replay fixture event streams and ensure UI renders expected summaries.

## 9. Open Questions
- Should issue creation default to `$EDITOR` or an inline Bubbletea form in MVP?
- How should multi-select operations (bulk status change) be exposed?
- What telemetry, if any, should the TUI emit for product analytics?

## 10. References
- `docs/TechSpec.md` for event model and storage contracts.
- `docs/reference/implementation-skeleton.md` for Bubbletea wiring examples.
