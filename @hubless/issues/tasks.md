# Hubless Tasks

> Source of truth: individual JSON files in `@hubless/issues/tasks/`. Task IDs follow `{project}/{milestone}/{type}/{number}`.\
> Regenerate this Markdown whenever tasks change (manual for now).

| ID | Title | Status | Owner | Labels | Badges | Updated |
| --- | --- | --- | --- | --- | --- | --- |
| [hubless/m0/task/0001](tasks/hubless-m0-task-0001.json) | Port progress updater from Python to Go | PLANNED | _unassigned_ | m0-foundations, prog | — | — |
| [hubless/m0/task/0004](tasks/hubless-m0-task-0004.json) | Structure @hubless planning artifacts | STARTED | _unassigned_ | m0-foundations, planning | — | 2025-09-18 |
| [hubless/m0/task/0005](tasks/hubless-m0-task-0005.json) | Evaluate markdown component library | PLANNED | _unassigned_ | m0-foundations, docs, automation | — | — |
| [hubless/m1/task/0002](tasks/hubless-m1-task-0002.json) | Introduce Fang-based CLI skeleton | PLANNED | _unassigned_ | m1-cli, cli | — | — |
| [hubless/m1/task/0003](tasks/hubless-m1-task-0003.json) | Implement Git event store adapter | PLANNED | _unassigned_ | m1-cli, event-store | — | — |
| [hubless/m1/task/0005](tasks/hubless-m1-task-0005.json) | Prototype Bubbletea TUI and Fang CLI wireframes with mocked data | STARTED | _unassigned_ | m1-cli, tui, cli | — | — |

## Task Dependency Graph

- Keep task dependencies in-sync with the JSON files.
- A future automation will render the DAG and feed the frontier worker pool.

## How to Update

1. Add or edit a JSON file under `@hubless/issues/tasks/` (see `../schema/task.schema.json`).
2. Update this table to reflect new statuses or metadata.
3. Once a task is `DONE` and has badges **Tested**, **Documented**, **Shipped**, move the JSON file to `tasks.archive/` and record it in `tasks.archive.md`.

## Anatomy of a Task

```json
{
  "id": "hubless/mX/task/0000",
  "title": "One-line summary",
  "status": "PLANNED | STARTED | BLOCKED | DONE",
  "owner": "github-handle or null",
  "description": "Paragraph with context and acceptance criteria.",
  "labels": ["epic", "area"],
  "links": [
    {"type": "doc", "url": "../../docs/...", "label": "Related spec"}
  ],
  "required_inputs": [
    {"resource": "../../path-or-url", "exclusivity": "read-only", "notes": "Constraints."}
  ],
  "expected_outputs": [
    {"resource": "../../artifact", "notes": "Acceptance hints."}
  ],
  "expertise": ["Go", "Charmbracelet"],
  "dependencies": ["hubless/mX/task/0000"],
  "badges": ["Tested", "Documented", "Shipped"],
  "created_at": "YYYY-MM-DD",
  "updated_at": "YYYY-MM-DD or null",
  "notes": [
    {"at": "YYYY-MM-DDThh:mm:ssZ", "author": "name", "body": "Short update."}
  ]
}
```
