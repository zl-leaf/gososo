package msg

type SearchResultMsg struct {
	Result int `json:"result"`
	Msg string `json:"msg"`
	Data []*SearchResultObj `json:"data"`
}

type SearchResultObj struct {
	URL string `json:"url"`
	Title string `json:"title"`
	Description string `json:"description"`
}