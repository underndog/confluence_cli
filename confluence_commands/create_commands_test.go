package confluence_commands

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreatePageCommand(t *testing.T) {
	cmd := CreatePageCommand()
	assert.Equal(t, "create", cmd.Name)
	assert.Equal(t, "Create a resource", cmd.Usage)
	assert.NotEmpty(t, cmd.Subcommands)
	pageCmd := cmd.Subcommands[0]
	assert.Equal(t, "page", pageCmd.Name)
	assert.Equal(t, "Create a page with title and content", pageCmd.Usage)
	assert.NotNil(t, pageCmd.Action)
	assert.NotEmpty(t, pageCmd.Flags)
}
