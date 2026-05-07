package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestContinue_CallsRebaseContinue(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"rebase --continue": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Continue(context.Background()); err != nil {
		t.Fatalf("Continue: %v", err)
	}
}

func TestAbort_CallsRebaseAbort(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"rebase --abort": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Abort(context.Background()); err != nil {
		t.Fatalf("Abort: %v", err)
	}
}
