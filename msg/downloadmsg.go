package msg

type DownloadResultMsg struct {
	URL string `json:"url"`
	StatusCode int `json:"StatusCode"`
	Path string `json:"path"`
	Redirects []string `json:"redirects"`
}