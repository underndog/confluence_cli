package confluence_actions

import (
	"confluence_cli/helper/http_request"
	"confluence_cli/log"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func UploadAttachmentAction(c *cli.Context) error {
	contentId := c.String("page-id")
	if contentId == "" {
		log.Error("Please provide --page-id xxxx")
		return fmt.Errorf("please provide --page-id xxxx")
	}

	filePath := c.String("file")
	if filePath == "" {
		log.Error("Please provide --file /path/to/file")
		return fmt.Errorf("please provide --file /path/to/file")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error("File does not exist:", filePath)
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Get file info for logging
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Error("Error getting file info:", err)
		return err
	}

	// Check if file is empty
	if fileInfo.Size() == 0 {
		log.Error("File is empty:", filePath)
		return fmt.Errorf("file is empty: %s", filePath)
	}

	log.Info("Uploading file:", filepath.Base(filePath), "Size:", fileInfo.Size(), "bytes")
	log.Info("Content ID:", contentId)

	resp, err := http_request.UploadConfluenceAttachment(contentId, filePath)
	if err != nil {
		log.Error("Error when uploading file via Confluence API:", err)
		return err
	}

	log.Info("Response status:", resp.Status())
	log.Info("Response code:", resp.StatusCode())

	if resp.StatusCode() == 200 {
		log.Info("File uploaded successfully!")
	} else {
		log.Error("Upload failed with status:", resp.Status())
		log.Error("Response body:", string(resp.Body()))
		return fmt.Errorf("upload failed with status: %s", resp.Status())
	}

	return nil
}
