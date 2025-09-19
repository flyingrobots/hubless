# Hubless Planning Data

This directory holds the human-authored JSON records plus generated markdown artifacts derived from them.

- `schema/` – JSON schemas (hand maintained).
- `roadmap/` – milestone & feature data.
  - `*.json` – source of truth (hand maintained).
  - `templates/` – Markdown templates that reference shared components (edit these when layout changes).
  - `generated/` – Output Markdown rendered from templates (do not edit, regenerated via `make docs`).
- `issues/` – stories & tasks following the task lifecycle.
  - `tasks/*.json` and `stories/*.json` – source of truth (hand maintained).
  - `templates/` – Markdown templates to document task/story tables (edit as needed).
  - `generated/` – Output Markdown rendered from templates (`tasks.md`, `archive.md`, etc.). Do not edit by hand.

Run `make docs` from the repository root to regenerate anything in `generated/` after changing source JSON or templates.
