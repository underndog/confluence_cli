package helper

import (
	"testing"
)

func TestGetCurrentFormattedDate(t *testing.T) {
	format := "2006-01-02"
	date := getCurrentFormattedDate(format)
	if len(date) != len("2006-01-02") {
		t.Errorf("expected date format length %d, got %d", len("2006-01-02"), len(date))
	}
}
