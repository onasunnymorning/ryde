package ryde

import (
	"os"
	"testing"
)

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
}
