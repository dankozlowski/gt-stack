package git

import (
	"context"
	"fmt"
	"strings"
)

type FakeResponse struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

// FakeRunner returns canned responses keyed by space-joined args.
// Calls are recorded for assertion in tests.
type FakeRunner struct {
	Responses map[string]FakeResponse
	Calls     [][]string
}

func (f *FakeRunner) Run(ctx context.Context, args ...string) ([]byte, []byte, error) {
	f.Calls = append(f.Calls, append([]string(nil), args...))
	key := strings.Join(args, " ")
	resp, ok := f.Responses[key]
	if !ok {
		return nil, nil, fmt.Errorf("FakeRunner: no response configured for %q", key)
	}
	if resp.ExitCode != 0 {
		return resp.Stdout, resp.Stderr, fmt.Errorf("git exit %d", resp.ExitCode)
	}
	return resp.Stdout, resp.Stderr, nil
}
