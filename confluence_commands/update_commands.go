package confluence_commands

import (
	"confluence_cli/confluence_actions"

	"github.com/urfave/cli/v2"
)

func UpdatePageCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update a resource",
		Subcommands: []*cli.Command{
			{
				Name:   "page",
				Usage:  "Update a page's body content",
				Action: confluence_actions.UpdatePageAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "page-id",
						Usage:    "ID of the page to update",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "body-value",
						Usage: "Content to update the page with",
					},
					&cli.StringFlag{
						Name:  "body-value-from-file",
						Usage: "File path to read content from",
					},
					&cli.StringFlag{
						Name:  "file",
						Usage: "Path to a file to upload as an attachment (automatically embedded)",
					},
				},
			},
		},
	}
}
