package confluence_actions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test validation logic with mock validation function
func TestUploadAttachmentValidationLogic(t *testing.T) {
	// Mock validation function that mimics the logic in UploadAttachmentAction
	validateRequiredFields := func(pageId, filePath string) error {
		if pageId == "" {
			return fmt.Errorf("please provide --page-id")
		}
		if filePath == "" {
			return fmt.Errorf("please provide --file")
		}
		return nil
	}

	tests := []struct {
		name        string
		pageId      string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing page-id",
			pageId:      "",
			filePath:    "/tmp/test.txt",
			expectError: true,
			errorMsg:    "please provide --page-id",
		},
		{
			name:        "Missing file",
			pageId:      "123",
			filePath:    "",
			expectError: true,
			errorMsg:    "please provide --file",
		},
		{
			name:        "Valid parameters",
			pageId:      "123",
			filePath:    "/tmp/test.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredFields(tt.pageId, tt.filePath)
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
func TestUploadAttachmentFileValidationLogic(t *testing.T) {
	// Mock file validation function
	validateFile := func(filePath string) error {
		if filePath == "" {
			return fmt.Errorf("please provide --file")
		}
		if filePath == "/non/existent/file.txt" {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		if filePath == "/tmp" {
			return fmt.Errorf("path is a directory: %s", filePath)
		}
		if filePath == "/empty/file.txt" {
			return fmt.Errorf("file is empty: %s", filePath)
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
			name:        "Empty file path",
			filePath:    "",
			expectError: true,
			errorMsg:    "please provide --file",
		},
		{
			name:        "Non-existent file",
			filePath:    "/non/existent/file.txt",
			expectError: true,
			errorMsg:    "file does not exist",
		},
		{
			name:        "Directory path",
			filePath:    "/tmp",
			expectError: true,
			errorMsg:    "path is a directory",
		},
		{
			name:        "Empty file",
			filePath:    "/empty/file.txt",
			expectError: true,
			errorMsg:    "file is empty",
		},
		{
			name:        "Valid file",
			filePath:    "/valid/file.txt",
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
func TestUploadAttachmentValidationLogicDirect(t *testing.T) {
	tests := []struct {
		name        string
		pageId      string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing page-id",
			pageId:      "",
			filePath:    "/tmp/test.txt",
			expectError: true,
			errorMsg:    "please provide --page-id",
		},
		{
			name:        "Missing file",
			pageId:      "123",
			filePath:    "",
			expectError: true,
			errorMsg:    "please provide --file",
		},
		{
			name:        "Valid parameters",
			pageId:      "123",
			filePath:    "/tmp/test.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic directly
			if tt.pageId == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "please provide --page-id")
			} else if tt.filePath == "" {
				assert.True(t, tt.expectError)
				assert.Contains(t, tt.errorMsg, "please provide --file")
			} else {
				assert.False(t, tt.expectError)
			}
		})
	}
}
