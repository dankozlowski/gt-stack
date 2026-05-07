package cli

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	var out bytes.Buffer
	root := NewRootCmd("0.0.1-test", "abc123", "go1.22")
	root.SetOut(&out)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	got := out.String()
	for _, want := range []string{"0.0.1-test", "abc123", "go1.22"} {
		if !bytes.Contains(out.Bytes(), []byte(want)) {
			t.Errorf("version output missing %q\ngot:\n%s", want, got)
		}
	}
}
