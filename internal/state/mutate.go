package state

import (
	"context"
	"strconv"

	"github.com/dankoz/gt-stacks/internal/git"
)

// SetParent records the parent of a branch in git config.
func SetParent(ctx context.Context, g *git.Git, branch, parent string) error {
	return g.ConfigSet(ctx, "branch."+branch+".gts-parent", parent)
}

// SetPR records the PR number for a branch.
func SetPR(ctx context.Context, g *git.Git, branch string, pr int) error {
	return g.ConfigSet(ctx, "branch."+branch+".gts-pr", strconv.Itoa(pr))
}

// Untrack removes both gts-parent and gts-pr config keys.
func Untrack(ctx context.Context, g *git.Git, branch string) error {
	if err := g.ConfigUnset(ctx, "branch."+branch+".gts-parent"); err != nil {
		return err
	}
	_ = g.ConfigUnset(ctx, "branch."+branch+".gts-pr") // ignore errors (key may be unset)
	return nil
}

// SetTrunk records the trunk branch.
func SetTrunk(ctx context.Context, g *git.Git, trunk string) error {
	return g.ConfigSet(ctx, "gts.trunk", trunk)
}
