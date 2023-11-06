package ryde

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"
)

// TestNewXMLAnalyzer tests the NewXMLAnalyzer function which creates a new XMLAnalyzer object and performs various tests on it.
// It creates a temporary file for testing, writes some data to the file, and checks that the file size and name were set correctly.
// It also checks if we can get a reader, get a decoder, and open a file that doesn't exist.
// Finally, it checks if we can close the file and get an error when trying to close it again.
func TestNewXMLAnalyzer(t *testing.T) {
	// Create a temporary file for testing
	f, err := os.CreateTemp("", "*test.xml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(f.Name())

	// Write some data to the file
	data := []byte("test data")
	if _, err := f.Write(data); err != nil {
		t.Fatalf("Failed to write data to file: %v", err)
	}

	// Close the file
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Call the function being tested
	a, err := NewXMLAnalyzer(f.Name())
	if err != nil {
		t.Fatalf("NewXMLAnalyzer failed with error: %v", err)
	}

	// Check that the file size and name were set correctly
	if a.XMLFile.FileSize != int64(len(data)) {
		t.Errorf("Expected file size to be %d, got %d", len(data), a.XMLFile.FileSize)
	}
	if a.XMLFile.FileName != f.Name() {
		t.Errorf("Expected file name to be %s, got %s", f.Name(), a.XMLFile.FileName)
	}
	if a.GetBaseXMLFileName() != f.Name()[:len(f.Name())-4] {
		t.Errorf("Expected base file name to be %s, got %s", f.Name()[:len(f.Name())-4], a.GetBaseXMLFileName())
	}

	// Try and open a file that doesn't exist
	_, err = NewXMLAnalyzer("does_not_exist.xml")
	if err == nil {
		t.Error("Expected NewXMLAnalyzer to return an error when opening a file that doesn't exist")
	}

	// Try and open a file that doesn't end with .xml
	_, err = NewXMLAnalyzer("test.txt")
	if err == nil {
		t.Error("Expected NewXMLAnalyzer to return an error when opening a file that doesn't end with .xml")
	}

	// Check if we can get a reader
	err = a.OpenXMLFile()
	if err != nil {
		t.Fatalf("Failed to open XML file: %v", err)
	}
	if a.XMLFile.osFile == nil {
		t.Error("Expected XMLFile.osFile to be set")
	}
	err = a.CloseXMLFile()
	if err != nil {
		t.Fatalf("Failed to close XML file: %v", err)
	}
	// close the file again
	err = a.CloseXMLFile()
	if err == nil {
		t.Error("Expected CloseXMLFile to return an error when the file is already closed")
	}
	// Break the fileName and try to open the file
	a.XMLFile.FileName = "does_not_exist.xml"
	err = a.OpenXMLFile()
	if err == nil {
		t.Error("Expected OpenXMLFile to return an error when the file doesn't exist")
	}

	// Try and get a decoder without opening the file
	err = a.CreateXMLDecoder()
	if err == nil {
		t.Error("Expected GetXMLDecoder to return an error when the file is not open")
	}
	// Fix the filename and try again
	a.XMLFile.FileName = f.Name()
	err = a.OpenXMLFile()
	if err != nil {
		t.Fatalf("Failed to open XML file: %v", err)
	}
	err = a.CreateXMLDecoder()
	if err != nil {
		t.Fatalf("Failed to get XML decoder: %v", err)
	}

}

// TestAnalyzeDepositTagFail tests the AnalyzeDepositTag method of the XMLAnalyzer struct
// when the deposit tag is not found in the XML file. It creates a temporary file for testing,
// initializes a new XMLAnalyzer with the file name, and calls AnalyzeDepositTag method twice,
// expecting an error the first time and an ErrNoDepositTagFound error the second time.
func TestAnalyzeDepositTagFail(t *testing.T) {
	// Create a temporary file for testing
	f, err := os.CreateTemp("", "*test.xml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(f.Name())

	a, err := NewXMLAnalyzer(f.Name())
	if err != nil {
		t.Fatalf("NewXMLAnalyzer returned an error: %s", err)
	}

	err = a.AnalyzeDepositTag()
	if err == nil {
		t.Fatalf("Expecting an error from AnalyzeDepositTag(), but got none")
	}

	err = a.OpenXMLFile()
	if err != nil {
		t.Fatalf("OpenXMLFile() returned an error: %s", err)
	}
	err = a.CreateXMLDecoder()
	if err != nil {
		t.Fatalf("CreateXMLDecoder() returned an error: %s", err)
	}

	err = a.AnalyzeDepositTag()
	if err == nil {
		t.Fatalf("Expected an ErrNoDepositTagFound error but got none")
	}
}

// TestAnalyzeDepositTagSuccess tests the AnalyzeDepositTag function of the XMLAnalyzer struct.
// It creates a temporary file for testing, initializes a new XMLAnalyzer with the file, opens the file, creates an XML decoder,
// analyzes the deposit tag, and checks if the deposit ID, type, and watermark match the expected values.
// If any of these steps fail, the test fails.
func TestAnalyzeDepositTagSuccess(t *testing.T) {
	// Create a temporary file for testing
	f, err := createValidXMLDepositTestFile()
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(f)

	a, err := NewXMLAnalyzer(f)
	if err != nil {
		t.Fatalf("NewXMLAnalyzer returned an error: %s", err)
	}

	err = a.OpenXMLFile()
	if err != nil {
		t.Fatalf("OpenXMLFile() returned an error: %s", err)
	}
	err = a.CreateXMLDecoder()
	if err != nil {
		t.Fatalf("CreateXMLDecoder() returned an error: %s", err)
	}

	err = a.AnalyzeDepositTag()
	if err != nil {
		t.Fatalf("AnalyzeDepositTag(), returned an error: %s", err)
	}

	if a.Deposit.ID != "20191017001" {
		t.Fatalf("AnalyzeDepositTag(), expected deposit to have id 20191017001, got: %s", a.Deposit.ID)
	}

	if a.Deposit.Type != "FULL" {
		t.Fatalf("AnalyzeDepositTag(), expected deposit to have type FULL, got: %s", a.Deposit.Type)
	}

	if a.Deposit.Watermark != "2019-10-17T00:00:00Z" {
		t.Fatalf("AnalyzeDepositTag(), expected deposit to have watermark 2019-10-17T00:00:00Z, got: %s", a.Deposit.Watermark)
	}
}

// Returns a filename to a testfile conataining a valid XML deposit. Or returns an error
func createValidXMLDepositTestFile() (string, error) {
	// Create a temporary file for testing
	f, err := os.CreateTemp("", "*test.xml")
	if err != nil {
		return "", fmt.Errorf("Failed to create temporary file: %v", err)
	}

	// Write some data to the file
	data := []byte(getValidFullDepositXMLString())
	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("Failed to write data to file: %v", err)
	}
	return f.Name(), nil
}
func TestCountLines(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "single line",
			input:    "hello world",
			expected: 0,
		},
		{
			name:     "two lines",
			input:    "hello\nworld",
			expected: 1,
		},
		{
			name:     "multiple lines",
			input:    "hello\nworld\nhow are you?",
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file for testing
			f, err := os.CreateTemp("", "*test.csv")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(f.Name())

			// Write some data to the file
			data := []byte(tc.input)
			if _, err := f.Write(data); err != nil {
				t.Fatalf("Failed to write data to file: %v", err)
			}

			file, err := os.OpenFile(f.Name(), os.O_RDONLY, 0444)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}
			actual, err := CountLines(file)
			if err != nil {
				t.Fatalf("CountLines returned an error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("CountLines(%q) = %d, expected %d", tc.input, actual, tc.expected)
			}
		})
	}
}
func TestCreateCSVFiles(t *testing.T) {
	// Create a temporary file for testing
	f, err := createValidXMLDepositTestFile()
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	a, err := NewXMLAnalyzer(f)
	if err != nil {
		t.Fatalf("NewXMLAnalyzer returned an error: %s", err)
	}
	err = a.CreateCSVFiles()
	if err != nil {
		t.Fatalf("CreateCSVFiles failed with error: %v", err)
	}

	for k, v := range a.CSVFiles {
		expectedFileName := a.GetBaseXMLFileName() + CSVFilesAndSuffixes[k]
		if v.FileName != expectedFileName {
			t.Errorf("Expected CSV file name to be %s, but got %s", expectedFileName, v.FileName)
		}
	}
}
func TestCreateCSVWriters(t *testing.T) {
	// Create a temporary file for testing
	f, err := os.CreateTemp("", "*test.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(f.Name())

	// Create a new XMLAnalyzer object
	a := &XMLAnalyzer{
		CSVFiles: map[string]CSVFile{
			"test": {
				FileName: f.Name(),
			},
		},
	}

	// Call the function being tested
	err = a.CreateCSVWriters()
	if err != nil {
		t.Fatalf("CreateCSVWriters failed with error: %v", err)
	}

	// Check that the file descriptor and CSV writer were set correctly
	csvFile := a.CSVFiles["test"]
	if csvFile.fileDescriptor == nil {
		t.Error("Expected file descriptor to be set")
	}
	if csvFile.CsvWriter == nil {
		t.Error("Expected CSV writer to be set")
	}

	// Write some data to the CSV file
	data := [][]string{
		{"header1", "header2", "header3"},
		{"data1", "data2", "data3"},
	}
	for _, row := range data {
		err = csvFile.CsvWriter.Write(row)
		if err != nil {
			t.Fatalf("Failed to write data to CSV file: %v", err)
		}
	}
	csvFile.CsvWriter.Flush()

	// Read the data from the CSV file and check that it matches the expected data
	file, err := os.Open(f.Name())
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	readData, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read data from CSV file: %v", err)
	}
	if !reflect.DeepEqual(readData, data) {
		t.Errorf("Expected data to be %v, got %v", data, readData)
	}
}
func TestFlushCSVWriters(t *testing.T) {
	// Create a temporary file for testing
	f, err := createValidXMLDepositTestFile()

	// Create a new CSV file and add it to the XMLAnalyzer's CSVFiles map
	a, err := NewXMLAnalyzer(f)
	a.CreateCSVFiles()
	a.CreateCSVWriters()

	// Flush the CSV writers
	err = a.FlushCSVWriters()
	if err != nil {
		t.Fatalf("FlushCSVWriters failed with error: %v", err)
	}

	// Check that the CSV writers were flushed
	for _, v := range a.CSVFiles {
		if v.CsvWriter != nil {
			t.Error("Expected CsvWriter to be nil after flushing")
		}
	}
}

func TestCloseCSVFiles(t *testing.T) {
	// Create temporary CSV files for testing
	file1, err := os.CreateTemp("", "*test1.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file1.Name())

	file2, err := os.CreateTemp("", "*test2.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file2.Name())

	// Create a test file
	f, err := createValidXMLDepositTestFile()
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Create a new XMLAnalyzer object and add the CSV files to its CSVFiles map
	a, err := NewXMLAnalyzer(f)
	if err != nil {
		t.Fatalf("NewXMLAnalyzer returned an error: %s", err)
	}
	a.CreateCSVFiles()
	a.CreateCSVWriters()

	// Close the CSV files
	err = a.CloseCSVFiles()
	if err != nil {
		t.Fatalf("CloseCSVFiles failed with error: %v", err)
	}

	// Check that the file descriptors are nil
	if a.CSVFiles["file1"].fileDescriptor != nil {
		t.Errorf("Expected file descriptor for file1 to be nil, got %v", a.CSVFiles["file1"].fileDescriptor)
	}
	if a.CSVFiles["file2"].fileDescriptor != nil {
		t.Errorf("Expected file descriptor for file2 to be nil, got %v", a.CSVFiles["file2"].fileDescriptor)
	}
}

func TestCountLinesInCSVFiles(t *testing.T) {
	f, err := createValidXMLDepositTestFile()
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	a, err := NewXMLAnalyzer(f)
	if err != nil {
		t.Fatalf("NewXMLAnalyzer returned an error: %s", err)
	}
	a.CreateCSVFiles()
	a.CreateCSVWriters()
	for _, v := range a.CSVFiles {
		for i := 0; i < 5; i++ {
			v.CsvWriter.Write([]string{"test"})
		}
	}
	a.FlushCSVWriters()
	a.CloseCSVFiles()

	err = a.CountLinesInCSVFiles()
	if err != nil {
		t.Fatalf("CountLinesInCSVFiles failed with error: %v", err)
	}

	for _, v := range a.CSVFiles {
		if v.LineCount != 5 {
			t.Errorf("Expected Linecount to be 5, got %d", v.LineCount)
		}
	}

}
