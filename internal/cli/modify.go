package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newModifyCmd(c *core.Core) *cobra.Command {
	var opts core.ModifyOpts
	cmd := &cobra.Command{
		Use:     "modify",
		Aliases: []string{"m"},
		Short:   "Amend the current commit, or create a new one",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Modify(cmd.Context(), opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.Amend, "amend", "a", false, "Amend the previous commit")
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Commit message")
	return cmd
}
