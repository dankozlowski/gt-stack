package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/dankoz/gt-stacks/internal/state"
)

func (c *Core) Status(ctx context.Context) (string, error) {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return "", err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return "", err
	}
	br, ok := s.Branches[cur]
	var b strings.Builder
	fmt.Fprintf(&b, "branch: %s\n", cur)
	fmt.Fprintf(&b, "trunk: %s\n", s.Trunk)
	if !ok || !br.Tracked {
		fmt.Fprintln(&b, "(not tracked — run `gts track` to add to a stack)")
		return b.String(), nil
	}
	fmt.Fprintf(&b, "parent: %s\n", br.Parent)
	if br.PR > 0 {
		fmt.Fprintf(&b, "pr:     #%d (%s)\n", br.PR, br.PRState)
	}
	return b.String(), nil
}
