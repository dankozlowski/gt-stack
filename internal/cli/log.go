package cli

import (
	"os"

	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newLogCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "log",
		Aliases: []string{"ls", "l"},
		Short:   "Print stack tree",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			color := term.IsTerminal(int(os.Stdout.Fd()))
			out, err := c.LogTree(cmd.Context(), color)
			if err != nil {
				return err
			}
			cmd.Print(out)
			return nil
		},
	}
}
