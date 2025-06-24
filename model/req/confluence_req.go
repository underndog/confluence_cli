package req

type Body struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type Page struct {
	SpaceID  string `json:"spaceId"`
	ParentId string `json:"parentId"`
	Status   string `json:"status"`
	Title    string `json:"title"`
	Body     Body   `json:"body"`
}

type UpdatePageBody struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type UpdatePageVersion struct {
	Number  int    `json:"number"`
	Message string `json:"message"`
}

type UpdatePagePayload struct {
	ID      string            `json:"id"`
	Status  string            `json:"status"`
	Title   string            `json:"title"`
	Body    UpdatePageBody    `json:"body"`
	Version UpdatePageVersion `json:"version"`
}
