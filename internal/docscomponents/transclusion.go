package docscomponents

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TransclusionOptions instructs the CLI invocation that renders templates.
type TransclusionOptions struct {
	Bin        string
	Args       []string
	BasePath   string
	InputPath  string
	OutputPath string
}

// RunTransclusion runs the markdown-transclusion CLI to render a document template
// using opts. It validates that opts.InputPath and opts.OutputPath are set, defaults
// opts.Bin to "markdown-transclusion" when empty, and uses opts.BasePath or the
// current working directory as the base. Paths are resolved to absolute values
// relative to the base, the output directory is created if necessary, and the CLI is
// executed with the provided context. On failure the returned error includes the
// CLI's combined output.
func RunTransclusion(ctx context.Context, opts TransclusionOptions) error {
	if opts.InputPath == "" {
		return errors.New("input path is required")
	}
	if opts.OutputPath == "" {
		return errors.New("output path is required")
	}

	bin := strings.TrimSpace(opts.Bin)
	if bin == "" {
		bin = "markdown-transclusion"
	}

	basePath := opts.BasePath
	if basePath == "" {
		detected, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("determine working directory: %w", err)
		}
		basePath = detected
	}

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Errorf("resolve base path: %w", err)
	}

	absInput, err := makeAbsoluteWithBase(opts.InputPath, absBasePath)
	if err != nil {
		return fmt.Errorf("resolve input path: %w", err)
	}

	absOutput, err := makeAbsoluteWithBase(opts.OutputPath, absBasePath)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(absOutput), 0o755); err != nil {
		return fmt.Errorf("ensure output directory: %w", err)
	}

	args := append([]string{}, opts.Args...)
	args = append(args, absInput)
	args = append(args, "--output", absOutput, "--base-path", absBasePath)

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = absBasePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run %s: %w\n%s", bin, err, strings.TrimSpace(string(output)))
	}

	return nil
}

// absolute path fails.
func makeAbsoluteWithBase(pathValue, base string) (string, error) {
	if filepath.IsAbs(pathValue) {
		return pathValue, nil
	}
	if strings.TrimSpace(pathValue) == "" {
		return "", errors.New("path cannot be empty")
	}
	candidate := filepath.Join(base, pathValue)
	return filepath.Abs(candidate)
}
