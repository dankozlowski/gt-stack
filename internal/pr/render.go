package pr

import (
	"fmt"
	"strings"

	"github.com/dankoz/gt-stacks/internal/state"
)

// RenderBlock builds the markdown block listing every tracked branch in the
// stack containing `current`, root-to-leaf. The current branch is highlighted.
func RenderBlock(s *state.Stack, current string) string {
	chain := chainContaining(s, current)
	var b strings.Builder
	fmt.Fprintf(&b, "**Stack** (this PR is **▶ %s**):\n\n", current)
	for _, name := range chain {
		br := s.Branches[name]
		glyph := stateGlyph(br, name == current)
		switch {
		case br.PR == 0:
			fmt.Fprintf(&b, "- %s %s (no PR yet)\n", glyph, name)
		case name == current:
			fmt.Fprintf(&b, "- %s **#%d · %s**  ← you are here\n", glyph, br.PR, name)
		default:
			fmt.Fprintf(&b, "- %s #%d · %s\n", glyph, br.PR, name)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func stateGlyph(b *state.Branch, isCurrent bool) string {
	switch {
	case isCurrent:
		return "▶"
	case b.PRState == "MERGED":
		return "✓"
	case b.PRState == "CLOSED":
		return "✗"
	case b.PR > 0:
		return "○"
	default:
		return "·"
	}
}

// chainContaining returns the linear path from root-of-stack down to (and
// including) leaf-most descendant of `current`. Diverging children are NOT
// included; a stack footer represents one PR's lineage.
func chainContaining(s *state.Stack, current string) []string {
	br, ok := s.Branches[current]
	if !ok {
		return nil
	}
	// Walk to root.
	var ancestors []string
	for cur := br; cur != nil && cur.Parent != "" && cur.Parent != s.Trunk; {
		ancestors = append([]string{cur.Parent}, ancestors...)
		cur = s.Branches[cur.Parent]
	}
	chain := append(ancestors, current)
	// Walk down: pick first child each level (deterministic via Children sort).
	cur := current
	for {
		kids := s.Children(cur)
		if len(kids) == 0 {
			break
		}
		cur = kids[0]
		chain = append(chain, cur)
	}
	return chain
}
