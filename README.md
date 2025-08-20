# Confluence CLI - Complete Command Reference

The `confluence_cli` is a powerful command-line tool designed for seamless interaction with Confluence's API. It provides comprehensive functionality for creating pages, updating content, managing attachments, and automatically embedding macros including attachment macros and action list macros.

## Table of Contents

- [Quick Start](#quick-start)
- [Commands Overview](#commands-overview)
- [Environment Setup](#environment-setup)
- [Create Page Command](#create-page-command)
- [Upload Attachment Command](#upload-attachment-command)
- [Update Page Command](#update-page-command)
- [Automatic Macro Features](#automatic-macro-features)
- [Complete Usage Examples](#complete-usage-examples)
- [API Reference](#api-reference)
- [Troubleshooting](#troubleshooting)
- [Notes and Best Practices](#notes-and-best-practices)

## 🆕 Latest Updates

**NEW FEATURES**:

1. **Attachment Macro Management**: Automatically enables attachment macros when files are uploaded
2. **Smart Macro Detection**: Detects missing attachment macros and action list macros and re-enables them automatically
3. **Page Version Management**: Handles page versions correctly during updates
4. **File Upload Integration**: Seamlessly uploads and embeds file attachments
5. **Helper Function System**: Comprehensive helper functions for macro creation and management

## Quick Start

### 1\. Environment Setup

```plaintext
# Load environment variables
source confluence_env.sh

# Or set manually
export CONFLUENCE_URL="https://sample.atlassian.net/"
export EMAIL="your-email@gmail.com"
export API_TOKEN="your-api-token"
```

### Required Environment Variables

#### `CONFLUENCE_URL`

*   **Description:** Your Confluence instance URL
*   **Format:** `https://your-domain.atlassian.net/`
*   **Required:** Yes

#### `EMAIL`

*   **Description:** Your Confluence account email
*   **Format:** `your-email@company.com`
*   **Required:** Yes

#### `API_TOKEN`

*   **Description:** Your Confluence API token
*   **Format:** `ATATT3x...` (long string)
*   **Required:** Yes

### 2\. Build the Tool

```plaintext
cd confluence_cli
go build -o confluence_cli
chmod +x confluence_cli
```

## 3\. Docker

Command format:

```plaintext
docker run -it --rm \
-e API_TOKEN="<token>" \
-e CONFLUENCE_URL='https://opswat.atlassian.net' \
-e EMAIL="<email>" \
--name confluence-cli quay.io/underndog/confluence_cli:latest create page --space-id "<space_id>" --parent-page-id "<parent_page_id>" --title "<Tile of Page>" --body-value <content>
```

## Commands Overview

The tool provides three main commands:

1. **`create page`** - Create new Confluence pages with optional content, attachments, and automatic macro embedding
2. **`update page`** - Update existing pages with new content and/or attachments, ensuring macros remain active
3. **`upload attachment`** - Upload attachments to existing pages without changing the page content

## Create Page Command

Creates new Confluence pages with optional content, attachments, and automatic macro embedding.

### Basic Syntax

```plaintext
confluence_cli create page \
  --space-id {space id} \
  --parent-page-id {parent page id} \
  --title {title of page} \
  [--body-value {content}] \
  [--body-value-from-file {file path}] \
  [--file {attachment file path}]
```

### Required Parameters

#### `--space-id {space id}`

*   **Description:** Specifies the ID of the space where the new page will be created
*   **Usage:** Replace `{space id}` with the actual space ID
*   **Required:** Yes
*   **How to find:** Space settings → Space details. Or call API

#### `--parent-page-id {parent page id}`

*   **Description:** Indicates the ID of the parent page under which the new page will be nested
*   **Usage:** Replace `{parent page id}` with the actual parent page ID
*   **Required:** Yes
*   **How to find:** Page URL or Page Information. Or call API

#### `--title {title of page}`

*   **Description:** Sets the title for the new Confluence page
*   **Usage:** Replace `{title of page}` with the desired page title
*   **Required:** Yes
*   **Note:** Tool automatically adds date prefixes: `[YYYY-MM]` and `[YYYY-MM-DD HH:MM:SS]`

### Optional Parameters

#### `--body-value {content}`

*   **Description:** Sets the content for the page directly
*   **Usage:** Replace `{content}` with the desired page content
*   **Required:** No (optional)
*   **Example:** `--body-value "This is my page content"`

#### `--body-value-from-file {file path}`

*   **Description:** Specifies the file path that contains the content for the page body
*   **Usage:** Replace `{file path}` with the actual path to the content file
*   **Required:** No (optional)
*   **Example:** `--body-value-from-file "/path/to/content.html"`

#### `--file {attachment file path}`

*   **Description:** Path to the file to upload as attachment. The file will always be embedded and displayed on the page.
*   **Usage:** Replace `{attachment file path}` with the actual path to the file (use absolute path for best results)
*   **Required:** No (optional)
*   **Note:** For create page, the attachment will be uploaded to the child page (nested page) that is created
*   **Example:** `--file "/home/user/documents/report.pdf"`

## Upload Attachment Command

Uploads attachments to existing Confluence pages without modifying the page content.

### Basic Syntax

```plaintext
confluence_cli upload attachment \
  --page-id {page id} \
  --file {attachment file path}
```

### Required Parameters

#### `--page-id {page id}`

*   **Description:** ID of the page to upload attachment to
*   **Usage:** Replace `{page id}` with the actual page ID
*   **Required:** Yes
*   **How to find:** Page URL or Page Information or call API

#### `--file {attachment file path}`

*   **Description:** Path to the file to upload as attachment
*   **Usage:** Replace `{attachment file path}` with the actual path to the file
*   **Required:** Yes
*   **Note:** File will be uploaded as attachment but not embedded in the page content
*   **Example:** `--file "/home/user/documents/report.pdf"`

## Update Page Command

Updates existing Confluence pages with new content and/or attachments, ensuring macros remain active.

### Basic Syntax

```plaintext
confluence_cli update page \
  --page-id {page id} \
  [--body-value {content}] \
  [--body-value-from-file {file path}] \
  [--file {attachment file path}]
```

### Required Parameters

#### `--page-id {page id}`

*   **Description:** ID of the page to update
*   **Usage:** Replace `{page id}` with the actual page ID
*   **Required:** Yes
*   **How to find:** Page URL or Page Information or call API

### Optional Parameters

#### `--body-value {content}`

*   **Description:** Sets the content for the page directly
*   **Usage:** Replace `{content}` with the desired page content
*   **Required:** No (optional)
*   **Example:** `--body-value "This is my page content"`

#### `--body-value-from-file {file path}`
- **Description:** Specifies the file path that contains the content for the page body
- **Usage:** Replace `{file path}` with the actual path to the content file
- **Required:** No (optional)
- **Example:** `--body-value-from-file "/path/to/updated-content.txt"`

#### `--file {attachment file path}`
- **Description:** Path to the file to upload as attachment
- **Usage:** Replace `{attachment file path}` with the actual path to the file
- **Required:** No (optional)
- **Example:** `--file "/path/to/new-attachment.pdf"`

## Complete Usage Examples

### Create Page Examples

#### 1\. Create page with empty content

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page"
```
**Result:** Page created with automatic action list macro and attachment macro

#### 2\. Create page with direct content

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value "This is the page content"
```
**Result:** Page created with content, automatic action list macro, and attachment macro

#### 3\. Create page with content from file

```plaintext
./confluence_cli create page \
  --space-id "SPACE" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value-from-file "/path/to/content.txt"
```
**Result:** Page created with file content, automatic action list macro, and attachment macro

#### 4\. Create page with content and attachment

```plaintext
./confluence_cli create page \
  --space-id "SPACE" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value "Content here" \
  --file "/path/to/attachment.txt"
```
**Result:** Page created with content, file attachment (embedded), action list macro, and attachment macro

#### 5\. Create page with content from file and attachment

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value-from-file "/path/to/content.html" \
  --file "/path/to/attachment.pdf"
```
**Result:** Page created with HTML content, PDF attachment (embedded), action list macro, and attachment macro

### Update Page Examples

#### 1\. Update page with direct content

```plaintext
./confluence_cli update page \
  --page-id "123456" \
  --body-value "Updated content"
```
**Result:** Page updated with new content, macros remain active

#### 2\. Update page with content from file

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value-from-file "/path/to/updated-content.txt"
```
**Result:** Page updated with file content, macros remain active

#### 3\. Update page with content and upload attachment

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value "Updated content" \
  --file "/path/to/new-attachment.pdf"
```
**Result:** Page updated with content, new attachment (embedded), macros remain active

#### 4\. Update page with content from file and upload attachment

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value-from-file "/path/to/report.html" \
  --file "/path/to/report.html"
```
**Result:** Page updated with HTML content, file attachment (embedded), macros remain active

### Upload Attachment Examples

#### 1\. Upload single attachment to existing page

```plaintext
./confluence_cli upload attachment \
  --page-id "123456" \
  --file "/path/to/document.pdf"
```

## API Reference

This tool uses the following Confluence REST APIs:

| Operation | Method | Endpoint | Description |
| --- | --- | --- | --- |
| **Create Page** | `POST` | `/wiki/api/v2/pages` | Creates new pages |
| **Update Page** | `PUT` | `/wiki/api/v2/pages/{pageId}` | Updates existing pages |
| **Get Page** | `GET` | `/wiki/api/v2/pages/{pageId}` | Retrieves page details |
| **Upload Attachment** | `POST` | `/wiki/rest/api/content/{pageId}/child/attachment` | Uploads files as attachments |
| **Get Pages by Title** | `GET` | `/wiki/api/v2/pages?title={title}` | Searches pages by title |

## Troubleshooting

### Common Issues

1. **Macros not visible after manual edit**
   - The CLI automatically detects and re-adds missing macros
   - Check logs for "Page has been manually edited, re-adding macros..." message

2. **Version conflicts**
   - The CLI handles version conflicts automatically
   - If you see "Version conflict detected" message, the CLI will retry

### Log Messages

- `"Adding attachment macro to ensure it's enabled"` - Attachment macro is being added
- `"Macros enabled successfully"` - All macros have been successfully added
- `"Page has been manually edited, re-adding macros..."` - Macros are being re-added after manual edit

## Notes and Best Practices

1. **Always use absolute paths** for file parameters to avoid path resolution issues
2. **Macros are automatically managed** - no need to manually add or remove them
3. **Test results are automatically parsed** from HTML content to determine action list status
4. **Version handling is automatic** - the CLI manages page versions transparently

## Support

For issues or questions:
1. Check the logs for detailed error messages
2. Ensure all required environment variables are set
3. Verify file paths are correct and accessible
4. Check Confluence API permissions for your account
