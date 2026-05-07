package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Runner abstracts execution of `git` so it can be faked in tests.
type Runner interface {
	Run(ctx context.Context, args ...string) (stdout, stderr []byte, err error)
}

type execRunner struct {
	dir string
}

// NewExecRunner returns a Runner that exec's the real `git` binary in dir.
func NewExecRunner(dir string) Runner {
	return &execRunner{dir: dir}
}

func (e *execRunner) Run(ctx context.Context, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = e.dir
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return out.Bytes(), errb.Bytes(),
			fmt.Errorf("git %s: %w (stderr: %s)", args[0], err, errb.String())
	}
	return out.Bytes(), errb.Bytes(), nil
}
