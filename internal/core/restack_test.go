package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

// Restack from current=feat/a should rebase children of feat/a (just feat/b here)
// onto feat/a using `git rebase --onto feat/a <oldBase> feat/b`.
func TestRestack_RebasesChildren(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("feat/a\n")},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/a\nfeat/b\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":              {Stdout: []byte("main\n")},
			"config --get branch.feat/b.gts-parent":              {Stdout: []byte("feat/a\n")},
			"config --get branch.feat/a.gts-pr":                  {ExitCode: 1},
			"config --get branch.feat/b.gts-pr":                  {ExitCode: 1},
			"merge-base feat/a feat/b":                           {Stdout: []byte("abc123\n")},
			"rebase --onto feat/a abc123 feat/b":                 {},
			"checkout feat/a":                                    {}, // restore current after rebase
		},
	}
	c := New(git.New(fr), nil)
	if err := c.Restack(context.Background()); err != nil {
		t.Fatalf("Restack: %v", err)
	}
}
