package git

import (
	"context"
	"strings"
)

type Git struct {
	r Runner
}

func New(r Runner) *Git { return &Git{r: r} }

func (g *Git) CurrentBranch(ctx context.Context) (string, error) {
	out, _, err := g.r.Run(ctx, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (g *Git) IsClean(ctx context.Context) (bool, error) {
	out, _, err := g.r.Run(ctx, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(out))) == 0, nil
}

func (g *Git) RepoRoot(ctx context.Context) (string, error) {
	out, _, err := g.r.Run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (g *Git) Branches(ctx context.Context) ([]string, error) {
	out, _, err := g.r.Run(ctx, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			names = append(names, line)
		}
	}
	return names, nil
}

// ConfigGet returns the value or empty string if the key is unset (exit 1 from `git config --get`).
func (g *Git) ConfigGet(ctx context.Context, key string) (string, error) {
	out, _, err := g.r.Run(ctx, "config", "--get", key)
	if err != nil {
		// git config exits 1 when key is missing; treat that as "no value".
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}

// ConfigGetAll returns all values for a key (multi-valued config).
func (g *Git) ConfigGetAll(ctx context.Context, key string) ([]string, error) {
	out, _, err := g.r.Run(ctx, "config", "--get-all", key)
	if err != nil {
		return nil, nil
	}
	var vals []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			vals = append(vals, line)
		}
	}
	return vals, nil
}

// MergeBase returns the merge-base of two refs.
func (g *Git) MergeBase(ctx context.Context, a, b string) (string, error) {
	out, _, err := g.r.Run(ctx, "merge-base", a, b)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (g *Git) ConfigSet(ctx context.Context, key, value string) error {
	_, _, err := g.r.Run(ctx, "config", "--local", key, value)
	return err
}

func (g *Git) ConfigUnset(ctx context.Context, key string) error {
	_, _, err := g.r.Run(ctx, "config", "--local", "--unset", key)
	return err
}

func (g *Git) Checkout(ctx context.Context, branch string) error {
	_, _, err := g.r.Run(ctx, "checkout", branch)
	return err
}

func (g *Git) BranchCreate(ctx context.Context, name, startPoint string) error {
	_, _, err := g.r.Run(ctx, "checkout", "-b", name, startPoint)
	return err
}

func (g *Git) BranchDelete(ctx context.Context, name string) error {
	_, _, err := g.r.Run(ctx, "branch", "-D", name)
	return err
}

// RebaseOnto runs: git rebase --onto <newParent> <oldParent> <branch>
// Used when a branch's parent has moved (e.g. after parent was amended).
func (g *Git) RebaseOnto(ctx context.Context, newParent, oldParent, branch string) error {
	_, _, err := g.r.Run(ctx, "rebase", "--onto", newParent, oldParent, branch)
	return err
}

func (g *Git) RebaseContinue(ctx context.Context) error {
	_, _, err := g.r.Run(ctx, "rebase", "--continue")
	return err
}

func (g *Git) RebaseAbort(ctx context.Context) error {
	_, _, err := g.r.Run(ctx, "rebase", "--abort")
	return err
}

func (g *Git) CommitAll(ctx context.Context, message string) error {
	_, _, err := g.r.Run(ctx, "commit", "-a", "-m", message)
	return err
}

func (g *Git) AmendNoEdit(ctx context.Context) error {
	_, _, err := g.r.Run(ctx, "commit", "--amend", "--no-edit")
	return err
}

func (g *Git) Fetch(ctx context.Context, remote string) error {
	_, _, err := g.r.Run(ctx, "fetch", remote, "--prune")
	return err
}
