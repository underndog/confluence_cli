package confluence_commands

import (
	"confluence_cli/confluence_actions"

	"github.com/urfave/cli/v2"
)

func CreatePageCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a resource",
		Subcommands: []*cli.Command{
			{
				Name:   "page",
				Usage:  "Create a page with title and content",
				Action: confluence_actions.CreatePageAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "space-id",
						Usage: "Specifies the ID of the space where the new page will be created",
					},
					&cli.StringFlag{
						Name:  "parent-page-id",
						Usage: "Indicates the ID of the parent page under which the new page will be nested",
					},
					&cli.StringFlag{
						Name:  "title",
						Usage: "Sets the title for the new Confluence page",
					},
					&cli.StringFlag{
						Name:  "body-value-from-file",
						Usage: "Specifies the file path that contains the content for the page body",
					},
					&cli.StringFlag{
						Name:  "body-value",
						Usage: "The content for the page body",
					},
				},
			},
			{
				Name:   "attachment",
				Usage:  "Upload a file as attachment to a Confluence page",
				Action: confluence_actions.UploadAttachmentAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "page-id",
						Usage: "The ID of the page where the file will be uploaded as attachment",
					},
					&cli.StringFlag{
						Name:  "file",
						Usage: "Path to the file to upload",
					},
				},
			},
		},
	}
}
