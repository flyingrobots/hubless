package release

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewRejectsMissingRepoRoot(t *testing.T) {
	t.Parallel()

	_, err := New(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected New to reject missing repository root")
	}
	if !strings.Contains(err.Error(), "repo root does not exist") {
		t.Fatalf("expected missing root error, got %q", err)
	}
}

func TestNewRejectsNonGitDirectory(t *testing.T) {
	t.Parallel()

	_, err := New(t.TempDir())
	if err == nil {
		t.Fatal("expected New to reject a non-Git directory")
	}
	if !strings.Contains(err.Error(), "repo root is not a Git working tree") {
		t.Fatalf("expected non-Git root error, got %q", err)
	}
}

func TestTagArgsUseSignedTagFlag(t *testing.T) {
	t.Parallel()

	got := tagArgs("v1.2.3", "release-notes.md", true)
	want := []string{"tag", "-s", "v1.2.3", "-F", "release-notes.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tagArgs signed mismatch:\nwant %#v\n got %#v", want, got)
	}
}

func TestTagArgsUseAnnotatedTagFlagByDefault(t *testing.T) {
	t.Parallel()

	got := tagArgs("v1.2.3", "release-notes.md", false)
	want := []string{"tag", "-a", "v1.2.3", "-F", "release-notes.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tagArgs annotated mismatch:\nwant %#v\n got %#v", want, got)
	}
}
