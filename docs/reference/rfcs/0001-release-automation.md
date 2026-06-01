# RFC 0001: Release Automation CLI

## Summary

Provide a Git-first release command that creates annotated tags from generated release notes, supports dry-runs, and remains safe to use in local or CI environments. The current implementation accepts an explicit `--version`, optional `--notes`, `--dry-run`, `--skip-checks`, and opt-in `--sign`; automatic bump inference, manifest updates, push, and GitHub Releases remain follow-up work.

## Goals

- Keep releases entirely Git-driven (tags and commits are the source of truth).
- Create annotated tags from rendered release notes, with optional GPG signing.
- Accept explicit versions today while preserving room for conventional-commit bump inference.
- Allow optional version manifest updates (`VERSION` file) in a future phase without making them mandatory.
- Support preflight checks (fmt/lint/test/docs) with configuration to skip when appropriate.
- Provide clear recovery guidance if any step fails.
- Make publishing to remotes (Git push, GitHub Releases in the future) opt-in.

## Non-goals

- Automatic publishing to GitHub Releases in v1 (captured as a follow-up adapter).
- Full semantic-release parity. We only need basic semantic bumping and release note templating.
- Replacing existing CI release workflows; the CLI should integrate, not dictate.

## CLI Design

```bash
hubless release \
  --version X.Y.Z \
  [--notes docs/reference/release-notes.md] \
  [--tag-prefix v] \
  [--sign] \
  [--dry-run] \
  [--skip-checks]
```

- `--version` is required in the current implementation.
- `--tag-prefix` defaults to `v` (for example, `v1.2.3`).
- `--notes` points to rendered release notes and defaults to `docs/reference/release-notes.md`.
- `--sign` runs `git tag -s`; the default is an unsigned annotated tag via `git tag -a`.
- `--dry-run` prints the tag command and note location without creating a tag.
- `--skip-checks` skips the configured preflight command sequence.

### Future CLI Surface

```bash
hubless release \
  [--patch | --minor | --major | --bump auto] \
  [--version-file VERSION] \
  [--push [<remote>]] \
  [--no-edit] \
  [--skip-verify | --skip-fmt | --skip-lint | --skip-test | --skip-docs]
```

- `--bump auto` would inspect commits since the previous tag using conventional commit semantics:
  - `BREAKING CHANGE` → major
  - `feat` → minor
  - otherwise → patch
  - `--patch|--minor|--major` override the automatic decision.
- `--version-file` would update the specified manifest and commit it as `chore(release): vX.Y.Z`.
- `--push` would push the tag and any manifest commit to the specified remote.
- `--no-edit` would skip launching `$EDITOR` for last-mile edits.
- `--skip-*` would control individual preflight steps.

## Workflow

1. **Preflight**

   - Ensure working tree is clean (`git diff --quiet` for tracked/untracked).
   - Run preflight checks unless skipped:
     - `make fmt`
     - `make lint`
     - `make test`
     - `make docs`
   - After each check, confirm the worktree is still clean or abort.

2. **Determine Version**

   - Current implementation requires `--version X.Y.Z`.
   - Future implementation can find the latest tag matching `^<prefix>\d+\.\d+\.\d+$`, infer a conventional-commit bump, and apply an override from `--patch|--minor|--major`.

3. **Assemble Notes**

   - Load release notes from `--notes` (default `docs/reference/release-notes.md`). If missing, fail unless future `--generate-notes` is provided.
   - Render notes through an optional Go template (future enhancement) or write raw notes to a temp file.
   - Unless `--no-edit`, open `$EDITOR` for final tweaks.

4. **Version Manifest (optional)**

   - If `--version-file` provided:
     - Update the file to the new version.
     - Create a commit `chore(release): vX.Y.Z` (signed commit is optional; default unsigned, with `--sign-commit` to opt-in).

5. **Create Tag**

   - Create annotated tag (`git tag -s` or `-a` depending on flags) using the prepared notes.
   - If the tag already exists, abort. This project does not support replacing release tags.

6. **Push (optional)**

   - If `--push` provided, push the tag and the manifest commit (if created) to the specified remote (default `origin`).

7. **Dry-run handling**

   - Skip tag/commit creation, but still report intended version, tag, and note location. Optionally skip expensive checks if `--skip-verify` used.

8. **Rollback Guidance**

   - Document recovery steps: `git tag -d <tag>` and reset/checkout for the version manifest commit.

## Implementation Plan

- **`internal/release` service**

  - Current: explicit version validation, configurable preflight checks, annotated tag creation, opt-in signed tags, and dry-run output.
  - Future: add commit parsing for conventional commits, version bump logic, optional manifest update, and push logic via `git push <remote> <tag>`.
  - Provide structured dry-run output for CI logs.

- **CLI `cmd/release`**

  - Map flags to service options.
  - Handle environment detection, signal-aware cancellation, and graceful error messaging.
  - Provide `--help` with detailed flag descriptions.

- **Makefile integration**

  - Targets `make release VERSION=X.Y.Z` and `make release-dry VERSION=X.Y.Z` (with optional `NOTES`, `VERSION_FILE`, etc.).

- **Tests & Validation**

  - Unit tests for version bump inference, manifest updates, tag creation commands (mock os/exec).
  - Docker-based integration test (`scripts/test-release-docker.sh`) to exercise the CLI in an isolated repo without remotes.
  - Optional GitHub Action to run the docker test on PRs touching release tooling.

## Documentation & Adoption

- Update `README.md` with release command usage, examples, and hooks.
- Expand `CONTRIBUTING.md` with release checklist (install hooks, ensure GPG configured, run `make release` flow).
- Provide recovery instructions in docs (`docs/reference/archive-structure.md` or new release guide).
- Track future adapters (GitHub Release, manifest syncing) in follow-up RFCs.

## Open Questions

- Should auto-generated notes include commit sections (feat/fix/etc.) or rely purely on docs/reference output?
- Do we want to support tagging multiple modules (monorepo scenario) in the future?
- Should we add guardrails for tag naming (prefix enforcement, semantic validation) beyond the current regex?

## Rollout

1. Implement bump inference, signing, manifest handling, and push flags.
2. Wire Docker release test into the quality GitHub Action.
3. Update docs and run the release command on a staged repo to build trust.
4. After team sign-off, mark the feature DONE and update the changelog via the new tooling.
