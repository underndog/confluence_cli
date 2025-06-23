package confluence_actions

import (
	"testing"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"confluence_cli/log"
)

// Mock context for minimal flag testing
func TestCreatePageAction_MissingFlags(t *testing.T) {
	log.InitLogger(true)
	app := cli.NewApp()
	set := flagSet([]string{})
	c := cli.NewContext(app, set, nil)
	err := CreatePageAction(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space-id")
}

// Helper to create a flag set
func flagSet(args []string) *flag.FlagSet {
	if len(args)%2 != 0 {
		panic("args must contain an even number of elements (key-value pairs)")
	}
	set := flag.NewFlagSet("test", 0)
	set.String("space-id", "", "")
	set.String("parent-page-id", "", "")
	set.String("title", "", "")
	set.String("body-value-from-file", "", "")
	set.String("body-value", "", "")
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break // or handle the error appropriately
		}
		set.Set(args[i], args[i+1])
	}
	return set
}

// More tests can be added with mocks for http_request and log
