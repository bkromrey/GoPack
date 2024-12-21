package tui

// this model orchestrates - it decides which sub-model is to be shown
// and routes messages to appropriate sub-models

import (
	fileHandling "GoPack/fileHandling"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// each TUI screen will be represented by an integer
type TUIScreen int

// list of all the TUI screens (models)
const (
	mainMenu TUIScreen = iota
	openScreen
	newScreen
	listScreen
	addItemScreen
	deleteScreen
)

const SaveDirectory = "saved_lists/"

const MenuText string = "               ** Go Pack! **    \n" +
	"An application to help you get packed for\nyour next trip and remind you what you packed \non your last trip.\n\n\n" +
	"All data entered within this application is \nprivate and stored as .json files on your computer.\n\n" +
	"Please press the corresponding key.\n\n" +
	"O: Open a Packing List\n" +
	"N: Create a New Packing List\n" +
	"D: Delete a Packing List\n" +
	"Q: Quit"

// MainScreen stores information about the model's state
type MainScreen struct {
	activeScreen     TUIScreen
	Functions        []MenuAction
	DisplayContent   string
	selectedFilePath string
}

func InitializeTUI() tea.Model {
	return MainScreen{
		Functions: []MenuAction{
			MenuAction{
				Description: "",
				OnSelect:    func() tea.Msg { return ViewMainMenuMsg{} }, // show the main menu
			},
			MenuAction{
				Description: "O: open list",
				OnSelect:    func() tea.Msg { return OpenListMsg{} }, // return a Msg with a specific type
			},
			MenuAction{
				Description: "N: new list",
				OnSelect:    func() tea.Msg { return NewListMsg{} },
			},
			MenuAction{}, //list screen
			MenuAction{}, //add item screen
			MenuAction{
				Description: "D: delete list",
				OnSelect:    func() tea.Msg { return DeleteListMsg{} },
			},
			//menuAction{ //
			//	description: "A: about this app",
			//	onSelect:    func() tea.Msg { return struct{}{} }, // TODO - what is this double struct? answer its a placeholder that returns an empty struct
			//},
		},
		DisplayContent: MenuText,
	}
}

var PackingListTUI = InitializeTUI()

// Init: INITIALIZE MODEL - returns command
func (m MainScreen) Init() tea.Cmd {
	return nil
}

// Update: UPDATE MODEL - responds to user input. takes in a Msg, returns an updated model & command
func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// MESSAGE - handle different types of messages differently
	switch msg := msg.(type) {

	// sub-model wants to return to the main menu
	case UpModelMsg:
		m.activeScreen = mainMenu

	// the Open menu selection
	case OpenListMsg:
		m.activeScreen = openScreen

	case fileSelectedMsg:
		m.activeScreen = listScreen
		m.selectedFilePath = msg.fileName

	case goToPackingListObjMsg:
		m.activeScreen = listScreen
		return InitListScreenModelByObj(msg.listObject), getPackedPercentage(msg.listObject) //nil

	case ViewMainMenuMsg:
		return m.ViewMainMenu(), nil

	case NewListMsg:
		m.activeScreen = newScreen

	case DeleteListMsg:
		m.activeScreen = deleteScreen

	// KEYPRESS - do different things based on what key is pressed
	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		// return the model unchanged and call the action's associated function
		case "o":
			m.activeScreen = openScreen
			return m, m.Functions[openScreen].OnSelect

		// goes back to the main menu
		case "esc":
			return m, m.Functions[mainMenu].OnSelect

		// new packing list
		case "n":
			return m, m.Functions[newScreen].OnSelect

		// delete a packing list
		case "d":
			return m, m.Functions[deleteScreen].OnSelect
		}

	}

	// ON ACTIVE SCREEN
	switch m.activeScreen {

	// pass the message off to the openScreen model's Update function
	case openScreen:
		// create a new model & pass the update to it to ensure we are always looking at updated directory structure
		openModel := InitOpenScreenModel()
		openModel.Update(msg)
		return openModel, nil

	case newScreen:
		return InitNewListModel(), nil

	case listScreen:
		viewThisListModel := InitListScreenModelByJSON(m.selectedFilePath)
		return viewThisListModel, getPackedPercentage(viewThisListModel.packingList)

	case addItemScreen:

	case deleteScreen:
		return InitDeleteListModel(), nil

	default:
		m.activeScreen = mainMenu
	}

	return m, nil

}

type UpModelMsg struct{}

type ViewMainMenuMsg struct{}

// ViewMainMenu is a Cmd that shows the main menu
func (m MainScreen) ViewMainMenu() tea.Model {
	m.DisplayContent = fmt.Sprintln(MenuText)
	return m
}

// we are going to use this to change the view to openScreen
type OpenListMsg struct{}

type NewListMsg struct{}

func (m MainScreen) NewList() tea.Model {
	m.DisplayContent = fmt.Sprintln("Create a new list!!")

	var newlist fileHandling.PackingList

	newlist.ListName = "test"
	m.DisplayContent += fmt.Sprintln(newlist)

	return m
}

type DeleteListMsg struct{}

// VIEW MODEL - presents data to user. everything that should be displayed is returned as one giant string
func (m MainScreen) View() string {
	var content string
	content = m.DisplayContent
	return content
}

// add menu options
type MenuAction struct {
	Description string
	OnSelect    func() tea.Msg
}
