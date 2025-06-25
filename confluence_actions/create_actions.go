package confluence_actions

import (
	"confluence_cli/helper/http_request"
	"confluence_cli/log"
	"confluence_cli/model/req"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

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
				respPages, err := http_request.GetConfluencePagesByTitle(fmt.Sprintf("[%s] %s", formattedDate, title))
				if err != nil {
					return err
				}
				var pagesInfo req.PagesInfo
				if err = json.Unmarshal(respPages.Body(), &pagesInfo); err == nil && len(pagesInfo.Results) > 0 {
					nextParentID = pagesInfo.Results[0].ID
				}
			}
		}
	} else {
		var createdResultPage req.CreatePageResult
		if err := json.Unmarshal(resp.Body(), &createdResultPage); err == nil {
			nextParentID = createdResultPage.ID
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

	filePath := c.String("file")
	if filePath != "" {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("attachment file does not exist: %s", filePath)
		}

		uploadResp, err := http_request.UploadConfluenceAttachment(createdPageId, filePath)
		if err != nil {
			return err
		}

		if uploadResp.StatusCode() < 300 {
			log.Info("File uploaded successfully as attachment!")

			// Enable macro by default when file is uploaded, unless explicitly disabled
			shouldEmbed := true // Always enable macro when file is present
			if shouldEmbed {
				log.Info("Embedding attached file into page content...")

				resp, err := http_request.GetConfluencePageByID(createdPageId + "?body-format=storage")
				if err != nil {
					return err
				}
				var currentPage req.CreatePageResult
				if err := json.Unmarshal(resp.Body(), &currentPage); err != nil {
					return err
				}

				macroXHTML := fmt.Sprintf(
					`<p><ac:structured-macro ac:name="attachments" ac:schema-version="1"></ac:structured-macro></p>`,
				)
				newBody := currentPage.Body.Storage.Value + macroXHTML

				payload := req.UpdatePagePayload{
					ID:      createdPageId,
					Status:  "current",
					Title:   currentPage.Title,
					Body:    req.UpdatePageBody{Representation: "storage", Value: newBody},
					Version: req.UpdatePageVersion{Number: currentPage.Version.Number + 1, Message: "Embedded attached file via CLI"},
				}

				embedResp, err := http_request.UpdateConfluencePage(createdPageId, payload)
				if err != nil || embedResp.StatusCode() >= 300 {
					return fmt.Errorf("embedding attachment failed")
				}
				log.Info("File embedded successfully into page.")
			}
		} else {
			return fmt.Errorf("upload failed with status: %s", uploadResp.Status())
		}
	}

	log.Info("All operations completed successfully.")
	return nil
}
