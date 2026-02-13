# GIT SYNC STATUS

Tool to verify whether a folder is a Git repository and report sync status between local and remote branches.

## Scope

- Check current branch by default
- Optional: check all local tracking branches

## Statuses

- `NOT_A_GIT_REPO`
  - Folder is not a Git repository.

- `NO_REMOTE`
  - Repository exists but no remote is configured.

- `NO_UPSTREAM`
  - Current branch exists but does not track a remote branch.
  - Optional check: verify whether this branch was already merged into the default branch (for example, `main` or `master`).
  - If already merged, suggest removing it from the local repository.

- `SYNCED`
  - Local and upstream point to the same commit (`ahead=0`, `behind=0`)
  - Working tree is clean (no staged, unstaged, or untracked files)

- `SYNC_PENDING`
  - Local branch is ahead of upstream (`ahead>0`, `behind=0`)

- `LATE`
  - Local branch is behind upstream (`ahead=0`, `behind>0`)

- `DIVERGED`
  - Local and upstream both have unique commits (`ahead>0`, `behind>0`)

## Working tree flags

- `WORKTREE_DIRTY`: staged, unstaged, or untracked files exist
- `DETACHED_HEAD`: repository is not on a branch
- `REMOTE_UNREACHABLE`: failed to fetch or compare remote (network or auth issue)

## Output recommendation

Always show:

- repository path
- current branch
- upstream branch (if any)
- ahead and behind counts
- sync status
- working tree flags

## Git command by status

- `NOT_A_GIT_REPO`
  - `git rev-parse --is-inside-work-tree`

- `NO_REMOTE`
  - `git remote`

- `NO_UPSTREAM`
  - `git rev-parse --abbrev-ref --symbolic-full-name @{u}`
  - Optional merged check with default branch:
    - `git branch --merged main`
    - `git branch --merged master`
  - Suggested local removal (when merged):
    - `git branch -d <branch>`

- `SYNCED`
  - `git fetch --prune`
  - `git rev-list --left-right --count @{u}...HEAD`
  - `git status --porcelain`

- `SYNC_PENDING`
  - `git fetch --prune`
  - `git rev-list --left-right --count @{u}...HEAD`

- `LATE`
  - `git fetch --prune`
  - `git rev-list --left-right --count @{u}...HEAD`

- `DIVERGED`
  - `git fetch --prune`
  - `git rev-list --left-right --count @{u}...HEAD`

## Working tree flag commands

- `WORKTREE_DIRTY`
  - `git status --porcelain`

- `DETACHED_HEAD`
  - `git symbolic-ref --quiet --short HEAD`

- `REMOTE_UNREACHABLE`
  - `git ls-remote --heads origin`

## Tiny parser for ahead and behind

Command:

- `git rev-list --left-right --count @{u}...HEAD`

Output format:

- `<behind> <ahead>`
  - first number = commits present in upstream and missing locally (`behind`)
  - second number = commits present locally and missing in upstream (`ahead`)

Examples:

- `0 0` -> `SYNCED` (if working tree is clean)
- `0 3` -> `SYNC_PENDING`
- `2 0` -> `LATE`
- `2 3` -> `DIVERGED`

## Development

### Build and run

- Build binary:
  - `make build`
- Run TUI (default):
  - `make run`
  - or `go run ./cmd/git-sync-status --path /path/to/repo`
- Plain text output:
  - `go run ./cmd/git-sync-status --plain --path /path/to/repo`
- JSON output:
  - `go run ./cmd/git-sync-status --json --path /path/to/repo`
- List local branches:
  - `go run ./cmd/git-sync-status --list-branches --path /path/to/repo`

### TUI keybinds

- `r`: refresh status
- `q`: quit

### Test and quality

- Run tests:
  - `make test`
- Run vet:
  - `make vet`
- Run race detector:
  - `make race`
- Run linter (requires `golangci-lint` installed):
  - `make lint`

## Release

- Local cross-platform archives are configured with GoReleaser:
  - `.goreleaser.yaml`
- Example local dry run:
  - `goreleaser release --snapshot --clean`
