If you want push content into a block code:   

```shell
# Set the title to "Integration Test Results" followed by the current date and time in YYYY-MM-DD hh:mm:ss format 
TITLE="Integration Test Results $(date '+%Y-%m-%d %H:%M:%S')"

FILE_PATH='result.txt' # path to your file

# Read the content of the file into a variable 
FILE_CONTENT=$(<"$FILE_PATH")

# Escape XML special characters in file content
XML_ESCAPED_CONTENT=$(echo "$FILE_CONTENT" | sed -e 's/&/\&amp;/g' -e 's/</\&lt;/g' -e 's/>/\&gt;/g')

# Wrap the XML escaped content within the Confluence structured macro syntax
ESCAPED_CONTENT="<ac:structured-macro ac:name=\"code\"><ac:plain-text-body><![CDATA[$XML_ESCAPED_CONTENT]]></ac:plain-text-body></ac:structured-macro>"

# Escape for JSON
JSON_ESCAPED_CONTENT=$(echo "$ESCAPED_CONTENT" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e ':a' -e 'N' -e '$!ba' -e 's/\n/\\n/g')

# cURL command 
curl --request POST \
  --url "https://nimtechnology.atlassian.net/wiki/api/v2/pages" \
  --user "your_email:your_api_token" \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{ 
    "spaceId": "98432", 
    "status": "current", 
    "title": "'"${TITLE}"'", 
    "parentId": "589825", 
    "body": { 
      "representation": "storage", 
      "value": "'"${JSON_ESCAPED_CONTENT}"'" 
    } 
  }'
```