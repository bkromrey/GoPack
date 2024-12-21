package tui

import (
	"GoPack/fileHandling"
	"GoPack/zmqClient"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type ExportModel struct {
	list           fileHandling.PackingList
	prompt         string
	username       string
	password       string
	activeInput    textInput
	usernamePrompt textinput.Model
	passwordPrompt textinput.Model
	results        string
}

const (
	usernameInput textInput = iota
	passwordInput
	doneInput
	exportDone
)

func InitExportModel(list fileHandling.PackingList) ExportModel {

	// initialize text inputs

	usernamePrompt := textinput.New()
	usernamePrompt.Placeholder = "youremail@gmail.com"
	usernamePrompt.Focus()

	passwordPrompt := textinput.New()
	passwordPrompt.Placeholder = "(app) password"

	return ExportModel{
		list: list,
		prompt: "Packing list will be automatically saved as a PDF. \nTo email this PDF " +
			"to yourself, type in your email address and app password.",
		activeInput:    0,
		usernamePrompt: usernamePrompt,
		passwordPrompt: passwordPrompt,
	}
}

func (m ExportModel) Init() tea.Cmd { return nil }

func (m ExportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case exportedMsg:
		m.activeInput = exportDone
		m.results = msg.response

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return PackingListTUI, GoToPackingListObj(m.list)

		// handle enter key based on which input is active
		case "enter":

			switch m.activeInput {

			case usernameInput:
				// save input to model
				m.username = m.usernamePrompt.Value()
				m.usernamePrompt.Blur()

				// update active prompt
				m.activeInput = passwordInput

				// display next prompt
				m.passwordPrompt.Focus()

			case passwordInput:
				// save input
				m.password = m.passwordPrompt.Value()
				m.passwordPrompt.Blur()

				// update active prompt
				m.activeInput = doneInput

				return m, exportList(m)
			}

		}

		// press any key to return to list if export is done
		if m.activeInput == exportDone {
			return PackingListTUI, GoToPackingListObj(m.list)
		}
	}

	// send update back to text input models
	m.usernamePrompt, cmd = m.usernamePrompt.Update(msg)
	m.passwordPrompt, cmd = m.passwordPrompt.Update(msg)
	return m, cmd
}

func (m ExportModel) View() string {

	// display information based on which input is active

	switch m.activeInput {

	case usernameInput:
		return m.prompt + "\n\n" + m.usernamePrompt.View()

	case passwordInput:
		return m.prompt + "\n\n" + m.usernamePrompt.View() + "\n\n" + m.passwordPrompt.View()

	case doneInput:
		return "Exporting packing list..."

	case exportDone:
		return m.results + "\n\nPress any key to go back to the packing list."
	}
	return "this is the export model"
}

// COMMANDS ----------------------------------------

func exportList(m ExportModel) tea.Cmd {
	return func() tea.Msg {

		// make request of the microservice - receive a string in response
		response := zmqClient.SendExportRequest(m.list, m.username, m.password)

		// clean up response - manually clean up long string response from microservice
		response = strings.TrimLeft(response, "Msg{Frames:{\"")
		responseSplit := strings.Split(response, "File")
		var formattedResponse string

		formattedResponse = responseSplit[0] + "\nFile" + responseSplit[1]

		if strings.Contains(formattedResponse, "Emailed to") {
			split2 := strings.Split(formattedResponse, "Emailed to")
			formattedResponse = split2[0] + "\nEmailed to" + split2[1]
		}

		formattedResponse = strings.TrimRight(formattedResponse, "\"}}")

		return exportedMsg{
			response: formattedResponse,
		}
	}
}

type exportedMsg struct {
	response string
}
