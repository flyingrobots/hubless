package release

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrVersionRequired = errors.New("version is required")
)

type Options struct {
	Version    string
	NotesPath  string
	DryRun     bool
	SkipChecks bool
}

type Releaser struct {
	repoRoot string
}

// New returns a Releaser for the repository located at repoRoot.
// repoRoot must be a non-empty path; it may be relative and will be resolved
// to an absolute path. Returns an error if repoRoot is empty or the path
// cannot be resolved.
func New(repoRoot string) (*Releaser, error) {
	if repoRoot == "" {
		return nil, errors.New("repo root is required")
	}

	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repo root: %w", err)
	}

	return &Releaser{repoRoot: absRoot}, nil
}

func (r *Releaser) Run(ctx context.Context, opts Options) error {
	if strings.TrimSpace(opts.Version) == "" {
		return ErrVersionRequired
	}
	version := normalizeVersion(opts.Version)

	notesPath := opts.NotesPath
	if notesPath == "" {
		notesPath = filepath.Join(r.repoRoot, "docs", "reference", "release-notes.md")
	} else if !filepath.IsAbs(notesPath) {
		notesPath = filepath.Join(r.repoRoot, notesPath)
	}

	notes, err := os.ReadFile(notesPath)
	if err != nil {
		return fmt.Errorf("read release notes: %w", err)
	}
	trimmedNotes := strings.TrimSpace(string(notes))
	if trimmedNotes == "" {
		return errors.New("release notes file is empty")
	}

	if !opts.SkipChecks {
		if err := r.runChecks(ctx); err != nil {
			return err
		}
	}

	if err := r.ensureClean(ctx); err != nil {
		return err
	}

	if opts.DryRun {
		fmt.Printf("[dry-run] Ready to tag %s using notes from %s\n", version, notesPath)
		fmt.Printf("[dry-run] Tag message preview:\n%s\n", trimmedNotes)
		fmt.Println("[dry-run] Next steps:")
		fmt.Printf("  git tag -a %s -F <notes>\n", version)
		fmt.Printf("  git push origin %s\n", version)
		fmt.Println("  Optionally: gh release create", version, "-F", notesPath)
		return nil
	}

	if err := r.ensureTagDoesNotExist(ctx, version); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "release-notes-*.md")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	if _, err := tempFile.WriteString(trimmedNotes + "\n"); err != nil {
		return fmt.Errorf("write temp release notes: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp release notes: %w", err)
	}

	if err := r.runCommand(ctx, "git", "tag", "-a", version, "-F", tempFile.Name()); err != nil {
		return fmt.Errorf("create tag: %w", err)
	}

	fmt.Printf("Created annotated tag %s.\n", version)
	fmt.Println("Next steps:")
	fmt.Printf("  git push origin %s\n", version)
	fmt.Println("  Optionally: gh release create", version, "-F", notesPath)
	return nil
}

func (r *Releaser) runChecks(ctx context.Context) error {
	commands := [][]string{
		{"make", "fmt"},
		{"make", "lint"},
		{"make", "test"},
		{"make", "docs"},
	}

	for _, args := range commands {
		if err := r.runCommand(ctx, args[0], args[1:]...); err != nil {
			return fmt.Errorf("run %s: %w", strings.Join(args, " "), err)
		}
		if err := r.ensureClean(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Releaser) ensureClean(ctx context.Context) error {
	output, err := r.capture(ctx, "git", "status", "--porcelain=v1", "--untracked-files=no")
	if err != nil {
		return fmt.Errorf("git status: %w", err)
	}
	if strings.TrimSpace(output) != "" {
		return fmt.Errorf("working tree has uncommitted changes:\n%s", output)
	}
	return nil
}

func (r *Releaser) ensureTagDoesNotExist(ctx context.Context, version string) error {
	if err := r.runCommand(ctx, "git", "fetch", "--tags", "--prune", "--quiet"); err != nil {
		return fmt.Errorf("fetch tags: %w", err)
	}
	output, err := r.capture(ctx, "git", "tag", "--list", version)
	if err != nil {
		return fmt.Errorf("check existing tags: %w", err)
	}
	if strings.TrimSpace(output) != "" {
		return fmt.Errorf("tag %s already exists", version)
	}
	return nil
}

func (r *Releaser) runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Releaser) capture(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.repoRoot
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// normalizeVersion trims whitespace from v and ensures it begins with a "v".
// If the trimmed version is empty it is returned unchanged; otherwise, a
// leading "v" is added when missing (e.g. "1.2.3" -> "v1.2.3").
func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return version
	}
	version = strings.TrimLeft(version, "vV")
	return "v" + version
}
