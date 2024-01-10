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
