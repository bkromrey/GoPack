package tui

import (
	"GoPack/fileHandling"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type newListModel struct {
	displayText      string
	activeInput      textInput
	listNameInput    textinput.Model
	listName         string
	departDateInput  textinput.Model
	departDate       string
	returnDateInput  textinput.Model
	returnDate       string
	destinationInput textinput.Model
	destination      string
	fileName         string
}

type textInput int

const (
	listNameInput textInput = iota
	departDateInput
	returnDateInput
	destinationInput
	confirmationPrompt
	finalSelection // TODO maybe delete this?
)

func InitNewListModel() newListModel {

	// initialize text inputs
	listname := textinput.New()
	listname.Placeholder = "packing list name"
	listname.Focus()

	departdate := textinput.New()
	departdate.Placeholder = "departure date (YYYY-MM-DD)"

	returndate := textinput.New()
	returndate.Placeholder = "return date (YYYY-MM-DD)"

	destination := textinput.New()
	destination.Placeholder = "destination"

	initalDisplay := "Please follow the prompts to create a new packing list. You will have a chance " +
		"\nto review your responses before the list is created.\n" // + listname.View()
	return newListModel{
		displayText:      initalDisplay,
		listNameInput:    listname,
		departDateInput:  departdate,
		returnDateInput:  returndate,
		destinationInput: destination,
		activeInput:      listNameInput,
	}
}

// INITIALIZE MODEL
func (m newListModel) Init() tea.Cmd { return nil }

// UPDATE MODEL
func (m newListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	// HANDLE USER INPUTS
	switch msg := msg.(type) {

	// HANDLE FILE CREATION
	case listCreatedMsg:
		return PackingListTUI, OpenSelectedFile(msg.fileName)

	// HANDLE USER INPUTS
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return InitializeTUI(), nil

		case "ctrl+c":

			return m, tea.Quit

		case "y":
			if m.activeInput == confirmationPrompt {
				m.activeInput = finalSelection
				return m, CreateJSON(m)
			}

		case "n":
			if m.activeInput == confirmationPrompt {
				return InitializeTUI(), nil
			}

		case "enter":

			// handle enter differently depending on which is active when enter was pressed
			switch m.activeInput {

			case listNameInput:
				// save input
				m.listName = m.listNameInput.Value()
				m.listNameInput.Blur() // cmd

				// update active prompt
				m.activeInput = departDateInput

				// display next prompt
				m.departDateInput.Focus() //cmd

			case departDateInput:
				// save input
				m.departDate = m.departDateInput.Value()
				m.departDateInput.Blur()

				// update active prompt
				m.activeInput = returnDateInput

				// display next prompt
				m.returnDateInput.Focus()

			case returnDateInput:
				// save input
				m.returnDate = m.returnDateInput.Value()
				m.returnDateInput.Blur()

				// update active prompt
				m.activeInput = destinationInput

				// display next prompt
				m.destinationInput.Focus() // TODO - this is where its running into an error

			case destinationInput:
				// save input from previous
				m.destination = m.destinationInput.Value()
				m.destinationInput.Blur()

				// update active prompt
				m.activeInput = confirmationPrompt

			case finalSelection:
				// TODO create the .json file and then send somethign back to main menu model to open this specific list. might need to move this out of this switch though.

			}

		}

	}

	// send update back to text input models
	m.listNameInput, cmd = m.listNameInput.Update(msg)
	m.departDateInput, cmd = m.departDateInput.Update(msg)
	m.returnDateInput, cmd = m.returnDateInput.Update(msg)
	m.destinationInput, cmd = m.destinationInput.Update(msg)
	return m, cmd
}

// VIEW MODEL
func (m newListModel) View() string {
	switch m.activeInput {

	case listNameInput:
		return m.displayText + "\n\n" + m.listNameInput.View()

	case departDateInput:
		return m.displayText + "\n\n" + m.listNameInput.View() + "\n\n" + m.departDateInput.View()

	case returnDateInput:
		return m.displayText + "\n\n" + m.listNameInput.View() + "\n\n" + m.departDateInput.View() + "\n\n" + m.returnDateInput.View()

	case destinationInput:
		return m.displayText + "\n\n" + m.listNameInput.View() + "\n\n" + m.departDateInput.View() + "\n\n" + m.returnDateInput.View() + "\n\n" + m.destinationInput.View()

	case confirmationPrompt:
		confirmation := fmt.Sprintf("Create a new packing list with the following characteristics?\n\nList Name:\t%s\nTrip Dates:\t%s-%s\nDestination:\t%s", m.listName, m.departDate, m.returnDate, m.destination)
		instruction := "\n\nY: Yes! Create new list & go to the list view to start adding items to the list!\nN/ESC: Maybe not. Return to the main menu.\n\nType in your selection now:\n> "
		return confirmation + instruction

	case finalSelection:
		final := fmt.Sprintf("Creating packing list...")
		return final
	}

	return ""
}

type listCreatedMsg struct {
	fileName string
}

// CreateJSON is a Cmd that takes in an argument (the current model) and from that
// information, creates a .json file on disk.
// Returns: a listCreatedMsg containing the name of the newly created .json file
func CreateJSON(m newListModel) tea.Cmd {

	return func() tea.Msg {

		// create packing list object from defined input
		var p fileHandling.PackingList
		p.ListName = m.listName
		p.DepartDate = m.departDate
		p.ReturnDate = m.returnDate
		p.Destination = m.destination

		// create filename
		m.fileName = m.listName + ".json"

		// write packing list object to file
		p.SaveList(SaveDirectory + m.fileName)

		return listCreatedMsg{
			fileName: m.fileName,
		}
	}
}
