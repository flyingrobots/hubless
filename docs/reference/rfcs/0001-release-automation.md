# RFC 0001: Release Automation CLI

## Summary
Provide a Git-first release command that bumps versions, produces annotated (signed) tags from the generated release notes, and optionally pushes/publishes the results. The tool must integrate with existing changelog automation, support dry-runs, and remain safe to use in local or CI environments.

## Goals
- Keep releases entirely Git-driven (tags and commits are the source of truth).
- Automate version bumping using conventional commits, with manual overrides.
- Default to signed annotated tags, with escape hatches when necessary.
- Allow optional version manifest updates (`VERSION` file) without making them mandatory.
- Support preflight checks (fmt/lint/test/docs) with configuration to skip when appropriate.
- Provide clear recovery guidance if any step fails.
- Make publishing to remotes (Git push, GitHub Releases in the future) opt-in.

## Non-goals
- Automatic publishing to GitHub Releases in v1 (captured as a follow-up adapter).
- Full semantic-release parity. We only need basic semantic bumping and release note templating.
- Replacing existing CI release workflows; the CLI should integrate, not dictate.

## CLI Design
```
hubless release \
  [--patch | --minor | --major | --bump auto] \
  [--notes docs/reference/release-notes.md] \
  [--version-file VERSION] \
  [--tag-prefix v] \
  [--sign | --no-sign] \
  [--push [<remote>]] \
  [--dry-run] [--no-edit] \
  [--skip-verify | --skip-fmt | --skip-lint | --skip-test | --skip-docs]
```

- `--bump auto` (default) inspects commits since the previous tag using conventional commit semantics:
  - `BREAKING CHANGE` → major
  - `feat` → minor
  - otherwise → patch
  - `--patch|--minor|--major` override the automatic decision.
- `--tag-prefix` defaults to `v` (e.g., `v1.2.3`).
- `--notes` points to the rendered release notes (defaults to `docs/reference/release-notes.md`). If the file is missing and `--allow-generate-notes` is set (future), the CLI can synthesize notes from commit messages.
- `--version-file` updates the specified manifest (e.g., `VERSION`) and commits it as `chore(release): vX.Y.Z` unless `--no-commit-version` is passed.
- `--sign` (default) runs `git tag -s`; `--no-sign` downgrades to `-a`.
- `--push` pushes the tag (and version commit if it exists). Default remote `origin`; allow overrides via `--push upstream`.
- `--dry-run` prints the resolved bump, notes, and tag command without executing changes.
- `--no-edit` skips launching `$EDITOR` for last-mile edits.
- `--skip-*` controls preflight checks. `--skip-verify` skips all; individual flags skip specific steps.

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
   - Find the latest tag matching `^<prefix>\d+\.\d+\.\d+$`.
   - If `--patch|--minor|--major` specified, use that bump.
   - Otherwise, compute bump using conventional commit analysis from previous tag to `HEAD`.
   - Apply bump to previous version; default to `0.1.0` if no prior tag.

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
   - If the tag already exists, abort unless `--force-replace` and `--confirm` specified (dangerous).

6. **Push (optional)**
   - If `--push` provided, push the tag and the manifest commit (if created) to the specified remote (default `origin`).

7. **Dry-run handling**
   - Skip tag/commit creation, but still report intended version, tag, and note location. Optionally skip expensive checks if `--skip-verify` used.

8. **Rollback Guidance**
   - Document recovery steps: `git tag -d <tag>` and reset/checkout for the version manifest commit.

## Implementation Plan

- **`internal/release` service**
  - Add commit parsing for conventional commits (lightweight parser).
  - Add version bump logic and optional manifest update.
  - Support signed tags (`git tag -s`) with `--no-sign` fallback.
  - Implement push logic via `git push <remote> <tag>` (and `git push <remote> HEAD` for manifest commit if present).
  - Provide structured dry-run output for CI logs.

- **CLI `cmd/release`**
  - Map flags to service options.
  - Handle `$EDITOR`, environment detection, and graceful error messaging.
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
