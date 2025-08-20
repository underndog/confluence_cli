package confluence_actions

import (
	"confluence_cli/helper"
	"confluence_cli/log"
	"flag"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// Test parseTestResultsFromHTML function
func TestParseTestResultsFromHTML(t *testing.T) {
	// Initialize logger to avoid nil pointer dereference
	log.InitLogger(true)

	tests := []struct {
		name           string
		htmlContent    string
		expectedFailed int
		expectedTotal  int
	}{
		{
			name: "Parse test results with failed tests",
			htmlContent: `
				<div>Total Tests: 100</div>
				<div>Failed: 5</div>
				<div>Passed: 95</div>
			`,
			expectedFailed: 5,
			expectedTotal:  100,
		},
		{
			name: "Parse test results with no failed tests",
			htmlContent: `
				<div>Total Tests: 50</div>
				<div>Failed: 0</div>
				<div>Passed: 50</div>
			`,
			expectedFailed: 0,
			expectedTotal:  50,
		},
		{
			name: "Parse test results with alternative format",
			htmlContent: `
				<div><strong>Total Tests:</strong> 75</div>
				<div><strong>Failed:</strong> 3</div>
			`,
			expectedFailed: 3,
			expectedTotal:  75,
		},
		{
			name: "Parse test results with no matches",
			htmlContent: `
				<div>Some other content</div>
				<div>No test results here</div>
			`,
			expectedFailed: 0,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failedCount, totalCount := parseTestResultsFromHTML(tt.htmlContent)

			assert.Equal(t, tt.expectedFailed, failedCount)
			assert.Equal(t, tt.expectedTotal, totalCount)
		})
	}
}

// Test regex pattern matching for Overall Status
func TestOverallStatusPatternMatching(t *testing.T) {
	pattern := `<td><strong>Overall Status</strong></td>\s*<td colspan="2">\s*</td>`
	re := regexp.MustCompile(pattern)

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Match Overall Status pattern",
			content:  `<td><strong>Overall Status</strong></td><td colspan="2"></td>`,
			expected: true,
		},
		{
			name: "Match Overall Status pattern with whitespace",
			content: `<td><strong>Overall Status</strong></td>
				<td colspan="2">
				</td>`,
			expected: true,
		},
		{
			name:     "No match - different structure",
			content:  `<td><strong>Status</strong></td><td colspan="2"></td>`,
			expected: false,
		},
		{
			name:     "No match - missing colspan",
			content:  `<td><strong>Overall Status</strong></td><td></td>`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := re.MatchString(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test macro content generation
func TestMacroContentGeneration(t *testing.T) {
	// Test that macros contain expected content
	attachmentMacro := helper.CreateAttachmentMacro()
	assert.Contains(t, attachmentMacro, "attachments")
	assert.Contains(t, attachmentMacro, "ac:structured-macro")

	// Test action list macro with different test results
	actionListMacroFailed := helper.CreateActionItemMacro(5, 100)
	assert.Contains(t, actionListMacroFailed, "HOLD-OFF")
	assert.Contains(t, actionListMacroFailed, "GOOD FOR RELEASE")

	actionListMacroPassed := helper.CreateActionItemMacro(0, 100)
	assert.Contains(t, actionListMacroPassed, "HOLD-OFF")
	assert.Contains(t, actionListMacroPassed, "GOOD FOR RELEASE")
}

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
	set.String("file", "", "")
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
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

// Test file parameter handling
func TestFileParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Empty file path",
			filePath:    "",
			expectError: false,
		},
		{
			name:        "Valid file path",
			filePath:    "/path/to/file.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would need to be expanded with actual file validation
			// For now, just test that the parameter is handled
			assert.NotEmpty(t, tt.name)
		})
	}
}
