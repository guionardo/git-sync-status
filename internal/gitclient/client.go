package gitclient

import "context"

type Client interface {
	IsGitRepo(ctx context.Context, path string) (bool, error)
	CurrentBranch(ctx context.Context, path string) (string, error)
	IsDetachedHead(ctx context.Context, path string) (bool, error)
	HasRemote(ctx context.Context, path string, remote string) (bool, error)
	RemoteReachable(ctx context.Context, path string, remote string) (bool, error)
	Upstream(ctx context.Context, path string) (string, error)
	UpstreamForBranch(ctx context.Context, path string, branch string) (string, error)
	FetchPrune(ctx context.Context, path string, remote string) error
	AheadBehind(ctx context.Context, path string) (behind int, ahead int, err error)
	AheadBehindRefs(ctx context.Context, path string, leftRef string, rightRef string) (behind int, ahead int, err error)
	IsWorktreeDirty(ctx context.Context, path string) (bool, error)
	DefaultBranch(ctx context.Context, path string, remote string) (string, error)
	IsBranchMergedInto(ctx context.Context, path string, branch string, base string) (bool, error)
	LocalBranches(ctx context.Context, path string) ([]string, error)
}
