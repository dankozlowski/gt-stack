package gh

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type FakeResponse struct {
	Stdout, Stderr []byte
	ExitCode       int
}

type FakeCall struct {
	Args  []string
	Stdin string
}

type FakeRunner struct {
	Responses map[string]FakeResponse
	Calls     []FakeCall
}

func (f *FakeRunner) Run(ctx context.Context, stdin io.Reader, args ...string) ([]byte, []byte, error) {
	var sin string
	if stdin != nil {
		b, _ := io.ReadAll(stdin)
		sin = string(b)
	}
	f.Calls = append(f.Calls, FakeCall{Args: append([]string(nil), args...), Stdin: sin})
	key := strings.Join(args, " ")
	resp, ok := f.Responses[key]
	if !ok {
		return nil, nil, fmt.Errorf("FakeRunner: no response for %q", key)
	}
	if resp.ExitCode != 0 {
		return resp.Stdout, resp.Stderr, fmt.Errorf("gh exit %d", resp.ExitCode)
	}
	return resp.Stdout, resp.Stderr, nil
}
