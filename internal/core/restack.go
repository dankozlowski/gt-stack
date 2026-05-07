package core

import (
	"context"
	"fmt"

	"github.com/dankoz/gt-stacks/internal/state"
)

// Restack walks descendants of the current branch and rebases each onto its
// recorded parent. Used after `modify` or after upstream parent changes.
//
// Algorithm: BFS from current branch's children. For each child:
//   - oldBase = merge-base(parent, child)
//   - git rebase --onto parent oldBase child
//
// On conflict, the rebase is left in progress; user runs `gts continue`.
// We checkout back to the original branch on success.
func (c *Core) Restack(ctx context.Context) error {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return err
	}
	if err := restackDescendants(ctx, c, s, cur); err != nil {
		return err
	}
	return c.Git.Checkout(ctx, cur)
}

func restackDescendants(ctx context.Context, c *Core, s *state.Stack, root string) error {
	queue := s.Children(root)
	for len(queue) > 0 {
		branch := queue[0]
		queue = queue[1:]

		parent := s.Branches[branch].Parent
		oldBase, err := c.Git.MergeBase(ctx, parent, branch)
		if err != nil {
			return fmt.Errorf("merge-base %s..%s: %w", parent, branch, err)
		}
		if err := c.Git.RebaseOnto(ctx, parent, oldBase, branch); err != nil {
			return fmt.Errorf("rebase %s onto %s: %w (run `gts continue` after resolving)", branch, parent, err)
		}
		queue = append(queue, s.Children(branch)...)
	}
	return nil
}
