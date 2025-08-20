package helper

import (
	"fmt"
	"os"
	"strings"
)

// function to format content for Confluence code macro
func FormatForConfluenceCodeMacro(filePath string) (string, error) {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Format as XML for Confluence code macro
	result := fmt.Sprintf(`<ac:structured-macro ac:name="code" ac:schema-version="1">
		<ac:parameter ac:name="language">text</ac:parameter>
		<ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body>
		</ac:structured-macro>`, string(content))

	return result, nil
}

// function to create attachment macro (for enabling existing macros)
func CreateAttachmentMacro() string {
	return `<p><ac:structured-macro ac:name="attachments" ac:schema-version="1"></ac:structured-macro></p>`
}

// function to check if content contains attachment macro
func HasAttachmentMacro(content string) bool {
	return strings.Contains(content, `ac:name="attachments"`)
}

// function to check if content contains action list macro
func HasActionListMacro(content string) bool {
	return strings.Contains(content, `ac:task-list`)
}

// function to enable macros if they don't exist
func EnableMacrosIfMissing(content string) string {
	result := content

	// Enable attachment macro if missing
	if !HasAttachmentMacro(content) {
		result += "\n" + CreateAttachmentMacro()
	}

	return result
}
