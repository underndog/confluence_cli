package http_request

import (
	"confluence_cli/helper"
	"confluence_cli/log"
	"confluence_cli/model/req"
	"github.com/go-resty/resty/v2"
)

type AuthSuccess struct {
	/* variables */
}
type AuthError struct {
	/* variables */
}

func CreateConfluencePage(spaceID, parentId, title, bodyValue string) (*resty.Response, error) {
	// Create a Resty Client
	client := resty.New()
	body := req.Page{
		SpaceID:  spaceID,
		ParentId: parentId,
		Status:   "current",
		Title:    title,
		Body: req.Body{
			Representation: "storage",
			Value:          bodyValue,
		},
	}
	email := helper.GetEnvOrDefault("EMAIL", "dc.nim94@gmail.com")
	apiToken := helper.GetEnvOrDefault("API_TOKEN", "ATATT3xFfGF0HhLy35h_ScihWNjoT8cABSw2f-GdCapur0eUVD1T54Ll1F4gLQhcAmpudC54lCBpDrvku68RTxor_4idlGEADw2TAzZqwYjs8a9iI2OhobhTILl5G7KX0lsEcbNXUj4Rq8EtYtZqBLIzH6WAhc10jn6jmx2TEFKVcc1a2g2GbiA=0A85027B")

	// POST JSON string
	// No need to set content type, if you have client level setting
	resp, err := client.R().
		SetBasicAuth(email, apiToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		Post(helper.GetEnvOrDefault("CONFLUENCE_URL", "https://nimtechnology.atlassian.net") + "/wiki/api/v2/pages")
	if err != nil {
		log.Error(err)
	}
	return resp, nil
}

func GetConfluencePagesByTitle(title string) (*resty.Response, error) {
	// Create a Resty Client
	client := resty.New()

	email := helper.GetEnvOrDefault("EMAIL", "dc.nim94@gmail.com")
	apiToken := helper.GetEnvOrDefault("API_TOKEN", "ATATT3xFfGF0HhLy35h_ScihWNjoT8cABSw2f-GdCapur0eUVD1T54Ll1F4gLQhcAmpudC54lCBpDrvku68RTxor_4idlGEADw2TAzZqwYjs8a9iI2OhobhTILl5G7KX0lsEcbNXUj4Rq8EtYtZqBLIzH6WAhc10jn6jmx2TEFKVcc1a2g2GbiA=0A85027B")

	resp, err := client.R().
		SetBasicAuth(email, apiToken).
		SetQueryParams(map[string]string{
			"title": title,
		}).
		SetHeader("Accept", "application/json").
		Get(helper.GetEnvOrDefault("CONFLUENCE_URL", "https://nimtechnology.atlassian.net") + "/wiki/api/v2/pages")
	if err != nil {
		log.Error(err)
	}
	return resp, err
}
