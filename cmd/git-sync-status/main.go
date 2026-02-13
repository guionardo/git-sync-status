package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/guionardo/git_sync_status/internal/gitclient"
	"github.com/guionardo/git_sync_status/internal/service"
	"github.com/guionardo/git_sync_status/internal/tui"
)

func main() {
	repoPath := flag.String("path", ".", "Repository path to inspect")
	remote := flag.String("remote", "origin", "Remote name to compare against")
	plain := flag.Bool("plain", false, "Print plain text status and exit")
	jsonOut := flag.Bool("json", false, "Print JSON status and exit")
	listBranches := flag.Bool("list-branches", false, "List local branches and exit")
	flag.Parse()

	client := gitclient.NewShellClient()
	analyzer := service.NewAnalyzer(client, *remote)

	if *listBranches {
		branches, err := analyzer.ScanLocalBranches(context.Background(), *repoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing branches: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(tui.RenderBranchList(branches))
		return
	}

	if *jsonOut {
		result := analyzer.Analyze(context.Background(), *repoPath)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *plain {
		result := analyzer.Analyze(context.Background(), *repoPath)
		fmt.Printf("path=%s\nbranch=%s\nupstream=%s\nstatus=%s\nahead=%d\nbehind=%d\n",
			result.RepoPath, result.Branch, result.Upstream, result.Status, result.Ahead, result.Behind)
		if len(result.Flags) > 0 {
			fmt.Printf("flags=%v\n", result.Flags)
		}
		if len(result.Actions) > 0 {
			fmt.Printf("actions=%v\n", result.Actions)
		}
		if result.Err != "" {
			fmt.Printf("error=%s\n", result.Err)
		}
		return
	}

	m := tui.NewModel(analyzer, *repoPath)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
}
