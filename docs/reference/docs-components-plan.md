# Docs Component Library Plan

## Integration Summary

- `cmd/docs-components` walks the `@hubless/` JSON (milestones, features, stories, tasks) and emits reusable Markdown fragments under `docs/components/`.
- Templates now live alongside the data (`@hubless/roadmap/templates/`, `@hubless/issues/templates/`) and transclude those fragments with the `markdown-transclusion` CLI.
- Rendered docs land in `@hubless/roadmap/generated/` and `@hubless/issues/generated/`; everything in `generated/` is overwritten by the pipeline, while JSON + templates remain human-edited.
- Shared snippets (progress, dependencies, tables, status summaries) can be referenced by any doc. Templates under `@hubless/**/templates/` are the transclusion sources; manually maintained docs should link to snippets unless they are wired into the render pipeline.
- Release notes (`docs/reference/release-notes.md`) are generated from the same archive/changelog snippets for easy copy/paste into changelogs, and the root `CHANGELOG.md` is rebuilt from `CHANGELOG.template.md` on every `make docs` run.

## Data Flow

1. **Generator** (`go run ./cmd/docs-components`) reads JSON, validates required fields, and writes Markdown snippets to `docs/components/...`.
2. **Transclusion** runs `markdown-transclusion` against each template, resolving `![[...]]` references into complete documents in `@hubless/**/generated/`.
3. Downstream docs can link to those snippets, or templates can transclude them when a generated output is desired.

## Directory Layout

- `@hubless/README.md` – quick guide on what’s generated vs. hand-edited.
- `@hubless/roadmap/`
  - `milestones/*.json`, `features/*.json` – source data (edit these).
  - `templates/` – Markdown shells with `![[...]]` includes (edit these).
  - `generated/` – rendered Markdown (do **not** edit).
- `@hubless/issues/`
  - `stories/*.json`, `tasks/*.json` – source data.
  - `templates/`, `generated/` – same pattern as above.
- `docs/components/`
  - `roadmap/milestones-table.md`, `.../stories-table.md` – tabular listings.
  - `roadmap/progress.md` – completion bars across milestones/features/stories/tasks.
  - `roadmap/dependencies.md` – dependency matrix for milestones → tasks.
  - `roadmap/dependencies-graph.md` – Mermaid graph of in-repo dependencies.
  - `issues/tasks-table.md` – task rollup.
  - `issues/status-summary.md` – task counts by status.
  - `issues/archived-stories.md`, `issues/archived-tasks.md` – completed work rollups feeding the archive template.
  - `@hubless/issues/templates/archive.md` / `generated/archive.md` – archive overview assembled from the archived snippets.
  - `issues/changelog.md` – release-ready bullet list of completed tasks (consumed by `docs/reference/release-notes.*`).
- `docs/reference/release-notes.template.md` / `docs/reference/release-notes.md` – generated release notes ready to drop into changelog entries.
- `docs/reference/palettes.json` (+ `docs/reference/palettes.schema.json`) – optional custom palette definitions merged with built-ins.
  - See `docs/reference/archive-structure.md` for a deeper explanation of archive/changelog relationships.
  - Downstream docs (PRD §9.2, TechSpec §10.2) embed the archived task rollup directly.

## Running the Pipeline

1. Install or clone [`markdown-transclusion`](https://github.com/flyingrobots/markdown-transclusion) (Node ≥20). Build it locally if not installed globally.
2. Configure how to invoke the CLI:

   ```bash
   export MARKDOWN_TRANSCLUSION_BIN=markdown-transclusion   # global install
   # or
   export MARKDOWN_TRANSCLUSION_BIN=node
   export MARKDOWN_TRANSCLUSION_SCRIPT=/path/to/markdown-transclusion/dist/cli.js
   ```

3. From the repo root run:

   ```bash
   make docs          # preferred entrypoint
   # or run the script directly
   ./scripts/render-docs.sh
   ```

   Optional variables:

   - `MARKDOWN_TRANSCLUSION_SCRIPT` – optional script path passed as the first CLI argument when `MARKDOWN_TRANSCLUSION_BIN=node`.
   - `MARKDOWN_TRANSCLUSION_BASE` – override the base path passed to the CLI.
   - `MARKDOWN_TRANSCLUSION_ARGS` – extra flags forwarded to the CLI.

   Useful CLI flags:

   - `--graph-direction` (LR/RL/TB/BT) tunes the dependency graph orientation.
   - `--graph-clusters` groups nodes by type using Mermaid subgraphs.
   - `--graph-palette` selects a Mermaid color palette (`evergreen`, `infrared`, `zerothrow`, or any palette declared in `docs/reference/palettes.json`).
   - `--palette-file` points to an alternate palette JSON document; when omitted, `docs/reference/palettes.json` is used only if it exists.
4. Validate outputs:

   ```bash
   make docs-test     # generator unit tests
   make docs-verify   # ensure transclusions fully resolved
   ```

5. Mirror CI locally with matching toolchain:

   ```bash
   ./scripts/ci-local.sh
   ```

   This builds `.ci/Dockerfile` (Go 1.26.3, Node 20, markdown-transclusion 1.2.0, golangci-lint v2.12.2), mounts the repo at `/workspace`, and executes the same sequence used in GitHub
   Actions (`fmt-check`, lightweight lint, `go vet`, tests, docs render, docs verification).

## Verification

- `go test ./internal/docscomponents` exercises the generator (fixtures & formatting guards).
- `go run ./cmd/docs-components --skip-transclusion` regenerates snippets without touching templates (helpful during development).
- `make docs` should leave only expected deltas under `docs/components/` and `@hubless/**/generated/`.
- `.github/workflows/docs.yml` enforces `make docs` + `make docs-test` on every push/PR and fails if generated files drift.

## Extending the Library

1. Model the new component in Go (add a generator method and unit tests).
2. Emit the Markdown snippet under an appropriate `docs/components/<area>/` path.
3. Reference it from the relevant template(s) in `@hubless/**/templates/`.
4. Add the snippet to any docs that should surface the data.

## Future Enhancements

- Render dependency graphs (Mermaid or Fang-based diagrams) once graph schemas stabilize.
- Expand the Makefile with grouped doc targets (e.g., `make docs/roadmap`, `make docs/tui`).
- Integrate the pipeline into CI so merges fail if generated docs drift from JSON.
