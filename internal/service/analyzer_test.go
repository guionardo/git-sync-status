package service

import (
	"context"
	"errors"
	"testing"
)

type fakeClient struct {
	isRepo           bool
	currentBranch    string
	detached         bool
	hasRemote        bool
	reachable        bool
	upstream         string
	upstreamErr      error
	fetchErr         error
	behind           int
	ahead            int
	aheadBehindErr   error
	dirty            bool
	defaultBranch    string
	defaultBranchErr error
	merged           bool
	mergedErr        error
	branches         []string
}

func (f *fakeClient) IsGitRepo(context.Context, string) (bool, error) { return f.isRepo, nil }
func (f *fakeClient) CurrentBranch(context.Context, string) (string, error) {
	return f.currentBranch, nil
}
func (f *fakeClient) IsDetachedHead(context.Context, string) (bool, error) { return f.detached, nil }
func (f *fakeClient) HasRemote(context.Context, string, string) (bool, error) {
	return f.hasRemote, nil
}
func (f *fakeClient) RemoteReachable(context.Context, string, string) (bool, error) {
	return f.reachable, nil
}
func (f *fakeClient) Upstream(context.Context, string) (string, error) {
	return f.upstream, f.upstreamErr
}
func (f *fakeClient) UpstreamForBranch(context.Context, string, string) (string, error) {
	return f.upstream, f.upstreamErr
}
func (f *fakeClient) FetchPrune(context.Context, string, string) error { return f.fetchErr }
func (f *fakeClient) AheadBehind(context.Context, string) (int, int, error) {
	return f.behind, f.ahead, f.aheadBehindErr
}
func (f *fakeClient) AheadBehindRefs(context.Context, string, string, string) (int, int, error) {
	return f.behind, f.ahead, f.aheadBehindErr
}
func (f *fakeClient) IsWorktreeDirty(context.Context, string) (bool, error) { return f.dirty, nil }
func (f *fakeClient) DefaultBranch(context.Context, string, string) (string, error) {
	return f.defaultBranch, f.defaultBranchErr
}
func (f *fakeClient) IsBranchMergedInto(context.Context, string, string, string) (bool, error) {
	return f.merged, f.mergedErr
}
func (f *fakeClient) LocalBranches(context.Context, string) ([]string, error) { return f.branches, nil }

func TestAnalyzerStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fc   *fakeClient
		want string
	}{
		{
			name: "not git repo",
			fc:   &fakeClient{isRepo: false},
			want: "NOT_A_GIT_REPO",
		},
		{
			name: "no remote",
			fc: &fakeClient{
				isRepo: true, currentBranch: "main", hasRemote: false,
			},
			want: "NO_REMOTE",
		},
		{
			name: "no upstream",
			fc: &fakeClient{
				isRepo: true, currentBranch: "feature", hasRemote: true, reachable: true,
				upstreamErr: errors.New("has no upstream branch"), defaultBranch: "main", merged: true,
			},
			want: "NO_UPSTREAM",
		},
		{
			name: "synced",
			fc: &fakeClient{
				isRepo: true, currentBranch: "main", hasRemote: true, reachable: true,
				upstream: "origin/main", behind: 0, ahead: 0,
			},
			want: "SYNCED",
		},
		{
			name: "sync pending",
			fc: &fakeClient{
				isRepo: true, currentBranch: "main", hasRemote: true, reachable: true,
				upstream: "origin/main", behind: 0, ahead: 2,
			},
			want: "SYNC_PENDING",
		},
		{
			name: "late",
			fc: &fakeClient{
				isRepo: true, currentBranch: "main", hasRemote: true, reachable: true,
				upstream: "origin/main", behind: 3, ahead: 0,
			},
			want: "LATE",
		},
		{
			name: "diverged",
			fc: &fakeClient{
				isRepo: true, currentBranch: "main", hasRemote: true, reachable: true,
				upstream: "origin/main", behind: 1, ahead: 1,
			},
			want: "DIVERGED",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			analyzer := NewAnalyzer(tc.fc, "origin")
			got := analyzer.Analyze(context.Background(), "/tmp/repo")
			if string(got.Status) != tc.want {
				t.Fatalf("got status %s, want %s", got.Status, tc.want)
			}
		})
	}
}

func TestAnalyzerAddsDirtyFlag(t *testing.T) {
	t.Parallel()
	fc := &fakeClient{
		isRepo: true, currentBranch: "main", hasRemote: true, reachable: true,
		upstream: "origin/main", behind: 0, ahead: 0, dirty: true,
	}
	analyzer := NewAnalyzer(fc, "origin")
	got := analyzer.Analyze(context.Background(), "/tmp/repo")
	if !got.HasFlag("WORKTREE_DIRTY") {
		t.Fatalf("expected WORKTREE_DIRTY flag")
	}
}
