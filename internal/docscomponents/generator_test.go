package docscomponents_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flyingrobots/hubless/internal/docscomponents"
)

func TestGeneratorGenerate(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "roadmap", "milestones", "sample-milestone.json"), map[string]any{
		"id":           "sample/milestone/0001",
		"title":        "Sample Milestone",
		"status":       "DONE",
		"dependencies": []any{},
		"features":     []any{},
		"tasks":        []any{},
		"notes":        []any{},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "roadmap", "features", "sample-feature.json"), map[string]any{
		"id":     "sample/feature/0001",
		"title":  "Sample Feature",
		"status": "PLANNED",
		"dependencies": []any{
			"sample/milestone/0001",
		},
		"stories": []any{},
		"tasks":   []any{},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "issues", "stories", "sample-story.json"), map[string]any{
		"id":     "sample/story/0001",
		"title":  "Sample Story",
		"status": "DONE",
		"dependencies": []any{
			"sample/feature/0001",
		},
		"tasks": []any{},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "issues", "tasks", "sample-task-1.json"), map[string]any{
		"id":           "sample/task/0001",
		"title":        "Sample Task Done",
		"status":       "DONE",
		"owner":        "dev",
		"labels":       []any{"docs"},
		"badges":       []any{"Tested"},
		"updated_at":   "2025-09-19",
		"dependencies": []any{"sample/story/0001"},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "issues", "tasks", "sample-task-2.json"), map[string]any{
		"id":           "sample/task/0002",
		"title":        "Sample Task Planned",
		"status":       "PLANNED",
		"labels":       []any{"docs"},
		"badges":       []any{},
		"updated_at":   nil,
		"dependencies": []any{},
	})

	componentsDir := filepath.Join(repoRoot, "docs", "components")
	gen, err := docscomponents.NewGenerator(repoRoot, componentsDir, docscomponents.GeneratorOptions{})
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}

	if err := gen.Generate(context.Background()); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	progress := readFile(t, filepath.Join(componentsDir, "roadmap", "progress.md"))
	if !strings.Contains(progress, "[##########] 100%") {
		t.Fatalf("expected progress bar to show 100%% completion, got:\n%s", progress)
	}

	dependencies := readFile(t, filepath.Join(componentsDir, "roadmap", "dependencies.md"))
	if !strings.Contains(dependencies, "sample/task/0001") {
		t.Fatalf("expected task dependency row in dependencies summary, got:\n%s", dependencies)
	}

	archivedStories := readFile(t, filepath.Join(componentsDir, "issues", "archived-stories.md"))
	if !strings.Contains(archivedStories, "sample/story/0001") {
		t.Fatalf("expected archived stories snippet to include completed story, got:\n%s", archivedStories)
	}

	archivedTasks := readFile(t, filepath.Join(componentsDir, "issues", "archived-tasks.md"))
	if !strings.Contains(archivedTasks, "Sample Task Done") {
		t.Fatalf("expected archived tasks snippet to include completed task, got:\n%s", archivedTasks)
	}

	changelog := readFile(t, filepath.Join(componentsDir, "issues", "changelog.md"))
	if !strings.Contains(changelog, "- 2025-09-19") {
		t.Fatalf("expected changelog snippet to include dated bullet, got:\n%s", changelog)
	}

	graph := readFile(t, filepath.Join(componentsDir, "roadmap", "dependencies-graph.md"))
	if !strings.Contains(graph, "graph LR") || !strings.Contains(graph, "Sample Task Done") {
		t.Fatalf("expected mermaid dependency graph to include task node label, got:\n%s", graph)
	}
	if !strings.Contains(graph, "classDef milestone") {
		t.Fatalf("expected mermaid graph to include class definitions, got:\n%s", graph)
	}
}

func TestGeneratorGenerateCustomGraphOptions(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "roadmap", "milestones", "m.json"), map[string]any{
		"id":     "custom/milestone",
		"title":  "Custom Milestone",
		"status": "DONE",
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "roadmap", "features", "f.json"), map[string]any{
		"id":           "custom/feature",
		"title":        "Custom Feature",
		"status":       "DONE",
		"dependencies": []any{"custom/milestone"},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "issues", "stories", "s.json"), map[string]any{
		"id":           "custom/story",
		"title":        "Custom Story",
		"status":       "DONE",
		"dependencies": []any{"custom/feature"},
	})

	writeJSON(t, filepath.Join(repoRoot, "@hubless", "issues", "tasks", "t.json"), map[string]any{
		"id":           "custom/task",
		"title":        "Custom Task",
		"status":       "DONE",
		"updated_at":   "2025-09-19",
		"dependencies": []any{"custom/story"},
	})

	componentsDir := filepath.Join(repoRoot, "docs", "components")
	gen, err := docscomponents.NewGenerator(repoRoot, componentsDir, docscomponents.GeneratorOptions{
		GraphDirection: "tb",
		GraphClusters:  true,
		GraphPalette:   "infrared",
	})
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}

	if err := gen.Generate(context.Background()); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	graph := readFile(t, filepath.Join(componentsDir, "roadmap", "dependencies-graph.md"))
	if !strings.Contains(graph, "graph TB") {
		t.Fatalf("expected dependency graph to honour direction TB, got:\n%s", graph)
	}
	if !strings.Contains(strings.ToLower(graph), "subgraph feature") {
		t.Fatalf("expected dependency graph to include clusters, got:\n%s", graph)
	}
	if !strings.Contains(graph, "classDef milestone fill:#1A1C23") {
		t.Fatalf("expected infrared palette colors in graph, got:\n%s", graph)
	}
}

func writeJSON(t *testing.T, path string, payload any) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", path, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(data)
}
