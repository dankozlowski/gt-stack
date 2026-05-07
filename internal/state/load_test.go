package state

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestLoad_BuildsTreeFromConfig(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"config --get gts.trunk": {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {
				Stdout: []byte("main\nfeat/a\nfeat/b\nfeat/c\nstandalone\n"),
			},
			"config --get branch.main.gts-parent":       {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":     {Stdout: []byte("main\n")},
			"config --get branch.feat/b.gts-parent":     {Stdout: []byte("feat/a\n")},
			"config --get branch.feat/c.gts-parent":     {Stdout: []byte("feat/b\n")},
			"config --get branch.standalone.gts-parent": {ExitCode: 1},
			"config --get branch.feat/a.gts-pr":         {Stdout: []byte("100\n")},
			"config --get branch.feat/b.gts-pr":         {ExitCode: 1},
			"config --get branch.feat/c.gts-pr":         {Stdout: []byte("102\n")},
		},
	}
	g := git.New(fr)

	s, err := Load(context.Background(), g)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.Trunk != "main" {
		t.Errorf("trunk = %q, want main", s.Trunk)
	}
	if !s.Branches["feat/a"].Tracked {
		t.Errorf("feat/a should be tracked")
	}
	if s.Branches["standalone"].Tracked {
		t.Errorf("standalone should NOT be tracked (no parent config)")
	}
	if s.Branches["feat/a"].PR != 100 {
		t.Errorf("feat/a PR = %d, want 100", s.Branches["feat/a"].PR)
	}
	if s.Branches["feat/b"].PR != 0 {
		t.Errorf("feat/b PR = %d, want 0", s.Branches["feat/b"].PR)
	}
}
