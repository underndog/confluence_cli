package helper

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write test content to temp file
	testContent := "This is test content for code macro"
	_, err = tempFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Test successful formatting with real file
	result, err := FormatForConfluenceCodeMacro(tempFile.Name())
	if err != nil {
		t.Fatalf("FormatForConfluenceCodeMacro failed: %v", err)
	}

	// Check that the result contains expected elements
	assert.Contains(t, result, "ac:structured-macro ac:name=\"code\"")
	assert.Contains(t, result, "ac:plain-text-body")
	assert.Contains(t, result, "<![CDATA[")
	assert.Contains(t, result, "]]>")
	assert.Contains(t, result, testContent)

	// Test with non-existent file
	_, err = FormatForConfluenceCodeMacro("/non/existent/file.txt")
	assert.Error(t, err)
}

func TestMacroDetectionFunctions(t *testing.T) {
	// Test HasAttachmentMacro
	contentWithAttachment := `<ac:structured-macro ac:name="attachments"></ac:structured-macro>`
	contentWithoutAttachment := `<p>Some content without attachment macro</p>`

	assert.True(t, HasAttachmentMacro(contentWithAttachment))
	assert.False(t, HasAttachmentMacro(contentWithoutAttachment))

	// Test HasActionListMacro
	contentWithActionList := `<ac:task-list>Some task list content</ac:task-list>`
	contentWithoutActionList := `<p>Some content without action list macro</p>`

	assert.True(t, HasActionListMacro(contentWithActionList))
	assert.False(t, HasActionListMacro(contentWithoutActionList))
}

func TestEnableMacrosIfMissing(t *testing.T) {
	// Test content without attachment macro
	contentWithoutMacro := `<p>Content without macros</p>`
	result := EnableMacrosIfMissing(contentWithoutMacro)

	// Should add attachment macro
	assert.Contains(t, result, "ac:name=\"attachments\"")
	assert.Contains(t, result, contentWithoutMacro)

	// Test content with existing attachment macro
	contentWithMacro := `<p>Content with <ac:structured-macro ac:name="attachments"></ac:structured-macro></p>`
	resultWithExisting := EnableMacrosIfMissing(contentWithMacro)

	// Should not add duplicate macro
	attachmentCount := strings.Count(resultWithExisting, "ac:name=\"attachments\"")
	assert.Equal(t, 1, attachmentCount, "Should not add duplicate attachment macro")
}

func TestMacroContentValidation(t *testing.T) {
	// Test attachment macro structure
	attachmentMacro := CreateAttachmentMacro()
	assert.Contains(t, attachmentMacro, "ac:structured-macro")
	assert.Contains(t, attachmentMacro, "ac:name=\"attachments\"")
	assert.Contains(t, attachmentMacro, "ac:schema-version=\"1\"")

	// Test code macro formatting with temp file
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString("Test content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	codeMacro, err := FormatForConfluenceCodeMacro(tempFile.Name())
	if err != nil {
		t.Fatalf("FormatForConfluenceCodeMacro failed: %v", err)
	}
	assert.Contains(t, codeMacro, "ac:structured-macro ac:name=\"code\"")
	assert.Contains(t, codeMacro, "ac:plain-text-body")
}
