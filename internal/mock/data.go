package mock

import (
	"time"
)

// Issue represents a simplified issue summary/detail payload used by mocked CLI/TUI flows.
type Issue struct {
	ID          string
	Title       string
	Status      string
	Priority    string
	Assignee    string
	LastUpdated time.Time
	Body        string
	Comments    []Comment
	Events      []TimelineEvent
}

// Comment represents a lightweight discussion entry.
type Comment struct {
	Author    string
	Body      string
	CreatedAt time.Time
}

// TimelineEvent captures key events rendered in the detail timeline.
type TimelineEvent struct {
	Label     string
	Actor     string
	Timestamp time.Time
	Note      string
}

// BoardColumn records kanban column membership with simple counts.
type BoardColumn struct {
	Name      string
	Limit     int
	Issues    []BoardCard
	Highlight string
}

// BoardCard is a compact representation of a card within the kanban view.
type BoardCard struct {
	ID       string
	Title    string
	Assignee string
	Priority string
}

// StatusSection summarises one status buffer section.
type StatusSection struct {
	Title   string
	Items   []string
	Hint    string
	Counter int
}

// MockCatalog returns a slice of sample Issue values used by list/detail wireframes.
// The provided now time is used to compute LastUpdated, Comment.CreatedAt, and TimelineEvent.Timestamp
// so callers can control the generated timestamps.
func MockCatalog(now time.Time) []Issue {
	return []Issue{
		{
			ID:          "hubless/m1/task/0005",
			Title:       "Prototype Bubbletea TUI and Fang CLI wireframes",
			Status:      "in-progress",
			Priority:    "high",
			Assignee:    "james",
			LastUpdated: now.Add(-2 * time.Hour),
			Body:        "- Validate layout decisions\n- Capture screenshots of primary buffers\n- Iterate on keyboard flow feedback",
			Comments: []Comment{
				{Author: "james", Body: "Kicking off mocked UI pass.", CreatedAt: now.Add(-90 * time.Minute)},
				{Author: "codex", Body: "Wireframes landed in docs/images.", CreatedAt: now.Add(-30 * time.Minute)},
			},
			Events: []TimelineEvent{
				{Label: "status:started", Actor: "james", Timestamp: now.Add(-3 * time.Hour), Note: "Moved from backlog"},
				{Label: "comment", Actor: "james", Timestamp: now.Add(-90 * time.Minute), Note: "Kicking off mocked UI pass."},
				{Label: "comment", Actor: "codex", Timestamp: now.Add(-30 * time.Minute), Note: "Wireframes landed in docs/images."},
			},
		},
		{
			ID:          "hubless/m1/task/0002",
			Title:       "Introduce Fang-based CLI skeleton",
			Status:      "planned",
			Priority:    "medium",
			Assignee:    "_unassigned_",
			LastUpdated: now.Add(-6 * time.Hour),
			Body:        "Skeleton CLI with Fang command tree and dependency injection for services.",
			Comments: []Comment{
				{Author: "codex", Body: "Waiting on wireframe validation.", CreatedAt: now.Add(-5 * time.Hour)},
			},
			Events: []TimelineEvent{
				{Label: "planned", Actor: "codex", Timestamp: now.Add(-7 * time.Hour), Note: "Pulled into milestone."},
			},
		},
		{
			ID:          "hubless/m0/task/0004",
			Title:       "Structure @hubless planning artifacts",
			Status:      "in-progress",
			Priority:    "medium",
			Assignee:    "james",
			LastUpdated: now.Add(-26 * time.Hour),
			Body:        "Maintain DAG of tasks and stories for upcoming milestones.",
			Comments: []Comment{
				{Author: "james", Body: "Docs refreshed with new sections.", CreatedAt: now.Add(-23 * time.Hour)},
			},
			Events: []TimelineEvent{
				{Label: "status:started", Actor: "james", Timestamp: now.Add(-27 * time.Hour), Note: "Board grooming"},
			},
		},
	}
}

// MockStatusSections returns sample status sections used by the home buffer view.
//
// It produces four static sections — "Focus", "Inbox", "Boards", and "Saved Filters" —
// each with example items, a short hint for keyboard interaction, and a counter.
//
// The `now` parameter is accepted for API consistency with other mock factories but is
// not used to compute the returned data.
func MockStatusSections(now time.Time) []StatusSection {
	return []StatusSection{
		{
			Title:   "Focus",
			Items:   []string{"hubless/m1/task/0005 · Mocked UI flows", "hubless/m0/task/0004 · Planning artifacts"},
			Hint:    "g i to drill into issues",
			Counter: 2,
		},
		{
			Title:   "Inbox",
			Items:   []string{"codex commented on hubless/m1/task/0005", "Merge request awaiting review"},
			Hint:    "enter to expand timeline",
			Counter: 4,
		},
		{
			Title:   "Boards",
			Items:   []string{"Open: 8", "In Progress: 5 (1 over WIP)", "Review: 2", "Done: 12"},
			Hint:    "g b to inspect columns",
			Counter: 5,
		},
		{
			Title:   "Saved Filters",
			Items:   []string{"1. @me", "2. High priority", "3. Needs review"},
			Hint:    "number keys to apply",
			Counter: 3,
		},
	}
}

// MockBoard returns a deterministic set of sample Kanban columns and cards used by mock UI flows.
//
// The returned slice contains three columns ("Open", "In Progress", "Review") populated with
// BoardCard entries (ID, Title, Assignee, Priority). Intended for development and testing of
// Kanban/board views—not for production data.
func MockBoard() []BoardColumn {
	return []BoardColumn{
		{
			Name:      "Open",
			Limit:     8,
			Highlight: "",
			Issues: []BoardCard{
				{ID: "task/0010", Title: "Capture sync diagnostics", Assignee: "alex", Priority: "medium"},
				{ID: "task/0011", Title: "Draft CLI usage guide", Assignee: "james", Priority: "high"},
			},
		},
		{
			Name:      "In Progress",
			Limit:     5,
			Highlight: "over",
			Issues: []BoardCard{
				{ID: "task/0005", Title: "Prototype Bubbletea wireframes", Assignee: "james", Priority: "high"},
				{ID: "task/0004", Title: "Planning DAG maintenance", Assignee: "james", Priority: "medium"},
				{ID: "task/0009", Title: "Sync command telemetry", Assignee: "sam", Priority: "low"},
			},
		},
		{
			Name:      "Review",
			Limit:     3,
			Highlight: "",
			Issues: []BoardCard{
				{ID: "task/0007", Title: "Git adapter tests", Assignee: "mia", Priority: "high"},
			},
		},
	}
}
