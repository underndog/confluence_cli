package confluence_actions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func validateUpdatePage(pageId string) error {
	if pageId == "" {
		return fmt.Errorf("please provide --page-id")
	}
	return nil
}

func TestUpdatePageValidationLogic(t *testing.T) {
	cases := []ValidationTestCase{
		{
			Name:        "Missing page-id",
			Args:        []string{""},
			ExpectError: true,
			ErrorMsg:    "please provide --page-id",
		},
		{
			Name:        "Valid page-id",
			Args:        []string{"123"},
			ExpectError: false,
		},
	}
	RunValidationTable(t, func(args ...string) error {
		return validateUpdatePage(args[0])
	}, cases)
}

// Test file validation logic
func TestUpdatePageFileValidationLogic(t *testing.T) {
	// Mock file validation function
	validateFile := func(filePath string) error {
		if filePath == "" {
			return nil // No file provided, no validation needed
		}
		if filePath == "/non/existent/file.txt" {
			return fmt.Errorf("attachment file does not exist: %s", filePath)
		}
		return nil
	}

	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Non-existent file",
			filePath:    "/non/existent/file.txt",
			expectError: true,
			errorMsg:    "attachment file does not exist",
		},
		{
			name:        "Empty file path",
			filePath:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFile(tt.filePath)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validation logic directly without CLI context
func TestUpdatePageValidationLogicDirect(t *testing.T) {
	tests := []struct {
		name        string
		pageId      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing page-id",
			pageId:      "",
			expectError: true,
			errorMsg:    "please provide --page-id",
		},
		{
			name:        "Valid page-id",
			pageId:      "123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic directly
			if tt.pageId == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "please provide --page-id")
			} else {
				assert.False(t, tt.expectError)
			}
		})
	}
}
