package main

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/onasunnymorning/ryde"
)

func main() {
	deposit, err := ryde.NewXMLDeposit("", "123", "456", 0, time.Now().UTC())
	if err != nil {
		panic(err)
	}

	bytes, err := xml.Marshal(deposit)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
