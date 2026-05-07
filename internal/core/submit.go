package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/pr"
	"github.com/dankoz/gt-stacks/internal/state"
)

type SubmitOpts struct {
	Draft bool
}

type SubmitReport struct {
	Created []SubmitOutcome
	Updated []SubmitOutcome
	Skipped []SubmitOutcome // body unchanged
	Failed  []SubmitOutcome
}

type SubmitOutcome struct {
	Branch string
	PR     int
	Err    error
}

// Submit creates PRs for any tracked branch in the current stack chain that
// lacks one, then rewrites the footer block of every PR in the chain.
func (c *Core) Submit(ctx context.Context, opts SubmitOpts) (*SubmitReport, error) {
	cur, err := c.Git.CurrentBranch(ctx)
	if err != nil {
		return nil, err
	}
	s, err := state.Load(ctx, c.Git)
	if err != nil {
		return nil, err
	}

	chain := submitChainContaining(s, cur)
	if len(chain) == 0 {
		return nil, fmt.Errorf("current branch %s is not tracked", cur)
	}

	rep := &SubmitReport{}

	// Phase 1: create missing PRs (sequentially; ordering matters because base depends on parent's PR).
	for _, name := range chain {
		br := s.Branches[name]
		if br.PR > 0 {
			continue
		}
		base := br.Parent
		if base == "" {
			base = s.Trunk
		}
		newPR, err := c.GH.PRCreate(ctx, gh.CreateOpts{
			Title: name,
			Body:  "", // empty body; footer added in phase 2
			Base:  base,
			Head:  name,
			Draft: opts.Draft,
		})
		if err != nil {
			rep.Failed = append(rep.Failed, SubmitOutcome{Branch: name, Err: err})
			return rep, err
		}
		br.PR = newPR.Number
		br.PRState = "OPEN"
		if err := state.SetPR(ctx, c.Git, name, newPR.Number); err != nil {
			return rep, err
		}
		rep.Created = append(rep.Created, SubmitOutcome{Branch: name, PR: newPR.Number})
	}

	// Phase 2: rewrite footers in parallel (bounded).
	const workers = 4
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, name := range chain {
		br := s.Branches[name]
		if br.PR == 0 {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(name string, num int) {
			defer wg.Done()
			defer func() { <-sem }()
			result := updateFooter(ctx, c, s, name, num)
			mu.Lock()
			defer mu.Unlock()
			switch {
			case result.Err != nil:
				rep.Failed = append(rep.Failed, result.SubmitOutcome)
			case result.changed:
				rep.Updated = append(rep.Updated, result.SubmitOutcome)
			default:
				rep.Skipped = append(rep.Skipped, result.SubmitOutcome)
			}
		}(name, br.PR)
	}
	wg.Wait()

	return rep, nil
}

type footerResult struct {
	SubmitOutcome
	changed bool
}

func updateFooter(ctx context.Context, c *Core, s *state.Stack, name string, num int) footerResult {
	res := footerResult{SubmitOutcome: SubmitOutcome{Branch: name, PR: num}}
	view, err := c.GH.PRView(ctx, num)
	if err != nil {
		res.Err = err
		return res
	}
	block := pr.RenderBlock(s, name)
	newBody, changed := pr.UpdateBody(view.Body, block)
	if !changed {
		return res
	}
	if err := c.GH.PREdit(ctx, num, gh.EditOpts{Body: &newBody}); err != nil {
		res.Err = err
		return res
	}
	res.changed = true
	return res
}

// submitChainContaining returns the linear root-to-leaf path through `current`.
// (Identical algorithm to internal/pr.chainContaining; duplicated to avoid
// import cycle pr → core → pr.)
func submitChainContaining(s *state.Stack, current string) []string {
	br, ok := s.Branches[current]
	if !ok {
		return nil
	}
	var anc []string
	for cur := br; cur != nil && cur.Parent != "" && cur.Parent != s.Trunk; {
		anc = append([]string{cur.Parent}, anc...)
		next, ok := s.Branches[cur.Parent]
		if !ok {
			break
		}
		cur = next
	}
	chain := append(anc, current)
	cur := current
	for {
		kids := s.Children(cur)
		if len(kids) == 0 {
			break
		}
		cur = kids[0]
		chain = append(chain, cur)
	}
	return chain
}
