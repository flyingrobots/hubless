package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/flyingrobots/hubless/internal/release"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("hubless-release: ")

	var (
		repoRoot   string
		version    string
		notesPath  string
		dryRun     bool
		skipChecks bool
	)

	flag.StringVar(&repoRoot, "repo", ".", "Repository root (defaults to current directory)")
	flag.StringVar(&version, "version", "", "Version to tag (required)")
	flag.StringVar(&notesPath, "notes", "docs/reference/release-notes.md", "Path to release notes markdown")
	flag.BoolVar(&dryRun, "dry-run", false, "Show actions without creating a tag")
	flag.BoolVar(&skipChecks, "skip-checks", false, "Skip fmt/lint/test/docs before tagging")
	flag.Parse()

	releaser, err := release.New(repoRoot)
	if err != nil {
		log.Fatalf("initialize releaser: %v", err)
	}

	if err := releaser.Run(context.Background(), release.Options{
		Version:    version,
		NotesPath:  notesPath,
		DryRun:     dryRun,
		SkipChecks: skipChecks,
	}); err != nil {
		if errors.Is(err, release.ErrVersionRequired) {
			fmt.Fprintln(os.Stderr, "--version is required")
			flag.Usage()
			os.Exit(2)
		}
		log.Fatalf("release failed: %v", err)
	}
}
