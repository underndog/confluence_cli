package confluence_actions

import (
	"confluence_cli/helper"
	"confluence_cli/helper/http_request"
	"confluence_cli/log"
	"confluence_cli/model/req"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

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

	// ALWAYS enable attachment macro when updating page (regardless of file upload)
	log.Info("Ensuring attachment macro and action list macro are enabled...")

	// Get the LATEST page content and version (not cached)
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

	log.Info("Current page version:", currentPage.Version.Number)

	// Check if macros already exist
	hasAttachmentMacro := strings.Contains(currentPage.Body.Storage.Value, "attachments")
	hasActionList := strings.Contains(currentPage.Body.Storage.Value, "ac:task-list")

	// Add file attachment macro only if it doesn't exist
	attachmentMacro := helper.CreateAttachmentMacro()

	// Parse test results from current page content
	failedCount, totalCount := parseTestResultsFromHTML(currentPage.Body.Storage.Value)

	// Debug logging
	log.Info("Parsed test results - Failed:", failedCount, "Total:", totalCount)

	// Add action list macro cho Overall Status
	actionListMacro := helper.CreateActionItemMacro(failedCount, totalCount)

	var newBody string
	var needsUpdate bool

	// Find and replace placeholder in Overall Status row if action list doesn't exist
	if !hasActionList {
		// Pattern to find Overall Status row
		overallStatusPattern := `<td><strong>Overall Status</strong></td>\s*<td colspan="2">\s*</td>`

		re := regexp.MustCompile(overallStatusPattern)
		if re.MatchString(currentPage.Body.Storage.Value) {
			log.Info("Adding action list macro to Overall Status")
			replacement := `<td><strong>Overall Status</strong></td>
				<td colspan="2">
					` + actionListMacro + `
				</td>`
			newBody = re.ReplaceAllString(currentPage.Body.Storage.Value, replacement)
			needsUpdate = true
		} else {
			// Fallback: append to end
			newBody = currentPage.Body.Storage.Value
			if !hasAttachmentMacro {
				newBody += attachmentMacro
				needsUpdate = true
			}
			newBody += actionListMacro
			needsUpdate = true
		}
	} else {
		// Action list already exists, only need to check attachment macro
		newBody = currentPage.Body.Storage.Value
		if !hasAttachmentMacro {
			newBody += attachmentMacro
			needsUpdate = true
		}
	}

	// Ensure attachment macro is always added to end
	if !strings.Contains(newBody, "attachments") {
		log.Info("Adding attachment macro to ensure it's enabled")
		newBody += attachmentMacro
		needsUpdate = true
	}

	// Helper function để tạo payload và update page
	updatePageWithMacros := func(newBody string) error {
		payload := req.UpdatePagePayload{
			ID:      pageId,
			Status:  "current",
			Title:   currentPage.Title,
			Body:    req.UpdatePageBody{Representation: "storage", Value: newBody},
			Version: req.UpdatePageVersion{Number: currentPage.Version.Number + 1, Message: "Enabled macros via CLI"},
		}

		log.Info("Updating page with new version:", currentPage.Version.Number+1)

		embedResp, err := http_request.UpdateConfluencePage(pageId, payload)
		if err != nil {
			log.Error("Update page failed with error:", err)
			return err
		}

		if embedResp.StatusCode() >= 300 {
			if embedResp.StatusCode() == 409 {
				log.Error("Version conflict detected. Page was modified by another user.")
				log.Error("Please try again or refresh the page content.")
			}
			return fmt.Errorf("enabling macros failed with status: %s", embedResp.Status())
		}
		log.Info("Macros enabled successfully.")
		return nil
	}

	// Update page if needed
	if needsUpdate {
		log.Info("Updating page content with macros...")
		if err := updatePageWithMacros(newBody); err != nil {
			return err
		}

		// After successful update, check if page needs additional update due to manual edits
		log.Info("Checking if page needs additional update after manual edits...")

		// Get the LATEST page content and version (not cached)
		latestResp, err := http_request.GetConfluencePageByID(pageId + "?body-format=storage")
		if err != nil {
			log.Warn("Could not get latest page info for version check")
		} else {
			var latestPage req.CreatePageResult
			if err := json.Unmarshal(latestResp.Body(), &latestPage); err == nil {
				log.Info("Latest page version after update:", latestPage.Version.Number)

				// If version increased more than expected, page was manually edited
				expectedVersion := currentPage.Version.Number + 1
				if latestPage.Version.Number > expectedVersion {
					log.Info("Page has been manually edited, re-adding macros...")

					// Check if macros still exist in the edited content
					hasMacrosInLatest := strings.Contains(latestPage.Body.Storage.Value, "attachments") &&
						strings.Contains(latestPage.Body.Storage.Value, "ac:task-list")

					if !hasMacrosInLatest {
						log.Info("Macros were lost during manual edit, re-adding them...")

						// Re-add macros to the latest content
						reupdateBody := latestPage.Body.Storage.Value

						// Add attachment macro if missing
						if !strings.Contains(reupdateBody, "attachments") {
							reupdateBody += attachmentMacro
						}

						// Add action list macro if missing
						if !strings.Contains(reupdateBody, "ac:task-list") {
							reupdateBody += actionListMacro
						}

						// Update page again with macros
						reupdatePayload := req.UpdatePagePayload{
							ID:      pageId,
							Status:  "current",
							Title:   latestPage.Title,
							Body:    req.UpdatePageBody{Representation: "storage", Value: reupdateBody},
							Version: req.UpdatePageVersion{Number: latestPage.Version.Number + 1, Message: "Re-added macros after manual edit"},
						}

						if reupdateResp, err := http_request.UpdateConfluencePage(pageId, reupdatePayload); err != nil {
							log.Warn("Failed to re-add macros after manual edit")
						} else if reupdateResp.StatusCode() >= 300 {
							log.Warn("Failed to re-add macros after manual edit, status:", reupdateResp.Status())
						} else {
							log.Info("Macros successfully re-added after manual edit")
						}
					} else {
						log.Info("Macros still exist after manual edit, no re-update needed")
					}
				}
			}
		}
	} else {
		log.Info("All macros already exist, no update needed.")
	}

	log.Info("All operations completed successfully.")
	return nil
}
