package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/guionardo/git_sync_status/internal/domain"
	"github.com/guionardo/git_sync_status/internal/service"
)

type resultMsg struct {
	result     domain.Result
	branchRows []service.BranchStatus
	err        error
}

type errMsg struct {
	err error
}

type Model struct {
	analyzer *service.Analyzer
	repoPath string
	keys     keyMap

	loading  bool
	result   domain.Result
	branches []service.BranchStatus
	lastErr  error
}

func NewModel(analyzer *service.Analyzer, repoPath string) Model {
	return Model{
		analyzer: analyzer,
		repoPath: repoPath,
		keys:     defaultKeyMap(),
		loading:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.refreshCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, m.keys.Quit):
			return m, tea.Quit
		case keyMatches(msg, m.keys.Refresh):
			m.loading = true
			m.lastErr = nil
			return m, m.refreshCmd()
		}
	case resultMsg:
		m.loading = false
		m.result = msg.result
		m.branches = msg.branchRows
		m.lastErr = msg.err
		return m, nil
	case errMsg:
		m.loading = false
		m.lastErr = msg.err
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	return m.renderView()
}

func (m Model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		result := m.analyzer.Analyze(context.Background(), m.repoPath)
		branchRows, err := m.analyzer.AnalyzeAllBranches(context.Background(), m.repoPath)
		return resultMsg{result: result, branchRows: branchRows, err: err}
	}
}

func keyMatches(msg tea.KeyMsg, binding key.Binding) bool {
	return key.Matches(msg, binding)
}
