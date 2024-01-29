package confluence_actions

import (
	"confluence_cli/helper/http_request"
	"confluence_cli/log"
	"confluence_cli/model/req"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"time"
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
	} else {
		return fmt.Errorf("please provide --body-value xxxx")
	}

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
	log.Info("Create Pages Successfully")
	return nil
}
