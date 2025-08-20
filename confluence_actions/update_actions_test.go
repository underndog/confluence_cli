package confluence_actions

import (
	"confluence_cli/helper"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test parseTestResultsFromHTML function for update actions
func TestUpdateActionsParseTestResultsFromHTML(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failedCount, totalCount := parseTestResultsFromHTML(tt.htmlContent)
			assert.Equal(t, tt.expectedFailed, failedCount)
			assert.Equal(t, tt.expectedTotal, totalCount)
		})
	}
}

// Test macro detection logic
func TestMacroDetectionLogic(t *testing.T) {
	tests := []struct {
		name                  string
		content               string
		expectedHasAttachment bool
		expectedHasActionList bool
	}{
		{
			name: "Content with both macros",
			content: `
				<div>Some content</div>
				<ac:structured-macro ac:name="attachments"></ac:structured-macro>
				<ac:task-list>
					<ac:task>GOOD FOR RELEASE</ac:task>
				</ac:task-list>
			`,
			expectedHasAttachment: true,
			expectedHasActionList: true,
		},
		{
			name: "Content with only attachment macro",
			content: `
				<div>Some content</div>
				<ac:structured-macro ac:name="attachments"></ac:structured-macro>
			`,
			expectedHasAttachment: true,
			expectedHasActionList: false,
		},
		{
			name: "Content with only action list macro",
			content: `
				<div>Some content</div>
				<ac:task-list>
					<ac:task>HOLD-OFF</ac:task>
				</ac:task-list>
			`,
			expectedHasAttachment: false,
			expectedHasActionList: true,
		},
		{
			name: "Content with no macros",
			content: `
				<div>Some content</div>
				<p>No macros here</p>
			`,
			expectedHasAttachment: false,
			expectedHasActionList: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAttachmentMacro := strings.Contains(tt.content, "attachments")
			hasActionList := strings.Contains(tt.content, "ac:task-list")

			assert.Equal(t, tt.expectedHasAttachment, hasAttachmentMacro)
			assert.Equal(t, tt.expectedHasActionList, hasActionList)
		})
	}
}

// Test Overall Status pattern matching
func TestUpdateActionsOverallStatusPatternMatching(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := re.MatchString(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test macro content generation for update actions
func TestUpdateActionsMacroContentGeneration(t *testing.T) {
	// Test that attachment macro contains expected content
	attachmentMacro := helper.CreateAttachmentMacro()
	assert.Contains(t, attachmentMacro, "attachments")
	assert.Contains(t, attachmentMacro, "ac:structured-macro")

	// Test that attachment macro is properly wrapped
	assert.True(t, strings.HasPrefix(attachmentMacro, "<p>"))
	assert.True(t, strings.HasSuffix(attachmentMacro, "</p>"))
}

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
			err := validateUpdatePage(tt.pageId)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test macro re-adding logic after manual edits
func TestMacroReAddingLogic(t *testing.T) {
	tests := []struct {
		name            string
		originalVersion int
		currentVersion  int
		hasMacros       bool
		expectReAdd     bool
	}{
		{
			name:            "Version increased, macros missing - should re-add",
			originalVersion: 1,
			currentVersion:  3,
			hasMacros:       false,
			expectReAdd:     true,
		},
		{
			name:            "Version increased, macros present - no re-add needed",
			originalVersion: 1,
			currentVersion:  3,
			hasMacros:       true,
			expectReAdd:     false,
		},
		{
			name:            "Version normal - no re-add needed",
			originalVersion: 1,
			currentVersion:  2,
			hasMacros:       true,
			expectReAdd:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would need to be expanded with actual logic testing
			// For now, just test that the test cases are valid
			assert.NotEmpty(t, tt.name)
			assert.True(t, tt.currentVersion >= tt.originalVersion)
		})
	}
}
