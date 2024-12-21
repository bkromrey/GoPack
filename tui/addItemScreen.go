package tui

import (
	"GoPack/fileHandling"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// this is the screen that is presented when the user decides to add an item to a packing list

// DEFINE OUR MODEL
type addItemScreenModel struct {
	displayText         string
	activeInput         textInput
	itemNameInput       textinput.Model
	itemName            string
	categoryInput       textinput.Model
	category            string
	packedLocationInput textinput.Model
	packedLocation      string
	listObject          fileHandling.PackingList
}

const (
	itemNameInput textInput = iota
	categoryInput
	packedLocationInput
	confirmationItem
	itemConfirmed
)

// FUNCTION TO CREATE NEW INSTANCE OF THIS MODEL
func InitAddItemScreenModel(list fileHandling.PackingList) addItemScreenModel {

	initialDisplay := "Add New Item To List\n\nCategory & packed location are optional, and if provided, will add additional context \nwhen viewing the packing list. Leaving these blank will not affect the functionality \nof the app."

	// initialize text inputs
	itemName := textinput.New()
	itemName.Prompt = "item name:\n"
	itemName.Focus()

	category := textinput.New()
	category.Prompt = "category (optional, hit ENTER to leave blank):\n"

	packedLocation := textinput.New()
	packedLocation.Prompt = "packed location (optional, hit ENTER to leave blank):\n"

	return addItemScreenModel{
		activeInput:         itemNameInput,
		listObject:          list,
		itemNameInput:       itemName,
		categoryInput:       category,
		packedLocationInput: packedLocation,
		displayText:         initialDisplay,
	}
}

// INITIALIZE MODEL
func (m addItemScreenModel) Init() tea.Cmd {
	return nil
}

// UPDATE MODEL
func (m addItemScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		//exit program
		case "ctrl+c":

			return m, tea.Quit

		// go back to packing list
		case "esc":
			//return // TODO - we need something like OpenSelectedFile, but that passes in an packingList Object and not a JSON file

		// only handle y if we are on the confirmation screen
		case "y":
			if m.activeInput == confirmationItem {
				//addToList(m)
				m.activeInput = itemConfirmed

				var item fileHandling.ListItem
				item.ItemName = m.itemName
				item.ItemCategory = m.category
				item.ItemLocation = m.packedLocation
				item.Packed = false

				// add new item to packing list object
				m.listObject.Contents = append(m.listObject.Contents, item)

				return PackingListTUI, GoToPackingListObj(m.listObject)
			}
		case "n":
			if m.activeInput == confirmationItem {
				return PackingListTUI, GoToPackingListObj(m.listObject)
			}

		// handle enter key based on which input is active
		case "enter":

			switch m.activeInput {

			case itemNameInput:
				// save input
				m.itemName = m.itemNameInput.Value()
				m.itemNameInput.Blur()

				// update active prompt
				m.activeInput = categoryInput

				// display next prompt
				m.categoryInput.Focus()

			case categoryInput:
				// save input
				m.category = m.categoryInput.Value()
				m.categoryInput.Blur()

				// update active prompt
				m.activeInput = packedLocationInput

				// display next prompt
				m.packedLocationInput.Focus()

			case packedLocationInput:
				// save input
				m.packedLocation = m.packedLocationInput.Value()
				m.packedLocationInput.Blur()

				// update active prompt
				m.activeInput = confirmationItem

				//case itemConfirmed:
				//	return PackingListTUI, GoToPackingListObj(m)
			}
		}

		//case itemAddedMsg:
		//	return m, GoToPackingListObj(m)

	}

	if m.activeInput == itemConfirmed {
		return PackingListTUI, GoToPackingListObj(m.listObject)
	}

	// send update back to text input models
	m.itemNameInput, cmd = m.itemNameInput.Update(msg)
	m.categoryInput, cmd = m.categoryInput.Update(msg)
	m.packedLocationInput, cmd = m.packedLocationInput.Update(msg)
	return m, cmd
}

// VIEW MODEL
func (m addItemScreenModel) View() string {

	// display information based on which input is active
	switch m.activeInput {

	case itemNameInput:
		return m.displayText + "\n\n" + m.itemNameInput.View()

	case categoryInput:
		return m.displayText + "\n\n" + m.itemNameInput.View() + "\n\n" + m.categoryInput.View()

	case packedLocationInput:
		return m.displayText + "\n\n" + m.itemNameInput.View() + "\n\n" + m.categoryInput.View() + "\n\n" + m.packedLocationInput.View()

	case confirmationItem:
		confirmation := fmt.Sprintf("Add the following item to the packing list?\n\nItem Name:\t\t%s\nCategory:\t\t%s\nPacked Location:\t%s\n\nConfirm item addition? Y/N\n> ", m.itemName, m.category, m.packedLocation)
		return confirmation

	case itemConfirmed:
		return "adding item..." + fmt.Sprintln(m.listObject.Contents)
	}
	return "this is the add item model"
}

// COMMANDS

type goToPackingListObjMsg struct {
	listObject fileHandling.PackingList
}

// GoToPackingListObj is a Cmd that will pass a message containing the packing
// list object back to the main screen model so the list view can be
// re-opened on this object.
// Used by functions in the add item screen as well as the view packing list screen.
// Returns a goToPackingListObjMsg
func GoToPackingListObj(list fileHandling.PackingList) tea.Cmd {
	return func() tea.Msg {
		return goToPackingListObjMsg{
			listObject: list,
		}
	}
}
