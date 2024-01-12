# Confluence CLI Command Explanation

The `confluence_cli` command is used for interacting with Confluence's API via the command line. Here's a breakdown of the command for creating a new page in Confluence:

```shell
confluence_cli create page \
--space-id {space id} \
--parent-page-id {parent page id} \
--title {tile of page} \
--body-from-value-file {file path to file}
```

# Command Components

- `confluence_cli`: The name of the CLI tool designed for Confluence operations.

- `create page`: The action to be performed. `create` indicates creation of a new entity, and `page` specifies that the entity is a page within Confluence.

## Options

### `--space-id {space id}`:
- **Description:** Specifies the ID of the space where the new page will be created.
- **Usage:** Replace `{space id}` with the actual space ID.

### `--parent-page-id {parent page id}`:
- **Description:** Indicates the ID of the parent page under which the new page will be nested.
- **Usage:** Replace `{parent page id}` with the actual parent page ID.

### `--title {title of page}`:
- **Description:** Sets the title for the new Confluence page.
- **Usage:** Replace `{title of page}` with the desired page title.

### `--body-from-value-file {file path to file}`:
- **Description:** Specifies the file path that contains the content for the page body.
- **Usage:** Replace `{file path to file}` with the actual path to the content file.

## Environment Variables.

In order to connect Your Confluence. You must configure the environments such as:   

`CONFLUENCE_URL`:   
- **Description:** your confluence link such as: `https://nimtechnology.atlassian.net`

# How build Binary file.

On windows
```shell
go build -o confluence_cli.exe

.\confluence_cli.exe create page --space-id 98432 --parent-page-id 589825 --title ahihi --body-value-from-file hihi1
```


