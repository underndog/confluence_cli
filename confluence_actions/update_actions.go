package confluence_actions

import (
	"confluence_cli/helper"
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
		log.Info("Processing file upload for page ID:", pageId)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("attachment file does not exist: %s", filePath)
		}

		log.Info("File exists, starting upload...")

		uploadResp, err := http_request.UploadConfluenceAttachment(pageId, filePath)
		if err != nil {
			log.Error("Upload failed with error:", err)
			return err
		}

		log.Info("Upload response status:", uploadResp.StatusCode())
		log.Info("Upload response body:", string(uploadResp.Body()))

		if uploadResp.StatusCode() < 300 {
			log.Info("File uploaded successfully as attachment!")
		} else {
			log.Error("Upload failed with status:", uploadResp.Status())
			log.Error("Upload response body:", string(uploadResp.Body()))
			return fmt.Errorf("upload failed with status: %s", uploadResp.Status())
		}
	} else {
		log.Info("No file attachment specified, but will still ensure attachment macro is enabled.")
	}

	// ALWAYS enable attachment macro when updating page
	log.Info("Ensuring attachment macro is enabled...")

	// Get the LATEST page content and version
	resp, err := http_request.GetConfluencePageByID(pageId + "?body-format=storage")
	if err != nil {
		log.Error("Failed to get page content:", err)
		return err
	}

	var currentPage req.CreatePageResult
	if err := json.Unmarshal(resp.Body(), &currentPage); err != nil {
		log.Error("Failed to unmarshal page content:", err)
		return err
	}

	// Check if content already has attachment macro
	hasAttachmentMacro := helper.HasAttachmentMacro(currentPage.Body.Storage.Value)
	hasActionListMacro := helper.HasActionListMacro(currentPage.Body.Storage.Value)

	// Log macro status for debugging
	log.Info("Macro status check - Attachment:", hasAttachmentMacro, "Action List:", hasActionListMacro)

	if !hasAttachmentMacro || !hasActionListMacro {
		log.Info("Some macros are missing, will re-enable them")

		// Helper function để tạo payload và update page
		updatePageWithMacros := func(newBody string) error {
			payload := req.UpdatePagePayload{
				ID:      pageId,
				Status:  "current",
				Title:   currentPage.Title,
				Body:    req.UpdatePageBody{Representation: "storage", Value: newBody},
				Version: req.UpdatePageVersion{Number: currentPage.Version.Number + 1, Message: "Re-enabled macros via CLI"},
			}

			log.Info("Updating page with new version:", currentPage.Version.Number+1)

			embedResp, err := http_request.UpdateConfluencePage(pageId, payload)
			if err != nil {
				log.Error("Failed to update page:", err)
				return err
			}

			if embedResp.StatusCode() < 300 {
				log.Info("Macros re-enabled successfully")
				return nil
			} else {
				log.Error("Failed to re-enable macros. Status:", embedResp.StatusCode())
				return fmt.Errorf("failed to re-enable macros")
			}
		}

		// Build new body with missing macros
		newBody := currentPage.Body.Storage.Value

		if !hasAttachmentMacro {
			log.Info("Adding attachment macro")
			newBody += "\n" + helper.CreateAttachmentMacro()
		}

		if !hasActionListMacro {
			log.Info("Action list macro is missing - this should be handled by Python logic")
			log.Info("Current page version:", currentPage.Version.Number)
		}

		if err := updatePageWithMacros(newBody); err != nil {
			return err
		}
	} else {
		log.Info("All required macros are already enabled")
	}

	log.Info("Page update completed successfully")
	return nil
}
