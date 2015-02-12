package msg

type DownloadResultMsg struct {
	URL string `json:"url"`
	Path string `json:"path"`
	Redirects []string `json:"redirects"`
}