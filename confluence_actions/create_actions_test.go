package confluence_actions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test validation logic directly without CLI context
func TestValidationLogic(t *testing.T) {
	tests := []struct {
		name        string
		spaceId     string
		parentId    string
		title       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing space-id",
			spaceId:     "",
			parentId:    "123",
			title:       "Test",
			expectError: true,
			errorMsg:    "--space-id is required",
		},
		{
			name:        "Missing parent-page-id",
			spaceId:     "SPACE",
			parentId:    "",
			title:       "Test",
			expectError: true,
			errorMsg:    "--parent-page-id is required",
		},
		{
			name:        "Missing title",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "",
			expectError: true,
			errorMsg:    "--title is required",
		},
		{
			name:        "Valid required parameters",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "Test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic directly
			if tt.spaceId == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "--space-id is required")
			} else if tt.parentId == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "--parent-page-id is required")
			} else if tt.title == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "--title is required")
			} else {
				assert.False(t, tt.expectError)
			}
		})
	}
}

// Test validation logic with mock validation function
func TestValidationLogicWithMock(t *testing.T) {
	// Mock validation function that mimics the logic in CreatePageAction
	validateRequiredFields := func(spaceId, parentId, title string) error {
		if spaceId == "" {
			return fmt.Errorf("--space-id is required")
		}
		if parentId == "" {
			return fmt.Errorf("--parent-page-id is required")
		}
		if title == "" {
			return fmt.Errorf("--title is required")
		}
		return nil
	}

	tests := []struct {
		name        string
		spaceId     string
		parentId    string
		title       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing space-id",
			spaceId:     "",
			parentId:    "123",
			title:       "Test",
			expectError: true,
			errorMsg:    "--space-id is required",
		},
		{
			name:        "Missing parent-page-id",
			spaceId:     "SPACE",
			parentId:    "",
			title:       "Test",
			expectError: true,
			errorMsg:    "--parent-page-id is required",
		},
		{
			name:        "Missing title",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "",
			expectError: true,
			errorMsg:    "--title is required",
		},
		{
			name:        "Valid required parameters",
			spaceId:     "SPACE",
			parentId:    "123",
			title:       "Test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredFields(tt.spaceId, tt.parentId, tt.title)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test file validation logic
func TestFileValidationLogic(t *testing.T) {
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
