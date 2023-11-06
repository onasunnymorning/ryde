package ryde

import (
	"bytes"
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
	CSVFiles map[string]CSVFile   `json:"csvFiles"` // The CSVs file generated during analysis.
	Deposit  XMLDepositUnMarshall `json:"deposit"`  // The struct for containing the UnMarshalled Deposit info
}

// CSVFile represents a CSV file with its metadata and read/write functionality.
type CSVFile struct {
	FileName       string      `json:"fileName"`  // The name of the CSV file.
	FileSize       int64       `json:"fileSize"`  // The size of the CSV file in bytes.
	LineCount      int         `json:"lineCount"` // The number of lines in the CSV file.
	fileDescriptor *os.File    // The file descriptor for the CSV file.
	CsvWriter      *csv.Writer // The CSV writer for the CSV file.
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
		// TODO: Should we return an error here or just continue?
		return err
	}
	a.XMLFile.osFile = nil
	a.XMLFile.Decoder = nil
	return nil
}

// Returns an XML decoder for the XMLFile.
func (a *XMLAnalyzer) CreateXMLDecoder() error {
	if a.XMLFile.osFile == nil {
		// TODO: Should we return an error here or open the file?
		return ErrNoXMLReader
	}
	a.XMLFile.Decoder = xml.NewDecoder(a.XMLFile.osFile)
	return nil
}

// returns an <rde:deposit> tag by reading the tokens from the decoder
func (a *XMLAnalyzer) AnalyzeDepositTag() error {
	if a.XMLFile.Decoder == nil {
		// TODO: Should we return an error here or create the decoder?
		return ErrNoXMLDecoder
	}

	err := a.XMLFile.Decoder.Decode(&a.Deposit)
	if err != nil {
		if err == io.EOF {
			return ErrNoDepositTagInFile
		}
		return fmt.Errorf("error decoding deposit: %s", err)
	}

	// found := false

	// for {
	// 	// Stop the loop if we already found and decoded the deposit tag
	// 	if found {
	// 		break
	// 	}
	// 	// Read the next token
	// 	t, tokenErr := a.XMLFile.Decoder.Token()
	// 	if tokenErr != nil {
	// 		if tokenErr == io.EOF {
	// 			return ErrNoDepositTagInFile
	// 		}
	// 		return fmt.Errorf("error decoding token: %s", tokenErr)
	// 	}

	// 	// Only process start elements of type deposit
	// 	switch se := t.(type) {
	// 	case xml.StartElement:
	// 		if se.Name.Local == "deposit" {
	// 			var d = XMLDepositUnMarshall{}
	// 			if err := a.XMLFile.Decoder.DecodeElement(&d, &se); err != nil {
	// 				return fmt.Errorf("error decoding deposit: %s", tokenErr)
	// 			}
	// 			found = true  // Mark as found so we exit as soon as possible
	// 			a.Deposit = d // Save to our XMLAnalyzer struct
	// 		}
	// 	}

	// }
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
	a.CSVFiles = make(map[string]CSVFile)
	return &a, nil
}

// Count the number of lines and save the fileSize for the set of CSV files. Use this to check against the number of objects in the header
func (a *XMLAnalyzer) CountLinesInCSVFilesAndSaveSize() error {
	for k, csvFile := range a.CSVFiles {
		file, err := os.OpenFile(csvFile.FileName, os.O_RDONLY, 0444)
		if err != nil {
			return err
		}
		defer file.Close()
		lineCount, err := CountLines(file)
		if err != nil {
			return err
		}
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		csvFile.FileSize = fi.Size()
		csvFile.LineCount = lineCount
		a.CSVFiles[k] = csvFile
	}
	return nil
}

// Count the number of lines in a file by looking for \n occurrences. Use this to check against the number of objects in the header
func CountLines(r io.Reader) (int, error) {
	var target []byte = []byte("\n")
	var err error

	count := 0

	buffer := make([]byte, 32*1024)

	for {
		read, err := r.Read(buffer)
		if err != nil {
			break
		}

		count += bytes.Count(buffer[:read], target)
	}

	if err == io.EOF {
		return count, nil
	}

	return count, err
}

// Creates CSV files for holding the data from the XML file.
func (a *XMLAnalyzer) CreateCSVFiles() error {
	for k, v := range CSVFilesAndSuffixes {
		csvFile := CSVFile{
			FileName: a.GetBaseXMLFileName() + v,
		}
		a.CSVFiles[k] = csvFile
	}
	return nil
}

// Create csvWriters for each CSV file.
func (a *XMLAnalyzer) CreateCSVWriters() error {
	for k, v := range a.CSVFiles {
		file, err := os.OpenFile(v.FileName, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		v.fileDescriptor = file
		v.CsvWriter = csv.NewWriter(file)
		a.CSVFiles[k] = v
	}
	return nil
}

// Flush csvWriters for each CSV file.
func (a *XMLAnalyzer) FlushCSVWriters() error {
	for k, v := range a.CSVFiles {
		v.CsvWriter.Flush()
		v.CsvWriter = nil
		a.CSVFiles[k] = v
	}
	return nil
}

// Close file descriptors for each CSV file.
func (a *XMLAnalyzer) CloseCSVFiles() error {
	for _, v := range a.CSVFiles {
		v.fileDescriptor.Close()
		v.fileDescriptor = nil
	}
	return nil
}
