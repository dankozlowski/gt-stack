package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dankoz/gt-stacks/internal/git"
)

// Load builds a Stack by reading git config for every local branch.
func Load(ctx context.Context, g *git.Git) (*Stack, error) {
	trunk, err := g.ConfigGet(ctx, "gts.trunk")
	if err != nil {
		return nil, err
	}
	if trunk == "" {
		// Default trunk detection: prefer "main", fall back to "master".
		trunk = "main"
	}

	names, err := g.Branches(ctx)
	if err != nil {
		return nil, fmt.Errorf("list branches: %w", err)
	}

	s := &Stack{Trunk: trunk, Branches: make(map[string]*Branch, len(names))}
	for _, name := range names {
		parent, err := g.ConfigGet(ctx, "branch."+name+".gts-parent")
		if err != nil {
			return nil, err
		}
		var pr int
		if parent != "" {
			// Only look up PR for tracked branches — avoids extra git calls
			// (and keeps test fixtures minimal).
			prStr, err := g.ConfigGet(ctx, "branch."+name+".gts-pr")
			if err != nil {
				return nil, err
			}
			if prStr != "" {
				pr, _ = strconv.Atoi(prStr)
			}
		}
		s.Branches[name] = &Branch{
			Name:    name,
			Parent:  parent,
			PR:      pr,
			Tracked: parent != "",
		}
	}
	return s, nil
}
