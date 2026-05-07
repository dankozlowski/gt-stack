package cli

import "github.com/spf13/cobra"

func NewRootCmd(version, commit, goVersion string) *cobra.Command {
	root := &cobra.Command{
		Use:           "gts",
		Short:         "Stacked PR workflow on top of git and gh",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newVersionCmd(version, commit, goVersion))
	return root
}
