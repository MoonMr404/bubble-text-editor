package main

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cdd6f4"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#fffb00")).
			Padding(1, 2)

	// Select file view style
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#0b415a")).
				Bold(true).
				PaddingLeft(1)

	normalItemStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#305023")).
			Bold(true).
			Foreground(lipgloss.Color("#f70000"))

	editingAreaStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#0b415a"))

	setTextToBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))

	setBorderToBox = lipgloss.NewStyle().Border(lipgloss.ThickBorder(), true, false)
	// columnLeftStyle
)
