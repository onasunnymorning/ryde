package main

import (
	"flag"
	"log"

	"github.com/onasunnymorning/ryde"
)

func main() {

	filename := flag.String("f", "", "(path to) filename")
	flag.Parse()

	if *filename == "" {
		log.Fatal("Please provide a filename")
	}

	a, err := ryde.NewXMLAnalyzer(*filename)
	if err != nil {
		log.Fatal(err)
	}

	err = a.AnalyzeDepositTag()
	if err != nil {
		log.Fatal(err)
	}

	err = a.AnalyzeTags()
	if err != nil {
		log.Fatal(err)
	}

	err = a.CountLinesInCSVFilesAndSaveSize()
	if err != nil {
		log.Fatal(err)
	}

	a.CheckValidationRules()

	a.CheckCSVColumnLength()

	err = a.WriteJSON()
	if err != nil {
		log.Fatal(err)
	}

}
