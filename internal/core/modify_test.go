package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestModify_AmendNoEdit(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"commit --amend --no-edit": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Modify(context.Background(), ModifyOpts{Amend: true}); err != nil {
		t.Fatalf("Modify: %v", err)
	}
}

func TestModify_NewCommit(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"commit -a -m progress": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Modify(context.Background(), ModifyOpts{Message: "progress"}); err != nil {
		t.Fatalf("Modify: %v", err)
	}
}
