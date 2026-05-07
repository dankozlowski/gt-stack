package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/dankoz/gt-stacks/internal/render"
	"github.com/dankoz/gt-stacks/internal/state"
)

type model struct {
	core    *core.Core
	ctx     context.Context
	stack   *state.Stack
	current string
	err     error
	width   int
	height  int
	quit    bool
}

type stackLoaded struct {
	stack   *state.Stack
	current string
}

type loadErr struct{ err error }

func (m model) Init() tea.Cmd { return loadStack(m.ctx, m.core) }

func loadStack(ctx context.Context, c *core.Core) tea.Cmd {
	return func() tea.Msg {
		cur, err := c.Git.CurrentBranch(ctx)
		if err != nil {
			return loadErr{err}
		}
		s, err := c.LoadStack(ctx)
		if err != nil {
			return loadErr{err}
		}
		return stackLoaded{stack: s, current: cur}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case stackLoaded:
		m.stack = msg.stack
		m.current = msg.current
	case loadErr:
		m.err = msg.err
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

var border = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2)

func (m model) View() string {
	if m.err != nil {
		return border.Render(fmt.Sprintf("error: %v\n\nq quit", m.err))
	}
	if m.stack == nil {
		return border.Render("loading…")
	}
	tree := render.StackTree(m.stack, m.current, true)
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[↑↓] nav  [enter] checkout  [s] submit  [r] restack  [c] create  [?] help  [q] quit")
	return border.Render(tree + "\n" + footer)
}
