package gitclient

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type ShellClient struct{}

func NewShellClient() *ShellClient {
	return &ShellClient{}
}

func (c *ShellClient) IsGitRepo(ctx context.Context, path string) (bool, error) {
	out, err := c.runGit(ctx, path, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, nil
	}
	return strings.TrimSpace(out) == "true", nil
}

func (c *ShellClient) CurrentBranch(ctx context.Context, path string) (string, error) {
	out, err := c.runGit(ctx, path, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *ShellClient) IsDetachedHead(ctx context.Context, path string) (bool, error) {
	if _, err := c.runGit(ctx, path, "symbolic-ref", "--quiet", "--short", "HEAD"); err != nil {
		return true, nil
	}
	return false, nil
}

func (c *ShellClient) HasRemote(ctx context.Context, path string, remote string) (bool, error) {
	out, err := c.runGit(ctx, path, "remote")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == remote {
			return true, nil
		}
	}
	return false, nil
}

func (c *ShellClient) RemoteReachable(ctx context.Context, path string, remote string) (bool, error) {
	if _, err := c.runGit(ctx, path, "ls-remote", "--heads", remote); err != nil {
		return false, nil
	}
	return true, nil
}

func (c *ShellClient) Upstream(ctx context.Context, path string) (string, error) {
	out, err := c.runGit(ctx, path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *ShellClient) UpstreamForBranch(ctx context.Context, path string, branch string) (string, error) {
	spec := fmt.Sprintf("%s@{upstream}", branch)
	out, err := c.runGit(ctx, path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", spec)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *ShellClient) FetchPrune(ctx context.Context, path string, remote string) error {
	_, err := c.runGit(ctx, path, "fetch", "--prune", remote)
	return err
}

func (c *ShellClient) AheadBehind(ctx context.Context, path string) (behind int, ahead int, err error) {
	out, err := c.runGit(ctx, path, "rev-list", "--left-right", "--count", "@{u}...HEAD")
	if err != nil {
		return 0, 0, err
	}
	return parseAheadBehind(out)
}

func (c *ShellClient) AheadBehindRefs(ctx context.Context, path string, leftRef string, rightRef string) (behind int, ahead int, err error) {
	out, err := c.runGit(ctx, path, "rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", leftRef, rightRef))
	if err != nil {
		return 0, 0, err
	}
	return parseAheadBehind(out)
}

func parseAheadBehind(out string) (behind int, ahead int, err error) {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("invalid ahead/behind output %q", out)
	}

	behind, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid behind count %q: %w", fields[0], err)
	}

	ahead, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid ahead count %q: %w", fields[1], err)
	}

	return behind, ahead, nil
}

func (c *ShellClient) IsWorktreeDirty(ctx context.Context, path string) (bool, error) {
	out, err := c.runGit(ctx, path, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

func (c *ShellClient) DefaultBranch(ctx context.Context, path string, remote string) (string, error) {
	out, err := c.runGit(ctx, path, "symbolic-ref", "--short", fmt.Sprintf("refs/remotes/%s/HEAD", remote))
	if err == nil {
		parts := strings.Split(strings.TrimSpace(out), "/")
		return parts[len(parts)-1], nil
	}

	for _, candidate := range []string{"main", "master"} {
		if _, checkErr := c.runGit(ctx, path, "rev-parse", "--verify", candidate); checkErr == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not determine default branch")
}

func (c *ShellClient) IsBranchMergedInto(ctx context.Context, path string, branch string, base string) (bool, error) {
	out, err := c.runGit(ctx, path, "branch", "--format=%(refname:short)", "--merged", base)
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == branch {
			return true, nil
		}
	}
	return false, nil
}

func (c *ShellClient) LocalBranches(ctx context.Context, path string) ([]string, error) {
	out, err := c.runGit(ctx, path, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	raw := strings.Split(strings.TrimSpace(out), "\n")
	branches := make([]string, 0, len(raw))
	for _, b := range raw {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (c *ShellClient) runGit(ctx context.Context, path string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}
