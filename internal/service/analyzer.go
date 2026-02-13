package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/guionardo/git_sync_status/internal/domain"
	"github.com/guionardo/git_sync_status/internal/gitclient"
)

type Analyzer struct {
	client gitclient.Client
	remote string
}

func NewAnalyzer(client gitclient.Client, remote string) *Analyzer {
	if remote == "" {
		remote = "origin"
	}
	return &Analyzer{client: client, remote: remote}
}

func (a *Analyzer) Analyze(ctx context.Context, repoPath string) domain.Result {
	result := domain.Result{RepoPath: repoPath}

	isRepo, err := a.client.IsGitRepo(ctx, repoPath)
	if err != nil {
		result.Status = domain.StatusNotAGitRepo
		result.Err = err.Error()
		result.Actions = []string{"Initialize repository with: git init"}
		return result
	}
	if !isRepo {
		result.Status = domain.StatusNotAGitRepo
		result.Actions = []string{"Initialize repository with: git init"}
		return result
	}

	detached, _ := a.client.IsDetachedHead(ctx, repoPath)
	if detached {
		result.Flags = append(result.Flags, "DETACHED_HEAD")
		result.Details = append(result.Details, "HEAD is detached")
	}

	branch, err := a.client.CurrentBranch(ctx, repoPath)
	if err == nil {
		result.Branch = branch
	}
	if result.Branch == "" && detached {
		result.Branch = "(detached)"
	}

	hasRemote, err := a.client.HasRemote(ctx, repoPath, a.remote)
	if err != nil {
		result.Err = err.Error()
	}
	if !hasRemote {
		result.Status = domain.StatusNoRemote
		result.Actions = []string{
			fmt.Sprintf("Add remote: git remote add %s <url>", a.remote),
		}
		a.enrichWorktreeState(ctx, repoPath, &result)
		return result
	}

	reachable, _ := a.client.RemoteReachable(ctx, repoPath, a.remote)
	if !reachable {
		result.Flags = append(result.Flags, "REMOTE_UNREACHABLE")
		result.Details = append(result.Details, fmt.Sprintf("Remote %q is unreachable", a.remote))
	}

	upstream, err := a.client.Upstream(ctx, repoPath)
	if err != nil && isNoUpstreamErr(err) {
		result.Status = domain.StatusNoUpstream
		result.Actions = []string{
			fmt.Sprintf("Set upstream: git push -u %s %s", a.remote, fallbackBranch(result.Branch)),
		}
		a.enrichNoUpstreamHints(ctx, repoPath, &result)
		a.enrichWorktreeState(ctx, repoPath, &result)
		return result
	}
	if err != nil {
		result.Status = domain.StatusNoUpstream
		result.Err = err.Error()
		result.Actions = []string{
			fmt.Sprintf("Set upstream: git push -u %s %s", a.remote, fallbackBranch(result.Branch)),
		}
		a.enrichWorktreeState(ctx, repoPath, &result)
		return result
	}
	result.Upstream = upstream

	if err := a.client.FetchPrune(ctx, repoPath, a.remote); err != nil {
		result.Flags = append(result.Flags, "REMOTE_UNREACHABLE")
		result.Details = append(result.Details, "Fetch failed; ahead/behind may be stale")
	}

	behind, ahead, err := a.client.AheadBehind(ctx, repoPath)
	if err != nil {
		result.Err = err.Error()
		result.Details = append(result.Details, "Could not compute ahead/behind")
	} else {
		result.Behind = behind
		result.Ahead = ahead
	}

	switch {
	case result.Behind > 0 && result.Ahead > 0:
		result.Status = domain.StatusDiverged
		result.Actions = []string{
			"Review incoming changes: git pull --rebase",
			"Resolve conflicts if needed, then push",
		}
	case result.Behind > 0:
		result.Status = domain.StatusLate
		result.Actions = []string{"Update branch: git pull --rebase"}
	case result.Ahead > 0:
		result.Status = domain.StatusSyncPending
		result.Actions = []string{"Push local commits: git push"}
	default:
		result.Status = domain.StatusSynced
		result.Actions = []string{"No sync action required"}
	}

	a.enrichWorktreeState(ctx, repoPath, &result)
	if result.HasFlag("WORKTREE_DIRTY") {
		result.Actions = append(result.Actions, "Review local changes: git status")
	}

	return result
}

func (a *Analyzer) enrichNoUpstreamHints(ctx context.Context, repoPath string, result *domain.Result) {
	branch := fallbackBranch(result.Branch)
	base, err := a.client.DefaultBranch(ctx, repoPath, a.remote)
	if err != nil || base == "" {
		base = "main"
	}

	merged, err := a.client.IsBranchMergedInto(ctx, repoPath, branch, base)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("Merged-branch check failed: %v", err))
		return
	}
	if merged {
		result.NOUpstreamWasMerged = true
		result.NOUpstreamMergeBase = base
		result.NOUpstreamSuggestion = fmt.Sprintf("Branch %q appears merged into %q; consider deleting it locally: git branch -d %s", branch, base, branch)
		result.Actions = append(result.Actions, result.NOUpstreamSuggestion)
	}
}

func (a *Analyzer) enrichWorktreeState(ctx context.Context, repoPath string, result *domain.Result) {
	dirty, err := a.client.IsWorktreeDirty(ctx, repoPath)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("Could not inspect worktree: %v", err))
		return
	}
	if dirty {
		result.Flags = append(result.Flags, "WORKTREE_DIRTY")
		result.Details = append(result.Details, "Working tree has staged, unstaged, or untracked changes")
	}
}

func isNoUpstreamErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "no upstream") ||
		strings.Contains(msg, "has no upstream branch") ||
		strings.Contains(msg, "not configured for branch")
}

func fallbackBranch(branch string) string {
	if strings.TrimSpace(branch) == "" || branch == "(detached)" {
		return "HEAD"
	}
	return branch
}
