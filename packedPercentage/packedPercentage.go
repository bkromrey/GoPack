package packedPercentage

import (
	"GoPack/fileHandling"
)

type Results struct {
	TotalItems       int            `json:"TotalItems"`
	PackedItems      int            `json:"PackedItems"`
	TotalInLocation  map[string]int `json:"TotalInLocation"`
	PackedByLocation map[string]int `json:"PackedByLocation"`
}

// CountPacked takes in a PackingList object and determines how many items
// are packed, both overall, and per category.
func CountPacked(packingList fileHandling.PackingList) Results {
	// first, we only care about the contents, so make it easier to reference
	list := packingList.Contents

	totalItems := len(list)

	var packedItems int
	totalInLocation := make(map[string]int)
	packedByLocation := make(map[string]int)

	// iterate over all items and tally up the number that are already packed
	item := 0
	for item < totalItems {

		// if item has no location, don't tally it in total # of items in a loc
		if list[item].ItemLocation != "" {
			totalInLocation[list[item].ItemLocation] += 1
		}

		if list[item].Packed {

			packedItems += 1
			packedByLocation[list[item].ItemLocation] += 1

		}
		item += 1
	}

	return Results{
		TotalItems:       totalItems,
		PackedItems:      packedItems,
		TotalInLocation:  totalInLocation,
		PackedByLocation: packedByLocation,
	}
}
