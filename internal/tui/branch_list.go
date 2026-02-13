package tui

import (
	"strings"

	"github.com/guionardo/git_sync_status/internal/service"
)

func RenderBranchList(branches []service.BranchSummary) string {
	if len(branches) == 0 {
		return "No local branches found."
	}
	var out []string
	for _, branch := range branches {
		out = append(out, "- "+branch.Name)
	}
	return strings.Join(out, "\n")
}
