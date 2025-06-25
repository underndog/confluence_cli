## Confluence CLI - Complete Command Reference

The `confluence_cli` is a powerful command-line tool designed for seamless interaction with Confluence's API. It provides comprehensive functionality for creating pages, updating content, managing attachments, and displaying file attachments with beautiful macros.

## Table of Contents

*   [Quick Start](#quick-start)
*   [Commands Overview](#commands-overview)
*   [Environment Setup](#environment-setup)
*   [Create Page Command](#create-page-command)
*   [Upload Attachment Command](#upload-attachment-command)
*   [Update Page Command](#update-page-command)
*   [Complete Usage Examples](#complete-usage-examples)
*   [API Reference](#api-reference)
*   [Troubleshooting](#troubleshooting)
*   [Notes and Best Practices](#notes-and-best-practices)

## 🆕 Latest Update

**NEW FEATURE**: Macro embedding is now automatically enabled when using the `--file` flag. This means:

*   When you upload a file with `--file`, it will automatically be embedded and displayed on the page
*   This provides a better user experience by showing attachments directly on the page
*   No additional flags needed - just use `--file` and the file will be embedded automatically

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

Example:

```plaintext
docker run -it --rm \
-e API_TOKEN="ATATT3xFfGF02YUScHu9GuuqycnLcSzwhorVf3sYwEwwAEG08c5LYGRg7ltuZmPnYUvoRhkkM1-bM4HpT9g9oNjOfv37P8jdOM0HfuSmsSwwd-8IEabKemAH0gattqzH7dt-6o-nP_vpAvN1aQI=77E0A659" \
-e CONFLUENCE_URL='https://opswat.atlassian.net' \
-e EMAIL="devops.nim@nimtechnology.com" \
--name confluence-cli quay.io/underndog/confluence_cli:latest create page --space-id "2703951026" --parent-page-id "4114481660" --title "Test report 11" --body-value test
```

## Commands Overview

The tool provides three main commands:

1.  `**create page**` - Create new Confluence pages with optional content and attachments (attachments are always embedded)
2.  `**update page**` - Update existing pages with new content and/or attachments (attachments are always embedded)
3.  `**upload attachment**` - Upload attachments to existing pages without changing the page content

## Create Page Command

Creates new Confluence pages with optional content and attachments.

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

Updates existing Confluence pages with new content and/or attachments.

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

*   **Description:** Specifies the file path that contains the content for the page body
*   **Usage:** Replace `{file path}` with the actual path to the content file
*   **Required:** No (optional)
*   **Example:** `--body-value-from-file "/path/to/content.html"`

#### `--file {attachment file path}`

*   **Description:** Path to the file to upload as attachment
*   **Usage:** Replace `{attachment file path}` with the actual path to the file
*   **Required:** No (optional)
*   **Example:** `--file "/home/user/documents/report.pdf"`

## Complete Usage Examples

### Create Page Examples

#### 1\. Create page with empty content

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page"
```

#### 2\. Create page with direct content

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value "This is the page content"
```

#### 3\. Create page with content from file

```plaintext
./confluence_cli create page \
  --space-id "SPACE" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value-from-file "/path/to/content.txt"
```

#### 4\. Create page with content and attachment

```plaintext
./confluence_cli create page \
  --space-id "SPACE" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value "Content here" \
  --file "/path/to/attachment.txt"
```

**Note:** File will be automatically embedded and displayed on the page.

#### 5\. Create page with content from file and attachment

```plaintext
./confluence_cli create page \
  --space-id "SPACE_ID" \
  --parent-page-id "123456" \
  --title "Test Page" \
  --body-value-from-file "/path/to/content.html" \
  --file "/path/to/attachment.pdf"
```

**Note:** File will be automatically embedded and displayed on the page.

### Update Page Examples

#### 1\. Update page with direct content

```plaintext
./confluence_cli update page \
  --page-id "123456" \
  --body-value "Updated content"
```

#### 2\. Update page with content from file

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value-from-file "/path/to/updated-content.txt"
```

#### 3\. Update page with content and upload attachment

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value "Updated content" \
  --file "/path/to/new-attachment.pdf"
```

**Note:** File will be automatically embedded and displayed on the page.

#### 4\. Update page with content from file and upload attachment

```plaintext
./confluence_cli update page \
  --page-id "4118873955" \
  --body-value-from-file "/path/to/report.html" \
  --file "/path/to/report.html"
```

**Note:** File will be automatically embedded and displayed on the page.

### Upload Attachment Examples

#### 1\. Upload single attachment to existing page

```plaintext
./confluence_cli upload attachment \
  --page-id "123456" \
  --file "/path/to/document.pdf"
```

**Note:** File will be uploaded as attachment but not embedded in the page content.

## API Reference

This tool uses the following Confluence REST APIs:

| Operation | Method | Endpoint | Description |
| --- | --- | --- | --- |
| **Create Page** | `POST` | `/wiki/api/v2/pages` | Creates new pages |
| **Update Page** | `PUT` | `/wiki/api/v2/pages/{pageId}` | Updates existing pages |
| **Get Page** | `GET` | `/wiki/api/v2/pages/{pageId}` | Retrieves page details |
| **Upload Attachment** | `POST` | `/wiki/rest/api/content/{pageId}/child/attachment` | Uploads files as attachments |
| **Get Pages by Title** | `GET` | `/wiki/api/v2/pages?title={title}` | Searches pages by title |