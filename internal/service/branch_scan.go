package service

import (
	"context"

	"github.com/guionardo/git_sync_status/internal/domain"
)

type BranchSummary struct {
	Name string
}

type BranchStatus struct {
	Branch   string
	Upstream string
	Status   domain.Status
	Behind   int
	Ahead    int
	Flags    []string
}

func (a *Analyzer) ScanLocalBranches(ctx context.Context, repoPath string) ([]BranchSummary, error) {
	branches, err := a.client.LocalBranches(ctx, repoPath)
	if err != nil {
		return nil, err
	}

	out := make([]BranchSummary, 0, len(branches))
	for _, b := range branches {
		out = append(out, BranchSummary{Name: b})
	}
	return out, nil
}

func (a *Analyzer) AnalyzeAllBranches(ctx context.Context, repoPath string) ([]BranchStatus, error) {
	branches, err := a.client.LocalBranches(ctx, repoPath)
	if err != nil {
		return nil, err
	}

	rows := make([]BranchStatus, 0, len(branches))
	for _, branch := range branches {
		row := BranchStatus{
			Branch: branch,
			Status: domain.StatusNoUpstream,
		}

		upstream, err := a.client.UpstreamForBranch(ctx, repoPath, branch)
		if err != nil {
			rows = append(rows, row)
			continue
		}
		row.Upstream = upstream

		behind, ahead, err := a.client.AheadBehindRefs(ctx, repoPath, upstream, branch)
		if err != nil {
			row.Flags = append(row.Flags, "REMOTE_UNREACHABLE")
			rows = append(rows, row)
			continue
		}
		row.Behind = behind
		row.Ahead = ahead

		switch {
		case row.Behind > 0 && row.Ahead > 0:
			row.Status = domain.StatusDiverged
		case row.Behind > 0:
			row.Status = domain.StatusLate
		case row.Ahead > 0:
			row.Status = domain.StatusSyncPending
		default:
			row.Status = domain.StatusSynced
		}

		rows = append(rows, row)
	}

	return rows, nil
}
