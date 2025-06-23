package helper

import (
	"testing"
	"time"
)

func TestgetCurrentFormattedDate(t *testing.T) {
	format := "2006-01-02"
	date := getCurrentFormattedDate(format)
	// Validate the date format using time.Parse
	if _, err := time.Parse(format, date); err != nil {
		t.Errorf("date '%s' does not match expected format '%s': %v", date, format, err)
	}
	
	// Additional check for expected length
	if len(date) != len(format) {
		t.Errorf("expected date format length %d, got %d", len(format), len(date))
 	}
}
