package generatePDF

import (
	"GoPack/fileHandling"
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"log"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	//"github.com/johnfercher/maroto/v2/pkg/components/line"

	"github.com/johnfercher/maroto/v2"

	//"github.com/johnfercher/maroto/v2/pkg/components/code"
	//"github.com/johnfercher/maroto/v2/pkg/components/image"
	//"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
)

const SAVE_DIR = "pdf_exports"

// The getMaroto function configures the parameter with which we wish to generate the PDF and returns a core.Maroto
// interface.
func getMaroto(list fileHandling.PackingList) core.Maroto {
	pdfConfig := config.NewBuilder().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	pdf := maroto.New(pdfConfig)

	pdf.AddRow(7,
		col.New(3),
		col.New(6).Add(
			text.New(list.ListName, props.Text{
				Size:  14,
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		),
		col.New(3))

	pdf.AddRow(6,
		col.New(3),
		col.New(6).Add(
			text.New(list.Destination, props.Text{
				Size:  12,
				Align: align.Center,
			}),
		),
		col.New(3))

	pdf.AddRow(4,
		col.New(3),
		col.New(6).Add(
			text.New(list.DepartDate+" to "+list.ReturnDate, props.Text{
				Size:  10,
				Align: align.Center,
			}),
		))

	// include packing list contents
	pdf.AddRow(5)
	pdf.AddRows(getItemList(list.Contents)...)

	return pdf
}

// GeneratePDF is a function used to parse a PackingList object into a PDF file
// for easier viewing and/or printing. Returns a string to indicate if successful
// or not.
//
// The PDF is saved in the pdf_exports folder.
func GeneratePDF(list fileHandling.PackingList) (string, string) {

	var saveLocation string

	// create the PDF object
	pdf := getMaroto(list)
	document, err := pdf.Generate()
	if err != nil {
		return saveLocation, fmt.Sprintf("PDF not exported properly: %v", err.Error())
		//log.Fatal(err.Error())
	}

	// save PDF object to disk
	saveLocation = fmt.Sprintf("%v/%v.pdf", SAVE_DIR, list.ListName)
	err = document.Save(saveLocation)
	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO can also save as a .txt file (or pass that text) to the email function, so that the email body can contain the list as well as PDF attachement

	return saveLocation, fmt.Sprintf("%v successfully exported.", list.ListName)
}

//// PDF CONTENTS -------------------------------------------------------------------------------------------------------

func getItemList(contents []fileHandling.ListItem) []core.Row {
	rows := []core.Row{

		// setup header columns for packing list
		row.New(5).Add(
			col.New(2),
			//text.NewCol(1),
			text.NewCol(4, "Item", props.Text{Style: fontstyle.Bold}),
			text.NewCol(3, "Category", props.Text{Style: fontstyle.Bold}),
			text.NewCol(3, "Packed Location", props.Text{Style: fontstyle.Bold}),
			//col.New(1),
		),
		row.New(3),
	}

	// iterate over items in the list and generate a row for each
	var itemRows []core.Row

	for _, currentItem := range contents {

		var status string

		if currentItem.Packed {
			status = "[x]"
		} else {
			status = "[  ]"
		}

		currentRow := row.New(5).Add(
			col.New(1),
			text.NewCol(1, status),
			text.NewCol(4, currentItem.ItemName),
			text.NewCol(3, currentItem.ItemCategory),
			text.NewCol(3, currentItem.ItemLocation),
			//col.New(1),
		)

		itemRows = append(itemRows, currentRow)
	}

	// combine header & content rows
	rows = append(rows, itemRows...)
	return rows
}
