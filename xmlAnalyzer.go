package ryde

import (
	"encoding/csv"
	"os"
	"strings"
)

// Defines an struct to hold all assets and information about the XML file being analyzed
type XMLAnalyzer struct {
	XMLFile  XMLFile   `json:"xmlFile"`  // The XML file being analyzed.
	CSVFiles []CSVFile `json:"csvFiles"` // The CSVs file generated during analysis.
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
