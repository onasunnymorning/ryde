package ryde

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// Defines an struct to hold all assets and information about the XML file being analyzed
type XMLAnalyzer struct {
	XMLFile  XMLFile              `json:"xmlFile"`  // The XML file being analyzed.
	CSVFiles []CSVFile            `json:"csvFiles"` // The CSVs file generated during analysis.
	Deposit  XMLDepositUnMarshall `json:"deposit"`  // The struct for containing the UnMarshalled Deposit info
}

// CSVFile represents a CSV file with its metadata and read/write functionality.
type CSVFile struct {
	FileName  string      `json:"fileName"`  // The name of the CSV file.
	FileSize  int64       `json:"fileSize"`  // The size of the CSV file in bytes.
	LineCount int64       `json:"lineCount"` // The number of lines in the CSV file.
	Writer    *csv.Writer // A pointer to the csv.Writer used to write to the CSV file.
	Reader    *csv.Reader // A pointer to the csv.Reader used to read from the CSV file.
}

// XMLFile represents an XML file with its name and size.
type XMLFile struct {
	FileName string `json:"fileName"` // The name of the XML file.
	FileSize int64  `json:"fileSize"` // The size of the XML file in bytes.
	osFile   *os.File
	Decoder  *xml.Decoder
}

// Opens the XMLFile and saves the os.File pointer to the XMLFile.osFile field.
func (a *XMLAnalyzer) OpenXMLFile() error {
	reader, err := os.Open(a.XMLFile.FileName)
	if err != nil {
		return err
	}
	a.XMLFile.osFile = reader
	return nil
}

// Closes the XMLFile.osFile and removes the pointers from the XMLFile.osFile and XMLFile.Decoder fields.
func (a *XMLAnalyzer) CloseXMLFile() error {
	err := a.XMLFile.osFile.Close()
	if err != nil {
		return err
	}
	a.XMLFile.osFile = nil
	a.XMLFile.Decoder = nil
	return nil
}

// Returns an XML decoder for the XMLFile.
func (a *XMLAnalyzer) CreateXMLDecoder() error {
	if a.XMLFile.osFile == nil {
		return ErrNoXMLReader
	}
	a.XMLFile.Decoder = xml.NewDecoder(a.XMLFile.osFile)
	return nil
}

// returns an <rde:deposit> tag by reading the tokens from the decoder
func (a *XMLAnalyzer) AnalyzeDepositTag() error {
	if a.XMLFile.Decoder == nil {
		return ErrNoXMLDecoder
	}

	found := false

	for {
		// Stop the loop if we already found and decoded the deposit tag
		if found {
			break
		}
		// Read the next token
		t, tokenErr := a.XMLFile.Decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				return ErrNoDepositTagInFile
			}
			return fmt.Errorf("error decoding token: %s", tokenErr)
		}

		// Only process start elements of type deposit
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "deposit" {
				var d = XMLDepositUnMarshall{}
				if err := a.XMLFile.Decoder.DecodeElement(&d, &se); err != nil {
					return fmt.Errorf("error decoding deposit: %s", tokenErr)
				}
				found = true  // Mark as found so we exit as soon as possible
				a.Deposit = d // Save to our XMLAnalyzer struct
			}
		}

	}
	return nil
}

// Returns the XMLFile.FileName without the file extension.
func (a *XMLAnalyzer) GetBaseXMLFileName() string {
	return strings.Join(strings.Split(a.XMLFile.FileName, ".")[0:len(strings.Split(a.XMLFile.FileName, "."))-1], ".")
}

// NewXMLAnalyzer creates a new instance of XMLAnalyzer and returns a pointer to it.
// It takes a filename string as input and returns an error if the file cannot be opened or its size cannot be determined.
// The function opens the file, checks its size, and saves the filename and size to the XMLFile field of the XMLAnalyzer struct.
func NewXMLAnalyzer(filename string) (*XMLAnalyzer, error) {
	if !strings.HasSuffix(filename, ".xml") {
		return nil, ErrInvalidDepositFileName
	}
	a := XMLAnalyzer{}
	// Try and open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Check and save the file size
	fi, _ := f.Stat() // We can ignore the error here because we know the file exists
	a.XMLFile.FileSize = fi.Size()
	a.XMLFile.FileName = filename
	return &a, nil
}
