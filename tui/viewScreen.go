package tui

import (
	"GoPack/fileHandling"
	"GoPack/zmqClient"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"sort"
	"strings"
)

// this is the screen that is presented when the user decides to open a packing list file
const checked = "■"   // U+25A0   // 2612           ✅ // U+2705
const unchecked = "□" // U+25A1  	// 2610
const selectIndicator = "→"

// DEFINE OUR MODEL
type ListScreenModel struct {
	// TODO should this contain all the JSON variables?
	// TODO or would it be more convienent to just store everythign as a single JSON, then
	// updates and edits wouldn't be scattered about
	fileName         string
	packingList      fileHandling.PackingList
	prompt           string
	controls         string
	selected         int
	header           string
	itemString       string
	saveSuccess      bool
	totalItems       int
	packedItems      int
	totalInLocation  map[string]int
	packedByLocation map[string]int
}

var menuOptions = fmt.Sprintf("%-35v\t", "  A: Add item to list") +
	fmt.Sprintf("%-35v\n", "C: Sort packing list by category") +
	fmt.Sprintf("%-35v\t", "  S: Save and close list") +
	fmt.Sprintf("%-35v\n", "L: Sort packing list by packed location") +
	fmt.Sprintf("%-35v\t", "  Q: Close without saving") +
	fmt.Sprintf("%-35v\n", "W: Check weather at your destination") +
	fmt.Sprintf("%-35v\t", "ESC: Back to main menu") +
	fmt.Sprintf("%-35v\n", "E: Export list as PDF & email to yourself")

// InitListScreenModelByJSON CREATES A NEW INSTANCE OF THIS MODEL BASED ON CONTENTS OF JSON FILE
func InitListScreenModelByJSON(file string) ListScreenModel {

	// load json file into memory
	var list fileHandling.PackingList
	list.LoadList(file)

	return ListScreenModel{
		fileName:    file,
		packingList: list,
		header:      fmt.Sprintf("      %-35s%-35s%-35s\n\n", "ITEM", "CATEGORY", "PACKED LOCATION"),
		//header:      "      ITEM\t\t\t\t\t  CATEGORY\t\t\t      PACKED LOCATION\n\n",
		prompt: "\nUse either ↑ or 'k' to move your selection up and ↓ or 'j' to move down. \nUse the either the spacebar or enter key to mark/unmark the selected item as packed.\n\n" +
			"Otherwise, type one of the following:\n\n",
		controls: bottomInstructions.Render(menuOptions),
	}
}

// InitListScreenModelByObj will create a new instance of this model based on
// the contents of a packing list object.
func InitListScreenModelByObj(list fileHandling.PackingList) ListScreenModel {

	return ListScreenModel{
		fileName:    SaveDirectory + list.ListName + ".json",
		packingList: list,
		//header:      "      ITEM\t\t  CATEGORY\t      PACKED LOCATION\n\n",
		header: fmt.Sprintf("      %-35s%-35s%-35s\n\n", "ITEM", "CATEGORY", "PACKED LOCATION"),
		prompt: "\nUse either ↑ or 'k' to move your selection up and ↓ or 'j' to move down. \nUse the either the spacebar or enter key to mark/unmark the selected item as packed.\n\n" +
			"Otherwise, type one of the following:\n\n",
		controls: bottomInstructions.Render(menuOptions),
	}
}

// INITIALIZE MODEL
func (m ListScreenModel) Init() tea.Cmd {
	//return nil
	return getPackedPercentage(m.packingList)
}

// UPDATE MODEL
func (m ListScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.saveSuccess == true {
		return PackingListTUI, nil
	} else {

		switch msg := msg.(type) {

		case packedPercentMsg:
			m.totalItems = msg.data.TotalItems
			m.packedItems = msg.data.PackedItems
			m.totalInLocation = msg.data.TotalInLocation
			m.packedByLocation = msg.data.PackedByLocation

		case sortedListMsg:
			return PackingListTUI, GoToPackingListObj(msg.sortedList)

		case tea.KeyMsg:

			// HANDLE KEYPRESS EVENTS
			switch msg.String() {
			// exit program
			case "ctrl+c":
				return m, tea.Quit

			// close list without saving and go back to main menu
			case "q", "esc":
				return PackingListTUI, nil

			// save and close list
			case "s":
				m.packingList.SaveList(m.fileName)
				m.saveSuccess = true

			// add new item to the list
			case "a":
				return InitAddItemScreenModel(m.packingList), nil

			// SORT LIST BY CATEGORY
			case "c":
				return m, sortListContents(m.packingList, "ItemCategory")

			// SORT LIST BY PACKED LOCATION
			case "l":
				return m, sortListContents(m.packingList, "ItemLocation")

			// go to weather screen
			case "w":
				return InitWeatherScreenModel(m.packingList), getWeather(m.packingList)

			case "e":
				return InitExportModel(m.packingList), nil

			// DOWN MOVEMENT
			case "down", "j":

				// only do something if there is at least 1 item in the list
				if len(m.packingList.Contents) > 0 {
					m.selected = m.selected + 1

					// logic to handle when the select goes out of bounds. wrap around.
					if m.selected >= len(m.packingList.Contents) {
						m.selected = m.selected % len(m.packingList.Contents)
					}
				}

			// UP MOVEMENT
			case "up", "k":
				// only do something if there is at least 1 item in list
				if len(m.packingList.Contents) > 0 {
					m.selected = m.selected - 1

					// logic to handle when the select goes out of bounds. wrap around.
					if m.selected < 0 {
						m.selected = len(m.packingList.Contents) - 1
					}
				}

			// TOGGLE CHECKED/UNCHECKED
			case " ", "enter":

				// only if list has at least 1 item
				if len(m.packingList.Contents) > 0 {

					if m.packingList.Contents[m.selected].Packed {
						m.packingList.Contents[m.selected].Packed = false
					} else {
						m.packingList.Contents[m.selected].Packed = true
					}
				}
				return m, getPackedPercentage(m.packingList)
			}

		}
	}
	return m, nil

}

// VIEW MODEL - this is what is rendered as text on the screen
func (m ListScreenModel) View() string {
	//return m.prompt + m.controls

	if m.saveSuccess == true {
		return "File successfully saved! Press any key to return to the main menu."

	}

	// use a string builder for less memory usage
	var displayBuilder strings.Builder

	// display basic list information
	displayBuilder.WriteString("List Name:\t\t" + m.packingList.ListName)
	displayBuilder.WriteString("\nTrip Dates:\t\t" + m.packingList.DepartDate + " to " + m.packingList.ReturnDate)
	displayBuilder.WriteString("\nDestination:\t\t" + m.packingList.Destination)
	displayBuilder.WriteString("\n\n" + m.header)

	// include the items within the list
	displayBuilder.WriteString(m.GenerateItemString())

	// if more than 0 items in the list, display packed percentages
	if m.totalItems > 0 {
		overallPacked := calculatePercentage(m.packedItems, m.totalItems)
		overall := fmt.Sprintf("\n%v of items packed!\n", overallPacked)
		displayBuilder.WriteString(overall)

		// alphabetize by item location
		var locations []string
		for currentLoc := range m.totalInLocation {
			locations = append(locations, currentLoc)
		}
		sort.Strings(locations)

		//iterate over location packed percentages and display in alphabetical order
		for _, locName := range locations {
			currentPercentage := calculatePercentage(m.packedByLocation[locName], m.totalInLocation[locName])
			current := fmt.Sprintf("%60v%v packed in [%v]\n", "", currentPercentage, locName)
			displayBuilder.WriteString(current)
		}
	}

	// display instructions & controls at bottom of list
	displayBuilder.WriteString("\n" + m.prompt)
	displayBuilder.WriteString(m.controls)

	display := displayBuilder.String()

	return display
}

// GenerateItemString generates list as string, and tracks which item is which
func (m ListScreenModel) GenerateItemString() string {

	// use a string builder for less memory usage
	var itemBuilder strings.Builder

	// populate the packing list items
	for i, item := range m.packingList.Contents {

		// ITEM SELECTED, so display selectIndicator
		if m.selected == i {

			// item packed so use the packed icon
			if item.Packed {
				temp := fmt.Sprintf("%s %s  %-35s%-35s%-35s\n", selectIndicator, checked, item.ItemName, item.ItemCategory, item.ItemLocation)
				itemBuilder.WriteString(temp)
			} else {
				temp := fmt.Sprintf("%s %s  %-35s%-35s%-35s\n", selectIndicator, unchecked, item.ItemName, item.ItemCategory, item.ItemLocation)
				itemBuilder.WriteString(temp)
			}

			// ITEM NOT SELECTED
		} else {
			// item packed so use the packed icon
			if item.Packed {
				temp := fmt.Sprintf("   %s  %-35s%-35s%-35s\n", checked, item.ItemName, item.ItemCategory, item.ItemLocation) //,selectIndicator + " " + checked + "  " + item.ItemName + "\t\t" + item.ItemCategory + "\t\t" + item.ItemLocation + "\n")
				itemBuilder.WriteString(temp)
			} else {
				temp := fmt.Sprintf("   %s  %-35s%-35s%-35s\n", unchecked, item.ItemName, item.ItemCategory, item.ItemLocation)
				itemBuilder.WriteString(temp)
			}
		}
	}
	return itemBuilder.String()
	//
}

// COMMANDS -------------------------------------------------------------------

// getPackedPercentage is a Cmd (with args) that will call the packed percentage
// microservice to determine how much of the list contents have already been
// packed
func getPackedPercentage(listObject fileHandling.PackingList) tea.Cmd {
	return func() tea.Msg {

		// make request of the microservice
		packedData := zmqClient.SendPackedPercentRequest(listObject)

		// return response if valid
		if packedData.TotalItems > 0 {
			// return a message to the TUI
			return packedPercentMsg{data: packedData}
		}
		return nil
	}
}

// calculatePercentage is a helper function of getPackedPercentage
func calculatePercentage(numerator int, denominator int) string {
	result := float64(numerator) / float64(denominator)
	result *= 100
	return fmt.Sprintf("%3.0f%%", result)
}

// packedPercentMsg is the Msg that the getPackedPercentage Cmd returns
type packedPercentMsg struct {
	data zmqClient.PackedPercentageResults
}

func sortListContents(list fileHandling.PackingList, sortMethod string) tea.Cmd {
	return func() tea.Msg {

		// make request of the microservice
		sortedList := zmqClient.SendListSortRequest(list, sortMethod)

		// return data
		return sortedListMsg{sortedList: sortedList}
	}
}

type sortedListMsg struct {
	sortedList fileHandling.PackingList
}
