package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dankoz/gt-stacks/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	if err := cli.NewRootCmd(version, commit, runtime.Version()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
