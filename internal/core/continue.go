package core

import "context"

func (c *Core) Continue(ctx context.Context) error {
	return c.Git.RebaseContinue(ctx)
}

func (c *Core) Abort(ctx context.Context) error {
	return c.Git.RebaseAbort(ctx)
}
