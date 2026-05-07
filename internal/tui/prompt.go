package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func newBranchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "feature/my-thing"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return ti
}
