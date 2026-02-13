package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guionardo/git_sync_status/internal/domain"
	"github.com/guionardo/git_sync_status/internal/gitclient"
)

func TestAnalyzerIntegrationNoRemote(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.name", "test")
	runGit(t, repo, "config", "user.email", "test@example.com")
	writeFile(t, filepath.Join(repo, "a.txt"), "hello")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "init")

	analyzer := NewAnalyzer(gitclient.NewShellClient(), "origin")
	got := analyzer.Analyze(context.Background(), repo)
	if got.Status != domain.StatusNoRemote {
		t.Fatalf("got %s, want %s", got.Status, domain.StatusNoRemote)
	}
}

func TestAnalyzerIntegrationSyncPending(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	work := filepath.Join(root, "work")

	runGit(t, root, "init", "--bare", remote)
	runGit(t, root, "clone", remote, seed)
	runGit(t, seed, "config", "user.name", "test")
	runGit(t, seed, "config", "user.email", "test@example.com")
	writeFile(t, filepath.Join(seed, "a.txt"), "hello")
	runGit(t, seed, "add", ".")
	runGit(t, seed, "commit", "-m", "seed")
	runGit(t, seed, "branch", "-M", "main")
	runGit(t, seed, "push", "-u", "origin", "main")
	runGit(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")

	runGit(t, root, "clone", remote, work)
	runGit(t, work, "config", "user.name", "test")
	runGit(t, work, "config", "user.email", "test@example.com")
	writeFile(t, filepath.Join(work, "b.txt"), "local")
	runGit(t, work, "add", ".")
	runGit(t, work, "commit", "-m", "local")

	analyzer := NewAnalyzer(gitclient.NewShellClient(), "origin")
	got := analyzer.Analyze(context.Background(), work)
	if got.Status != domain.StatusSyncPending {
		t.Fatalf("got %s, want %s", got.Status, domain.StatusSyncPending)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
}

func writeFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}
