# gts

A Graphite-style CLI for managing stacked git pull requests, plus a playful TUI. Wraps `git` and `gh` — no GitHub token to manage, no daemon, no remote service.

```
┌─ gts ─ myrepo · main ──────────────────────────────────-┐
│   main                                                  │
│   ├─● feat/auth-1     #4520   ✓ merged                  │
│   └─● feat/auth-2     #4521   ✓ approved                │
│     └─▶ feat/auth-3   #4522   ○ in review     ← here    │
│       └─○ feat/oauth  (no PR yet)                       │
│                                                         │
│  [↑↓] nav  [⏎] checkout  [s] submit  [r] restack  [q]   │
└─────────────────────────────────────────────────────────┘
```

## Features

- **Stacked PRs**: track parent/child relationships between branches; `gts up`/`gts down` walk the chain.
- **Auto-managed PR descriptions**: every PR in a stack gets a footer block listing every other PR in the stack, with the current PR highlighted. Updates are idempotent — a no-op run produces zero `gh pr edit` calls.
- **One-shot stack ops**: `gts sync` fetches trunk, prunes merged branches, restacks survivors. `gts submit` opens or updates every PR in your current stack in parallel.
- **Conflict-aware restack**: rebase conflicts are surfaced as state, not error — `gts continue` resumes after you resolve.
- **Playful TUI**: run `gts` with no arguments to launch a Bubble Tea stack browser.

## Requirements

- `git` (any reasonably recent version)
- [`gh`](https://cli.github.com) authenticated via `gh auth login`
- Go 1.25+ to build from source

## Install

### Homebrew (once tagged)

```sh
brew install dankoz/gts/gts
```

### From source

```sh
go install github.com/dankoz/gt-stacks/cmd/gts@latest
```

### From clone

```sh
git clone https://github.com/dankoz/gt-stacks.git
cd gt-stacks
make build
./gts version
```

## Quick start

```sh
# In any git repo with `gh auth login` configured:
gts trunk main                     # set the trunk branch (one-time)
git checkout -b feat/auth-1        # start a feature branch the normal way
gts track                          # tell gts this branch is part of a stack

# Make some commits, then stack a child branch:
gts create feat/auth-2 -m "more"   # creates child off feat/auth-1, commits staged changes

# When you're ready to open PRs for the whole chain:
gts submit                         # creates PRs, writes the stack footer

# Edit the parent and propagate:
gts down                           # checkout feat/auth-1
gts modify --amend                 # amend the commit
gts restack                        # replay descendants on top
gts submit                         # update PRs (idempotent — only edits what changed)

# Trunk moved? Branches merged? One command:
gts sync                           # fetch, prune merged branches, restack survivors
```

Run `gts` with no arguments to drop into the TUI.

## Command reference

| Command | Aliases | Description |
| --- | --- | --- |
| `gts` | — | Launch TUI |
| `gts trunk [branch]` | — | Show or set the trunk branch (default: `main`) |
| `gts track [parent]` | — | Track current branch as a child of `[parent]` (default: trunk) |
| `gts untrack` | — | Stop tracking the current branch |
| `gts create <name>` | `c` | Create a new branch as a child of the current branch (use `-m` to commit staged changes) |
| `gts modify` | `m` | Amend the current commit (`--amend`) or create a new one (`-m <msg>`) |
| `gts log` | `ls`, `l` | Print the stack tree |
| `gts status` | `st` | Show current branch and stack position |
| `gts up [n]` | `u` | Checkout child (interactive picker on multiple) |
| `gts down [n]` | `d` | Checkout parent |
| `gts checkout [branch]` | `co` | Checkout a branch (interactive picker if no arg) |
| `gts restack` | `r` | Replay descendants onto current parents |
| `gts continue` | — | Resume a paused rebase (or `--abort`) |
| `gts sync` | — | Fetch trunk, prune merged branches, restack survivors |
| `gts submit` | `s` | Create or update PRs for the current stack |
| `gts version` | — | Print version information |

## How it works

**State storage**: Stack metadata lives in `git config`:

- `branch.<name>.gts-parent` — the parent branch
- `branch.<name>.gts-pr` — PR number once known
- `gts.trunk` — the trunk branch

This is inspectable via `git config --list` and survives normal git operations. A small JSON cache lives in `.git/gts/state.json` for transient bookkeeping.

**PR description footer**: When `gts submit` runs, it inserts (or rewrites) a managed block in every PR in the stack:

```
<!-- gts:stack-start -->
**Stack** (this PR is **▶ feat/auth-2**):

- ✓ #4520 · feat/auth-1
- **▶ #4521 · feat/auth-2**  ← you are here
- ○ #4522 · feat/auth-3
<!-- gts:stack-end -->
```

Anything you wrote outside the markers is preserved exactly. Re-running `gts submit` with no stack changes results in **zero** `gh pr edit` calls.

**No daemon, no token storage**: `gts` shells out to `gh` for every GitHub operation. If `gh` is authenticated, `gts` is authenticated. Nothing else to configure.

## Development

```sh
make build       # build ./gts
make test        # go test -race -count=1 ./...
make lint        # golangci-lint run
make fmt         # gofmt -s -w .
```

The codebase is laid out:

```
cmd/gts/         # entry: with args → cobra; no args → TUI
internal/git/    # typed wrapper around `git` subprocess
internal/gh/     # typed wrapper around `gh` subprocess
internal/state/  # stack metadata: load from git config + JSON cache
internal/pr/     # PR footer marker parsing + idempotent rewrite
internal/core/   # operations called by both CLI and TUI
internal/cli/    # cobra commands
internal/tui/    # Bubble Tea program
internal/render/ # shared tree rendering (used by `gts log` and TUI)
```

Tests use a `FakeRunner` in place of real `git`/`gh` invocations — the suite runs in well under a second and never touches the network.

## Status

Early. v1 covers the core stack workflow (~14 commands). `gts absorb` is intentionally deferred. Full graphite parity (split, fold, squash, dash, etc.) is tracked for v1.x.

## Acknowledgements

Inspired by the [Graphite CLI](https://graphite.com/docs/command-reference). Built with [Cobra](https://github.com/spf13/cobra) and [Charmbracelet](https://charm.sh)'s Bubble Tea / Lipgloss / Bubbles.

## License

MIT
