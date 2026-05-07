package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newContinueCmd(c *core.Core) *cobra.Command {
	var abort bool
	cmd := &cobra.Command{
		Use:   "continue",
		Short: "Resume a paused rebase (or --abort)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if abort {
				return c.Abort(cmd.Context())
			}
			return c.Continue(cmd.Context())
		},
	}
	cmd.Flags().BoolVar(&abort, "abort", false, "Abort the in-progress rebase")
	return cmd
}
