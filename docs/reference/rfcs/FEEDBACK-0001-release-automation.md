# Feedback: RFC 0001 – Release Automation
Love the direction. A couple of places to sharpen before we merge:

1. **Version source of truth** – Let's stick with tags as canonical. Optional VERSION file is fine, but release tooling should update it automatically when present so we don't drift. This also means the CLI needs a flag to skip the file update if a team doesn't use it.
2. **Bump inference** – Yes to conventional commits. I'd suggest `--bump auto` defaulting to conventional-commit scan over `prevTag..HEAD`, with `--patch|--minor|--major` as overrides. Call out that without conventional commits, users need to specify the bump explicitly.
3. **Signed tags** – Default to `git tag -s`. Provide `--no-sign` for edge cases but make sure we surface the GPG requirement up front. Maybe add a helper command to check GPG config or a doc link.
4. **Push semantics** – Agree tag creation should be local by default. A `--push` flag that pushes both the tag and optional version bump commit sounds right. Maybe `--push-origin <remote>` for flexibility.
5. **Skip controls / dry-run** – Current design always runs fmt/lint/test/docs. Keep that as default but add `--skip-verify`, `--skip-docs`, etc., so large repos can iterate faster. Dry run should still validate cleanliness but skip expensive steps when requested.
6. **Failure / rollback story** – Add a section describing how to recover if tag creation fails halfway (e.g., tag exists, GPG failure). A simple `git tag -d <tag>` plus re-run is likely enough, but let's document it.

Address those notes and I’m +1 on landing the RFC.
