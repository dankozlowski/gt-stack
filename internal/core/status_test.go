package core

import (
	"context"
	"strings"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestStatus_TrackedBranch(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("feat/b\n")},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/a\nfeat/b\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":              {Stdout: []byte("main\n")},
			"config --get branch.feat/b.gts-parent":              {Stdout: []byte("feat/a\n")},
			"config --get branch.feat/a.gts-pr":                  {Stdout: []byte("100\n")},
			"config --get branch.feat/b.gts-pr":                  {ExitCode: 1},
		},
	}
	c := New(git.New(fr), nil)
	out, err := c.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	for _, want := range []string{"feat/b", "parent: feat/a", "trunk: main"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
}
