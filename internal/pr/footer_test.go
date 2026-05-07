package pr

import (
	"strings"
	"testing"

	"github.com/dankoz/gt-stacks/internal/state"
)

func TestParseBlock_Present(t *testing.T) {
	body := "Hello\n\n<!-- gts:stack-start -->\nold content\n<!-- gts:stack-end -->\n\nbye"
	got, ok := ParseBlock(body)
	if !ok {
		t.Fatalf("expected to find block")
	}
	if got != "old content\n" {
		t.Errorf("got %q, want %q", got, "old content\n")
	}
}

func TestParseBlock_Absent(t *testing.T) {
	body := "Hello\nNo markers here"
	if _, ok := ParseBlock(body); ok {
		t.Errorf("expected no block")
	}
}

func TestReplaceOrAppend_AppendsWhenAbsent(t *testing.T) {
	body := "Summary\n\nDetails"
	got := ReplaceOrAppend(body, "RENDERED")
	want := "Summary\n\nDetails\n\n<!-- gts:stack-start -->\nRENDERED\n<!-- gts:stack-end -->\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestReplaceOrAppend_ReplacesWhenPresent(t *testing.T) {
	body := "Hi\n<!-- gts:stack-start -->\nOLD\n<!-- gts:stack-end -->\n"
	got := ReplaceOrAppend(body, "NEW")
	want := "Hi\n<!-- gts:stack-start -->\nNEW\n<!-- gts:stack-end -->\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderBlock_LinearStack(t *testing.T) {
	s := &state.Stack{
		Trunk: "main",
		Branches: map[string]*state.Branch{
			"feat/a": {Name: "feat/a", Parent: "main", PR: 100, PRState: "MERGED", Tracked: true},
			"feat/b": {Name: "feat/b", Parent: "feat/a", PR: 101, PRState: "OPEN", Tracked: true},
			"feat/c": {Name: "feat/c", Parent: "feat/b", PR: 102, PRState: "OPEN", Tracked: true},
		},
	}
	got := RenderBlock(s, "feat/b")

	for _, want := range []string{
		"**Stack**",
		"#100",
		"feat/a",
		"#101",
		"feat/b",
		"← you are here",
		"#102",
		"feat/c",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rendered block missing %q\ngot:\n%s", want, got)
		}
	}
}

func TestRenderBlock_BranchWithoutPR(t *testing.T) {
	s := &state.Stack{
		Trunk: "main",
		Branches: map[string]*state.Branch{
			"feat/a": {Name: "feat/a", Parent: "main", PR: 100, PRState: "OPEN", Tracked: true},
			"feat/b": {Name: "feat/b", Parent: "feat/a", PR: 0, Tracked: true},
		},
	}
	got := RenderBlock(s, "feat/a")
	if !strings.Contains(got, "feat/b (no PR yet)") {
		t.Errorf("expected '(no PR yet)' for feat/b\ngot:\n%s", got)
	}
}
