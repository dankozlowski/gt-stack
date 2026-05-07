package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newTrackCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "track [parent]",
		Short: "Track current branch as a child of [parent] (default: trunk)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parent := ""
			if len(args) == 1 {
				parent = args[0]
			}
			return c.Track(cmd.Context(), parent)
		},
	}
}

func newUntrackCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "untrack",
		Short: "Stop tracking current branch",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Untrack(cmd.Context())
		},
	}
}
