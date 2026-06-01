package main

import (
	"reflect"
	"testing"
)

func TestParseArgsHonorsShellQuoting(t *testing.T) {
	t.Parallel()

	got := parseArgs(`--flag 'two words' "three words" bare`)
	want := []string{"--flag", "two words", "three words", "bare"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseArgs mismatch:\nwant %#v\n got %#v", want, got)
	}
}
