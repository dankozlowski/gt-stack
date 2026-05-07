package cli

import (
	"fmt"

	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func newSubmitCmd(c *core.Core) *cobra.Command {
	var opts core.SubmitOpts
	cmd := &cobra.Command{
		Use:     "submit",
		Aliases: []string{"s"},
		Short:   "Create or update PRs for the current stack",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rep, err := c.Submit(cmd.Context(), opts)
			out := cmd.OutOrStdout()
			if rep != nil {
				for _, o := range rep.Created {
					fmt.Fprintf(out, "✓ created #%d (%s)\n", o.PR, o.Branch)
				}
				for _, o := range rep.Updated {
					fmt.Fprintf(out, "✓ updated #%d (%s)\n", o.PR, o.Branch)
				}
				for _, o := range rep.Skipped {
					fmt.Fprintf(out, "· #%d (%s) unchanged\n", o.PR, o.Branch)
				}
				for _, o := range rep.Failed {
					fmt.Fprintf(out, "✗ #%d (%s): %v\n", o.PR, o.Branch, o.Err)
				}
			}
			return err
		},
	}
	cmd.Flags().BoolVar(&opts.Draft, "draft", false, "Open new PRs as draft")
	return cmd
}
