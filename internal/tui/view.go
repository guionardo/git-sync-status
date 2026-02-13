package tui

import (
	"fmt"
	"strings"

	"github.com/guionardo/git_sync_status/internal/domain"
	"github.com/guionardo/git_sync_status/internal/service"
)

func (m Model) renderView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Git Sync Status"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(m.repoPath))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading repository status...\n\n")
		b.WriteString(mutedStyle.Render("Press q to quit"))
		return b.String()
	}

	if m.lastErr != nil {
		b.WriteString(errStyle.Render("Error: " + m.lastErr.Error()))
		b.WriteString("\n")
	}

	b.WriteString(boxStyle.Render(m.renderStatusCard()))
	if len(m.branches) > 0 {
		b.WriteString("\n\n")
		b.WriteString(boxStyle.Render(m.renderAllBranchesCard()))
	}
	b.WriteString("\n\n")

	help := []string{"r refresh", "q quit"}
	b.WriteString(mutedStyle.Render(strings.Join(help, " • ")))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderStatusCard() string {
	r := m.result
	lines := []string{
		headerStyle.Render("Repository"),
		fmt.Sprintf("Path: %s", r.RepoPath),
		fmt.Sprintf("Branch: %s", fallback(r.Branch, "(unknown)")),
		fmt.Sprintf("Upstream: %s", fallback(r.Upstream, "(none)")),
		"",
		headerStyle.Render("Status"),
		m.renderStatusValue(r.Status),
		fmt.Sprintf("Ahead/Behind: %d/%d", r.Ahead, r.Behind),
	}

	if len(r.Flags) > 0 {
		lines = append(lines, "", headerStyle.Render("Flags"))
		for _, flag := range r.Flags {
			lines = append(lines, "- "+flag)
		}
	}

	if len(r.Actions) > 0 {
		lines = append(lines, "", headerStyle.Render("Suggested actions"))
		for _, action := range r.Actions {
			lines = append(lines, "- "+action)
		}
	}

	if len(r.Details) > 0 {
		lines = append(lines, "", headerStyle.Render("Diagnostics"))
		for _, detail := range r.Details {
			lines = append(lines, "- "+detail)
		}
	}

	if r.Err != "" {
		lines = append(lines, "", errStyle.Render("Last git error"), r.Err)
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderStatusValue(status domain.Status) string {
	switch status {
	case domain.StatusSynced:
		return okStyle.Render(string(status))
	case domain.StatusSyncPending, domain.StatusLate, domain.StatusNoUpstream:
		return warnStyle.Render(string(status))
	case domain.StatusDiverged, domain.StatusNoRemote, domain.StatusNotAGitRepo:
		return errStyle.Render(string(status))
	default:
		return string(status)
	}
}

func (m Model) renderAllBranchesCard() string {
	lines := []string{
		headerStyle.Render("All Branches"),
		"",
		m.renderBranchTable(m.branches),
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderBranchTable(rows []service.BranchStatus) string {
	branchW := len("BRANCH")
	upstreamW := len("UPSTREAM")
	statusW := len("STATUS")
	abW := len("A/B")
	flagsW := len("FLAGS")

	for _, row := range rows {
		if len(row.Branch) > branchW {
			branchW = len(row.Branch)
		}
		up := fallback(row.Upstream, "-")
		if len(up) > upstreamW {
			upstreamW = len(up)
		}
		if len(string(row.Status)) > statusW {
			statusW = len(string(row.Status))
		}
		ab := fmt.Sprintf("%d/%d", row.Ahead, row.Behind)
		if len(ab) > abW {
			abW = len(ab)
		}
		flags := "-"
		if len(row.Flags) > 0 {
			flags = strings.Join(row.Flags, ",")
		}
		if len(flags) > flagsW {
			flagsW = len(flags)
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%-*s  %-*s  %-*s  %-*s  %-*s\n",
		branchW, "BRANCH", upstreamW, "UPSTREAM", statusW, "STATUS", abW, "A/B", flagsW, "FLAGS")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", branchW+upstreamW+statusW+abW+flagsW+8))
	for _, row := range rows {
		flags := "-"
		if len(row.Flags) > 0 {
			flags = strings.Join(row.Flags, ",")
		}
		ab := fmt.Sprintf("%d/%d", row.Ahead, row.Behind)
		fmt.Fprintf(&b, "%-*s  %-*s  %-*s  %-*s  %-*s\n",
			branchW, row.Branch,
			upstreamW, fallback(row.Upstream, "-"),
			statusW, string(row.Status),
			abW, ab,
			flagsW, flags,
		)
	}
	return strings.TrimRight(b.String(), "\n")
}

func fallback(value string, alt string) string {
	if strings.TrimSpace(value) == "" {
		return alt
	}
	return value
}
