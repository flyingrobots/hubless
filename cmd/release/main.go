package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/flyingrobots/hubless/internal/release"
)

// main is the entry point for the hubless release CLI.
//
// It parses command-line flags to configure a release and invokes the releaser:
// - repo: repository root (default ".")
// - version: version to tag (required)
// - notes: path to release notes markdown (default "docs/reference/release-notes.md")
// - dry-run: show actions without creating a tag
// - skip-checks: skip fmt/lint/test/docs before tagging
//
// If --version is omitted the program prints a short message, shows usage and exits with code 2.
// Any other initialization or run error is logged and the program exits non‑zero.
func main() {
	log.SetFlags(0)
	log.SetPrefix("hubless-release: ")

	var (
		repoRoot   string
		version    string
		notesPath  string
		dryRun     bool
		skipChecks bool
		sign       bool
	)

	flag.StringVar(&repoRoot, "repo", ".", "Repository root (defaults to current directory)")
	flag.StringVar(&version, "version", "", "Version to tag (required)")
	flag.StringVar(&notesPath, "notes", "docs/reference/release-notes.md", "Path to release notes markdown")
	flag.BoolVar(&dryRun, "dry-run", false, "Show actions without creating a tag")
	flag.BoolVar(&skipChecks, "skip-checks", false, "Skip fmt/lint/test/docs before tagging")
	flag.BoolVar(&sign, "sign", false, "Create a GPG-signed release tag")
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected positional arguments: %s\n", strings.Join(flag.Args(), " "))
		flag.Usage()
		os.Exit(2)
	}

	releaser, err := release.New(repoRoot)
	if err != nil {
		log.Fatalf("initialize releaser: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := releaser.Run(ctx, release.Options{
		Version:    version,
		NotesPath:  notesPath,
		DryRun:     dryRun,
		SkipChecks: skipChecks,
		Sign:       sign,
	}); err != nil {
		if errors.Is(err, release.ErrVersionRequired) {
			fmt.Fprintln(os.Stderr, "--version is required")
			flag.Usage()
			os.Exit(2)
		}
		log.Fatalf("release failed: %v", err)
	}
}
