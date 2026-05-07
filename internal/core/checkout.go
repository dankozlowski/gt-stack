package core

import (
	"context"
	"fmt"

	"github.com/dankoz/gt-stacks/internal/state"
)

// Picker chooses among multiple branches when ambiguity exists. Tests pass nil.
type Picker func(prompt string, options []string) (string, error)

// CheckoutDown moves to parent (n levels). Returns the branch checked out.
func (c *Core) CheckoutDown(ctx context.Context, n int) (string, error) {
	if n < 1 {
		n = 1
	}
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return "", err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return "", err
	}
	target := cur
	for i := 0; i < n; i++ {
		br, ok := s.Branches[target]
		if !ok || br.Parent == "" {
			return "", fmt.Errorf("cannot go down %d levels from %s (reached trunk)", n, cur)
		}
		target = br.Parent
	}
	if err := c.Git.Checkout(ctx, target); err != nil {
		return "", err
	}
	return target, nil
}

// CheckoutUp moves to a child (n levels). When multiple children exist, picker is consulted.
func (c *Core) CheckoutUp(ctx context.Context, n int, pick Picker) (string, error) {
	if n < 1 {
		n = 1
	}
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return "", err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return "", err
	}
	target := cur
	for i := 0; i < n; i++ {
		kids := s.Children(target)
		switch len(kids) {
		case 0:
			return "", fmt.Errorf("no children above %s", target)
		case 1:
			target = kids[0]
		default:
			if pick == nil {
				return "", fmt.Errorf("multiple children of %s: %v", target, kids)
			}
			chosen, err := pick(fmt.Sprintf("which child of %s?", target), kids)
			if err != nil {
				return "", err
			}
			target = chosen
		}
	}
	if err := c.Git.Checkout(ctx, target); err != nil {
		return "", err
	}
	return target, nil
}

// Checkout switches to the named branch. Errors if branch does not exist locally.
func (c *Core) Checkout(ctx context.Context, branch string) error {
	return c.Git.Checkout(ctx, branch)
}
