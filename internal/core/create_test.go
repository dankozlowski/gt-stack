package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestCreate_FromTrackedParent(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                      {Stdout: []byte("feat/a\n")},
			"status --porcelain":                             {Stdout: []byte(" M file.go\n")}, // staged
			"checkout -b feat/b feat/a":                      {},
			"commit -a -m WIP":                               {},
			"config --local branch.feat/b.gts-parent feat/a": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Create(context.Background(), "feat/b", "WIP"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if len(fr.Calls) < 4 {
		t.Fatalf("expected at least 4 git calls, got %d", len(fr.Calls))
	}
}

func TestCreate_NoStagedChanges_NoCommit(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                      {Stdout: []byte("feat/a\n")},
			"status --porcelain":                             {Stdout: []byte("")},
			"checkout -b feat/b feat/a":                      {},
			"config --local branch.feat/b.gts-parent feat/a": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Create(context.Background(), "feat/b", ""); err != nil {
		t.Fatalf("Create: %v", err)
	}
	for _, call := range fr.Calls {
		if call[0] == "commit" {
			t.Errorf("did not expect a commit call when worktree is clean: %v", call)
		}
	}
}
