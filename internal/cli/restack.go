package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newRestackCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "restack",
		Aliases: []string{"r"},
		Short:   "Replay descendants onto current parents",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Restack(cmd.Context())
		},
	}
}
