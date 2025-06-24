package helper

import (
	"os"
	"testing"
)

func TestFormatForConfluenceCodeMacro(t *testing.T) {
	// Test with valid file
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `package main
	
	func main() {
		println("Hello, World!")
	}`
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	result, err := FormatForConfluenceCodeMacro(tmpFile.Name())
	if err != nil {
		t.Errorf("unexpected error for valid file: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result for valid file")
	}
	// Test with nonexistent file
	_, err = FormatForConfluenceCodeMacro("nonexistent_file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
