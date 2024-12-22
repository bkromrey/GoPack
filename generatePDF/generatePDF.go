package generatePDF

import (
	"GoPack/fileHandling"
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
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

	mrt := maroto.New(pdfConfig)

	err := mrt.RegisterHeader(getPageHeader(list))
	if err != nil {
		log.Fatal(err)
	}

	return mrt
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

// PDF CONTENTS -------------------------------------------------------------------------------------------------------
func getPageHeader(list fileHandling.PackingList) core.Row {
	return row.New(20).Add(
		col.New(3),
		col.New(6).Add(
			text.New(list.ListName, props.Text{
				Size:  8,
				Align: align.Center,
			}),
			text.New(list.DepartDate+" - "+list.ReturnDate, props.Text{
				Size:  5,
				Align: align.Center,
			}),
		),
		col.New(3),
	)
}
