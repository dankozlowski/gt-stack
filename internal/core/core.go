package core

import (
	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
)

// Core holds wired-up dependencies for all stack operations.
type Core struct {
	Git *git.Git
	GH  *gh.GH
}

func New(g *git.Git, h *gh.GH) *Core {
	return &Core{Git: g, GH: h}
}
