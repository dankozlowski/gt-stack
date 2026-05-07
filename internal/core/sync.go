package core

import (
	"context"
	"fmt"
	"sort"

	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/state"
)

type SyncReport struct {
	Fetched   bool
	Deleted   []string
	Restacked []string
}

// Sync fetches trunk's remote, refreshes PR states, deletes branches whose PRs are merged,
// and restacks remaining tracked branches onto the up-to-date trunk.
//
// Conflicts during restack are surfaced as errors with the standard `gts continue` hint.
func (c *Core) Sync(ctx context.Context) (*SyncReport, error) {
	rep := &SyncReport{}

	if _, err := c.Git.CurrentBranch(ctx); err != nil {
		return nil, err
	}
	if err := c.Git.Fetch(ctx, "origin"); err != nil {
		return rep, fmt.Errorf("fetch origin: %w", err)
	}
	rep.Fetched = true

	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return rep, err
	}

	// Collect tracked branch names for a single bulk PR query.
	// Sort for deterministic gh-args ordering (so tests and logs are stable).
	var heads []string
	for name, br := range s.Branches {
		if br.Tracked && br.PR > 0 {
			heads = append(heads, name)
		}
	}
	sort.Strings(heads)

	prs, err := c.GH.PRList(ctx, gh.ListOpts{Heads: heads, State: "all"})
	if err != nil {
		return rep, fmt.Errorf("gh pr list: %w", err)
	}
	prByHead := map[string]gh.PR{}
	for _, p := range prs {
		prByHead[p.HeadRefName] = p
	}

	// Delete branches whose PR is merged. Iterate sorted for stable Deleted order.
	var trackedSorted []string
	for name, br := range s.Branches {
		if br.Tracked && br.PR > 0 {
			trackedSorted = append(trackedSorted, name)
		}
		_ = br // appease lint
	}
	sort.Strings(trackedSorted)
	for _, name := range trackedSorted {
		p, ok := prByHead[name]
		if !ok {
			continue
		}
		if p.State == "MERGED" {
			if err := c.Git.BranchDelete(ctx, name); err != nil {
				return rep, fmt.Errorf("delete %s: %w", name, err)
			}
			_ = state.Untrack(ctx, c.Git, name)
			rep.Deleted = append(rep.Deleted, name)
			delete(s.Branches, name)
		}
	}

	// Restack what remains. Children of trunk get rebased onto trunk; their
	// children get rebased onto their (now-restacked) parents in turn.
	if err := restackDescendants(ctx, c, s, s.Trunk); err != nil {
		return rep, err
	}
	for name := range s.Branches {
		if s.Branches[name].Tracked {
			rep.Restacked = append(rep.Restacked, name)
		}
	}
	sort.Strings(rep.Restacked)
	return rep, nil
}
