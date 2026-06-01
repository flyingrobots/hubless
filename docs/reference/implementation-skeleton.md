# Hubless Implementation Skeleton

## Document Control

- Version: 0.1
- Last updated: 2025-09-18
- Maintainer: Platform Engineering

## 1. Purpose

This reference collects scaffolding snippets for implementing Hubless using Go. It mirrors the architecture described in `docs/TechSpec.md` and provides minimal examples to accelerate prototyping. The code is illustrative and may omit imports, helpers, error handling, and tests for brevity.

## 2. Project Layout

```bash
hubless/
├─ cmd/
│  ├─ docs-components/
│  │  └─ main.go               # generated docs component renderer
│  ├─ release/
│  │  └─ main.go               # release automation CLI
│  └─ hubless/
│     └─ main.go               # planned product composition root
├─ internal/
│  ├─ docscomponents/          # docs component generator
│  ├─ release/                 # release service
│  ├─ domain/                  # pure domain types and logic
│  │  ├─ events.go
│  │  └─ issue.go
│  ├─ application/             # use cases
│  │  └─ services.go
│  ├─ ports/                   # boundary interfaces (in/out)
│  │  └─ repository.go
│  ├─ adapters/
│  │  └─ gitstore/
│  │     └─ git_store.go
│  └─ ui/
│     └─ tui/                  # Bubbletea implementation
│        ├─ model.go
│        ├─ listview.go
│        └─ styles.go
├─ go.mod
├─ scripts/
│  ├─ render-docs.sh
│  ├─ verify-docs.sh
│  └─ test-release-docker.sh
└─ Makefile                    # build/test helpers
```

## 3. Module Definition (`go.mod`)

```go
module github.com/flyingrobots/hubless

go 1.26.3

require (
    github.com/charmbracelet/bubbles v1.0.0
    github.com/charmbracelet/bubbletea v1.3.10
    github.com/charmbracelet/glamour v1.0.0
    github.com/charmbracelet/lipgloss v1.1.0
)
```

## 4. Domain Layer

### 4.1 Events (`internal/domain/events.go`)

```go
package domain

import "time"

type EventType string

const (
    EventIssueCreated        EventType = "issue:created"
    EventIssueEdited         EventType = "issue:edited"
    EventIssueStatusChanged  EventType = "issue:status_changed"
    EventIssueAssigned       EventType = "issue:assigned"
    EventIssueCommented      EventType = "issue:commented"
    EventPullRequestOpened   EventType = "pr:opened"
    EventPullRequestMerged   EventType = "pr:merged"
    EventBoardMoved          EventType = "board:moved"
)

type Event struct {
    Type    EventType
    Issue   IssueID
    Actor   string
    TS      time.Time
    Lamport int
    Payload map[string]any
    EventID string
}
```

### 4.2 Issue Aggregate (`internal/domain/issue.go`)

```go
package domain

import "time"

type IssueID string

type Priority string

const (
    PriorityHigh   Priority = "high"
    PriorityMedium Priority = "medium"
    PriorityLow    Priority = "low"
)

type Issue struct {
    ID          IssueID
    Title       string
    Body        string
    Status      string
    Priority    Priority
    Assignee    string
    LastUpdated time.Time
    EventCount  int
}

func Replay(id IssueID, events []Event) Issue {
    issue := Issue{ID: id, Status: "open", Priority: PriorityMedium}
    for _, evt := range events {
        issue.LastUpdated = evt.TS
        issue.EventCount++
        switch evt.Type {
        case EventIssueCreated:
            issue.Title = getString(evt.Payload, "title", issue.Title)
            issue.Body = getString(evt.Payload, "body", issue.Body)
            if p := getString(evt.Payload, "priority", string(issue.Priority)); p != "" {
                issue.Priority = Priority(p)
            }
        case EventIssueEdited:
            issue.Title = getString(evt.Payload, "title", issue.Title)
            issue.Body = getString(evt.Payload, "body", issue.Body)
        case EventIssueStatusChanged:
            issue.Status = getString(evt.Payload, "to", issue.Status)
        case EventIssueAssigned:
            issue.Assignee = getString(evt.Payload, "assignee", issue.Assignee)
        }
    }
    return issue
}

func getString(payload map[string]any, key, fallback string) string {
    value, ok := payload[key].(string)
    if !ok || value == "" {
        return fallback
    }
    return value
}
```

## 5. Application Layer (`internal/application/services.go`)

```go
package application

import (
    "context"
    "sort"

    "github.com/flyingrobots/hubless/internal/domain"
    "github.com/flyingrobots/hubless/internal/ports"
)

type Service struct {
    store ports.EventStore
}

func NewService(store ports.EventStore) *Service {
    return &Service{store: store}
}

type IssueSummary struct {
    ID         domain.IssueID
    Title      string
    Status     string
    Priority   domain.Priority
    Assignee   string
    LastUpdate int64
}

func (s *Service) List(ctx context.Context) ([]IssueSummary, error) {
    ids, err := s.store.ListIssues(ctx)
    if err != nil {
        return nil, err
    }
    summaries := make([]IssueSummary, 0, len(ids))
    for _, id := range ids {
        events, err := s.store.LoadEvents(ctx, id)
        if err != nil {
            return nil, err
        }
        issue := domain.Replay(id, events)
        summaries = append(summaries, IssueSummary{
            ID:         issue.ID,
            Title:      issue.Title,
            Status:     issue.Status,
            Priority:   issue.Priority,
            Assignee:   issue.Assignee,
            LastUpdate: issue.LastUpdated.Unix(),
        })
    }
    sort.SliceStable(summaries, func(i, j int) bool {
        if pi, pj := domain.PriorityOrder(summaries[i].Priority), domain.PriorityOrder(summaries[j].Priority); pi != pj {
            return pi < pj
        }
        return summaries[i].LastUpdate > summaries[j].LastUpdate
    })
    return summaries, nil
}
```

## 6. Ports (`internal/ports/repository.go`)

```go
package ports

import (
    "context"
    "time"

    "github.com/flyingrobots/hubless/internal/domain"
)

type EventStore interface {
    ListIssues(ctx context.Context) ([]domain.IssueID, error)
    LoadEvents(ctx context.Context, id domain.IssueID) ([]domain.Event, error)
    AppendEvent(ctx context.Context, event domain.Event) (string, error)
    Now() time.Time
}
```

## 7. Git Adapter (`internal/adapters/gitstore/git_store.go`)

```go
package gitstore

import (
    "context"
    "fmt"
    "strings"

    "github.com/flyingrobots/hubless/internal/domain"
)

func (s *Store) AppendEvent(ctx context.Context, evt domain.Event) (string, error) {
    tree, err := s.gitWithInput("", "mktree")
    if err != nil {
        return "", err
    }
    msg := buildCommitMessage(evt)
    parent := s.currentRefHead(evt.Issue)
    oid, err := s.gitWithInput(msg, "commit-tree", tree, "-p", parent, "--author", s.author(), "--date", s.timestamp(evt.TS))
    if err != nil {
        return "", err
    }
    ref := fmt.Sprintf("refs/hubless/issues/%s", evt.Issue)
    if err := s.updateRef(ref, strings.TrimSpace(string(oid)), parent); err != nil {
        return "", err
    }
    return strings.TrimSpace(string(oid)), nil
}
```

## 8. TUI Wiring (`internal/ui/tui/model.go`)

```go
package tui

import (
    "context"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/flyingrobots/hubless/internal/application"
)

type Model struct {
    ctx    context.Context
    svc    *application.Service
    list   list.Model
    detail viewport.Model
}

func New(ctx context.Context, svc *application.Service, width, height int) Model {
    issues, _ := svc.List(ctx)
    items := make([]list.Item, len(issues))
    for i, summary := range issues {
        items[i] = issueItem{IssueSummary: summary}
    }
    l := newList(items, width/3, height)
    v := viewport.New(width-width/3-2, height)
    return Model{ctx: ctx, svc: svc, list: l, detail: v}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            return m.loadDetail(), nil
        }
    }
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}
```

## 9. Next Steps

- Flesh out unit tests for domain replay and adapters.
- Expand the Git adapter with catalog and feed updates.
- Integrate TUI commands with mutation operations (`CreateIssue`, `ChangeStatus`, `Comment`).
