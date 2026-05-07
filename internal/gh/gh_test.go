package gh

import (
	"context"
	"testing"
)

func TestPRView_ParsesJSON(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"pr view 4521 --json number,state,body,baseRefName,headRefName,title,url": {
				Stdout: []byte(`{
					"number": 4521,
					"state": "OPEN",
					"body": "Hello",
					"baseRefName": "feat/auth-1",
					"headRefName": "feat/auth-2",
					"title": "Auth: device flow",
					"url": "https://github.com/o/r/pull/4521"
				}`),
			},
		},
	}
	g := New(fr)
	pr, err := g.PRView(context.Background(), 4521)
	if err != nil {
		t.Fatalf("PRView: %v", err)
	}
	if pr.Number != 4521 || pr.State != "OPEN" || pr.HeadRefName != "feat/auth-2" {
		t.Errorf("unexpected PR: %+v", pr)
	}
}

func TestPRCreate_ArgsAndParsing(t *testing.T) {
	fr := &FakeRunner{
		Responses: map[string]FakeResponse{
			"pr create --title T --body-file - --base main --head feat/x": {
				Stdout: []byte("https://github.com/o/r/pull/99\n"),
			},
		},
	}
	g := New(fr)
	pr, err := g.PRCreate(context.Background(), CreateOpts{
		Title: "T", Body: "B", Base: "main", Head: "feat/x",
	})
	if err != nil {
		t.Fatalf("PRCreate: %v", err)
	}
	if pr.Number != 99 {
		t.Errorf("got PR #%d, want 99", pr.Number)
	}
}
