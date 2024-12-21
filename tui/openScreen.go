package tui

// this is the screen prompts the user to make a file selection (when opening a list)

import (
	"GoPack/fileHandling"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
)

type OpenScreenModel struct {
	fileList       []os.DirEntry
	selectedFile   os.DirEntry
	prompt         string
	controls       string
	displayContent string
}

func InitOpenScreenModel() OpenScreenModel {
	files := fileHandling.ListFiles()
	fileDisplay := fileHandling.NumberedDirectoryList(files)
	prompt := "Press the corresponding number to open a list, or use the escape key to return to the main menu.\n\n"
	controls := "\nESC: Return to main menu."

	return OpenScreenModel{
		prompt:         prompt,
		controls:       controls,
		fileList:       files,
		displayContent: fmt.Sprintf("%v\n%v\n%v", prompt, fileDisplay, controls),
	}
}

// INITIALIZE MODEL
func (m OpenScreenModel) Init() tea.Cmd {
	return nil // no i/o needed
}

// UPDATE MODEL
func (m OpenScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// MESSAGE - handle various types of message
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// go back to the main menu
		case "esc":
			return InitializeTUI(), nil // TODO should this be PackingListTUI instead of creating a new model?

		// easy ways to quit program
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// check if valid integers were entered

		// convert string to int using built-in libs
		selectedFileNum, err := strconv.Atoi(msg.String())

		// check to make sure that input is numbers and input is within accepted range

		if err == nil && selectedFileNum >= 1 && selectedFileNum <= len(m.fileList) {
			m.selectedFile = m.fileList[selectedFileNum-1]
			m.displayContent = fmt.Sprintf("You selected file: %v\n", SaveDirectory+m.selectedFile.Name())

			// return this file selection in a Cmd to the main model
			return PackingListTUI, OpenSelectedFile(m.selectedFile.Name())
		}

	}

	return m, nil
}

// VIEW MODEL
func (m OpenScreenModel) View() string {
	//something here
	return m.displayContent
}

// COMMANDS

// fileSelectedMsg is used to change the main view to the list screen on the
// provided list
type fileSelectedMsg struct {
	fileName string
}

// OpenSelectedFile is a Cmd that will pass a message containing information about
// the file selection back to the main model so the selected list can be opened
func OpenSelectedFile(selected string) tea.Cmd {
	return func() tea.Msg {
		return fileSelectedMsg{
			fileName: SaveDirectory + selected,
		}
	}
}
