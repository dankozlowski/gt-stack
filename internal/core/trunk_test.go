package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestSetTrunk_WritesConfig(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"config --local gts.trunk develop": {},
		},
	}
	c := New(git.New(fr), nil)
	if err := c.SetTrunk(context.Background(), "develop"); err != nil {
		t.Fatalf("SetTrunk: %v", err)
	}
	if len(fr.Calls) != 1 {
		t.Fatalf("expected 1 git call, got %d: %v", len(fr.Calls), fr.Calls)
	}
}
