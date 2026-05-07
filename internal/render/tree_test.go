package render

import (
	"strings"
	"testing"

	"github.com/dankoz/gt-stacks/internal/state"
)

func TestStackTree_LinearStack(t *testing.T) {
	s := &state.Stack{
		Trunk: "main",
		Branches: map[string]*state.Branch{
			"feat/a": {Name: "feat/a", Parent: "main", PR: 100, PRState: "MERGED", Tracked: true},
			"feat/b": {Name: "feat/b", Parent: "feat/a", PR: 101, PRState: "OPEN", Tracked: true},
		},
	}
	got := StackTree(s, "feat/b", false)
	for _, want := range []string{"main", "feat/a", "feat/b", "#100", "#101"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}
