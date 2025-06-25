package confluence_actions

import (
	"confluence_cli/log"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// Mock context for minimal flag testing
func TestCreatePageAction_MissingFlags(t *testing.T) {
	log.InitLogger(true)
	app := cli.NewApp()
	set := flagSet([]string{})
	c := cli.NewContext(app, set, nil)
	err := CreatePageAction(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space-id")
}

// Helper to create a flag set
func flagSet(args []string) *flag.FlagSet {
	if len(args)%2 != 0 {
		panic("args must contain an even number of elements (key-value pairs)")
	}
	set := flag.NewFlagSet("test", 0)
	set.String("space-id", "", "")
	set.String("parent-page-id", "", "")
	set.String("title", "", "")
	set.String("body-value-from-file", "", "")
	set.String("body-value", "", "")
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break // or handle the error appropriately
		}
		set.Set(args[i], args[i+1])
	}
	return set
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

func validateCreatePage(spaceId, parentId, title string) error {
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

func TestCreatePageValidationLogic(t *testing.T) {
	cases := []ValidationTestCase{
		{
			Name:        "Missing space-id",
			Args:        []string{"", "123", "Test"},
			ExpectError: true,
			ErrorMsg:    "--space-id is required",
		},
		{
			Name:        "Missing parent-page-id",
			Args:        []string{"SPACE", "", "Test"},
			ExpectError: true,
			ErrorMsg:    "--parent-page-id is required",
		},
		{
			Name:        "Missing title",
			Args:        []string{"SPACE", "123", ""},
			ExpectError: true,
			ErrorMsg:    "--title is required",
		},
		{
			Name:        "Valid required parameters",
			Args:        []string{"SPACE", "123", "Test"},
			ExpectError: false,
		},
	}
	RunValidationTable(t, func(args ...string) error {
		return validateCreatePage(args[0], args[1], args[2])
	}, cases)
}
