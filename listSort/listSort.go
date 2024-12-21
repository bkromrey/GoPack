package listSort

import (
	"GoPack/fileHandling"
	"sort"
)

func SortCategory(list fileHandling.PackingList) fileHandling.PackingList {
	var sortedList []fileHandling.ListItem
	items := list.Contents

	listsByCat := make(map[string][]fileHandling.ListItem)

	// go through each item in the list, and add item to a new slice for that category
	index := 0
	for index < len(items) {
		currentCategory := items[index].ItemCategory
		listsByCat[currentCategory] = append(listsByCat[currentCategory], items[index])
		index += 1
	}

	// alphabetize by category name
	var categories []string
	for currentCat, _ := range listsByCat {
		categories = append(categories, currentCat)
	}
	sort.Strings(categories)

	// then recombine these sub-lists into one list
	for _, sublist := range categories {
		sortedList = append(sortedList, listsByCat[sublist]...)
	}

	list.Contents = sortedList
	return list
}

func SortLocation(list fileHandling.PackingList) fileHandling.PackingList {
	var sortedList []fileHandling.ListItem
	items := list.Contents

	listsByLoc := make(map[string][]fileHandling.ListItem)

	// go through each item in the list, and add item to new slice for that location
	index := 0
	for index < len(items) {
		currentLocation := items[index].ItemLocation
		listsByLoc[currentLocation] = append(listsByLoc[currentLocation], items[index])
		index += 1
	}

	// alphabetize by item location
	var locations []string
	for currentLoc, _ := range listsByLoc {
		locations = append(locations, currentLoc)
	}
	sort.Strings(locations)

	// then recombine sub-lists into one list
	for _, sublist := range locations {
		sortedList = append(sortedList, listsByLoc[sublist]...)
	}

	list.Contents = sortedList
	return list
}
