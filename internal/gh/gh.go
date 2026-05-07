package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type PR struct {
	Number      int    `json:"number"`
	State       string `json:"state"`
	Body        string `json:"body"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type GH struct{ r Runner }

func New(r Runner) *GH { return &GH{r: r} }

const prFields = "number,state,body,baseRefName,headRefName,title,url"

func (g *GH) PRView(ctx context.Context, number int) (*PR, error) {
	out, _, err := g.r.Run(ctx, nil, "pr", "view", strconv.Itoa(number), "--json", prFields)
	if err != nil {
		return nil, err
	}
	var pr PR
	if err := json.Unmarshal(out, &pr); err != nil {
		return nil, fmt.Errorf("decode pr view: %w", err)
	}
	return &pr, nil
}

type ListOpts struct {
	Heads []string // filter by head ref names; empty = all
	State string   // "open"|"closed"|"merged"|"all" (default: open)
}

func (g *GH) PRList(ctx context.Context, opts ListOpts) ([]PR, error) {
	args := []string{"pr", "list", "--json", prFields, "--limit", "200"}
	if opts.State != "" {
		args = append(args, "--state", opts.State)
	}
	for _, h := range opts.Heads {
		args = append(args, "--head", h)
	}
	out, _, err := g.r.Run(ctx, nil, args...)
	if err != nil {
		return nil, err
	}
	var prs []PR
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, fmt.Errorf("decode pr list: %w", err)
	}
	return prs, nil
}

type CreateOpts struct {
	Title, Body, Base, Head string
	Draft                   bool
}

var prURLRe = regexp.MustCompile(`/pull/(\d+)`)

func (g *GH) PRCreate(ctx context.Context, opts CreateOpts) (*PR, error) {
	args := []string{"pr", "create",
		"--title", opts.Title,
		"--body-file", "-",
		"--base", opts.Base,
		"--head", opts.Head,
	}
	if opts.Draft {
		args = append(args, "--draft")
	}
	out, _, err := g.r.Run(ctx, strings.NewReader(opts.Body), args...)
	if err != nil {
		return nil, err
	}
	m := prURLRe.FindStringSubmatch(string(out))
	if m == nil {
		return nil, fmt.Errorf("could not parse PR number from gh output: %q", string(out))
	}
	n, _ := strconv.Atoi(m[1])
	return &PR{
		Number:      n,
		URL:         strings.TrimSpace(string(out)),
		HeadRefName: opts.Head,
		BaseRefName: opts.Base,
	}, nil
}

type EditOpts struct {
	Body *string // nil = leave unchanged
	Base *string // nil = leave unchanged
}

func (g *GH) PREdit(ctx context.Context, number int, opts EditOpts) error {
	args := []string{"pr", "edit", strconv.Itoa(number)}
	var in io.Reader
	if opts.Body != nil {
		args = append(args, "--body-file", "-")
		in = strings.NewReader(*opts.Body)
	}
	if opts.Base != nil {
		args = append(args, "--base", *opts.Base)
	}
	if len(args) == 3 {
		return nil // nothing to do
	}
	_, _, err := g.r.Run(ctx, in, args...)
	return err
}
