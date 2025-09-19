# Update Progress Algorithm

## Document Control

- Version: 0.1
- Last updated: 2025-09-18
- Source: Transcribed from `update_progress.py`

## 1. Purpose

This document captures the exact behavior of the original Python script that synchronized GitMind’s Features Ledger and README progress indicators. It exists so the algorithm can be reimplemented in other languages (e.g., Go) without referring back to the deleted script.

## 2. Inputs and Outputs

- **Inputs**:
  - Markdown ledger at `<root>/docs/features/Features_Ledger.md`.
  - Optional `<root>/README.md` status section.
  - Optional `--root` CLI flag or `GITMIND_ROOT` environment variable to locate the GitMind checkout.
- **Outputs**:
  - Updated progress blocks embedded in the ledger file.
  - Updated README progress block when the document contains a `## 📊 Status` section.

## 3. Root Resolution

1. Accept an optional `--root` flag.
2. Look for `HUBLESS_ROOT` in the environment.
3. `git rev-parse --show-toplevel`

## 4. Markdown Block Patterns

The script uses compiled regular expressions to locate fenced blocks:
- `<!-- group-progress:<slug>:begin --> … <!-- group-progress:<slug>:end -->`
- `<!-- progress-overall:begin --> … <!-- progress-overall:end -->`
- Milestone blocks for `mvp`, `alpha`, `beta`, `v1`.
- Task blocks embedded inside Markdown block quotes: `> <!-- tasks-mvp:begin -->` … `> <!-- tasks-mvp:end -->`.
- README guard block: `<!-- features-progress:begin -->` … `<!-- features-progress:end -->`.

## 5. Progress Bar Renderer

Given a percentage in the range `[0, 1]`, the script:

1. Clamps the value.
2. Computes `filled = round(pct * width)` with `width = 40`.
3. Inserts an edge character `▓` if there is a fractional remainder and space for an additional cell.
4. Pads the remainder of the bar with `░`.
5. Appends a textual percentage (e.g., ` 72%`).

## 6. Group Progress Calculation

For each group block in the ledger:

1. Identify the next Markdown table that starts with a header cell containing `Emoji`.
2. Parse the table header to locate the `Progress`, `KLoC`, and `Milestone` columns (case-insensitive substring match).
3. Iterate rows until a non-table line is hit.
4. Extract a numeric percentage from the `Progress` column. If missing, scan all cells for a trailing `%` or any digits.
5. Determine row weight:
   - Default: `0.1`.
   - If `KLoC` column contains a positive float, use that value instead.
6. Accumulate weighted averages and count rows.
7. Optionally append `(milestone_label, pct, weight)` tuples for downstream milestone aggregation.
8. Compute the weighted percentage or fallback to arithmetic mean if all weights are zero.

### Feature Tagging

The replacement block includes a `features=<count>` footer where `<count>` is the number of table rows processed for that group.

## 7. Milestone Aggregation

1. Map milestone labels to canonical keys using:
   - `MVP → mvp`
   - `Alpha → alpha`
   - `Beta → beta`
   - `v1.0.0` / `V1 → v1`
2. Sum `pct * weight` per milestone using the weights recorded for each feature.
3. Compute average percentage per milestone (`total / weight`).
4. Enforce gating order: `mvp ≥ alpha ≥ beta ≥ v1` by cascading `min(previous, current)` down the chain.
5. Calculate overall progress as the weighted sum of milestone percentages using fixed weights `{mvp: 0.3, alpha: 0.3, beta: 0.2, v1: 0.2}`.

## 8. Ledger Mutation Steps

For every ledger run:

1. Replace each group block with a fenced code block containing the new progress bar, a legend row, and the `features=<count>` footer.
2. Update the overall section (`<!-- progress-overall -->`) with a progress bar and inline legend `MVP 70% | Alpha 55% | …`.
3. Replace each milestone block (`progress-mvp`, etc.) with the gated percentage.
4. Extract outstanding tasks via the Tasklist parser (Section 9) and rewrite each block quote to contain either the list of tasks (each line prefixed with `> `) or the placeholder `> _All tracked tasks complete_`.
5. Write the ledger back to disk only if changes were detected.

## 9. Tasklist Parsing

1. Locate the first occurrence of `## Tasklist` (case-insensitive).
2. Inspect subsequent lines for unchecked tasks (`- [ ] …`).
3. Detect milestone tags in either `[tag]` or `(tag)` prefixes. Tags may include dotted forms (e.g., `MVP.core`); every segment is checked against the milestone label map.
4. Default untagged tasks to `mvp`.
5. Return a mapping `{mvp|alpha|beta|v1 → list[str]}` where list entries omit the tag wrapper but preserve the checkbox syntax.

## 10. README Update Logic

1. Skip if `<root>/README.md` does not exist or lacks the `## 📊 Status` heading.
2. Ensure the status section contains a guard block. If missing, insert the placeholder:

   ````markdown
   <!-- features-progress:begin -->
   ```text
   Feature progress to be updated via hubless/update_progress.py
   ```
   <!-- features-progress:end -->
   ````

3. When an overall percentage is available, replace the guard block with the rendered progress bar.
4. Collapse runs of more than two blank lines to keep the document tidy.

## 11. Execution Flow

```bash
parse_args()
root = resolve_root(args.root)
configure_paths(root)  # sets global ROOT, LEDGER, README
overall, milestone_progress, tasks = update_ledger()
update_readme(overall)
```

The script exits with code `0` on success and prints errors to stderr before exiting non-zero when path resolution fails or file operations raise exceptions.

## 12. Reimplementation Checklist

- Re-create the regex guards exactly as listed above to ensure idempotent updates.
- Preserve task defaulting to MVP when no tag is provided.
- Honor the gating rule for milestone percentages before calculating the overall score.
- Ensure the progress bar renderer replicates the Unicode characters and rounding behavior.
- Only write modified files to avoid unnecessary git diffs.
