package confluence_actions

import (
	"confluence_cli/helper/http_request"
	"confluence_cli/log"
	"confluence_cli/model/req"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func UpdatePageAction(c *cli.Context) error {
	pageId := c.String("page-id")
	if pageId == "" {
		log.Error("Please provide --page-id")
		return fmt.Errorf("please provide --page-id")
	}

	bodyValueFromFile := c.String("body-value-from-file")
	bodyValue := c.String("body-value")
	filePath := c.String("file")

	// Check if we need to update content (new content OR file upload with auto-embed)
	needsContentUpdate := bodyValueFromFile != "" || bodyValue != "" || filePath != ""

	if needsContentUpdate {
		var contentPage string

		// Get current page content first
		resp, err := http_request.GetConfluencePageByID(pageId)
		if err != nil {
			return err
		}
		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			return fmt.Errorf("failed to get page: %s", resp.Status())
		}
		var currentPage req.CreatePageResult
		if err := json.Unmarshal(resp.Body(), &currentPage); err != nil {
			return err
		}

		// Start with current content
		contentPage = currentPage.Body.Storage.Value

		// Update content if new content is provided
		if bodyValueFromFile != "" {
			data, err := os.ReadFile(bodyValueFromFile)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", bodyValueFromFile, err)
			}
			contentPage = string(data)
		} else if bodyValue != "" {
			contentPage = bodyValue
		}

		// Add embed macro if file is present (always enabled)
		if filePath != "" {
			if contentPage != "" {
				contentPage += "\n\n"
			}
			contentPage += `<ac:structured-macro ac:name="attachments" ac:schema-version="1"></ac:structured-macro>`
		}

		// Update page with new content
		payload := req.UpdatePagePayload{
			ID:     pageId,
			Status: "current",
			Title:  currentPage.Title,
			Body:   req.UpdatePageBody{Representation: "storage", Value: contentPage},
			Version: req.UpdatePageVersion{
				Number:  currentPage.Version.Number + 1,
				Message: "Updated content via CLI",
			},
		}

		updateResp, err := http_request.UpdateConfluencePage(pageId, payload)
		if err != nil {
			log.Error("Failed to update page:", err)
			return fmt.Errorf("failed to update page: %w", err)
		}
		if updateResp.StatusCode() < 200 || updateResp.StatusCode() >= 300 {
			log.Error("Update failed with status:", updateResp.Status())
			return fmt.Errorf("failed to update page with status: %s", updateResp.Status())
		}
		log.Info("Page content updated successfully!")
	} else {
		log.Info("No content or file upload requested. Keeping current page content unchanged.")
	}

	// Upload Attachment (if provided)
	if filePath != "" {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("attachment file does not exist: %s", filePath)
		}
		uploadResp, err := http_request.UploadConfluenceAttachment(pageId, filePath)
		if err != nil {
			return fmt.Errorf("attachment upload failed: %w", err)
		}
		if uploadResp.StatusCode() >= 200 && uploadResp.StatusCode() < 300 {
			log.Info("File uploaded successfully as attachment!")
		} else {
			return fmt.Errorf("attachment upload failed with status %s: %s", uploadResp.Status(), string(uploadResp.Body()))
		}
	}

	log.Info("All operations completed successfully.")
	return nil
}
