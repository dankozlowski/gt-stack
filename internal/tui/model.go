package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/dankoz/gt-stacks/internal/render"
	"github.com/dankoz/gt-stacks/internal/state"
)

type model struct {
	core     *core.Core
	ctx      context.Context
	stack    *state.Stack
	current  string
	err      error
	width    int
	height   int
	quit     bool
	branches []string // flattened, in display order
	cursor   int
	op       string // "" | "submit" | "restack" | "sync"
	opMsg    string
	spinner  spinner.Model
}

type stackLoaded struct {
	stack   *state.Stack
	current string
}

type loadErr struct{ err error }

type checkedOut struct{ branch string }

func (m model) Init() tea.Cmd {
	return tea.Batch(loadStack(m.ctx, m.core), m.spinner.Tick)
}

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

func doCheckout(ctx context.Context, c *core.Core, branch string) tea.Cmd {
	return func() tea.Msg {
		if err := c.Checkout(ctx, branch); err != nil {
			return loadErr{err}
		}
		return checkedOut{branch: branch}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case stackLoaded:
		m.stack = msg.stack
		m.current = msg.current
		m.branches = flattenBranches(m.stack)
		for i, n := range m.branches {
			if n == m.current {
				m.cursor = i
				break
			}
		}
	case loadErr:
		m.err = msg.err
	case checkedOut:
		m.current = msg.branch
	case opFinished:
		m.op = ""
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.opMsg = msg.summary
		}
		// Reload stack after op (PR states may have changed).
		return m, loadStack(m.ctx, m.core)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if m.op != "" {
			// While an op is running, only allow quit.
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				m.quit = true
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.stack == nil || len(m.branches) == 0 {
				return m, nil
			}
			target := m.branches[m.cursor]
			if target == m.current || target == m.stack.Trunk {
				return m, nil
			}
			return m, doCheckout(m.ctx, m.core, target)
		case "s":
			m.op = "submit"
			m.opMsg = ""
			return m, runSubmit(m.ctx, m.core)
		case "r":
			m.op = "restack"
			m.opMsg = ""
			return m, runRestack(m.ctx, m.core)
		case "y":
			m.op = "sync"
			m.opMsg = ""
			return m, runSync(m.ctx, m.core)
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
	if m.op != "" {
		body := m.spinner.View() + " " + m.op + "…"
		return border.Render(body + "\n\n[q] cancel")
	}
	if m.stack == nil {
		return border.Render("loading…")
	}
	rows := splitLines(render.StackTree(m.stack, m.current, true))
	for i, line := range rows {
		if i == m.cursor {
			rows[i] = "❯ " + line
		} else {
			rows[i] = "  " + line
		}
	}
	out := strings.Join(rows, "\n")
	if m.opMsg != "" {
		out += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Render(m.opMsg)
	}
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[↑↓] nav  [enter] checkout  [s] submit  [r] restack  [y] sync  [c] create  [?] help  [q] quit")
	return border.Render(out + "\n" + footer)
}

func splitLines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

func flattenBranches(s *state.Stack) []string {
	out := []string{s.Trunk}
	var walk func(name string)
	walk = func(name string) {
		out = append(out, name)
		for _, k := range s.Children(name) {
			walk(k)
		}
	}
	for _, root := range s.Roots() {
		walk(root)
	}
	return out
}
