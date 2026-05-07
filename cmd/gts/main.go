package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dankoz/gt-stacks/internal/cli"
	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/dankoz/gt-stacks/internal/gh"
	"github.com/dankoz/gt-stacks/internal/git"
	"github.com/dankoz/gt-stacks/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// No args (besides the program name) → launch TUI.
	if len(os.Args) == 1 {
		c := core.New(git.New(git.NewExecRunner(".")), gh.New(gh.NewExecRunner()))
		if err := tui.Run(c); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}
	if err := cli.NewRootCmd(version, commit, runtime.Version()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
