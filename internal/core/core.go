package core

import (
	"context"

	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
	"github.com/dankoz/gt-stacks/internal/state"
)

// Core holds wired-up dependencies for all stack operations.
type Core struct {
	Git *git.Git
	GH  *gh.GH
}

func New(g *git.Git, h *gh.GH) *Core {
	return &Core{Git: g, GH: h}
}

func (c *Core) LoadStack(ctx context.Context) (*state.Stack, error) {
	return state.Load(ctx, c.Git)
}
