package git

import (
	"context"
	"strings"
	"testing"
)

func TestExecRunner_Run(t *testing.T) {
	r := NewExecRunner(".")
	stdout, _, err := r.Run(context.Background(), "rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(strings.TrimSpace(string(stdout))) == 0 {
		t.Fatalf("expected toplevel path on stdout")
	}
}

func TestExecRunner_NonZeroExit(t *testing.T) {
	r := NewExecRunner(".")
	_, stderr, err := r.Run(context.Background(), "no-such-subcommand")
	if err == nil {
		t.Fatalf("expected error for unknown subcommand")
	}
	if len(stderr) == 0 {
		t.Fatalf("expected stderr to be populated on failure")
	}
}
