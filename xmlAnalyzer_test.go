package ryde

import (
	"fmt"
	"os"
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
