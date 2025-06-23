# Confluence CLI Command Explanation

The `confluence_cli` command is used for interacting with Confluence's API via the command line. Here's a breakdown of the command for creating a new page in Confluence:

### 1. Create Page
```shell
confluence_cli create page \
--space-id {space id} \
--parent-page-id {parent page id} \
--title {tile of page} \
[--body-value {content}] \
[--body-value-from-file {file path}] \
[--file {attachment file path}]
```

### 2. Upload Attachment
```shell
confluence_cli create attachment \
--page-id {page id} \
--file {file path}
```

# Command Components

- `confluence_cli`: The name of the CLI tool designed for Confluence operations.

## Create Page Command

- `create page`: The action to be performed. `create` indicates creation of a new entity, and `page` specifies that the entity is a page within Confluence.

### Options for Create Page

### `--space-id {space id}`:
- **Description:** Specifies the ID of the space where the new page will be created.
- **Usage:** Replace `{space id}` with the actual space ID.
- **Required:** Yes

### `--parent-page-id {parent page id}`:
- **Description:** Indicates the ID of the parent page under which the new page will be nested.
- **Usage:** Replace `{parent page id}` with the actual parent page ID.
- **Required:** Yes

### `--title {title of page}`:
- **Description:** Sets the title for the new Confluence page.
- **Usage:** Replace `{title of page}` with the desired page title.
- **Required:** Yes

### `--body-value {content}`:
- **Description:** Sets the content for the new Confluence page directly.
- **Usage:** Replace `{content}` with the desired page content.
- **Required:** No (optional)

### `--body-value-from-file {file path to file}`:
- **Description:** Specifies the file path that contains the content for the page body.
- **Usage:** Replace `{file path to file}` with the actual path to the content file.
- **Required:** No (optional)

#### `--file {attachment file path}`:
- **Description:** Path to the file to upload as attachment to the created page.
- **Usage:** Replace `{attachment file path}` with the actual path to the file.
- **Required:** No (optional)
- **Note:** The attachment will be uploaded to the child page (nested page) that is created.

## Upload Attachment Command

- `create attachment`: The action to upload a file as an attachment to an existing Confluence page.

### Options for Upload Attachment

#### `--page-id {page id}`:
- **Description:** The Content ID (Page ID) of the Confluence page where the file will be uploaded as attachment.
- **Usage:** Replace `{page id}` with the actual page ID.
- **Required:** Yes
- **How to find Page ID:**
  - From URL: `https://<CONFLUENCE_URL>/wiki/spaces/SPACE/pages/589825/Page+Title`
  - Page ID is: `589825`
  - Or use Confluence UI: Page → "..." → "Page Information" → Copy "Page ID"

#### `--file {file path}`:
- **Description:** Path to the file to upload as attachment.
- **Usage:** Replace `{file path}` with the actual path to the file.
- **Required:** Yes
- **Supported:** All file types supported by Confluence (txt, pdf, images, etc.)

## Environment Variables.

In order to connect Your Confluence. You must configure the environments such as:   

`CONFLUENCE_URL`:   
- **Description:** your confluence link such as: `https://nimtechnology.atlassian.net`
- **Required:** Yes
`EMAIL`:
- **Description:** your email to access Confluence API such as: `mr.nim@nimtechnology.com`
- **Required:** Yes
`API_TOKEN`:
- **Description:** The Token is used to access Confluence API such as: `XXXXXXXXXXXXXX`
- **Refer to:** https://nimtechnology.com/2024/01/05/confluence-integrate-with-confluence-by-api/
- **Required:** Yes

## How build Binary file.

### On Windows
```shell
go build -o confluence_cli.exe
```

### On Linux/Mac:
```shell
go build -o confluence_cli
```

## Usage Examples

### Create a Page with Content:
```shell
# On Windows
.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value "This is the page content"

# On Linux/Mac
./confluence_cli create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value "This is the page content"
```

### Create a Page with Content from File:
```shell
# On Windows
.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value-from-file result.txt

# On Linux/Mac
./confluence_cli create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value-from-file result.txt
```

### Create a Page with Empty Content:
```shell
# On Windows
.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title "Test Page"

# On Linux/Mac
./confluence_cli create page --space-id 98432 --parent-page-id 589825 --title "Test Page"
```

### Create a Page with Content and Attachment:
```shell
# On Windows
.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value "Content" --file /path/to/attachment.txt

# On Linux/Mac
./confluence_cli create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --body-value "Content" --file /path/to/attachment.txt
```

### Create a Page with Empty Content and Attachment:
```shell
# On Windows
.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --file /path/to/attachment.txt

# On Linux/Mac
./confluence_cli create page --space-id 98432 --parent-page-id 589825 --title "Test Page" --file /path/to/attachment.txt
```

### Upload Attachment to Existing Page:
```shell
# On Windows
.\confluence_cli.exe create attachment --page-id 589825 --file /path/to/file

# On Linux/Mac
./confluence_cli create attachment --page-id 589825 --file /path/to/file
```

## API Referenc

This tool uses the following Confluence REST APIs:
- **Create Page:** `POST /wiki/api/v2/pages`
- **Upload Attachment:** `POST /wiki/rest/api/content/{pageId}/child/attachment`
- **Get Pages by Title:** `GET /wiki/api/v2/pages?title={title}`

## Notes

- When creating a page with `--file` flag, the attachment will be uploaded to the child page (nested page) that is created, not the parent page.
- Content flags (`--body-value` and `--body-value-from-file`) are optional. If neither is provided, the page will be created with empty content.
- The tool automatically adds date prefixes to page titles in the format `[YYYY-MM]` and `[YYYY-MM-DD HH:MM:SS]` for better organization.