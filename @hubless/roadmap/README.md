# Hubless Roadmap Data

Structured roadmap data lives alongside JSON schemas so automation can generate schedules, dependency graphs, and dashboards.

## Milestones

JSON records under `milestones/` follow `@hubless/schema/milestone.schema.json`.

| ID | Title | Status |
|----|-------|--------|
| [hubless/milestone/m0-foundations](milestones/hubless-milestone-m0-foundations.json) | Repository foundations | IN_PROGRESS |
| [hubless/milestone/m0-5-cli-proof](milestones/hubless-milestone-m0-5-cli-proof.json) | CLI proof of concept | PLANNED |
| [hubless/milestone/m1-mvp](milestones/hubless-milestone-m1-mvp.json) | MVP release | PLANNED |
| [hubless/milestone/m1-5-enhancements](milestones/hubless-milestone-m1-5-enhancements.json) | Filters and activity enhancements | PLANNED |
| [hubless/milestone/m2-github-sync](milestones/hubless-milestone-m2-github-sync.json) | GitHub synchronization | PLANNED |
| [hubless/milestone/m3-ide](milestones/hubless-milestone-m3-ide.json) | IDE integrations | PLANNED |

## Features

Feature records under `features/` follow `@hubless/schema/feature.schema.json`.

| ID | Title | Status |
|----|-------|--------|
| [hubless/feature/repo-foundations](features/hubless-feature-repo-foundations.json) | Repository foundations and automation | IN_PROGRESS |
| [hubless/feature/event-store-cli](features/hubless-feature-event-store-cli.json) | Event store and CLI foundations | PLANNED |
| [hubless/feature/tui-experience](features/hubless-feature-tui-experience.json) | Magit-grade TUI experience | PLANNED |
| [hubless/feature/docs-components](features/hubless-feature-docs-components.json) | Markdown component library integration | PLANNED |
| [hubless/feature/git-sync](features/hubless-feature-git-sync.json) | Robust Git synchronization | PLANNED |
| [hubless/feature/github-projection](features/hubless-feature-github-projection.json) | GitHub projection integration | PLANNED |
| [hubless/feature/ide-integration](features/hubless-feature-ide-integration.json) | Editor integrations via LSP | PLANNED |

## Stories

Stories reside in `../issues/stories/` following `@hubless/schema/story.schema.json`.

| ID | Title | Status |
|----|-------|--------|
| [hubless/story/0001](../issues/stories/hubless-story-0001.json) | As a maintainer I have documented workflows and automation | IN_PROGRESS |
| [hubless/story/0002](../issues/stories/hubless-story-0002.json) | As a developer I can manage issues via CLI commands | PLANNED |
| [hubless/story/0003](../issues/stories/hubless-story-0003.json) | As a developer I can browse and update work in the TUI | PLANNED |
| [hubless/story/0004](../issues/stories/hubless-story-0004.json) | As a contributor I can sync Hubless work via Git remotes | PLANNED |
| [hubless/story/0005](../issues/stories/hubless-story-0005.json) | As a maintainer I can mirror events to GitHub | PLANNED |
| [hubless/story/0006](../issues/stories/hubless-story-0006.json) | As a developer I can interact with Hubless from my IDE | PLANNED |
| [hubless/story/0007](../issues/stories/hubless-story-0007.json) | As a documentarian I can compose docs from reusable components | PLANNED |

Keep these tables in sync with the JSON records. Automation will eventually consume the JSON directly to render dependency graphs and schedule projections.
