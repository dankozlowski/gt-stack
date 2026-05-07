package cli

import (
	"fmt"

	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newSyncCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Fetch trunk, prune merged branches, restack survivors",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rep, err := c.Sync(cmd.Context())
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if rep.Fetched {
				fmt.Fprintln(out, "✓ fetched origin")
			}
			for _, d := range rep.Deleted {
				fmt.Fprintf(out, "✓ deleted %s (merged)\n", d)
			}
			for _, r := range rep.Restacked {
				fmt.Fprintf(out, "✓ restacked %s\n", r)
			}
			return nil
		},
	}
}
