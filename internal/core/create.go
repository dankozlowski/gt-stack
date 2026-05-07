package core

import (
	"context"
	"fmt"

	"github.com/dankoz/gt-stacks/internal/state"
)

// Create creates a new branch as a child of the current branch.
// If the worktree has changes and message != "", they're committed on the new branch.
// If the worktree has changes and message == "", returns an error.
func (c *Core) Create(ctx context.Context, name, message string) error {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return err
	}
	clean, err := c.Git.IsClean(ctx)
	if err != nil {
		return err
	}
	if !clean && message == "" {
		return fmt.Errorf("worktree has changes; pass a commit message with `gts create <name> -m <msg>` or stash first")
	}
	if err := c.Git.BranchCreate(ctx, name, cur); err != nil {
		return err
	}
	if !clean {
		if err := c.Git.CommitAll(ctx, message); err != nil {
			return err
		}
	}
	return state.SetParent(ctx, c.Git, name, cur)
}
