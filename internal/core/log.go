package core

import (
	"context"

	"github.com/dankoz/gt-stacks/internal/render"
	"github.com/dankoz/gt-stacks/internal/state"
)

// LogTree returns the rendered stack tree with `current` highlighted.
func (c *Core) LogTree(ctx context.Context, color bool) (string, error) {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return "", err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return "", err
	}
	return render.StackTree(s, cur, color), nil
}
