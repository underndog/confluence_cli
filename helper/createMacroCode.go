package helper

import (
	"fmt"
	"os"
)

/***
Refer to links: https://community.atlassian.com/t5/Confluence-questions/Code-Macro-via-Confluence-REST-API/qaq-p/2097123
*/

// function to format content for Confluence code macro
func FormatForConfluenceCodeMacro(filePath string) (string, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Format as XML for Confluence code macro
	xmlContent := fmt.Sprintf(
		"<ac:structured-macro ac:name=\"code\">\n"+
			"  <ac:plain-text-body><![CDATA[\n"+
			"    %s\n"+ // you can see some space at begin of the content, you can adjust at here
			"  ]]></ac:plain-text-body>\n"+
			"</ac:structured-macro>",
		string(data),
	)

	return xmlContent, nil
}

const (
	// Macro templates
	ActionListMacroTemplate = `<ac:task-list>
		<ac:task>
			<ac:task-status>%s</ac:task-status>
			<ac:task-body><ac:structured-macro ac:name="status" ac:schema-version="1">
				<ac:parameter ac:name="colour">Green</ac:parameter>
				<ac:parameter ac:name="title">GOOD FOR RELEASE</ac:parameter>
			</ac:structured-macro></ac:task-body>
		</ac:task>
		<ac:task>
			<ac:task-status>%s</ac:task-status>
			<ac:task-body><ac:structured-macro ac:name="status" ac:schema-version="1">
				<ac:parameter ac:name="colour">Red</ac:parameter>
				<ac:parameter ac:name="title">HOLD-OFF</ac:parameter>
			</ac:structured-macro></ac:task-body>
		</ac:task>
	</ac:task-list>`

	AttachmentMacroTemplate = `<p><ac:structured-macro ac:name="attachments" ac:schema-version="1"></ac:structured-macro></p>`
)

// function to create Task List macro for Overall Status with dynamic status based on test results
func CreateActionItemMacro(failedCount, totalCount int) string {
	if failedCount > 0 {
		// Tests failed - HOLD-OFF should be checked
		return fmt.Sprintf(ActionListMacroTemplate, "incomplete", "complete")
	} else {
		// All tests passed - GOOD FOR RELEASE should be checked
		return fmt.Sprintf(ActionListMacroTemplate, "complete", "incomplete")
	}
}

// function to create attachment macro
func CreateAttachmentMacro() string {
	return AttachmentMacroTemplate
}
