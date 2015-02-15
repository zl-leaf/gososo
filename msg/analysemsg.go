package msg

type AnalyseOrderMsg struct {
	URL string `json:"url"`
	StatusCode int `json:"StatusCode"`
	Path string `json:"path"`
}