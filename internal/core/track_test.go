package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestTrack_RecordsParent(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                    {Stdout: []byte("feat/x\n")},
			"config --local branch.feat/x.gts-parent main": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Track(context.Background(), "main"); err != nil {
		t.Fatalf("Track: %v", err)
	}
}

func TestUntrack_RemovesConfig(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                       {Stdout: []byte("feat/x\n")},
			"config --local --unset branch.feat/x.gts-parent": {},
			"config --local --unset branch.feat/x.gts-pr":     {ExitCode: 5}, // unset on missing key
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Untrack(context.Background()); err != nil {
		t.Fatalf("Untrack: %v", err)
	}
}
