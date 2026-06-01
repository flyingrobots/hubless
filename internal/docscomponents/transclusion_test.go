package docscomponents

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunTransclusionFailsFastOnMissingBinary(t *testing.T) {
	t.Parallel()

	err := RunTransclusion(context.Background(), TransclusionOptions{
		Bin:        "definitely-not-a-real-markdown-transclusion-binary",
		BasePath:   t.TempDir(),
		InputPath:  "input.md",
		OutputPath: "output.md",
	})
	if err == nil {
		t.Fatal("expected missing binary to fail")
	}
	if !strings.Contains(err.Error(), "resolve transclusion binary") {
		t.Fatalf("expected missing binary resolution error, got %q", err)
	}
}

func TestRunTransclusionFailsBeforeCommandWhenInputMissing(t *testing.T) {
	t.Parallel()

	err := RunTransclusion(context.Background(), TransclusionOptions{
		Bin:        "true",
		BasePath:   t.TempDir(),
		InputPath:  "missing.md",
		OutputPath: "output.md",
	})
	if err == nil {
		t.Fatal("expected missing input to fail before invoking command")
	}
	if !strings.Contains(err.Error(), "input path does not exist") {
		t.Fatalf("expected missing input error, got %q", err)
	}
}

func TestRunTransclusionNormalizesOutputNewline(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	bin := filepath.Join(dir, "fake-transclusion")
	script := `#!/usr/bin/env sh
set -eu
out=""
while [ "$#" -gt 0 ]; do
  if [ "$1" = "--output" ]; then
    shift
    out="$1"
  fi
  shift || true
done
printf rendered > "$out"
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake transclusion binary: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "input.md"), []byte("source\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	err := RunTransclusion(context.Background(), TransclusionOptions{
		Bin:        bin,
		BasePath:   dir,
		InputPath:  "input.md",
		OutputPath: "output.md",
	})
	if err != nil {
		t.Fatalf("RunTransclusion: %v", err)
	}

	output, err := os.ReadFile(filepath.Join(dir, "output.md"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(output) != "rendered\n" {
		t.Fatalf("expected exactly one trailing newline, got %q", output)
	}
}
