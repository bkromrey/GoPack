package main

import (
	"GoPack/tui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	// create bubbletea program by intializing the main view model
	program := tea.NewProgram(tui.InitializeTUI(), tea.WithAltScreen())

	// catch any errors thrown, otherwise run program
	if _, err := program.Run(); err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
}
