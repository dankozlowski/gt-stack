package core

import (
	"context"
	"fmt"

	"github.com/dankoz/gt-stacks/internal/state"
)

// Track records `parent` as the parent of the currently checked-out branch.
// If parent is empty, defaults to the trunk.
func (c *Core) Track(ctx context.Context, parent string) error {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return err
	}
	if cur == "" {
		return fmt.Errorf("not on a branch")
	}
	if parent == "" {
		parent, err = c.GetTrunk(ctx)
		if err != nil {
			return err
		}
	}
	if parent == cur {
		return fmt.Errorf("a branch cannot be its own parent")
	}
	return state.SetParent(ctx, c.Git, cur, parent)
}

// Untrack removes the parent / PR config for the current branch.
func (c *Core) Untrack(ctx context.Context) error {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return err
	}
	return state.Untrack(ctx, c.Git, cur)
}
