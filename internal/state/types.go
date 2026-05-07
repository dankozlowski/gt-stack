package state

import "sort"

type Branch struct {
	Name    string
	Parent  string // "" if parent is trunk
	PR      int    // 0 if no PR yet
	PRState string // "OPEN"|"MERGED"|"CLOSED"|""
	Tracked bool
}

type Stack struct {
	Trunk    string
	Branches map[string]*Branch
}

// Children returns names of branches whose Parent is `name`. Sorted for stability.
func (s *Stack) Children(name string) []string {
	var out []string
	for _, b := range s.Branches {
		if b.Parent == name && b.Tracked {
			out = append(out, b.Name)
		}
	}
	sort.Strings(out)
	return out
}

// Ancestors returns the chain from `name`'s parent up to (but not including) trunk.
func (s *Stack) Ancestors(name string) []string {
	var out []string
	cur, ok := s.Branches[name]
	if !ok {
		return nil
	}
	for cur.Parent != "" && cur.Parent != s.Trunk {
		out = append(out, cur.Parent)
		next, ok := s.Branches[cur.Parent]
		if !ok {
			break
		}
		cur = next
	}
	return out
}

// Roots returns tracked branches whose parent is the trunk. Sorted.
func (s *Stack) Roots() []string {
	var out []string
	for _, b := range s.Branches {
		if b.Tracked && (b.Parent == "" || b.Parent == s.Trunk) {
			out = append(out, b.Name)
		}
	}
	sort.Strings(out)
	return out
}
