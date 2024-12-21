package tui

import "github.com/charmbracelet/lipgloss"

var bottomInstructions = lipgloss.NewStyle().
	//Align(lipgloss.Center).
	Foreground(lipgloss.Color("33")).
	TabWidth(lipgloss.NoTabConversion)
