package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dankoz/gt-stacks/internal/core"
)

// Run launches the TUI program. Returns when the user quits.
func Run(c *core.Core) error {
	p := tea.NewProgram(New(c), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func New(c *core.Core) tea.Model {
	return initialModel(c)
}

func initialModel(c *core.Core) model {
	return model{
		core: c,
		ctx:  context.Background(),
	}
}
