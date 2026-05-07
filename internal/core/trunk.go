package core

import (
	"context"

	"github.com/dankoz/gt-stacks/internal/state"
)

func (c *Core) SetTrunk(ctx context.Context, branch string) error {
	return state.SetTrunk(ctx, c.Git, branch)
}

func (c *Core) GetTrunk(ctx context.Context) (string, error) {
	v, err := c.Git.ConfigGet(ctx, "gts.trunk")
	if err != nil {
		return "", err
	}
	if v == "" {
		return "main", nil
	}
	return v, nil
}
