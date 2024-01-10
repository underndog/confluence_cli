package req

import "time"

type ErrorResponse struct {
	Errors []struct {
		Status int    `json:"status"`
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail any    `json:"detail"`
	} `json:"errors"`
}

type PagesInfo struct {
	Results []struct {
		ParentType string    `json:"parentType"`
		CreatedAt  time.Time `json:"createdAt"`
		AuthorID   string    `json:"authorId"`
		ID         string    `json:"id"`
		Version    struct {
			Number    int       `json:"number"`
			Message   string    `json:"message"`
			MinorEdit bool      `json:"minorEdit"`
			AuthorID  string    `json:"authorId"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"version"`
		Position int    `json:"position"`
		Title    string `json:"title"`
		Status   string `json:"status"`
		Body     struct {
		} `json:"body"`
		ParentID    string `json:"parentId"`
		SpaceID     string `json:"spaceId"`
		OwnerID     string `json:"ownerId"`
		LastOwnerID any    `json:"lastOwnerId"`
		Links       struct {
			Editui string `json:"editui"`
			Webui  string `json:"webui"`
			Tinyui string `json:"tinyui"`
		} `json:"_links"`
	} `json:"results"`
	Links struct {
	} `json:"_links"`
}

type CreatePageResult struct {
	ParentType string    `json:"parentType"`
	CreatedAt  time.Time `json:"createdAt"`
	AuthorID   string    `json:"authorId"`
	ID         string    `json:"id"`
	Version    struct {
		Number    int       `json:"number"`
		Message   string    `json:"message"`
		MinorEdit bool      `json:"minorEdit"`
		AuthorID  string    `json:"authorId"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"version"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	LastOwnerID any    `json:"lastOwnerId"`
	Body        struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
	ParentID string `json:"parentId"`
	SpaceID  string `json:"spaceId"`
	OwnerID  any    `json:"ownerId"`
	Position int    `json:"position"`
	Links    struct {
		Editui string `json:"editui"`
		Webui  string `json:"webui"`
		Tinyui string `json:"tinyui"`
	} `json:"_links"`
}
