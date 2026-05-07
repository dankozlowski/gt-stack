package git

import (
	"context"
	"testing"
)

func TestCurrentBranch(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"symbolic-ref --short HEAD": {Stdout: []byte("feat/auth-2\n")},
		},
	}
	g := New(fr)
	got, err := g.CurrentBranch(context.Background())
	if err != nil {
		t.Fatalf("CurrentBranch: %v", err)
	}
	if got != "feat/auth-2" {
		t.Errorf("got %q, want %q", got, "feat/auth-2")
	}
}

func TestIsClean(t *testing.T) {
	cases := []struct {
		name      string
		stdout    string
		wantClean bool
	}{
		{"clean", "", true},
		{"dirty", " M file.go\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fr := &FakeRunner{
				Responses: map[string]FakeResponse{
					"status --porcelain": {Stdout: []byte(tc.stdout)},
				},
			}
			g := New(fr)
			got, err := g.IsClean(context.Background())
			if err != nil {
				t.Fatalf("IsClean: %v", err)
			}
			if got != tc.wantClean {
				t.Errorf("got %v, want %v", got, tc.wantClean)
			}
		})
	}
}

func TestConfigGet_Missing(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"config --get branch.feat.gts-parent": {ExitCode: 1},
		},
	}
	g := New(fr)
	got, err := g.ConfigGet(context.Background(), "branch.feat.gts-parent")
	if err != nil {
		t.Fatalf("ConfigGet: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string for missing key", got)
	}
}

func TestConfigSet_RecordsArgs(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"config --local branch.feat.gts-parent main": {},
		},
	}
	g := New(fr)
	if err := g.ConfigSet(context.Background(), "branch.feat.gts-parent", "main"); err != nil {
		t.Fatalf("ConfigSet: %v", err)
	}
	if len(fr.Calls) != 1 {
		t.Fatalf("want 1 call, got %d", len(fr.Calls))
	}
}

func TestRebaseOnto_ArgOrder(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"rebase --onto main old-parent feat/x": {},
		},
	}
	g := New(fr)
	if err := g.RebaseOnto(context.Background(), "main", "old-parent", "feat/x"); err != nil {
		t.Fatalf("RebaseOnto: %v", err)
	}
}
