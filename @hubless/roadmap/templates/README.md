# Hubless Roadmap Data

Structured roadmap data lives alongside JSON schemas so automation can generate schedules, dependency graphs, and dashboards.

> Regenerate this document with `make docs` after updating roadmap JSON.

## Snapshot

![[docs/components/roadmap/progress.md]]

## Dependencies

![[docs/components/roadmap/dependencies.md]]

## Dependency Graph

![[docs/components/roadmap/dependencies-graph.md]]

## Milestones

JSON records under `milestones/` follow `@hubless/schema/milestone.schema.json`.

![[docs/components/roadmap/milestones-table.md]]

## Features

Feature records under `features/` follow `@hubless/schema/feature.schema.json`.

![[docs/components/roadmap/features-table.md]]

## Stories

Stories reside in `../issues/stories/` following `@hubless/schema/story.schema.json`.

![[docs/components/roadmap/stories-table.md]]

Keep these tables in sync with the JSON records. Automation will eventually consume the JSON directly to render dependency graphs and schedule projections.
