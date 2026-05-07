package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *core.Core) *cobra.Command {
	var msg string
	cmd := &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c"},
		Short:   "Create a new branch as a child of the current branch",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Create(cmd.Context(), args[0], msg)
		},
	}
	cmd.Flags().StringVarP(&msg, "message", "m", "", "Commit message for staged changes")
	return cmd
}
