package gh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
)

type Runner interface {
	Run(ctx context.Context, stdin io.Reader, args ...string) (stdout, stderr []byte, err error)
}

type execRunner struct{}

func NewExecRunner() Runner { return &execRunner{} }

func (e *execRunner) Run(ctx context.Context, stdin io.Reader, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, "gh", args...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return out.Bytes(), errb.Bytes(),
			fmt.Errorf("gh %s: %w (stderr: %s)", args[0], err, errb.String())
	}
	return out.Bytes(), errb.Bytes(), nil
}
