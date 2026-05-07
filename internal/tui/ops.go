package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dankoz/gt-stacks/internal/core"
)

type opFinished struct {
	name    string
	summary string
	err     error
}

func runSubmit(ctx context.Context, c *core.Core) tea.Cmd {
	return func() tea.Msg {
		rep, err := c.Submit(ctx, core.SubmitOpts{})
		if err != nil {
			return opFinished{name: "submit", err: err}
		}
		var b strings.Builder
		for _, o := range rep.Created {
			fmt.Fprintf(&b, "created #%d %s\n", o.PR, o.Branch)
		}
		for _, o := range rep.Updated {
			fmt.Fprintf(&b, "updated #%d %s\n", o.PR, o.Branch)
		}
		return opFinished{name: "submit", summary: strings.TrimSpace(b.String())}
	}
}

func runRestack(ctx context.Context, c *core.Core) tea.Cmd {
	return func() tea.Msg {
		if err := c.Restack(ctx); err != nil {
			return opFinished{name: "restack", err: err}
		}
		return opFinished{name: "restack", summary: "ok"}
	}
}

func runSync(ctx context.Context, c *core.Core) tea.Cmd {
	return func() tea.Msg {
		_, err := c.Sync(ctx)
		if err != nil {
			return opFinished{name: "sync", err: err}
		}
		return opFinished{name: "sync", summary: "ok"}
	}
}

func newSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return sp
}
