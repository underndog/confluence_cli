package confluence_commands

import (
	"confluence_cli/confluence_actions"

	"github.com/urfave/cli/v2"
)

func UploadCommand() *cli.Command {
	return &cli.Command{
		Name:  "upload",
		Usage: "Upload a resource",
		Subcommands: []*cli.Command{
			{
				Name:   "attachment",
				Usage:  "Upload an attachment to an existing page without changing the page content",
				Action: confluence_actions.UploadAttachmentAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "page-id",
						Usage:    "ID of the page to upload attachment to",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "file",
						Usage:    "Path to the file to upload as attachment",
						Required: true,
					},
				},
			},
		},
	}
}
