# gts — Design Spec

**Date:** 2026-05-07
**Author:** Dan Kozlowski

## 1. Summary

`gts` is a Go-based clone of the [Graphite CLI](https://graphite.com/docs/command-reference) for managing stacked git PRs. It wraps `git` and the `gh` GitHub CLI and provides:

- A CLI that mirrors graphite's command vocabulary as closely as possible.
- A standalone TUI (launched by running `gts` with no args) for browsing and operating on the stack.
- Automatic management of a "stack" footer block in each PR's description, showing every PR in the stack with the current PR highlighted.

Distribution is a single static binary via GoReleaser, available through a Homebrew tap.

## 2. Goals & Non-Goals

### Goals

- Faithful CLI parity with the **core stack workflow** (~14 commands; see §5).
- Idempotent, safe PR description updates that preserve user-authored content.
- Simple, clean, playful TUI built on Bubble Tea / Lipgloss.
- Shell out to `gh` for all GitHub interactions — no token storage in `gts`.
- Cross-platform: macOS (arm64/amd64), Linux (arm64/amd64).

### Non-Goals (v1)

- Full graphite parity (split, fold, squash, dash, ssh-key, completions for every alias). Tracked for v1.x.
- Native GitHub OAuth / token management. We require `gh auth login` to be configured.
- Mouse support in the TUI.
- GitLab / Bitbucket / other forge support.
- Auto-rollback of partially failed multi-step operations (re-run is the recovery path; ops are designed to be idempotent).

## 3. Language & Stack

**Go**, chosen over Rust because:

- 90% of work is orchestrating subprocesses (`git`, `gh`); Go's `os/exec` + goroutines fit naturally.
- Charmbracelet's Bubble Tea + Lipgloss + Bubbles is the gold standard for "playful" TUIs.
- Cobra is the de facto framework for tools shaped like graphite's CLI.
- GoReleaser provides single-config multi-platform builds + Homebrew tap.
- `gh` itself is in Go — reference patterns are directly available.

Key dependencies:

- `github.com/spf13/cobra` — CLI framework
- `github.com/charmbracelet/bubbletea` — TUI runtime
- `github.com/charmbracelet/lipgloss` — styling
- `github.com/charmbracelet/bubbles` — list, spinner, viewport, help, textinput
- `github.com/charmbracelet/x/exp/teatest` — TUI testing

## 4. Architecture

### Package Layout

```
gt-stacks/
├── cmd/gts/main.go           # entry: args → cobra, no args → TUI
├── internal/
│   ├── git/                  # git subprocess wrapper (typed)
│   ├── gh/                   # gh subprocess wrapper (typed, JSON in/out)
│   ├── state/                # stack metadata: git config + .git/gts/state.json
│   ├── pr/                   # PR description footer parsing + rendering
│   ├── core/                 # operations: Create, Modify, Submit, Sync, Restack…
│   ├── cli/                  # cobra commands (one file per command group)
│   ├── tui/                  # bubble tea program (model/update/view/keymap)
│   └── render/               # shared lipgloss styles + tree rendering
└── go.mod
```

**Layering rule:** `cli/` and `tui/` are presentation layers and call only into `core/`. `core/` calls `git/`, `gh/`, `state/`, `pr/`. The deeper packages are unaware of CLI or TUI concerns. `render/` is a leaf used by both `cli/` (for `gts log`) and `tui/`.

### Data Model

```go
// internal/state/types.go

type Stack struct {
    Trunk    string              // e.g. "main"
    Branches map[string]*Branch  // keyed by branch name
}

type Branch struct {
    Name    string
    Parent  string  // "" if parent is trunk
    PR      int     // 0 if no PR yet
    PRState string  // "open" | "merged" | "closed" | ""
    Tracked bool    // false = ignore (regular non-stacked branch)
}
```

The tree shape is implicit: walk via `Parent`. Children-of is computed on demand.

### Persistence

| What | Where | Source of truth? |
|---|---|---|
| Parent branch | `git config branch.<name>.gts-parent` | yes |
| PR number | `git config branch.<name>.gts-pr` | yes |
| Tracked flag | presence of `gts-parent` config key | yes |
| Trunk choice | `git config gts.trunk` | yes |
| PR states cache | `.git/gts/state.json` | no (rebuildable from `gh pr list`) |
| Last-fetched timestamps | `.git/gts/state.json` | no |

Putting authoritative metadata in `git config` means it lives in the repo's `.git/config`, survives normal git operations, and is inspectable via `git config --list`.

### Data Flow (example: `gts submit`)

1. CLI/TUI invokes `core.Submit(ctx, opts)`.
2. `core` reads stack via `state.Load()` and walks branches from current to trunk.
3. For each branch missing a PR: `gh.PR.Create(...)`. For each existing: refresh via `gh.PR.View(...)`.
4. `pr.RenderFooter(stack, currentBranch)` builds the markdown block.
5. For each PR whose body's managed block differs from the rendered block, `gh.PR.Edit(num, body=mergedBody)` (parallelized, bounded worker pool).
6. `state.Save()` persists any new PR numbers.

## 5. Command Surface (v1)

| Command | Aliases | Behavior |
|---|---|---|
| `gts` | — | Launch TUI |
| `gts create <name>` | `c` | Create branch as child of current; commits any staged changes |
| `gts modify` | `m` | Amend current commit OR create new commit; restack descendants |
| `gts submit` | `s` | Create/update PRs for the current stack; rewrite footers |
| `gts sync` | — | Fetch trunk, restack tracked branches, delete merged |
| `gts restack` | `r` | Replay current stack onto current parents |
| `gts continue` | — | Resume after rebase conflict |
| `gts log` | `ls`, `l` | Print stack tree |
| `gts status` | `st` | Current branch + position in stack |
| `gts up [n]` | `u` | Checkout child (interactive picker if multiple) |
| `gts down [n]` | `d` | Checkout parent |
| `gts checkout` | `co` | Interactive branch picker |
| `gts track [parent]` | — | Start tracking current branch |
| `gts untrack` | — | Stop tracking current branch |
| `gts trunk <branch>` | — | Set trunk (default: detected `main`/`master`) |
| `gts absorb` | — | **Stretch** — distribute staged hunks to right branches |

**Global flags:** `--no-color`, `--json` (where applicable), `-v/--verbose`.

`gts absorb` is marked stretch — hunk-level analysis adds significant complexity. Implemented last; may slip to v1.1.

## 6. PR Footer Mechanics

The marker block is the marquee feature.

### Markers

```
<!-- gts:stack-start -->
…managed content…
<!-- gts:stack-end -->
```

HTML comments are invisible in rendered markdown.

### Render Format

Walked root-to-leaf (trunk → leaf). Each tracked branch with a PR rendered as:

```
- [glyph] #<num> · <branch>
```

Glyphs: `✓` merged, `▶` current PR (also bold), `○` open, `✗` closed.

Example block:

```markdown
**Stack** (this PR is **▶ feat/auth-2**):

- ✓ #4520 · feat/auth-1
- **▶ #4521 · feat/auth-2**  ← you are here
- ○ #4522 · feat/auth-3
```

### Update Algorithm

1. Fetch existing body via `gh pr view <num> --json body`.
2. If start/end markers present → replace content between them. If absent → append (preceded by a blank line).
3. Render the block from the current `Stack` and the PR's branch.
4. Diff the rendered block against the existing parsed block. If unchanged, **skip the API call** to avoid noise in the PR activity feed.
5. Otherwise call `gh.PR.Edit(num, body=mergedBody)`.

### When Footers Are Rewritten

- After every `gts submit` (full stack iteration).
- After `gts sync` if any PR shape changed (parent reassignments, deletions).
- After `gts restack` if it altered which branches exist or their order.
- Never on read-only operations (`log`, `status`, `up`, `down`).

### Idempotency

Two consecutive runs with no stack change must produce zero `gh pr edit` calls. Guaranteed by the diff check in step 4.

### Concurrency

`gh pr edit` calls across a stack run in parallel with a bounded worker pool (default 4). A 6-branch stack should not take 6× single-call latency.

### Failure Handling

If `gh pr edit` fails mid-stack, per-branch outcomes are recorded and reported at the end:

```
✓ #4520 updated
✓ #4521 updated
✗ #4522 failed: rate limited
```

No automatic rollback. Re-running is the recovery path (idempotent).

## 7. TUI Design

Bubble Tea state machine. Lipgloss for styling. Bubbles for list / spinner / viewport / help / textinput.

### Entry

`gts` with no args launches the TUI. Every keystroke that mutates state calls the same `core/` function the corresponding CLI command would.

### Primary View

```
┌─ gts ─ myrepo · main ──────────────────────────────────┐
│                                                          │
│   main                                                   │
│   ├─● feat/auth-1     #4520   ✓ merged                  │
│   └─● feat/auth-2     #4521   ✓ approved                │
│     └─▶ feat/auth-3   #4522   ○ in review     ← here   │
│       ├─○ feat/oauth  (no PR yet)                        │
│       └─○ feat/sso    (no PR yet)                        │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  ↑↓  navigate    ⏎  checkout    s  submit    r  restack │
│  c   create      m  modify      ?  help      q  quit    │
└──────────────────────────────────────────────────────────┘
```

### States

- `browsing` — tree visible, key bindings active (default).
- `confirming` — modal dialog for destructive ops.
- `running` — operation in flight; spinner + streamed log lines.
- `prompting` — text input (e.g., new branch name).
- `error` — error toast at the bottom; dismiss with any key.

### Key Bindings (browsing)

| Key | Action |
|---|---|
| `↑` / `↓` / `j` / `k` | Move selection |
| `enter` | Checkout selected branch |
| `s` | Submit (full stack from current) |
| `r` | Restack |
| `y` | Sync |
| `c` | Create child of selected (prompts for name) |
| `m` | Modify (amend / new commit; opens `$EDITOR` for message) |
| `t` / `T` | Track / untrack selected |
| `?` | Help overlay |
| `q` / `ctrl-c` | Quit |

### Visual Language

- Soft rounded borders (Lipgloss `RoundedBorder`).
- Branch glyphs: `●` tracked, `○` no PR yet, `▶` current, `✓` merged.
- Restrained color palette. Trunk in muted gray. Current branch in accent (cyan default, configurable). Merged branches dimmed. Errors red. No rainbow.
- Charm-style dot-cycle spinner during ops.
- Emoji opt-in via `--emoji` (off by default).

### Async Work

Long ops (submit, sync) run in a goroutine that pushes `tea.Msg` values back to the program: progress, completion, error. UI never blocks. `q` aborts via context cancellation; in-flight `gh` calls finish or get killed via SIGTERM.

### Output Streaming

When an op runs, a collapsible viewport shows live `git`/`gh` output. On success it auto-collapses; on error it stays open with the failing line highlighted.

### Out of Scope (v1)

- Mouse support.
- Pluggable themes (we ship one, configurable accent color only).
- Inline diff preview.

## 8. Subprocess Strategy

### `internal/git/`

Typed wrapper around `git`. Single `Runner` interface:

```go
type Runner interface {
    Run(ctx context.Context, args ...string) (stdout, stderr []byte, err error)
}
```

Real impl uses `os/exec`. Test impl is a fake that returns canned output. All git interaction routes through this interface.

Methods return Go types. Examples:

```go
func (g *Git) CurrentBranch(ctx context.Context) (string, error)
func (g *Git) Branches(ctx context.Context) ([]Branch, error)
func (g *Git) RebaseOnto(ctx, upstream, branch string) error
func (g *Git) ConfigGet(ctx, key string) (string, error)
func (g *Git) ConfigSet(ctx, key, value string) error
```

Use porcelain commands (`-z`, `--porcelain`) where available so parsing is stable across git versions.

### `internal/gh/`

Same pattern. Always `--json` for reads.

```go
func (g *GH) PRView(ctx, num int) (*PR, error)         // gh pr view N --json …
func (g *GH) PRCreate(ctx, opts CreateOpts) (*PR, error) // gh pr create …
func (g *GH) PREdit(ctx, num int, opts EditOpts) error  // gh pr edit N --body -
func (g *GH) PRList(ctx, opts ListOpts) ([]PR, error)   // gh pr list --json …
```

`PRList` is used for stack-wide refresh in a single API call.

### Cancellation

Every call accepts `context.Context`. TUI quit / Ctrl-C cancels the context → child process gets SIGTERM. No zombies.

## 9. Error Handling

Three tiers:

### 1. User errors

Uncommitted changes, dirty worktree, branch already exists, not a git repo. Surfaced as a one-line message + a hint. No stack trace.

```
✗ Cannot create branch: working tree has uncommitted changes
  Hint: stash them with `git stash`, or commit first with `gts modify`
```

### 2. Subprocess failures

`git rebase` conflict, `gh` rate-limited, `gh` not authenticated. Show the underlying tool's stderr verbatim plus our context.

### 3. Bugs

Nil deref, parse failure on git output we expected. Print full stack trace + a "please file an issue" footer with version/commit/`runtime.Version()`.

### Conflicts

Rebase conflicts are surfaced as state, not error. `gts continue` resumes; `gts restack --abort` bails out. Mirrors graphite's UX.

### Atomicity

Multi-step operations (e.g., submit touches N PRs) record per-step outcomes. On failure, we print which steps succeeded and which need retry. We do not auto-rollback — re-running is the path forward (idempotent by design).

## 10. Testing Strategy

### Unit (fast, no network, no real git)

- `internal/pr/footer_test.go` — golden-file tests of footer rendering across varied stack shapes (empty, single branch, deep stack, branches without PRs, mid-stack merged).
- `internal/state/` — fake `Runner` returning canned `git config` output.
- `internal/core/` — fake `git.Runner` and `gh.Runner`; assert sequence of subprocess calls and resulting state transitions.

### Integration (real git, fake gh)

- Test harness creates a throwaway `git init` repo in `t.TempDir()`, runs real `git` commands.
- `gh` is replaced via a shim binary placed earlier in `PATH`. The shim reads request args, returns recorded JSON fixtures.
- Covers: create-stack, modify, restack, sync-with-merged-branches, footer-update happy path.

### TUI

- Bubble Tea's `teatest` package: drive a sequence of keystrokes, assert on rendered frames.

### E2E (manual, gated)

- `TestE2E_*` tests behind `-tags=e2e` that hit a real test repo and real GitHub. Run before releases, not on every CI commit.

### Coverage Target

80% on `core/`, `pr/`, `state/`. Lower on `cli/` and `tui/` (covered by integration).

## 11. Distribution

- **GoReleaser** (`.goreleaser.yaml`): builds for `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64` on every tag. Generates checksums + cosign-signed artifacts.
- **Homebrew tap**: GoReleaser auto-publishes a formula to `homebrew-gts`. Install via `brew install dankoz/gts/gts`.
- **Single static binary** — version embedded via `-ldflags="-X main.version=…"`. `gts version` prints version + commit + Go version.
- **Shell completions**: Cobra-generated. `gts completion bash|zsh|fish` ships in the binary.
- **CI**: GitHub Actions matrix across the build OSes. Runs `go test ./...`, `golangci-lint run`, `gofmt -l`. Tag push triggers GoReleaser.

## 12. Open Questions / Future Work

- **`gts absorb`** — design hunk-distribution algorithm. Look at `git absorb` (rust) for prior art.
- **Stack sharing across clones** — `git config` lives in `.git/`, not the working tree. If teammates pull a branch, they don't get its `gts-parent`. Future option: `gts export`/`gts import` writing a `.gts/stack.yaml` to the working tree.
- **Conflict resolution UX in the TUI** — v1 punts to CLI (`gts continue`); a future TUI mode could walk the user through conflicts.
- **Themes** — config file (`~/.config/gts/config.toml`) for accent color, glyphs, footer template.
- **Notifications** — desktop notification on async submit completion.
