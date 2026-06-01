package docscomponents

import (
	"context"
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
