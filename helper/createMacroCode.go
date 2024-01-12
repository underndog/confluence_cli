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
