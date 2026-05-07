package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newStatusCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Aliases: []string{"st"},
		Short:   "Show current branch and stack position",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := c.Status(cmd.Context())
			if err != nil {
				return err
			}
			cmd.Print(out)
			return nil
		},
	}
}
