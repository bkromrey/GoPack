package tui

import (
	"GoPack/fileHandling"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
	"strings"
)

// }
type deleteScreenModel struct {
	allFiles       []os.DirEntry
	fileToDelete   os.DirEntry
	prompt         textinput.Model
	displayContent string
	activeScreen   screen
}

type screen int

const (
	makeSelection = iota
	fileSelected
)

func InitDeleteListModel() deleteScreenModel {

	prompt := "Press the corresponding number to delete a list, or use the escape key to return to the main menu.\n\nList deletion is PERMANENT and cannot be undone.\n\n"

	files := fileHandling.ListFiles()
	fileDisplay := fileHandling.NumberedDirectoryList(files)

	// intialize text input confirmation
	confirmPrompt := textinput.New()
	confirmPrompt.Prompt = "Type 'DELETE' to confirm file deletion, or anything else to cancel operation.\n> "

	return deleteScreenModel{
		displayContent: prompt + fileDisplay,
		prompt:         confirmPrompt,
		allFiles:       files,
		activeScreen:   makeSelection,
	}
}

// INITIALIZE MODEL
func (m deleteScreenModel) Init() tea.Cmd { return nil }

// UPDATE MODEL
func (m deleteScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case SelectedMsg:
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return PackingListTUI, nil

		case "q", "ctrl+c":

			if m.activeScreen == fileSelected && msg.String() == "q" {
				break
			}
			return m, tea.Quit

		case "enter":
			if m.activeScreen == fileSelected {

				// if 'delete' was typed in delete the file
				if m.prompt.Value() == "delete" || m.prompt.Value() == "DELETE" {
					// perform file deletion
					fileHandling.DeleteList(m.fileToDelete)
				}

				// return to main menu
				return PackingListTUI, nil

			}
		}

		switch m.activeScreen {

		// allow any keypress
		case makeSelection:
			// check if valid integer was entered
			selection, err := strconv.Atoi(msg.String())

			if err == nil && selection >= 1 && selection <= len(m.allFiles) {
				m.fileToDelete = m.allFiles[selection-1]
				m.activeScreen = fileSelected
				m.displayContent = m.displayContent + fmt.Sprintf("\nReally delete %s?\n\n", strings.TrimSuffix(m.fileToDelete.Name(), ".json"))
				m.prompt.Focus()
				return m, nil
			}
		case fileSelected:

			//m.prompt.Focus()
		}

	}

	m.prompt, cmd = m.prompt.Update(msg)
	return m, cmd
}

// VIEW MODEL
func (m deleteScreenModel) View() string {

	if m.fileToDelete != nil {
		return m.displayContent + m.prompt.View()
	} else {
		return m.displayContent
	}
}

type SelectedMsg struct{}

func FileDeletion() tea.Msg {

	return SelectedMsg{}
}
