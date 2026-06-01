# Hubless Tasks

> Source of truth: individual JSON files in `@hubless/issues/tasks/`. Task IDs follow `{project}/{milestone}/{type}/{number}`.\\
> Regenerate this Markdown with `make docs` after updating task JSON files.

![[docs/components/issues/tasks-table.md]]

## Status Breakdown

![[docs/components/issues/status-summary.md]]

## Task Dependency Graph

- Keep task dependencies in-sync with the JSON files.
- A future automation will render the DAG and feed the frontier worker pool.

## How to Update

1. Add or edit a JSON file under `@hubless/issues/tasks/` (see `../schema/task.schema.json`).
2. Do not edit generated tables. Update the JSON under `@hubless/issues/tasks/` and re-run `make docs`.
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
