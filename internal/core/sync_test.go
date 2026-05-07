package core

import (
	"context"
	"testing"

	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
)

func TestSync_FetchesAndPrunesMerged(t *testing.T) {
	gitFR := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("main\n")},
			"fetch origin --prune":                               {},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/old\nfeat/new\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/old.gts-parent":            {Stdout: []byte("main\n")},
			"config --get branch.feat/new.gts-parent":            {Stdout: []byte("main\n")},
			"config --get branch.feat/old.gts-pr":                {Stdout: []byte("100\n")},
			"config --get branch.feat/new.gts-pr":                {Stdout: []byte("101\n")},
			"branch -D feat/old":                                 {},
			"config --local --unset branch.feat/old.gts-parent":  {},
			"config --local --unset branch.feat/old.gts-pr":      {},
			"merge-base main feat/new":                           {Stdout: []byte("base123\n")},
			"rebase --onto main base123 feat/new":                {},
		},
	}
	ghFR := &gh.FakeRunner{
		Responses: map[string]gh.FakeResponse{
			"pr list --json number,state,body,baseRefName,headRefName,title,url --limit 200 --state all --head feat/new --head feat/old": {
				Stdout: []byte(`[
					{"number":100,"state":"MERGED","headRefName":"feat/old","baseRefName":"main"},
					{"number":101,"state":"OPEN","headRefName":"feat/new","baseRefName":"main"}
				]`),
			},
		},
	}
	c := New(git.New(gitFR), gh.New(ghFR))
	report, err := c.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if len(report.Deleted) != 1 || report.Deleted[0] != "feat/old" {
		t.Errorf("expected feat/old deleted, got %v", report.Deleted)
	}
}
