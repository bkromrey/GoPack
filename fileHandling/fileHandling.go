package fileHandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const SaveDirectory = "saved_lists/"

type PackingList struct {
	ListName    string     `json:"ListName"`
	DepartDate  string     `json:"DepartDate"`
	ReturnDate  string     `json:"ReturnDate"`
	Destination string     `json:"Destination"`
	Contents    []ListItem `json:"ListContents"`
}

type ListItem struct {
	ItemName     string `json:"ItemName"`
	ItemCategory string `json:"ItemCategory"`
	ItemLocation string `json:"ItemLocation"`
	Packed       bool   `json:"Packed"`
}

//// create a new packing list
//func (p *PackingList) NewList(name string) {
//	p.ListName = name
//
//}

// load packing list from specified file
func (p *PackingList) LoadList(filepath string) {

	// open and read the file as bytes (read only)
	// ReadFile will close the file when its done
	jsonFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("file successfully opened...")

	// convert the bytes of json to a struct
	err2 := json.Unmarshal(jsonFile, p)
	if err2 != nil {
		//fmt.Println("error opening file:", err)
		log.Fatal(err)
	}
}

// method to save a packing list to a .json file in the default directory
func (p *PackingList) SaveList(filepath string) {

	// first marshall the data into json byte slice
	var rawdata []byte

	rawdata, err1 := json.Marshal(p)
	if err1 != nil {
		log.Fatal(err1)
	}

	// then write bytes to file
	err := os.WriteFile(filepath, rawdata, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// ListFiles shows a list of .json files within a specified directory
func ListFiles() []os.DirEntry {

	// create variable for storing file names
	files, err := os.ReadDir("saved_lists")

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("directory scanned successfully")

	return files

}

// NumberedDirectoryList creates a menu-appropriately formatted numbered list of all the files in directory
func NumberedDirectoryList(directory []os.DirEntry) string {
	// use the Builder to loop through file list and generate a string to display
	var fileBuilder strings.Builder
	for i, file := range directory {
		_, err := fileBuilder.WriteString(fmt.Sprintf("%d: ", i+1))
		if err != nil {
			log.Fatal(err)
		}

		_, err2 := fileBuilder.WriteString(strings.TrimSuffix(file.Name(), ".json"))
		if err2 != nil {
			log.Fatal(err2)
		}

		_, err3 := fileBuilder.WriteString("\n")
		if err3 != nil {
			log.Fatal(err3)
		}
	}
	return fileBuilder.String()
}

// DeleteList deletes the specified file from the directory
func DeleteList(file os.DirEntry) {

	fullpath := SaveDirectory + file.Name()

	err := os.Remove(fullpath)

	if err != nil {
		log.Fatal(err)
	}

}
