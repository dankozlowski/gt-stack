package state

import "testing"

func TestStack_Children(t *testing.T) {
	s := &Stack{
		Trunk: "main",
		Branches: map[string]*Branch{
			"a": {Name: "a", Parent: "main", Tracked: true},
			"b": {Name: "b", Parent: "a", Tracked: true},
			"c": {Name: "c", Parent: "a", Tracked: true},
			"d": {Name: "d", Parent: "b", Tracked: true},
		},
	}
	got := s.Children("a")
	if len(got) != 2 {
		t.Fatalf("want 2 children of a, got %d", len(got))
	}
	gotMap := map[string]bool{}
	for _, c := range got {
		gotMap[c] = true
	}
	if !gotMap["b"] || !gotMap["c"] {
		t.Errorf("expected children b, c; got %v", got)
	}
}

func TestStack_Ancestors(t *testing.T) {
	s := &Stack{
		Trunk: "main",
		Branches: map[string]*Branch{
			"a": {Name: "a", Parent: "main", Tracked: true},
			"b": {Name: "b", Parent: "a", Tracked: true},
			"c": {Name: "c", Parent: "b", Tracked: true},
		},
	}
	got := s.Ancestors("c")
	want := []string{"b", "a"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("ancestor[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
