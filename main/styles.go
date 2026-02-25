package main

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1D2837")).
			Foreground(lipgloss.Color("#cdd6f4"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89b4fa"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#89b4fa")).
			Padding(1, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#89b4fa")).
				Bold(true).
				PaddingLeft(1)

	normalItemStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	setTextToBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))

	setBorderToBox = lipgloss.NewStyle().Border(lipgloss.ThickBorder(), true, false)
	// columnLeftStyle
)
