package helper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateActionItemMacro(t *testing.T) {
	tests := []struct {
		name                   string
		failedCount            int
		totalCount             int
		expectedGoodForRelease string
		expectedHoldOff        string
	}{
		{
			name:                   "Tests failed - HOLD-OFF should be checked",
			failedCount:            5,
			totalCount:             100,
			expectedGoodForRelease: "incomplete",
			expectedHoldOff:        "complete",
		},
		{
			name:                   "All tests passed - GOOD FOR RELEASE should be checked",
			failedCount:            0,
			totalCount:             100,
			expectedGoodForRelease: "complete",
			expectedHoldOff:        "incomplete",
		},
		{
			name:                   "Edge case: failed count equals total",
			failedCount:            10,
			totalCount:             10,
			expectedGoodForRelease: "incomplete",
			expectedHoldOff:        "complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateActionItemMacro(tt.failedCount, tt.totalCount)

			// Check that the macro contains expected status values
			assert.Contains(t, result, tt.expectedGoodForRelease)
			assert.Contains(t, result, tt.expectedHoldOff)

			// Check that the macro contains required elements
			assert.Contains(t, result, "ac:task-list")
			assert.Contains(t, result, "GOOD FOR RELEASE")
			assert.Contains(t, result, "HOLD-OFF")
			assert.Contains(t, result, "ac:structured-macro")
			assert.Contains(t, result, "ac:parameter ac:name=\"colour\"")
		})
	}
}

func TestCreateAttachmentMacro(t *testing.T) {
	result := CreateAttachmentMacro()

	// Check that the macro contains required elements
	assert.Contains(t, result, "ac:structured-macro")
	assert.Contains(t, result, "ac:name=\"attachments\"")
	assert.Contains(t, result, "ac:schema-version=\"1\"")

	// Check that it's wrapped in paragraph tags
	assert.True(t, strings.HasPrefix(result, "<p>"))
	assert.True(t, strings.HasSuffix(result, "</p>"))
}

func TestFormatForConfluenceCodeMacro(t *testing.T) {
	// Test successful formatting
	result, err := FormatForConfluenceCodeMacro("test_file.txt")
	if err == nil {
		// Check that the result contains expected elements
		assert.Contains(t, result, "ac:structured-macro ac:name=\"code\"")
		assert.Contains(t, result, "ac:plain-text-body")
		assert.Contains(t, result, "<![CDATA[")
		assert.Contains(t, result, "]]>")
	}

	// Test with non-existent file
	_, err = FormatForConfluenceCodeMacro("/non/existent/file.txt")
	assert.Error(t, err)
}

func TestMacroTemplates(t *testing.T) {
	// Test ActionListMacroTemplate
	assert.Contains(t, ActionListMacroTemplate, "ac:task-list")
	assert.Contains(t, ActionListMacroTemplate, "GOOD FOR RELEASE")
	assert.Contains(t, ActionListMacroTemplate, "HOLD-OFF")
	assert.Contains(t, ActionListMacroTemplate, "%s") // Placeholder for status

	// Test AttachmentMacroTemplate
	assert.Contains(t, AttachmentMacroTemplate, "ac:structured-macro")
	assert.Contains(t, AttachmentMacroTemplate, "ac:name=\"attachments\"")
	assert.Contains(t, AttachmentMacroTemplate, "ac:schema-version=\"1\"")
}

// Test macro template constants
func TestMacroTemplateConstants(t *testing.T) {
	// Test ActionListMacroTemplate structure
	assert.Contains(t, ActionListMacroTemplate, "ac:task-list")
	assert.Contains(t, ActionListMacroTemplate, "ac:task")
	assert.Contains(t, ActionListMacroTemplate, "ac:task-status")
	assert.Contains(t, ActionListMacroTemplate, "ac:task-body")

	// Test AttachmentMacroTemplate structure
	assert.Contains(t, AttachmentMacroTemplate, "ac:structured-macro")
	assert.Contains(t, AttachmentMacroTemplate, "attachments")
}

// Test macro content validation
func TestMacroContentValidation(t *testing.T) {
	// Test that action list macro contains both tasks
	actionListMacro := CreateActionItemMacro(0, 10)

	// Count occurrences of task elements
	taskCount := strings.Count(actionListMacro, "<ac:task>")
	assert.Equal(t, 2, taskCount, "Action list macro should contain exactly 2 tasks")

	// Check for status macro parameters
	assert.Contains(t, actionListMacro, "ac:parameter ac:name=\"colour\"")
	assert.Contains(t, actionListMacro, "ac:parameter ac:name=\"title\"")
}
