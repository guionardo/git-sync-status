package tui

import (
	"strings"
	"testing"

	"github.com/guionardo/git_sync_status/internal/domain"
	"github.com/guionardo/git_sync_status/internal/service"
)

func TestRenderStatusCard(t *testing.T) {
	t.Parallel()

	m := Model{
		repoPath: "/tmp/repo",
		result: domain.Result{
			RepoPath: "/tmp/repo",
			Branch:   "main",
			Upstream: "origin/main",
			Status:   domain.StatusSyncPending,
			Ahead:    2,
			Behind:   0,
			Flags:    []string{"WORKTREE_DIRTY"},
			Actions:  []string{"Push local commits: git push"},
		},
	}

	out := m.renderStatusCard()
	wantContains := []string{"Path: /tmp/repo", "Status", "SYNC_PENDING", "Ahead/Behind: 2/0", "WORKTREE_DIRTY"}
	for _, want := range wantContains {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q. output: %s", want, out)
		}
	}
}

func TestRenderBranchTable(t *testing.T) {
	t.Parallel()

	m := Model{
		branches: []service.BranchStatus{
			{Branch: "main", Upstream: "origin/main", Status: domain.StatusSynced, Ahead: 0, Behind: 0},
			{Branch: "feature/a", Upstream: "origin/feature/a", Status: domain.StatusSyncPending, Ahead: 2, Behind: 0},
			{Branch: "feature/no-upstream", Status: domain.StatusNoUpstream, Ahead: 0, Behind: 0},
		},
	}

	out := m.renderAllBranchesCard()
	wantContains := []string{
		"All Branches",
		"BRANCH",
		"UPSTREAM",
		"STATUS",
		"A/B",
		"main",
		"origin/main",
		"feature/a",
		"SYNC_PENDING",
		"feature/no-upstream",
		"NO_UPSTREAM",
	}

	for _, want := range wantContains {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q. output: %s", want, out)
		}
	}
}
