package tui

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
)

func newTestCore(t *testing.T) *core.Core {
	t.Helper()
	gitFR := &git.FakeRunner{
		Responses: map[string]git.FakeResponse{
			"symbolic-ref --short HEAD":                          {Stdout: []byte("feat/a\n")},
			"config --get gts.trunk":                             {Stdout: []byte("main\n")},
			"for-each-ref --format=%(refname:short) refs/heads/": {Stdout: []byte("main\nfeat/a\n")},
			"config --get branch.main.gts-parent":                {ExitCode: 1},
			"config --get branch.feat/a.gts-parent":              {Stdout: []byte("main\n")},
			"config --get branch.feat/a.gts-pr":                  {ExitCode: 1},
		},
	}
	return core.New(git.New(gitFR), gh.New(&gh.FakeRunner{}))
}

func TestTUI_RendersStackOnLoad(t *testing.T) {
	m := New(newTestCore(t))
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	teatest.WaitFor(t, tm.Output(),
		func(out []byte) bool { return bytes.Contains(out, []byte("feat/a")) },
		teatest.WithDuration(2*time.Second),
	)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t)
}
