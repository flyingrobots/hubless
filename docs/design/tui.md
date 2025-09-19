# Hubless TUI Design Specification

## Document Control

- Version: 0.2
- Last updated: 2025-09-19
- Owner: Developer Experience Team

## 1. Overview

The Hubless text user interface (TUI) delivers the primary interactive experience for browsing and managing work items. It must feel immediate, keyboard-centric, and familiar to developers who rely on Magit, Tig, or other terminal-first tools. This document outlines framework choices, view layouts, navigation flows, and stylistic guidelines.

## 2. Framework Stack

| Library | Purpose | Notes |
|---------|---------|-------|
| [Bubbletea](https://github.com/charmbracelet/bubbletea) | Elm-inspired TUI state machine | Provides update loop, message handling, and window resizing support |
| [Bubbles](https://github.com/charmbracelet/bubbles) | Reusable UI widgets (list, viewport, text input, progress) | Accelerates prototyping; customize delegates for domain data |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Styling primitives (colors, borders, spacing) | Use for consistent theming and status-specific highlights |
| [Stickers](https://github.com/76creates/stickers) | Flexbox-style layout manager | Define responsive containers and breakpoints that adapt to terminal width |
| [Glamour](https://github.com/charmbracelet/glamour) | Markdown rendering | Applied to issue descriptions and comments |
| [Huh](https://github.com/charmbracelet/huh) (optional) | Form builder for interactive dialogs | Useful for guided issue creation beyond `$EDITOR` |
| [Wish](https://github.com/charmbracelet/wish) (stretch) | SSH-hosted multi-user sessions | Enables collaborative board sessions in future phases |

## 3. Primary Buffers & Views

Magit builds confidence through a "status" buffer that fan-outs into focused surfaces. Hubless mirrors that rhythm: the application always returns to a rich home buffer, and every other view is a focused derivative that keeps the same keyboard grammar.

### 3.1 Status Buffer (Home)

- **Layout**: Top statusline (repo, current branch, sync health, clock skew) with right-aligned quick help (`?` to expand). The body is a vertical stack of collapsible sections rendered with Lipgloss headings and `bubbles/list` delegates. Footer hosts a mini log of background jobs.
- **Sections**: Focus (assigned or bookmarked issues), Inbox (recent events grouped by project/milestone), Boards (columns with counts and SLA indicators), Saved Filters (user-defined queries), and System (sync warnings, detached HEAD alerts).
- **Interactions**: `tab`/`shift+tab` cycle sections; number keys `1-9` jump directly; `enter`/`l` drills into the section's default focused view; `h` or `esc` collapses back; `space` opens an item-specific transient menu; `.` folds/unfolds a section; `g` `g` replays the underlying queries without leaving the buffer.
- **Dynamics**: The buffer renders immediately with skeleton rows, replacing them with hydrated data streams as they arrive (<200 ms target). Each section displays its last refresh timestamp and pending background jobs so the user always feels oriented.

### 3.2 Issue List Focus

- **Layout**: Split-pane; left column uses `bubbles/list` for grouped issues, right column uses `bubbles/viewport` for a preview of the selected issue. Section headers in the list can be toggled with `tab` (Magit-style folding) to isolate status or priority buckets.
- **Displayed Fields**: Issue ID, title, priority glyph, status chip, assignee initials, relative last activity.
- **Interactions**: `↑/↓` or `k/j` change selection; `enter`/`l` opens detail view; `x` toggles multi-select; `A` opens the assignment transient; `S` opens the status-change transient; `c` starts a comment (defaults to `$EDITOR`); `o` expands the preview into full detail without leaving the buffer.
- **Transitions**: Accessible from the status buffer with `g` `i` or by drilling into any section item. Exiting with `b`, `h`, or `ctrl+o` returns to the status buffer and restores scroll position.

### 3.3 Issue Detail View

- **Layout**: Full-width viewport rendering Glamour-formatted markdown with sticky header (issue metadata, active timers, quick actions). Side panel (optional) shows related pull requests or dependencies when terminal width allows.
- **Content**: Entire event timeline (creation, edits, comments, status changes, assignments) with Magit-style section folding for comments versus system events.
- **Actions**: `c` add comment (opens `$EDITOR` or inline form); `s` change status via transient; `a` assign/unassign; `r` replay events to verify state; `b`/`h` return to the previous buffer; `m` opens dependency map; `f` toggles follow mode to auto-scroll as new events stream in.

### 3.4 Kanban View

- **Layout**: Configurable columns rendered with Lipgloss borders and headers that surface WIP limits and SLA warnings. Rows use compact cards with priority accent bars and assignee avatars (two-character initials).
- **Navigation**: Arrow keys or `h/l` switch columns; `k/j` move within a column; `tab` cycles focus between columns and action bars; `ctrl+f`/`ctrl+b` page through tall columns without losing selection.
- **Actions**: `space` or `enter` triggers the move transient pre-populated with legal status transitions; `M` toggles auto-leveling mode to rebalance WIP; `,` and `.` reorder columns on the fly.

### 3.5 Command Palette & Transient Menus

- **Trigger**: `space` on any focused entity opens a contextual transient (Magit popup) listing verbs with one-letter shortcuts. `:` opens a global command palette with fuzzy search over every action and saved workflow.
- **Presentation**: Transients show short/long keys (`s`/`S`), the command description, and whether it will drop into `$EDITOR`. Palette entries display the scope (buffer, global, experimental) and remember last invocation for quick repeat via `.`.
- **Extensibility**: Actions register through an application service so adapters can contribute verbs without touching the UI core. Popups share confirmation plumbing and respect dry-run mode.

### 3.6 Filter & Search View

- **Trigger**: `/` opens a `bubbles/textinput` component at the bottom of the screen with immediate filtering when typing. `ctrl+s` cycles between saved filters.
- **Query Language**: `status:open assignee:me priority:high` with fuzzy text search fallback. Queries can be saved (`alt+enter`) and surface inside the status buffer's Saved Filters section.
- **Persistence**: Named filters appear with numeric shortcuts (`1`, `2`, `3`) and are stored in the CLI config. Holding `shift` while applying a filter spawns a detachable buffer you can keep alongside the home buffer.

### 3.7 Sync Progress Overlay

- **Trigger**: `hubless sync`, pressing `S` inside the TUI, or when background watchers detect divergence.
- **Display**: `bubbles/progress` bar with stages (fetch, apply, project) and counters for events pushed/pulled. Idle view shows the last sync age and next scheduled sync.
- **Feedback**: Success and error notifications render as temporary toast panels at the bottom. Failures link to a troubleshooting transient with recommended Git commands.

### 3.8 Responsive Layout Profiles

- **Layout Engine**: Every buffer composes its view through a Stickers `FlexLayout`. The root model owns named regions (`statusline`, `rail`, `content`, `footer`) and passes them to child views so they can render proportionally without duplicating layout math.
- **Breakpoints**: Hubless defines three profiles: `sm` (`<100` columns) collapses to a single column with context overlays; `md` (`100-139`) keeps a narrow navigation rail and stackable content; `lg` (`≥140`) renders full side-by-side panes with generous gutters. Heights flex based on available rows but guarantee the footer and statusline stay pinned.
- **Implementation Sketch**:

```go
func layoutFor(width int) stickers.LayoutProfile {
    switch {
    case width < 100:
        return mobileProfile
    case width < 140:
        return mediumProfile
    default:
        return desktopProfile
    }
}

func (m Model) View() string {
    profile := layoutFor(m.viewport.Width)
    layout := profile.Layout(stickers.Size{Width: m.viewport.Width, Height: m.viewport.Height})
    return layout.Compose(
        stickers.Region("statusline", m.statusline.View()),
        stickers.Region("rail", m.sidebar.View()),
        stickers.Region("content", m.body.View()),
        stickers.Region("footer", m.footer.View()),
    )
}
```

- **Behavior**: Resize messages trigger a debounced recalculation of the active profile. Buffers only rerender when their assigned region dimensions change, keeping redraw minimal on aggressive resizes or when running inside `tmux` panes. Profiles live in `internal/ui/tui/layout` so new surfaces can reuse the same responsive logic.

## 4. Interaction Model

### 4.1 Boot Sequence & Initialization

- Launching `hubless` composes the root Bubbletea model with the status buffer already mounted so the user never lands on an empty screen.
- The initial Stickers profile defaults to `md` until the first `tea.WindowSizeMsg` arrives; the resulting layout broadcast ensures buffers render correctly even on terminals that hide their dimensions during boot.
- `Init` kicks off concurrent commands: catalog hydrate, snapshot replay, sync health probe, and configuration load. Each command streams partial results back to the status buffer sections so they can hydrate independently.
- While loading, the status buffer renders skeleton rows and a subtle spinner in the statusline. By the 200 ms mark we show at least one actionable section; if data is still pending the footer lists which jobs remain.
- Any fatal boot error drops the app into a dedicated diagnostic buffer that surfaces the offending Git command and recovery guidance.

### 4.2 Buffer Stack & Focus Management

- Every surface (status, issue list, detail, board, filter result) implements `tea.Model`. The root keeps a stack so buffers can push/pull like Magit's transient buffers.
- `ctrl+o` pops back (previous buffer), `ctrl+i` moves forward, and `g` `s` always returns home. The stack stores scroll offsets so returning restores context.
- Focus traversal follows a consistent order: main content → sidebar → footer hints. Within the status buffer, `tab`/`shift+tab` cycle sections; inside lists, `h/j/k/l` mirrors arrow keys.
- Buffers declare their affordances (supports multi-select, offers bulk actions) so global commands can decide whether to enable features such as the command palette.

### 4.3 Keybinding Strategy

| Keys | Scope | Action |
|------|-------|--------|
| `g` `s` | Global | Jump to status buffer |
| `g` `i` | Global | Jump to issue list focus |
| `g` `b` | Global | Jump to Kanban board |
| `g` `/` | Global | Open filter/search buffer |
| `?` | Global | Toggle inline help overlay with context-specific bindings |
| `space` | Focused item | Open transient menu of verbs |
| `tab` / `shift+tab` | Status buffer | Next/previous section |
| `↑/↓` or `k/j` | Lists & menus | Move selection |
| `h/l` | Lists & menus | Collapse/expand or move horizontal focus |
| `ctrl+o` / `ctrl+i` | Global | Back/forward buffer stack |

Bindings always offer both Vim-style and arrow-key equivalents to stay welcoming. Prefixes use short sequences (`g` `*`, `ctrl+*`) that Magit users expect, and the inline help overlay keeps the chords discoverable without leaving the keyboard.

### 4.4 Transient Actions & Editor Handoff

- `space` opens a transient menu anchored to the focused entity, showing one-character accelerators, full descriptions, and whether the action is destructive. Capital variants act on multi-selection when available.
- Confirmations happen inline (no `yes/no` prompts) by requiring a second keypress, e.g., `space` → `d` then `d` to delete, mirroring Magit's double-tap safety.
- Commands that require rich text hand off to `$EDITOR` but keep the transient state alive so returning allows `ctrl+enter` to submit immediately.
- Undo is intentionally absent for MVP; instead, the transient footer surfaces the equivalent Git command (e.g., `git revert <event>`) so users know the escape hatch.

### 4.5 Message & Effect Handling

- Child models expose strongly typed messages (e.g., `IssueSelected`, `StatusChanged`) that bubble to the root. The root orchestrates side effects with `tea.Batch` to minimize redraws.
- Long-running operations stream progress messages so views can diff and redraw incrementally rather than repainting full buffers.
- Global notifications (sync success/failure, background job completion) publish through a centralized toast manager that respects the current buffer's available screen real estate.
- Resizing events are debounced (magically 16 ms budget) to avoid thrashing on terminal resize, and rate-limited watchers prevent redundant diff computation when Git emits bursts of events. Each accepted resize emits a `LayoutChanged` message with the new Stickers profile so buffers can adjust only if their region footprint shifts.

## 5. Visual Styling Guidelines

- **Color Palette**:
  - Status Open: blue (`#1E90FF` equivalent).
  - Status In Progress: amber (`#FFBE00`).
  - Status Closed: green (`#32CD32`).
  - Snapshots or derived data: muted gray italics.
- **Priority Indicators**: 🔥 (high), ● (medium), · (low) prefixed to issue titles.
- **Borders**: Use Lipgloss rounded borders sparingly to delineate columns without overwhelming text density.
- **Typography**: Monospaced fonts by default; rely on padding rather than heavy separators.
- **Statusline**: Two-tone bar pinned to the top; left shows repo/branch/sync health, right shows context-aware key hints. Falls back to monochrome when colors are unavailable.
- **Transients**: Command popups use a subtle vertical gradient and drop shadow when supported; otherwise they render as plain bordered boxes aligned to the focused row.
- **Responsive Breakpoints**: Stickers profiles (`sm`, `md`, `lg`) control spacing, typography scale, and gutter width so smaller terminals prioritize readability while larger canvases surface secondary context.

## 6. Performance Considerations

- Avoid recomputing derived views in the render loop; memoize catalog-derived summaries.
- Cache status buffer aggregates (Focus, Inbox, Boards) and diff them before emitting messages so partial updates do not repaint the entire buffer.
- Precompute Stickers layout templates per breakpoint so resize handling only swaps profiles rather than re-allocating nodes each frame.
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
- How much end-user configuration should exists for the status buffer sections and key chords before we risk diluting the shared mental model?

## 10. References

- `docs/TechSpec.md` for event model and storage contracts.
- `docs/reference/implementation-skeleton.md` for Bubbletea wiring examples.

## 11. Wireframes

Desktop (`lg`) SVG mockups live under `docs/images/`:
- `docs/images/status-buffer-desktop.svg`
- `docs/images/issue-list-desktop.svg`
- `docs/images/issue-detail-desktop.svg`
- `docs/images/kanban-desktop.svg`
- `docs/images/command-transient-desktop.svg`
- `docs/images/filter-search-desktop.svg`
- `docs/images/sync-overlay-desktop.svg`

Mobile (`sm`) counterparts illustrate stacked layouts and overlays:
- `docs/images/status-buffer-mobile.svg`
- `docs/images/issue-list-mobile.svg`
- `docs/images/issue-detail-mobile.svg`
- `docs/images/kanban-mobile.svg`
- `docs/images/command-transient-mobile.svg`
- `docs/images/filter-search-mobile.svg`
- `docs/images/sync-overlay-mobile.svg`
