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
		log.Error("Please provide --space-id xxxx")
		return fmt.Errorf("please provide --space-id xxxx")
	}

	parentPageId := c.String("parent-page-id")
	if parentPageId == "" {
		log.Error("Please provide --parent-page-id xxxx")
		return fmt.Errorf("please provide --parent-page-id xxxx")
	}

	title := c.String("title")
	if title == "" {
		log.Error("Please provide --title xxxx")
		return fmt.Errorf("please provide --title xxxx")
	}

	bodyValueFromFile := c.String("body-value-from-file")
	bodyValue := c.String("body-value")
	var contentPage string

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

	var nextParentID string
	// Get the current date
	currentTime := time.Now()
	// Format the date to "YYYY-MM"
	formattedDate := currentTime.Format("2006-01")
	formattedDateTime := currentTime.Format("2006-01-02 15:04:05")

	resp, err := http_request.CreateConfluencePage(spaceId, parentPageId, fmt.Sprintf("[%s] %s", formattedDate, title), "INTEGRATION TESTING RESULT")
	if err != nil {
		log.Error("Error when Create Page via Confluence API: ", err)
		return err
	}
	if resp.Status() == "400 Bad Request" {
		var errResponse req.ErrorResponse
		err := json.Unmarshal(resp.Body(), &errResponse)
		if err != nil {
			log.Error("Error parsing JSON:", err)
			return err
		}
		if errResponse.Errors[0].Title == "A page with this title already exists: A page already exists with the same TITLE in this space" {
			respPages, err := http_request.GetConfluencePagesByTitle(fmt.Sprintf("[%s] %s", formattedDate, title))
			if err != nil {
				log.Error("Error when Get Page via Confluence API: ", err)
				return err
			}
			var pagesInfo req.PagesInfo
			err = json.Unmarshal(respPages.Body(), &pagesInfo)
			if err != nil {
				log.Error("Error parsing JSON:", err)
				return err
			}
			nextParentID = pagesInfo.Results[0].ID
		}
	} else if resp.Status() == "200 OK" {
		var createdResultPage req.CreatePageResult
		err := json.Unmarshal(resp.Body(), &createdResultPage)
		if err != nil {
			log.Error("Error parsing JSON: ", err)
			return err
		}
		nextParentID = createdResultPage.ID
	} else {
		log.Info(resp.Status())
		return fmt.Errorf("status Response from Confluence: %s", resp.Status())
	}

	//blockCode, err := helper.FormatForConfluenceCodeMacro(bodyValueFromFile)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}

	_, err = http_request.CreateConfluencePage(spaceId, nextParentID, fmt.Sprintf("[%s] %s", formattedDateTime, title), contentPage)
	if err != nil {
		log.Error("Error when Create Page via Confluence API: ", err)
		return err
	}

	// Check if file attachment is requested
	filePath := c.String("file")
	if filePath != "" {
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Error("File does not exist:", filePath)
			return fmt.Errorf("file does not exist: %s", filePath)
		}

		// Get the created child page ID for attachment
		// Use a more specific search to avoid race conditions
		childPageTitle := fmt.Sprintf("[%s] %s", formattedDateTime, title)
		respPages, err := http_request.GetConfluencePagesByTitle(childPageTitle)
		if err != nil {
			log.Error("Error when Get Page via Confluence API: ", err)
			return err
		}

		var pagesInfo req.PagesInfo
		err = json.Unmarshal(respPages.Body(), &pagesInfo)
		if err != nil {
			log.Error("Error parsing JSON: ", err)
			return err
		}

		// Find the most recently created page with the exact title
		var createdPageId string
		found := false
		for _, page := range pagesInfo.Results {
			if page.Title == childPageTitle {
				createdPageId = page.ID
				found = true
				break
			}
		}

		if !found {
			log.Error("Could not find created page for attachment")
			return fmt.Errorf("could not find created page for attachment")
		}

		// Upload attachment to the created child page
		log.Info("Uploading attachment to created page:", createdPageId)
		log.Info("Uploading attachment:", filePath)
		uploadResp, err := http_request.UploadConfluenceAttachment(createdPageId, filePath)
		if err != nil {
			log.Error("Error when uploading file via Confluence API:", err)
			return err
		}

		if uploadResp.StatusCode() >= 200 && uploadResp.StatusCode() < 300 {
			log.Info("File uploaded successfully as attachment!")
		} else {
			log.Error("Upload failed with status:", uploadResp.Status())
			log.Error("Response body:", string(uploadResp.Body()))
			return fmt.Errorf("upload failed with status: %s", uploadResp.Status())
		}
	}

	log.Info("Create Pages Successfully")
	return nil
}
