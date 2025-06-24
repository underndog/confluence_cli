package confluence_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestCreatePageCommand_Structure(t *testing.T) {
	command := CreatePageCommand()

	assert.Equal(t, "create", command.Name)
	assert.Equal(t, "Create a resource", command.Usage)
	assert.Len(t, command.Subcommands, 1)

	subcommand := command.Subcommands[0]
	assert.Equal(t, "page", subcommand.Name)
	assert.Equal(t, "Create a page with title and content", subcommand.Usage)
	assert.NotNil(t, subcommand.Action)

	// Check required flags
	requiredFlags := []string{"space-id", "parent-page-id", "title"}
	for _, flagName := range requiredFlags {
		found := false
		for _, flag := range subcommand.Flags {
			if flag.Names()[0] == flagName {
				found = true
				break
			}
		}
		assert.True(t, found, "Required flag %s not found", flagName)
	}

	// Check optional flags
	optionalFlags := []string{"body-value", "body-value-from-file", "file"}
	for _, flagName := range optionalFlags {
		found := false
		for _, flag := range subcommand.Flags {
			if flag.Names()[0] == flagName {
				found = true
				break
			}
		}
		assert.True(t, found, "Optional flag %s not found", flagName)
	}
}

func TestUpdatePageCommand_Structure(t *testing.T) {
	command := UpdatePageCommand()

	assert.Equal(t, "update", command.Name)
	assert.Equal(t, "Update a resource", command.Usage)
	assert.Len(t, command.Subcommands, 1)

	subcommand := command.Subcommands[0]
	assert.Equal(t, "page", subcommand.Name)
	assert.Equal(t, "Update a page's body content", subcommand.Usage)
	assert.NotNil(t, subcommand.Action)

	// Check required flags
	requiredFlags := []string{"page-id"}
	for _, flagName := range requiredFlags {
		found := false
		for _, flag := range subcommand.Flags {
			if flag.Names()[0] == flagName {
				found = true
				break
			}
		}
		assert.True(t, found, "Required flag %s not found", flagName)
	}

	// Check optional flags
	optionalFlags := []string{"body-value", "body-value-from-file", "file"}
	for _, flagName := range optionalFlags {
		found := false
		for _, flag := range subcommand.Flags {
			if flag.Names()[0] == flagName {
				found = true
				break
			}
		}
		assert.True(t, found, "Optional flag %s not found", flagName)
	}
}

func TestUploadCommand_Structure(t *testing.T) {
	command := UploadCommand()

	assert.Equal(t, "upload", command.Name)
	assert.Equal(t, "Upload a resource", command.Usage)
	assert.Len(t, command.Subcommands, 1)

	subcommand := command.Subcommands[0]
	assert.Equal(t, "attachment", subcommand.Name)
	assert.Equal(t, "Upload an attachment to an existing page without changing the page content", subcommand.Usage)
	assert.NotNil(t, subcommand.Action)

	// Check required flags
	requiredFlags := []string{"page-id", "file"}
	for _, flagName := range requiredFlags {
		found := false
		for _, flag := range subcommand.Flags {
			if flag.Names()[0] == flagName {
				found = true
				break
			}
		}
		assert.True(t, found, "Required flag %s not found", flagName)
	}
}

func TestCommands_Integration(t *testing.T) {
	// Test that all commands can be added to an app without errors
	app := &cli.App{
		Commands: []*cli.Command{
			CreatePageCommand(),
			UpdatePageCommand(),
			UploadCommand(),
		},
	}

	assert.Len(t, app.Commands, 3)
	assert.Equal(t, "create", app.Commands[0].Name)
	assert.Equal(t, "update", app.Commands[1].Name)
	assert.Equal(t, "upload", app.Commands[2].Name)
}
