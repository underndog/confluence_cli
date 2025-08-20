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
	"time"

	"github.com/urfave/cli/v2"
)

// Parse test results from HTML content to determine Overall Status
func parseTestResultsFromHTML(htmlContent string) (failedCount, totalCount int) {
	// Default values
	failedCount = 0
	totalCount = 0

	// Simple regex to find test counts
	// Look for patterns like "Failed: 47" or "Total Tests: 438"
	failedMatch := regexp.MustCompile(`Failed:\s*(\d+)`)
	totalMatch := regexp.MustCompile(`Total Tests:\s*(\d+)`)

	if failedMatches := failedMatch.FindStringSubmatch(htmlContent); len(failedMatches) > 1 {
		fmt.Sscanf(failedMatches[1], "%d", &failedCount)
		log.Info("Found failed count:", failedCount)
	}

	if totalMatches := totalMatch.FindStringSubmatch(htmlContent); len(totalMatches) > 1 {
		fmt.Sscanf(totalMatches[1], "%d", &totalCount)
		log.Info("Found total count:", totalCount)
	}

	// If no matches found, try alternative patterns
	if failedCount == 0 && totalCount == 0 {
		// Try to find in different format
		failedMatch2 := regexp.MustCompile(`<strong>Failed:</strong>\s*(\d+)`)
		totalMatch2 := regexp.MustCompile(`<strong>Total Tests:</strong>\s*(\d+)`)

		if failedMatches := failedMatch2.FindStringSubmatch(htmlContent); len(failedMatches) > 1 {
			fmt.Sscanf(failedMatches[1], "%d", &failedCount)
			log.Info("Found failed count (alt):", failedCount)
		}

		if totalMatches := totalMatch2.FindStringSubmatch(htmlContent); len(totalMatches) > 1 {
			fmt.Sscanf(totalMatches[1], "%d", &totalCount)
			log.Info("Found total count (alt):", totalCount)
		}
	}

	log.Info("Final parsed results - Failed:", failedCount, "Total:", totalCount)
	return failedCount, totalCount
}

func CreatePageAction(c *cli.Context) error {
	spaceId := c.String("space-id")
	if spaceId == "" {
		return fmt.Errorf("--space-id is required")
	}

	parentPageId := c.String("parent-page-id")
	if parentPageId == "" {
		return fmt.Errorf("--parent-page-id is required")
	}

	title := c.String("title")
	if title == "" {
		return fmt.Errorf("--title is required")
	}

	bodyValueFromFile := c.String("body-value-from-file")
	bodyValue := c.String("body-value")
	var contentPage string

	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01")
	formattedDateTime := currentTime.Format("2006-01-02 15:04:05")

	if bodyValueFromFile != "" {
		data, err := os.ReadFile(bodyValueFromFile)
		if err != nil {
			log.Error("Error when Create Page via Confluence API: ", err)
			return err
		}
		contentPage = string(data)
	} else if bodyValue != "" {
		contentPage = bodyValue
	}
	// If neither is provided, contentPage remains empty string

	// Log content trước khi tạo page
	log.Info("Content length:", len(contentPage))

	// Create the page
	log.Info("Creating page...")
	resp, err := http_request.CreateConfluencePage(spaceId, parentPageId, fmt.Sprintf("[%s] %s", formattedDate, title), contentPage)
	if err != nil {
		return err
	}
	var nextParentID string
	if resp.StatusCode() >= 400 {
		var errResponse req.ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errResponse); err == nil && len(errResponse.Errors) > 0 {
			if errResponse.Errors[0].Title == "A page with this title already exists: A page already exists with the same TITLE in this space" {
				log.Info("Page with title already exists, getting existing page ID...")
				respPages, err := http_request.GetConfluencePagesByTitle(fmt.Sprintf("[%s] %s", formattedDate, title))
				if err != nil {
					return err
				}
				var pagesInfo req.PagesInfo
				if err = json.Unmarshal(respPages.Body(), &pagesInfo); err == nil && len(pagesInfo.Results) > 0 {
					nextParentID = pagesInfo.Results[0].ID
					log.Info("Found existing page with ID:", nextParentID)
				}
			}
		}
	} else {
		var createdResultPage req.CreatePageResult
		if err := json.Unmarshal(resp.Body(), &createdResultPage); err == nil {
			nextParentID = createdResultPage.ID
			log.Info("Created parent page with ID:", nextParentID)
		}
	}
	if nextParentID == "" {
		return fmt.Errorf("could not determine parent page ID for content page")
	}

	childResp, err := http_request.CreateConfluencePage(spaceId, nextParentID, fmt.Sprintf("[%s] %s", formattedDateTime, title), contentPage)
	if err != nil {
		return err
	}
	if childResp.StatusCode() >= 300 {
		return fmt.Errorf("failed to create content page with status %s: %s", childResp.Status(), string(childResp.Body()))
	}
	var childPageResult req.CreatePageResult
	if err := json.Unmarshal(childResp.Body(), &childPageResult); err != nil {
		return err
	}
	createdPageId := childPageResult.ID
	log.Info("Successfully created content page with ID:", createdPageId)

	// Create and log the full URL using the webui link from API response
	if childPageResult.Links.Webui != "" {
		// Get base URL from the same source as http_request functions
		baseURL := helper.GetEnvOrDefault("CONFLUENCE_URL", "https://nimtechnology.atlassian.net")

		// Construct full URL
		fullURL := baseURL + "wiki" + childPageResult.Links.Webui
		log.Info("Page URL:", fullURL)
	} else {
		// Fallback if webui link is not available
		log.Info("Page created with ID:", createdPageId)
		log.Info("WebUI link not available in response")
	}

	filePath := c.String("file")
	if filePath != "" {
		log.Info("Processing file upload for page ID:", createdPageId)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("attachment file does not exist: %s", filePath)
		}

		log.Info("File exists, starting upload...")
		uploadResp, err := http_request.UploadConfluenceAttachment(createdPageId, filePath)
		if err != nil {
			return err
		}

		if uploadResp.StatusCode() < 300 {
			log.Info("File uploaded successfully as attachment!")

			// Enable macro by default when file is uploaded, unless explicitly disabled
			shouldEmbed := true // Always enable macro when file is present
			if shouldEmbed {
				log.Info("Ensuring attachment macro is enabled...")

				resp, err := http_request.GetConfluencePageByID(createdPageId + "?body-format=storage")
				if err != nil {
					return err
				}
				var currentPage req.CreatePageResult
				if err := json.Unmarshal(resp.Body(), &currentPage); err != nil {
					return err
				}

				// Add file attachment macro
				attachmentMacro := helper.CreateAttachmentMacro()

				// Action list macro will be set by generate_report.py
				// confluence_cli only enables attachment macro

				// Find and replace placeholder in Overall Status row
				currentBody := currentPage.Body.Storage.Value

				// Check macro status for debugging
				hasAttachmentMacro := helper.HasAttachmentMacro(currentBody)
				hasActionListMacro := helper.HasActionListMacro(currentBody)
				log.Info("Macro status check - Attachment:", hasAttachmentMacro, "Action List:", hasActionListMacro)

				// Pattern to find Overall Status row - support td with content
				// Match td có content bên trong, không cần closing tag
				overallStatusPattern := `<td><strong>Overall Status</strong></td>\s*<td colspan=['"]2['"]>`

				re := regexp.MustCompile(overallStatusPattern)

				if re.MatchString(currentBody) {
					log.Info("Overall Status pattern matched - adding attachment macro")

					// Chỉ thêm attachment macro vào cuối, không thay thế Overall Status
					newBody := currentPage.Body.Storage.Value + attachmentMacro

					// Update page với newBody
					payload := req.UpdatePagePayload{
						ID:      createdPageId,
						Status:  "current",
						Title:   currentPage.Title,
						Body:    req.UpdatePageBody{Representation: "storage", Value: newBody},
						Version: req.UpdatePageVersion{Number: currentPage.Version.Number + 1, Message: "Added attachment macro"},
					}

					log.Info("Updating page with new version:", currentPage.Version.Number+1)
					embedResp, err := http_request.UpdateConfluencePage(createdPageId, payload)
					if err != nil || embedResp.StatusCode() >= 300 {
						return fmt.Errorf("adding attachment macro failed")
					}
					log.Info("Attachment macro enabled successfully")
					return nil
				} else {
					log.Info("Overall Status pattern not matched - adding attachment macro to end")

					// Fallback: thêm attachment macro vào cuối
					newBody := currentPage.Body.Storage.Value + attachmentMacro

					// Update page với newBody
					payload := req.UpdatePagePayload{
						ID:      createdPageId,
						Status:  "current",
						Title:   currentPage.Title,
						Body:    req.UpdatePageBody{Representation: "storage", Value: newBody},
						Version: req.UpdatePageVersion{Number: currentPage.Version.Number + 1, Message: "Added attachment macro via fallback"},
					}

					log.Info("Updating page with new version:", currentPage.Version.Number+1)
					embedResp, err := http_request.UpdateConfluencePage(createdPageId, payload)
					if err != nil || embedResp.StatusCode() >= 300 {
						return fmt.Errorf("adding attachment macro failed")
					}
					log.Info("Attachment macro enabled successfully via fallback")
					return nil
				}
			}
		} else {
			return fmt.Errorf("upload failed with status: %s", uploadResp.Status())
		}
	} else {
		log.Info("No file attachment specified, but page creation completed successfully")
	}

	log.Info("Page creation completed successfully")
	return nil
}
