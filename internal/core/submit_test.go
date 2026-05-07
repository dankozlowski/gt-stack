package core

import (
	"context"
	"strings"
	"testing"

	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
)

func TestSubmit_CreatesMissingPR_AndUpdatesFooter(t *testing.T) {
	gitFR := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("feat/a\n")},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/a\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":              {Stdout: []byte("main\n")},
			"config --get branch.feat/a.gts-pr":                  {ExitCode: 1},
			"config --local branch.feat/a.gts-pr 200":            {},
		},
	}
	ghFR := &gh.FakeRunner{
		Responses: map[string]gh.FakeResponse{
			"pr create --title feat/a --body-file - --base main --head feat/a": {
				Stdout: []byte("https://github.com/o/r/pull/200\n"),
			},
			"pr view 200 --json number,state,body,baseRefName,headRefName,title,url": {
				Stdout: []byte(`{"number":200,"state":"OPEN","body":"","baseRefName":"main","headRefName":"feat/a","title":"feat/a","url":"https://github.com/o/r/pull/200"}`),
			},
			"pr edit 200 --body-file -": {},
		},
	}
	c := New(git.New(gitFR), gh.New(ghFR))
	rep, err := c.Submit(context.Background(), SubmitOpts{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if len(rep.Created) != 1 || rep.Created[0].Branch != "feat/a" {
		t.Errorf("expected creation of feat/a, got %+v", rep.Created)
	}
	// Confirm a body containing the marker block was sent on edit.
	for _, call := range ghFR.Calls {
		if len(call.Args) > 1 && call.Args[1] == "edit" {
			if !strings.Contains(call.Stdin, "gts:stack-start") {
				t.Errorf("edit body missing marker:\n%s", call.Stdin)
			}
		}
	}
}
