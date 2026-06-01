# Hubless Changelog

> This file is generated. Edit `CHANGELOG.template.md` or underlying JSON, then run `make docs`.

## Latest Release Snapshot

![[docs/reference/release-notes.md]]

## Unreleased

- Added `hubless release --sign` for GPG-signed release tags and stricter release CLI argument validation.
- Split `MARKDOWN_TRANSCLUSION_SCRIPT` from extra `MARKDOWN_TRANSCLUSION_ARGS` and added transclusion preflight validation.
- Pinned CI and container markdown-transclusion execution to version 1.2.0 and made the release-test container run as an unprivileged user.
- Tightened docs hygiene, palette schema validation, and generated Markdown newline normalization.
- Filtered shipped tasks out of active generated rollups only after all completion badges are present.

## Historical Archives

See `@hubless/issues/generated/archive.md` for the full backlog of completed stories and tasks.
