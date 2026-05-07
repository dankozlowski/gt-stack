package core

import (
	"context"
	"fmt"
)

type ModifyOpts struct {
	Amend   bool   // amend the previous commit
	Message string // commit message (used when !Amend)
}

// Modify either amends the current commit or creates a new commit.
// Descendants of the current branch will need restacking afterwards
// (caller should run Restack if desired — Modify itself does not restack).
func (c *Core) Modify(ctx context.Context, opts ModifyOpts) error {
	if opts.Amend && opts.Message != "" {
		return fmt.Errorf("--amend cannot be combined with -m")
	}
	if opts.Amend {
		return c.Git.AmendNoEdit(ctx)
	}
	if opts.Message == "" {
		return fmt.Errorf("commit message required (use --amend to amend without editing)")
	}
	return c.Git.CommitAll(ctx, opts.Message)
}
