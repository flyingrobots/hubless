package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/flyingrobots/hubless/internal/docscomponents"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// main is the CLI entrypoint for generating documentation components and optionally
// rendering documentation templates via the markdown-transclusion tool.
//
// It parses command-line flags to configure repository and output paths, generator
// options (graph direction, clusters, palette and palette file), and transclusion
// settings (binary, base path, and additional args). It constructs a docs
// generator, runs component generation, and—unless -skip-transclusion is set—resolves
// the transclusion binary and arguments (from flags or MARKDOWN_TRANSCLUSION_* env
// vars) and invokes markdown-transclusion to render a set of templates to their
// configured outputs. Any initialization, generation, or rendering error causes the
// program to log a fatal error and exit.
func main() {
	log.SetFlags(0)
	log.SetPrefix("docs-components: ")

	var (
		repoRoot          string
		componentsDir     string
		roadmapTemplate   string
		roadmapOutput     string
		tasksTemplate     string
		tasksOutput       string
		archiveTemplate   string
		archiveOutput     string
		releaseTemplate   string
		releaseOutput     string
		changelogTemplate string
		changelogOutput   string
		transclusionBin   string
		transclusionBase  string
		skipTransclusion  bool
		graphDirection    string
		graphClusters     bool
		graphPalette      string
		paletteFile       string
		transclusionArgs  stringSliceFlag
	)

	flag.StringVar(&repoRoot, "repo", ".", "Repository root (defaults to current directory)")
	flag.StringVar(&componentsDir, "components", "", "Components output directory (defaults to docs/components under repo)")
	flag.StringVar(&roadmapTemplate, "roadmap-template", "@hubless/roadmap/templates/README.md", "Template path for roadmap documentation")
	flag.StringVar(&roadmapOutput, "roadmap-output", "@hubless/roadmap/generated/README.md", "Output path for generated roadmap documentation")
	flag.StringVar(&tasksTemplate, "tasks-template", "@hubless/issues/templates/tasks.md", "Template path for tasks overview")
	flag.StringVar(&tasksOutput, "tasks-output", "@hubless/issues/generated/tasks.md", "Output path for generated tasks overview")
	flag.StringVar(&archiveTemplate, "archive-template", "@hubless/issues/templates/archive.md", "Template path for archive overview")
	flag.StringVar(&archiveOutput, "archive-output", "@hubless/issues/generated/archive.md", "Output path for generated archive overview")
	flag.StringVar(&releaseTemplate, "release-template", "docs/reference/release-notes.template.md", "Template path for release notes")
	flag.StringVar(&releaseOutput, "release-output", "docs/reference/release-notes.md", "Output path for generated release notes")
	flag.StringVar(&changelogTemplate, "changelog-template", "CHANGELOG.template.md", "Template path for root changelog")
	flag.StringVar(&changelogOutput, "changelog-output", "CHANGELOG.md", "Output path for generated changelog")
	flag.StringVar(&transclusionBin, "transclusion-bin", "", "Executable for markdown-transclusion CLI (defaults to MARKDOWN_TRANSCLUSION_BIN env or markdown-transclusion)")
	flag.StringVar(&transclusionBase, "transclusion-base", "", "Base path passed to markdown-transclusion (defaults to repo root)")
	flag.BoolVar(&skipTransclusion, "skip-transclusion", false, "Skip rendering templates with markdown-transclusion")
	flag.StringVar(&graphDirection, "graph-direction", "LR", "Direction for Mermaid dependency graph (LR, RL, TB, BT)")
	flag.BoolVar(&graphClusters, "graph-clusters", false, "Group dependency graph nodes by type using Mermaid subgraphs")
	flag.StringVar(&graphPalette, "graph-palette", "evergreen", "Mermaid palette for dependency graph (evergreen, infrared, zerothrow)")
	flag.StringVar(&paletteFile, "palette-file", "docs/reference/palettes.json", "Optional palette definition file (JSON)")
	flag.Var(&transclusionArgs, "transclusion-args", "Additional argument passed to markdown-transclusion (repeatable)")
	flag.Parse()

	ctx := context.Background()

	generator, err := docscomponents.NewGenerator(repoRoot, componentsDir, docscomponents.GeneratorOptions{
		GraphDirection: graphDirection,
		GraphClusters:  graphClusters,
		GraphPalette:   graphPalette,
		PaletteFile:    paletteFile,
	})
	if err != nil {
		log.Fatalf("initialise generator: %v", err)
	}

	if err := generator.Generate(ctx); err != nil {
		log.Fatalf("generate components: %v", err)
	}

	if skipTransclusion {
		return
	}

	if transclusionBin == "" {
		if envValue := os.Getenv("MARKDOWN_TRANSCLUSION_BIN"); envValue != "" {
			transclusionBin = envValue
		} else {
			transclusionBin = "markdown-transclusion"
		}
	}

	if len(transclusionArgs) == 0 {
		if envValue := os.Getenv("MARKDOWN_TRANSCLUSION_ARGS"); envValue != "" {
			transclusionArgs = append(transclusionArgs, parseArgs(envValue)...)
		}
	}

	if transclusionBase == "" {
		transclusionBase = generator.RepoRoot()
	}

	documents := []struct {
		template string
		output   string
	}{
		{template: roadmapTemplate, output: roadmapOutput},
		{template: tasksTemplate, output: tasksOutput},
		{template: archiveTemplate, output: archiveOutput},
		{template: releaseTemplate, output: releaseOutput},
		{template: changelogTemplate, output: changelogOutput},
	}

	for _, doc := range documents {
		if doc.template == "" || doc.output == "" {
			continue
		}

		opts := docscomponents.TransclusionOptions{
			Bin:        transclusionBin,
			Args:       []string(transclusionArgs),
			BasePath:   transclusionBase,
			InputPath:  doc.template,
			OutputPath: doc.output,
		}

		if err := docscomponents.RunTransclusion(ctx, opts); err != nil {
			log.Fatalf("render %s -> %s: %v", doc.template, doc.output, err)
		}
	}
}

// parseArgs splits raw into whitespace-separated fields (using strings.Fields)
// and returns a newly allocated slice containing those fields. An empty or
// all-whitespace input yields an empty slice.
func parseArgs(raw string) []string {
	fields := strings.Fields(raw)
	return append([]string(nil), fields...)
}
