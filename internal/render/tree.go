package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dankoz/gt-stacks/internal/state"
)

var (
	StyleTrunk   = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	StyleCurrent = lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Bold(true)
	StyleMerged  = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	StyleNoPR    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	StyleDefault = lipgloss.NewStyle()
)

// StackTree renders a tree of all tracked branches plus the trunk, with
// `current` highlighted. When color=false, no ANSI escapes are emitted.
func StackTree(s *state.Stack, current string, color bool) string {
	var b strings.Builder
	fmt.Fprintln(&b, styleIf(color, StyleTrunk, s.Trunk))
	roots := s.Roots()
	for i, root := range roots {
		isLast := i == len(roots)-1
		writeBranch(&b, s, root, current, "", isLast, color)
	}
	return b.String()
}

func writeBranch(b *strings.Builder, s *state.Stack, name, current, prefix string, isLast, color bool) {
	connector := "├─"
	childPrefix := prefix + "│ "
	if isLast {
		connector = "└─"
		childPrefix = prefix + "  "
	}
	br := s.Branches[name]
	glyph := branchGlyph(br, name == current)
	line := fmt.Sprintf("%s%s%s %s", prefix, connector, glyph, name)
	if br.PR > 0 {
		line += fmt.Sprintf("  #%d", br.PR)
	}
	if name == current {
		line += "  ← here"
		line = styleIf(color, StyleCurrent, line)
	} else if br.PRState == "MERGED" {
		line = styleIf(color, StyleMerged, line)
	} else if br.PR == 0 {
		line = styleIf(color, StyleNoPR, line)
	}
	fmt.Fprintln(b, line)

	kids := s.Children(name)
	sort.Strings(kids)
	for i, k := range kids {
		writeBranch(b, s, k, current, childPrefix, i == len(kids)-1, color)
	}
}

func branchGlyph(b *state.Branch, isCurrent bool) string {
	switch {
	case isCurrent:
		return "▶"
	case b.PRState == "MERGED":
		return "✓"
	case b.PRState == "CLOSED":
		return "✗"
	case b.PR > 0:
		return "●"
	default:
		return "○"
	}
}

func styleIf(color bool, s lipgloss.Style, text string) string {
	if !color {
		return text
	}
	return s.Render(text)
}
