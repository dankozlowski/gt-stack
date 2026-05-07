package cli

import (
	"fmt"

	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newTrunkCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "trunk [branch]",
		Short: "Show or set the trunk branch",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) == 0 {
				t, err := c.GetTrunk(ctx)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), t)
				return nil
			}
			return c.SetTrunk(ctx, args[0])
		},
	}
}
