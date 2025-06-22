package helper

import (
	"testing"
)

func TestFormatForConfluenceCodeMacro(t *testing.T) {
	_, err := FormatForConfluenceCodeMacro("nonexistent_file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
