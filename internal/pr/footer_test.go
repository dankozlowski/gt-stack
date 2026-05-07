package pr

import "testing"

func TestParseBlock_Present(t *testing.T) {
	body := "Hello\n\n<!-- gts:stack-start -->\nold content\n<!-- gts:stack-end -->\n\nbye"
	got, ok := ParseBlock(body)
	if !ok {
		t.Fatalf("expected to find block")
	}
	if got != "old content\n" {
		t.Errorf("got %q, want %q", got, "old content\n")
	}
}

func TestParseBlock_Absent(t *testing.T) {
	body := "Hello\nNo markers here"
	if _, ok := ParseBlock(body); ok {
		t.Errorf("expected no block")
	}
}

func TestReplaceOrAppend_AppendsWhenAbsent(t *testing.T) {
	body := "Summary\n\nDetails"
	got := ReplaceOrAppend(body, "RENDERED")
	want := "Summary\n\nDetails\n\n<!-- gts:stack-start -->\nRENDERED\n<!-- gts:stack-end -->\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestReplaceOrAppend_ReplacesWhenPresent(t *testing.T) {
	body := "Hi\n<!-- gts:stack-start -->\nOLD\n<!-- gts:stack-end -->\n"
	got := ReplaceOrAppend(body, "NEW")
	want := "Hi\n<!-- gts:stack-start -->\nNEW\n<!-- gts:stack-end -->\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
