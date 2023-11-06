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

	err = a.OpenXMLFile()
	if err != nil {
		log.Fatal(err)
	}

	err = a.CreateXMLDecoder()
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

}
