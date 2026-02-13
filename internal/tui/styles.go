package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	errStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)
