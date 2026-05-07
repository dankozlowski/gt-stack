package cli

import (
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
	"github.com/spf13/cobra"
)

func NewRootCmd(version, commit, goVersion string) *cobra.Command {
	root := &cobra.Command{
		Use:           "gts",
		Short:         "Stacked PR workflow on top of git and gh",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newVersionCmd(version, commit, goVersion))

	// Lazy core: real runners only constructed when a subcommand actually runs.
	mkCore := func() *core.Core {
		return core.New(git.New(git.NewExecRunner(".")), gh.New(gh.NewExecRunner()))
	}
	addStackCommands(root, mkCore)
	return root
}

func addStackCommands(root *cobra.Command, mkCore func() *core.Core) {
	root.AddCommand(newTrunkCmd(mkCore()))
	root.AddCommand(newTrackCmd(mkCore()))
	root.AddCommand(newUntrackCmd(mkCore()))
	root.AddCommand(newLogCmd(mkCore()))
	root.AddCommand(newStatusCmd(mkCore()))
	root.AddCommand(newUpCmd(mkCore()))
	root.AddCommand(newDownCmd(mkCore()))
	root.AddCommand(newCheckoutCmd(mkCore()))
	root.AddCommand(newCreateCmd(mkCore()))
	root.AddCommand(newModifyCmd(mkCore()))
	root.AddCommand(newRestackCmd(mkCore()))
	root.AddCommand(newContinueCmd(mkCore()))
}
