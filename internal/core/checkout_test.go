package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/git"
)

func TestCheckoutDown_GoesToParent(t *testing.T) {
	fr := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("feat/b\n")},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/a\nfeat/b\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":              {Stdout: []byte("main\n")},
			"config --get branch.feat/b.gts-parent":              {Stdout: []byte("feat/a\n")},
			"config --get branch.feat/a.gts-pr":                  {ExitCode: 1},
			"config --get branch.feat/b.gts-pr":                  {ExitCode: 1},
			"checkout feat/a":                                    {},
		},
	}
	c := New(git.New(fr), nil)
	got, err := c.CheckoutDown(context.Background(), 1)
	if err != nil {
		t.Fatalf("CheckoutDown: %v", err)
	}
	if got != "feat/a" {
		t.Errorf("got %q, want feat/a", got)
	}
}

func TestCheckoutUp_SingleChild(t *testing.T) {
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
			"checkout feat/b":                                    {},
		},
	}
	c := New(git.New(fr), nil)
	got, err := c.CheckoutUp(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("CheckoutUp: %v", err)
	}
	if got != "feat/b" {
		t.Errorf("got %q, want feat/b", got)
	}
}
